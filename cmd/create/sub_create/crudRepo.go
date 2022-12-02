/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>
*/
package subcreate

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"os"
	"regexp"
	"strings"

	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
	"golang.org/x/tools/imports"
)

var (
	outputFileName string
)

// crudRepoCmd represents the crudRepo command
var CrudRepoCmd = &cobra.Command{
	Use:   "crudRepo",
	Short: "generate crud in repository",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		all, err := cmd.Flags().GetBool("all")
		if err != nil {
			fmt.Println(err)
			return
		}
		if all {
			if err := genRead(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}
			if err := genCreate(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}
			if err := genUpdate(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}
			if err := genDelete(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}

			return
		}
		flagRead, err := cmd.Flags().GetBool("read")
		if err != nil {
			fmt.Println(err)
			return
		}
		if flagRead {
			if err := genRead(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}

		}
		flagCreate, err := cmd.Flags().GetBool("create")
		if err != nil {
			fmt.Println(err)
			return
		}
		if flagCreate {
			if err := genCreate(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}

		}
		flagUpdate, err := cmd.Flags().GetBool("update")
		if err != nil {
			fmt.Println(err)
			return
		}
		if flagUpdate {
			if err := genUpdate(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}

		}
		flagDelete, err := cmd.Flags().GetBool("delete")
		if err != nil {
			fmt.Println(err)
			return
		}
		if flagDelete {
			if err := genDelete(fileName, outputFileName); err != nil {
				fmt.Println(err)
			}
		}
		return
	},
}

func init() {
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// crudRepoCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	CrudRepoCmd.Flags().StringVarP(&fileName, "file", "f", "", "path file model")
	if err := CrudRepoCmd.MarkFlagRequired("file"); err != nil {
		fmt.Println(err)
	}
	CrudRepoCmd.Flags().StringVarP(&outputFileName, "output", "o", "", "path file repository")
	if err := CrudRepoCmd.MarkFlagRequired("output"); err != nil {
		fmt.Println(err)
	}
	CrudRepoCmd.Flags().BoolP("all", "a", false, "generate sql script (insert ,query ,update ,delete)")
	CrudRepoCmd.Flags().BoolP("create", "c", false, "generate sql script insert")
	CrudRepoCmd.Flags().BoolP("read", "r", false, "generate sql script query")
	CrudRepoCmd.Flags().BoolP("update", "u", false, "generate sql script update")
	CrudRepoCmd.Flags().BoolP("delete", "d", false, "generate sql script delete")
}

func getTag(s string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.TrimSpace(s), "db:", ""), `"`, "")
}

func genCreate(fn string, fnOut string) error {
	fmt.Println("Start generate create")
	m, tagsDB, shouldReturn, returnValue := readModel(fn)
	if shouldReturn {
		return returnValue
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := `
	func (p {{.OutPut.LowerCamelCase}}) Create{{.Input.CamelCase}}(ctx context.Context ,{{.Input.LowerCamelCase}} *models.{{.Input.CamelCase}}) error {
		tx, err := p.client.GetClient().Beginx()
		if err != nil {
			return err
		}
		
		if err := p.create{{.Input.CamelCase}}(ctx, tx, {{.Input.LowerCamelCase}}); err != nil {
			return err
		}
	
		return tx.Commit() 
	}
	`
	subTmp := `
	func (p {{.OutPut.LowerCamelCase}}) create{{.Input.CamelCase}}(ctx context.Context, tx *sqlx.Tx ,{{.Input.LowerCamelCase}} *models.{{.Input.CamelCase}}) error {
	sql :=` + "`" + `INSERT INTO {{snackCase .Input.Name}} ({{range $i, $a := .Params}}{{if $i}}, {{end}} "{{$a}}" {{end}})
	VALUES 
	(
		{{range $i, $a := .Params}}{{if $i}},
		{{end}}${{add $i 1}}{{end}}
	)
	` + "`" + `
	stmt, err := tx.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	{{$name := .Input.LowerCamelCase}}
	if _, err := stmt.ExecContext(
		ctx,{{range .Params}}
		{{$name}}.{{ camelCase .}},{{end}}
	); err != nil {
		return err
	}
	return nil
	}
	`
	tmp = strings.Join([]string{string(fo), tmp, subTmp}, "\n")

	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	var funcMap = template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("create").Funcs(funcMap).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, e)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}
	fmt.Println("success")
	return nil
}

func genRead(fn string, fnOut string) error {
	fmt.Println("Start generate read")
	m, tagsDB, shouldReturn, returnValue := readModel(fn)
	if shouldReturn {
		return returnValue
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := `
	func (p {{.OutPut.LowerCamelCase}}) Fetch{{.Input.CamelCase}}s(ctx context.Context,args *sync.Map , paginator *helperModel.Paginator ,{{.Input.LowerCamelCase}}s []*models.{{.Input.CamelCase}}) ([]*models.{{.Input.CamelCase}},error) {
		if args == nil {
			args = new(sync.Map)
		}
		var conds []string
		var valArgs []interface{}

		var where string
		if len(conds) > 0 {
		where = "WHERE " + strings.Join(conds, " AND ")
		}
		var paginatorSql string
		if paginator != nil {
		var limit = int(paginator.PerPage)
		var skipItem = (int(paginator.Page) - 1) * int(paginator.PerPage)
		paginatorSql = fmt.Sprintf(` + "`" + `
			LIMIT %d
			OFFSET %d
			` + "`" + `,
			limit,
			skipItem,
		)
	}
	sql := fmt.Sprintf(` + "`" + `
	SELECT
		%s,
		count(*) OVER() as total_row
	FROM
		{{snackCase .Input.Name}}
	%s
	%s	
	` + "`" + `,
		orm.GetSelector(models.{{.Input.CamelCase}}{}),
		where,
		paginatorSql,
	)
	sql = sqlx.Rebind(sqlx.DOLLAR, sql)
	myHelper.Println(sql)
	stmt, err := p.client.GetClient().PreparexContext(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Queryx(valArgs...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return p.orm{{.Input.CamelCase}}(ctx, rows, paginator, []string{})
	}
	`
	subTmp := `
	func (p {{.OutPut.LowerCamelCase}}) orm{{.Input.CamelCase}}(ctx context.Context, rows *sqlx.Rows, paginator *helperModel.Paginator, relationships []string) ([]*models.{{.Input.CamelCase}},error) {
		var ptrs = make([]*models.{{.Input.CamelCase}}, 0)
		var mapper, err = orm.NewRowsScan(rows)
		if err != nil {
			return nil, err
		}
	
		if mapper.TotalRows() > 0 {
			if paginator != nil {
				paginator.SetPaginatorByAllRows(mapper.PaginateTotal())
			}
			for _, row := range mapper.RowsValues() {
				var ptr = new(models.{{.Input.CamelCase}})
	
				ptr, err := orm.Orm{{.Input.CamelCase}}(ptr, row)
				if err != nil {
					return nil, err
				}
				if ptr != nil {
					exists, err := orm.IsDuplicateByPK(ptrs, ptr)
					if err != nil {
						return nil, err
					}
					if !exists {
						ptrs = append(ptrs, ptr)
					}
	
				}
			}
	
		}
	
		return ptrs, nil
	}
	`
	tmp = strings.Join([]string{string(fo), tmp, subTmp}, "\n")

	var buf bytes.Buffer
	var funcMap = template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
	}
	t, err := template.New("read").Funcs(funcMap).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, e)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}

func genUpdate(fn string, fnOut string) error {
	fmt.Println("Start generate update")
	m, tagsDB, shouldReturn, returnValue := readModel(fn)
	if shouldReturn {
		return returnValue
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := `
	func (p {{.OutPut.LowerCamelCase}}) Update{{.Input.CamelCase}}(ctx context.Context,id *uuid.UUID ,{{.Input.LowerCamelCase}} *models.{{.Input.CamelCase}}) error {
		tx, err := p.client.GetClient().Beginx()
		if err != nil {
			return err
		}
		
		if err := p.update{{.Input.CamelCase}}(ctx, tx, id, {{.Input.LowerCamelCase}}); err != nil {
			return err
		}
	
		return tx.Commit() 
	}
	`
	subTmp := `
	func (p {{.OutPut.LowerCamelCase}}) update{{.Input.CamelCase}}(ctx context.Context, tx *sqlx.Tx, id *uuid.UUID ,{{.Input.LowerCamelCase}} *models.{{.Input.CamelCase}}) error {
	sql :=` + "`" + `
		UPDATE 
			{{snackCase .Input.Name}}
		SET
			{{range $i, $a := .Params}}{{if $i}}, 
			{{end}} {{$a}} = ${{add $i 1}} {{end}}
		WHERE
			id = $1
			
	` + "`" + `
	stmt, err := tx.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	{{$name := .Input.LowerCamelCase}}
	if _, err := stmt.ExecContext(
		ctx,{{range .Params}}
		{{$name}}.{{ camelCase .}},{{end}}
	); err != nil {
		return err
	}
	return nil
	}
	`
	tmp = strings.Join([]string{string(fo), tmp, subTmp}, "\n")

	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	var funcMap = template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("update").Funcs(funcMap).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, e)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}

func genDelete(fn string, fnOut string) error {
	fmt.Println("Start generate delete")
	m, tagsDB, shouldReturn, returnValue := readModel(fn)
	if shouldReturn {
		return returnValue
	}

	mOut, err := getModelFromRepo(fnOut)
	if err != nil {
		return err
	}

	e := struct {
		Input  *Model
		Params []string
		OutPut *Model
	}{
		m,
		tagsDB,
		mOut,
	}
	fo, _ := os.ReadFile(fnOut)
	tmp := `
	func (p {{.OutPut.LowerCamelCase}}) Delete{{.Input.CamelCase}}(ctx context.Context,id *uuid.UUID) error {
		tx, err := p.client.GetClient().Beginx()
		if err != nil {
			return err
		}
		
		if err := p.delete{{.Input.CamelCase}}(ctx, tx, id); err != nil {
			return err
		}
	
		return tx.Commit() 
	}
	`
	subTmp := `
	func (p {{.OutPut.LowerCamelCase}}) delete{{.Input.CamelCase}}(ctx context.Context, tx *sqlx.Tx, id *uuid.UUID) error {
	sql :=` + "`" + `
		DELETE FROM 
			{{snackCase .Input.Name}}
		WHERE
			id = $1
			
	` + "`" + `
	stmt, err := tx.Preparex(sql)
	if err != nil {
		return err
	}
	defer stmt.Close()
	{{$name := .Input.LowerCamelCase}}
	if _, err := stmt.ExecContext(
		ctx,
		id,
	); err != nil {
		return err
	}
	return nil
	}
	`
	tmp = strings.Join([]string{string(fo), tmp, subTmp}, "\n")

	addFunc := func(x, y int) int {
		return x + y
	}
	var buf bytes.Buffer
	var funcMap = template.FuncMap{
		"camelCase": strcase.ToCamel,
		"snackCase": strcase.ToSnake,
		"add":       addFunc,
	}
	t, err := template.New("update").Funcs(funcMap).Parse(tmp)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, e)
	code, err := imports.Process("", buf.Bytes(), &imports.Options{
		Comments: true,
	})
	if err != nil {
		return err
	}

	if err := os.WriteFile(fnOut, code, os.ModeAppend); err != nil {
		return err
	}

	fmt.Println("success")
	return nil
}

func readModel(fn string) (*Model, []string, bool, error) {
	var m *Model
	var tagsDB []string

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		return nil, nil, true, err
	}
	reg := regexp.MustCompile(`db:"(.+)" `)
	for range f.Scope.Objects {
		for _, decl := range f.Decls {
			switch decl := decl.(type) {
			case *ast.GenDecl:
				if decl.Tok == token.IMPORT {
					continue
				}
				for _, spec := range decl.Specs {
					switch spec := spec.(type) {
					case *ast.TypeSpec:
						m = NewModel(spec.Name.String())
						switch st := spec.Type.(type) {
						case *ast.StructType:
							for _, field := range st.Fields.List {
								tag := field.Tag.Value
								tags := reg.FindAllString(tag, 1)
								switch field.Type.(type) {
								case *ast.StarExpr:
									t := getTag(tags[0])
									if t != "-" && t != "" {
										tagsDB = append(tagsDB, t)
									}
								case *ast.ArrayType:
									arr := field.Type.(*ast.ArrayType)
									if arr.Lbrack.IsValid() {
										switch arr.Elt.(type) {
										case *ast.StarExpr:
											t := getTag(tags[0])
											if t != "-" && t != "" {
												tagsDB = append(tagsDB, t)
											}
										case *ast.Ident:
											t := getTag(tags[0])
											if t != "-" && t != "" {
												tagsDB = append(tagsDB, t)
											}
										}
									}

								case *ast.Ident:
									t := getTag(tags[0])
									if t != "-" && t != "" {
										tagsDB = append(tagsDB, t)
									}
								}
							}
						}
					}
				}
			}
		}
		break
	}
	return m, tagsDB, false, nil
}

func getModelFromRepo(fn string) (*Model, error) {
	var m *Model

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fn, nil, 0)
	if err != nil {
		return nil, err
	}

	for _, decl := range f.Decls {
		switch decl := decl.(type) {
		case *ast.GenDecl:
			if decl.Tok == token.IMPORT {
				continue
			}
			for _, spec := range decl.Specs {
				switch spec := spec.(type) {
				case *ast.TypeSpec:
					m = NewModel(spec.Name.String())
					return m, nil
				}

			}
		}
	}

	return nil, fmt.Errorf("not found struct")
}
