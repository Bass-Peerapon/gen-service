/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package subcreate

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/Bass-Peerapon/gen-service/json2struct"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

var (
	modelName string
	fileName  string
)

const (
	MODEL_PATH = "./models/"
	ORM_PATH   = "./orm/"
)

type Variable string

var (
	typeUUID       Variable = "uuid"
	typeZeroUUID   Variable = "zerouuid"
	typeString     Variable = "string"
	typeInt32      Variable = "int32"
	typeInt64      Variable = "int64"
	typeFloat64    Variable = "float64"
	typeTimeStamp  Variable = "timestamp"
	typeDate       Variable = "date"
	typeZeroString Variable = "zerostring"
	typeZeroInt    Variable = "zeroint"
	typeZeroFloat  Variable = "zerofloat"
	typeZeroBool   Variable = "zerobool"
	typeDuration   Variable = "duration"
	TypeBool       Variable = "bool"
	TypeJson       Variable = "json"
	TypeInterface  Variable = "interface"
)

var (
	TYPE_MAP = map[Variable]string{
		typeUUID:       "*uuid.UUID",
		typeTimeStamp:  "*helperModel.Timestamp",
		typeZeroUUID:   "zero.UUID",
		typeString:     "string",
		typeInt32:      "int32",
		typeInt64:      "int64",
		typeFloat64:    "float64",
		typeDate:       "date",
		typeZeroString: "zero.string",
		typeZeroInt:    "zero.Int",
		typeZeroFloat:  "zero.Float",
		typeZeroBool:   "zero.Bool",
		typeDuration:   "duration",
		TypeBool:       "bool",
		TypeJson:       "map[string]interface{}",
		TypeInterface:  "interface{}",
	}
)

type Model struct {
	Name           string
	CamelCase      string
	LowerCamelCase string
	modelPath      string
	ormPath        string
}

type Table struct {
	Feild map[string]*Feild
	order []string
}

type Feild struct {
	Name string
	Type string
	Tag  string
}

func NewTable(name string) *Table {
	var t string
	t = fmt.Sprintf("`json:%s db:%s pk:%s`", `"-"`, strconv.Quote(name), strconv.Quote("Id"))
	return &Table{
		Feild: map[string]*Feild{"table_name": {Name: "TableName", Type: "struct{}", Tag: t}},
		order: []string{"table_name"},
	}
}

func (t *Table) NewCol(name string, typ Variable) {
	var tag string
	tag = fmt.Sprintf("`json:%s db:%s type:%s`", strconv.Quote(name), strconv.Quote(name), strconv.Quote(string(typ)))
	t.Feild[name] = &Feild{
		Name: strcase.ToCamel(name),
		Type: TYPE_MAP[typ],
		Tag:  tag,
	}
	t.order = append(t.order, name)
}

func (f Feild) PrintFeild() string {
	return fmt.Sprintf(`%s %s %s`, f.Name, f.Type, f.Tag)
}

func (t Table) PrintTable() string {
	var vals []string
	for _, key := range t.order {
		vals = append(vals, t.Feild[key].PrintFeild())
	}
	return strings.Join(vals, "\n")
}

func NewModel(name string) *Model {
	return &Model{
		Name:           name,
		CamelCase:      strcase.ToCamel(name),
		LowerCamelCase: strcase.ToLowerCamel(name),
		modelPath:      MODEL_PATH + name,
		ormPath:        ORM_PATH + name,
	}
}

func (s Model) GetModelPath() string {
	return s.modelPath
}

func (s Model) GetOrmPath() string {
	return s.ormPath
}

func generateModel(model *Model) error {
	fmt.Printf("start generate %s model...\n", model.Name)
	table := NewTable(model.Name)
	table.NewCol("id", typeUUID)
	table.NewCol("created_at", typeTimeStamp)
	table.NewCol("updated_at", typeTimeStamp)
	tmp := fmt.Sprintf(`
	package models
	import (
		"github.com/gofrs/uuid"
		helperModel "git.innovasive.co.th/backend/models"
	)	
	type {{.CamelCase}} struct {
		%s
	}
	`,
		table.PrintTable(),
	)
	if err := os.MkdirAll(MODEL_PATH, os.ModePerm); err != nil {
		return err
	}
	fn := model.GetModelPath() + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("models").Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, model)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	fmt.Println("success!!")
	return nil
}

func generateOrm(model *Model) error {
	fmt.Printf("start generate %s orm...\n", model.Name)
	tmp := `
	package orm	
	func Orm{{.CamelCase}}(ptr *models.{{.CamelCase}}, currentRow RowValue) (*models.{{.CamelCase}}, error) {
		v, err := fillValue(ptr, currentRow)
		if v != nil {
			return v.(*models.{{.CamelCase}}), nil
		}

		return nil, err
	}
	`
	if err := os.MkdirAll(ORM_PATH, os.ModePerm); err != nil {
		return err
	}
	fn := model.GetOrmPath() + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("orm").Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, model)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, 0644); err != nil {
		return err
	}

	fmt.Println("success!!")

	return nil
}

// modelCmd represents the model command
var ModelCmd = &cobra.Command{
	Use:   "model",
	Short: "generate model and generate orm",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		model := NewModel(modelName)
		if modelName != "" && fileName != "" {
			if err := Json2Struct(model.CamelCase, fileName, MODEL_PATH+modelName+SER_NAME_GO); err != nil {
				fmt.Println(err)
				return
			}
			if err := generateOrm(model); err != nil {
				fmt.Println(err)
				return
			}
			return
		}
		if modelName != "" {
			if err := generateModel(model); err != nil {
				fmt.Println(err)
				return
			}
			if err := generateOrm(model); err != nil {
				fmt.Println(err)
				return
			}
			return
		}

	},
}

func init() {
	ModelCmd.Flags().StringVarP(&modelName, "model_name", "n", "", "model name ")
	if err := ModelCmd.MarkFlagRequired("model_name"); err != nil {
		fmt.Println(err)
	}
	ModelCmd.Flags().StringVarP(&fileName, "file_json", "f", "", "file json ")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// modelCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// modelCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func Json2Struct(modelName, fileName string, outputName string) error {
	fmt.Println("start generate model form json")
	if err := os.MkdirAll(MODEL_PATH, os.ModePerm); err != nil {
		return err
	}
	var input io.Reader
	input = os.Stdin
	f, err := os.Open(fileName)
	if err != nil {
		return fmt.Errorf("reading input file: %s", err)
	}
	defer f.Close()
	input = f

	output, err := json2struct.Generate(input, modelName, "models")
	if err != nil {
		fmt.Fprintln(os.Stderr, "error parsing", err)
		return err
	}
	if err := ioutil.WriteFile(outputName, output, 0644); err != nil {
		return err

	}
	fmt.Println("success!!")
	return nil
}
