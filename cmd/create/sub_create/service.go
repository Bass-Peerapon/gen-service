/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package subcreate

import (
	"bytes"
	"fmt"
	"html/template"
	"os"
	"strconv"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

// serviceCmd represents the service command
var ServiceCmd = &cobra.Command{
	Use:   "service",
	Short: "generate service",
	Long:  ``,
	Run: func(_ *cobra.Command, _ []string) {
		fmt.Println("start create service" + serviceName)
		Generate(serviceName)
		fmt.Println("created service " + serviceName)

	},
}

func init() {
	ServiceCmd.Flags().StringVarP(&serviceName, "service_name", "n", "", "service name ")
	if err := ServiceCmd.MarkFlagRequired("service_name"); err != nil {
		fmt.Println(err)
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serviceCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serviceCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

type Service struct {
	Name           string
	CamelCase      string
	LowerCamelCase string
	defultPath     string
}

func NewService(name string) *Service {
	return &Service{
		Name:           name,
		CamelCase:      strcase.ToCamel(name),
		LowerCamelCase: strcase.ToLowerCamel(name),
		defultPath:     DEFULT_PATH + name,
	}
}

func (s *Service) SetDefultPath() {
	s.defultPath = DEFULT_PATH + s.Name
}

func (s Service) GetDefultPath() string {
	return s.defultPath
}

const (
	REPO        = "repository"
	USECASE     = "usecase"
	HTTP        = "http"
	HANDLER     = "handler"
	VALIDATOR   = "validator"
	DEFULT_PATH = "./service/"
	SER_NAME_GO = ".go"
)

const (
	FILE_MAIN = "./main.go"
)

var (
	GOPATH       = os.Getenv("GOPATH")
	SERVICE_PATH = "/service/"
)

var tmpRepoAdapter = `
package {{.Name}}
type {{ .CamelCase }}Repository interface {

}
`

var tmpUsecaseAdapter = `
package {{.Name}}
type {{ .CamelCase }}Usecase interface {

}
`

var tmpHttpAdapter = `
package {{.Name}} 
type {{ .CamelCase }}Handler interface {

}
`

var tmpRepo = `
package repository
type {{.LowerCamelCase}}Repository struct {
	client             *psql.Client
}

func New{{.LowerCamelCase}}Repository(client *psql.Client) {{.Name}}.{{.CamelCase}}Repository {
	return &{{.LowerCamelCase}}Repository{
		client:             client,	
	}
}
`
var tmpUsecase = `
package usecase
type {{.LowerCamelCase}}Usecase struct {
	{{.LowerCamelCase}}Repo {{.Name}}.{{.CamelCase}}Repository
}

func New{{.LowerCamelCase}}Usecase({{.LowerCamelCase}}Repo {{.Name}}.{{.CamelCase}}Repository) {{.Name}}.{{.CamelCase}}Usecase {
	return &{{.LowerCamelCase}}Usecase{
		{{.LowerCamelCase}}Repo : {{.LowerCamelCase}}Repo,	
	}
}
`
var tmpHttp = `
package http
type {{.LowerCamelCase}}Handler struct {
	{{.LowerCamelCase}}Us {{.Name}}.{{.CamelCase}}Usecase
}

func New{{.LowerCamelCase}}Handler({{.LowerCamelCase}}Us {{.Name}}.{{.CamelCase}}Usecase) {{.Name}}.{{.CamelCase}}Repository {
	return &{{.LowerCamelCase}}Handler{
		{{.LowerCamelCase}}Us:   {{.LowerCamelCase}}Us,	
	}
}
`

var tmpVal = `
package validator
type Validation struct{

}
`

var (
	serviceName string
)

func Generate(serviceName string) error {
	s := NewService(serviceName)

	if err := s.generateServiceDir(); err != nil {
		return err
	}
	if err := s.generateReposiroryAdapter(); err != nil {
		return err
	}
	if err := s.generateUsecaseAdapter(); err != nil {
		return err
	}
	if err := s.generateHandlerAdapter(); err != nil {
		return err
	}
	if err := s.generateHandler(); err != nil {
		return err
	}
	if err := s.generateUsecase(); err != nil {
		return err
	}
	if err := s.generateReposirory(); err != nil {
		return err
	}

	if err := s.generateValidator(); err != nil {
		return err
	}

	if err := InjectServiceInMain(s); err != nil {

	}

	return nil
}

func (s Service) generateServiceDir() error {
	if err := os.MkdirAll(s.GetDefultPath(), os.ModePerm); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/repository/{serivce_name}_repository.go
func (s Service) generateReposirory() error {
	dir := "./" + s.GetDefultPath() + "/" + REPO
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + REPO + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("repository").Parse(tmpRepo)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/usecase/{serivce_name}_usecase.go
func (s Service) generateUsecase() error {
	dir := "./" + s.GetDefultPath() + "/" + USECASE
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + USECASE + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("usecase").Parse(tmpUsecase)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/http/{serivce_name}_handler.go
func (s Service) generateHandler() error {
	dir := "./" + s.GetDefultPath() + "/" + HTTP
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + s.Name + "_" + HTTP + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("http").Parse(tmpHttp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// ./service/{service_name}/repository.go
func (s Service) generateReposiroryAdapter() error {
	fn := s.GetDefultPath() + "/" + REPO + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("repository").Parse(tmpRepoAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/usecase.go
func (s Service) generateUsecaseAdapter() error {
	fn := s.GetDefultPath() + "/" + USECASE + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("usecase").Parse(tmpUsecaseAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}

	return nil
}

// ./service/{service_name}/http.go
func (s Service) generateHandlerAdapter() error {
	fn := s.GetDefultPath() + "/" + HTTP + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("http").Parse(tmpHttpAdapter)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func (s Service) generateValidator() error {
	dir := "./" + s.GetDefultPath() + "/" + VALIDATOR
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	fn := dir + "/" + VALIDATOR + SER_NAME_GO
	var buf bytes.Buffer
	t, err := template.New("validator").Parse(tmpVal)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, s)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fn, code, os.ModePerm); err != nil {
		return err
	}

	return nil
}

func InjectServiceInMain(s *Service) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	fmt.Println("-- script import main.go--")
	dir = strings.Replace(dir, GOPATH+"/src/", "", -1)
	repo_import_alias := fmt.Sprintf("_%s_%s", s.Name, REPO)
	us_import_alias := fmt.Sprintf("_%s_%s", s.Name, USECASE)
	http_import_alias := fmt.Sprintf("_%s_%s", s.Name, HTTP)
	validator_import_alias := fmt.Sprintf("_%s_%s", s.Name, VALIDATOR)
	repo_import := fmt.Sprintf("%s %s", repo_import_alias, strconv.Quote(dir+SERVICE_PATH+s.Name+"/"+REPO))
	us_import := fmt.Sprintf("%s %s", us_import_alias, strconv.Quote(dir+SERVICE_PATH+s.Name+"/"+USECASE))
	http_import := fmt.Sprintf("%s %s", http_import_alias, strconv.Quote(dir+SERVICE_PATH+s.Name+"/"+HTTP))
	validate_import := fmt.Sprintf("%s %s", validator_import_alias, strconv.Quote(dir+SERVICE_PATH+s.Name+"/"+VALIDATOR))
	fmt.Println(repo_import)
	fmt.Println(us_import)
	fmt.Println(http_import)
	fmt.Println(validate_import)
	fmt.Println()
	fmt.Println("-- script inject repo --")
	fmt.Printf("%sRepo := %s.New%sRepository()\n", s.LowerCamelCase, repo_import_alias, s.CamelCase)
	fmt.Println("-- script inject usecase --")
	fmt.Printf("%sUs := %s.New%sUsecase(%sRepo)\n", s.LowerCamelCase, repo_import_alias, s.CamelCase, s.LowerCamelCase)
	fmt.Println("-- script inject http --")
	fmt.Printf("%sHandler := %s.New%sHandler(%sUs)\n", s.LowerCamelCase, repo_import_alias, s.CamelCase, s.LowerCamelCase)
	fmt.Println()
	return nil
}
