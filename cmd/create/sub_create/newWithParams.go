/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package subcreate

import (
	"bufio"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"regexp"
	"strings"

	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

// newWithParamsCmd represents the newWithParams command
var NewWithParamsCmd = &cobra.Command{
	Use:   "newWithParams",
	Short: "generate function new with params",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		if err := Generator(fileName); err != nil {
			fmt.Println(err)
			return
		}
	},
}

func init() {
	NewWithParamsCmd.Flags().StringVarP(&fileName, "file", "f", "", "path file golang model")
	if err := NewWithParamsCmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newWithParamsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newWithParamsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

const (
	v    = "val"
	valM = "valM"

	SUB_TYPE_STRING = "string"
	SUB_TYPE_TIME   = "time"
)

func WriteHeader(w *bufio.Writer, f *ast.File) error {
	if f != nil {

		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						switch st := spec.Type.(type) {
						case *ast.StructType:
							for _, field := range st.Fields.List {
								switch field.Type.(type) {
								case *ast.ArrayType:
									arr := field.Type.(*ast.ArrayType)
									if arr.Lbrack.IsValid() {
										switch arr.Elt.(type) {
										case *ast.StarExpr:
											fmt.Fprintf(w, "\t\t%v := make([]*%v, 0)\n", field.Names[0], field.Type.(*ast.ArrayType).Elt.(*ast.StarExpr).X)
										case *ast.Ident:
											fmt.Fprintf(w, "\t\t%v := []%v{}\n", field.Names[0], field.Type.(*ast.ArrayType).Elt)

										}
									}

								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func WriteTail(w *bufio.Writer, f *ast.File) error {
	if f != nil {

		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						switch st := spec.Type.(type) {
						case *ast.StructType:
							for _, field := range st.Fields.List {
								switch field.Type.(type) {
								case *ast.ArrayType:
									arr := field.Type.(*ast.ArrayType)
									if arr.Lbrack.IsValid() {
										fmt.Fprintf(w, "\tptr.%v = %v\n", field.Names[0], field.Names[0])
									}

								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}

func Generator(fileName string) error {
	fmt.Println("start generate new with params")
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		return err
	}
	if f != nil {
		fo, _ := os.ReadFile(fileName)
		file, _ := os.Create(fileName)
		w := bufio.NewWriter(file)
		fmt.Fprintf(w, string(fo))
		reg := regexp.MustCompile(`\S+`)

		for val := range f.Scope.Objects {
			for _, decl := range f.Decls {
				switch decl := decl.(type) {
				case *ast.GenDecl:
					if decl.Tok == token.IMPORT {
						continue
					}
					fmt.Fprintf(w, "\n")
					fmt.Fprintf(w, "func New%vWithParams(params map[string]interface{}, ptr *%v) *%v {\n", val, val, val)
					fmt.Fprintf(w, "\tif ptr == nil {\n")
					fmt.Fprintf(w, "\t\tptr = new(%v)\n", val)
					fmt.Fprintf(w, "\t}\n")
					fmt.Fprintf(w, "\n")
					if err := WriteHeader(w, f); err != nil {
						return err
					}
					fmt.Fprintf(w, "\n")
					fmt.Fprintf(w, "\tfor key, val := range params {\n")
					fmt.Fprintf(w, "\t    switch key {\n")
					for _, spec := range decl.Specs {
						switch spec := spec.(type) {
						case *ast.TypeSpec:
							switch st := spec.Type.(type) {
							case *ast.StructType:
								for _, field := range st.Fields.List {
									tag := field.Tag.Value
									tags := reg.FindAllString(tag, 1)
									switch field.Type.(type) {
									case *ast.StarExpr:
										WriteStarExpr(w, tags[0], field)
									case *ast.ArrayType:
										arr := field.Type.(*ast.ArrayType)
										if arr.Lbrack.IsValid() {
											switch arr.Elt.(type) {
											case *ast.StarExpr:
												WriteArrayStarExpr(w, tags[0], arr, field)
											case *ast.Ident:
												WriteArrayIdent(w, tags[0], field)
											}
										}

									case *ast.Ident:
										for _, name := range field.Names {
											WriteIdent(w, tags[0], field, name)
										}
									}
								}
							}
						}
					}
				}
			}
			fmt.Fprintf(w, "\t}\n")
			fmt.Fprintf(w, "}\n")
			fmt.Fprintf(w, "\n")
			if err := WriteTail(w, f); err != nil {
				return err
			}
			fmt.Fprintf(w, "\n")
			fmt.Fprintf(w, "\treturn ptr\n")
			fmt.Fprintf(w, "}\n")
			// Flush.

			break
		}
		w.Flush()
	}
	code, err := imports.Process(fileName, nil, &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fileName, code, 0644); err != nil {
		return err
	}
	fmt.Println("success")
	return nil
}

func WriteStarExpr(w *bufio.Writer, tag string, field *ast.Field) {
	fmt.Fprintf(w, "\t    "+`case %v:`+"\n", strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(tag, `json:`, ""), "`"), "`"))
	fmt.Fprintf(w, "\tif %s != nil {\n", v)
	star := field.Type.(*ast.StarExpr)
	switch star.X.(type) {
	case *ast.SelectorExpr:
		switch fmt.Sprintf("%v", star.X.(*ast.SelectorExpr).X) {
		case "uuid":
			fmt.Fprintf(w, "\t\tptr.%v, _ = %s\n", field.Names[0].Name, GetCastTypeFunc(star.X.(*ast.SelectorExpr).X, v))
		case "helperModels":
			fallthrough
		case "helperModel":
			fmt.Fprintf(w, "\t\t\tif reflect.TypeOf(%s).Kind() == reflect.String {\n", v)
			fmt.Fprintf(w, "\t\t\tti := %s\n", GetCastTypeFunc(star.X.(*ast.SelectorExpr).Sel, v, SUB_TYPE_STRING))
			fmt.Fprintf(w, "\t\t\tptr.%v = &ti\n", field.Names[0])
			fmt.Fprintf(w, "\t\t\t"+`}else if reflect.TypeOf(%s).String() == "time.Time" {`+"\n", v)
			fmt.Fprintf(w, "\t\t\tti := %s\n", GetCastTypeFunc(star.X.(*ast.SelectorExpr).Sel, v, SUB_TYPE_TIME))
			fmt.Fprintf(w, "\t\t\tptr.%v = &ti\n", field.Names[0])
			fmt.Fprintf(w, "\t\t\t}\n")
		}
	case *ast.Ident:
		fmt.Fprintf(w, "\t\tptr.%v = New%vWithParams(val.(map[string]interface{}), nil)\n", field.Names[0], star.X.(*ast.Ident).Name)
	}
	fmt.Fprintf(w, "\t\t}\n")
}

func WriteArrayStarExpr(w *bufio.Writer, tag string, arr *ast.ArrayType, field *ast.Field) {
	fmt.Fprintf(w, "\t    "+`case %v:`+"\n", strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(tag, `json:`, ""), "`"), "`"))
	fmt.Fprintf(w, "\tif %s != nil && len(%s.([]interface{})) > 0 {\n", v, v)
	fmt.Fprintf(w, "\t\tfor _, %s := range %s.([]interface{}) {\n", valM, v)
	fmt.Fprintf(w, "\t\t\t"+`%v := New%vWithParams(%s.(map[string]interface{}), nil)`+"\n", arr.Elt.(*ast.StarExpr).X, arr.Elt.(*ast.StarExpr).X, valM)
	fmt.Fprintf(w, "\t\t\tif %v != nil {\n", arr.Elt.(*ast.StarExpr).X)
	fmt.Fprintf(w, "\t\t\t\t"+`%v = append(%v, %v)`+"\n", field.Names[0], field.Names[0], arr.Elt.(*ast.StarExpr).X)
	fmt.Fprintf(w, "\t\t\t}\n")
	fmt.Fprintf(w, "\t\t\t\t}\n")
	fmt.Fprintf(w, "\t\t}\n")
}

func WriteArrayIdent(w *bufio.Writer, tag string, field *ast.Field) {
	fmt.Fprintf(w, "\t    "+`case %v:`+"\n", strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(tag, `json:`, ""), "`"), "`"))
	fmt.Fprintf(w, "\tif %s != nil && len(%s.([]interface{})) > 0 {\n", v, v)
	fmt.Fprintf(w, "\t\tfor _, %s := range %s.([]interface{}) {\n", valM, v)
	fmt.Fprintf(w, "\t\t\t"+`%v = append(%v, %s)`+"\n", field.Names[0], field.Names[0], GetCastTypeFunc(field.Type.(*ast.ArrayType).Elt, valM))
	fmt.Fprintf(w, "\t\t\t}\n")
	fmt.Fprintf(w, "\t\t}\n")
}

func WriteIdent(w *bufio.Writer, tag string, field *ast.Field, name *ast.Ident) {
	fmt.Fprintf(w, "\t    "+`case %v:`+"\n", strings.TrimRight(strings.TrimLeft(strings.ReplaceAll(tag, `json:`, ""), "`"), "`"))
	switch fmt.Sprintf("%v", field.Type) {
	case "string":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))
	case "int64":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))
	case "float64":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))
	case "int32":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))
	case "int":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))
	case "bool":
		fmt.Fprintf(w, "\t\t"+`ptr.%v = %s`+"\n", name.Name, GetCastTypeFunc(field.Type, v))

	}
}

func GetCastTypeFunc(expr ast.Expr, name string, subType ...string) string {
	var t string
	switch fmt.Sprintf("%v", expr) {
	case "string":
		t = fmt.Sprintf("cast.ToString(%s)", name)
	case "int64":
		t = fmt.Sprintf("cast.ToInt64(%s)", name)
	case "float64":
		t = fmt.Sprintf("cast.ToFloat64(%s)", name)
	case "int32":
		t = fmt.Sprintf("cast.ToInt32(%s)", name)
	case "int":
		t = fmt.Sprintf("cast.ToInt(%s)", name)
	case "bool":
		t = fmt.Sprintf("cast.ToBool(%s)", name)
	case "uuid":
		t = fmt.Sprintf("helper.ConvertToUUIDAndBinary(%s)", name)
	case "Timestamp":
		if len(subType) > 0 {
			if subType[0] == "string" {
				t = fmt.Sprintf("helperModel.NewTimestampFromString(%s.(string))", name)
			} else if subType[0] == "time" {
				t = fmt.Sprintf("helperModel.NewTimestampFromTime(%s.(time.Time))", name)
			}
		}

	}

	return t
}
