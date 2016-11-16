package generate

import (
	"bytes"
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"html/template"
	"io/ioutil"
	"os"
	"path"
	"regexp"
	"strings"

	"github.com/serenize/snaker"
	"github.com/urfave/cli"
)

// Regex to parse db struct tag like `db:"field_name"`
var dbFieldRegex, _ = regexp.Compile("`db:\"([a-zA-Z0-9_-]+),?.*\"`")

// ModelField represents a field from the struct we are parsing
type ModelField struct {
	Name   string
	DbName string
}

// Model represent a model to generate
type Model struct {
	Name        string
	Table       string
	FileName    string
	PackageName string
	Path        string
	FullPath    string
	Fields      []*ModelField
}

// ModelTemplateVars contains all the variable needed to render the new file
type ModelTemplateVars struct {
	ModelName      string
	ModelNameLC    string
	TableName      string
	ModelVar       string
	PackageName    string
	CreateStmt     string
	CreateStmtArgs string
	UpdateStmt     string
	UpdateStmtArgs string
}

// setDefault control what has been set in the model, and set default values where needed
func (m *Model) setDefault() error {
	if m.Name == "" {
		return errors.New("model name missing")
	}

	if m.FileName == "" {
		return errors.New("filename missing. use -f to specify one")
	}

	if m.PackageName == "" {
		return errors.New("package name missing. use -p to specify one")
	}

	if m.Table == "" {
		m.Table = snaker.CamelToSnake(m.Name)
	}

	pwd, err := os.Getwd()
	if err != nil {
		return err
	}
	m.Path = pwd
	m.FullPath = path.Join(m.Path, m.FileName)
	return nil
}

// Parse parses and render a model
func (m *Model) Parse() error {
	if err := m.setDefault(); err != nil {
		return err
	}

	// Open the file
	file, err := os.Open(m.FullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Put the content of the file in a string
	fileStr, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	// Parse the file
	astFile, err := parser.ParseFile(token.NewFileSet(), "", fileStr, parser.AllErrors)
	if err != nil {
		return err
	}
	if err := m.parseTarget(astFile); err != nil {
		return err
	}

	return m.generate()
}

// generate generates the new file
func (m *Model) generate() error {
	vars := &ModelTemplateVars{
		ModelName:   m.Name,
		ModelNameLC: strings.ToLower(m.Name),
		TableName:   m.Table,
		ModelVar:    string(strings.ToLower(m.Name)[0]),
		PackageName: m.PackageName,
	}

	// Create Statement
	createFields := make([]string, len(m.Fields))
	createValues := make([]string, len(m.Fields))
	createArgs := make([]string, len(m.Fields))
	for i, field := range m.Fields {
		createFields[i] = field.DbName
		createValues[i] = fmt.Sprintf("$%d", i+1)
		createArgs[i] = fmt.Sprintf("%s.%s", vars.ModelVar, field.Name)
	}
	vars.CreateStmtArgs = strings.Join(createArgs, ", ")
	vars.CreateStmt = fmt.Sprintf(
		"INSERT INTO %s (%s) VALUES (%s)",
		vars.TableName,
		strings.Join(createFields, ", "),
		strings.Join(createValues, ", "),
	)

	// Update Statement
	updateFields := make([]string, len(m.Fields))
	updateArgs := make([]string, len(m.Fields)+1)
	var i int
	for _, field := range m.Fields {
		updateFields[i] = fmt.Sprintf("%s = $%d", field.DbName, i+1)
		updateArgs[i] = fmt.Sprintf("%s.%s", vars.ModelVar, field.Name)
		i++
	}
	updateArgs[i] = fmt.Sprintf("%s.ID", vars.ModelVar)
	vars.UpdateStmtArgs = strings.Join(updateArgs, ", ")
	vars.UpdateStmt = fmt.Sprintf(
		"UPDATE %s SET %s WHERE id=$%d",
		vars.TableName,
		strings.Join(updateFields, ", "),
		i+1,
	)

	// Get the template and parse it with the variables we have
	t, err := template.New("model").Parse(modelTpl)
	if err != nil {
		fmt.Println(err)
		return err
	}
	buf := &bytes.Buffer{}
	if err := t.Execute(buf, vars); err != nil {
		fmt.Println(err)
		return err
	}
	output := strings.TrimSpace(buf.String())

	// Write the new file to the disk
	newFile, err := os.Create(m.generatedFileName())
	if err != nil {
		return err
	}
	defer newFile.Close()
	if _, err := newFile.WriteString(output); err != nil {
		return err
	}

	return nil
}

// generatedFileName returns the file name of the new file
func (m *Model) generatedFileName() string {
	return strings.TrimSuffix(m.FullPath, ".go") + "_generated.go"
}

// parseTarget parses the source file to get the Model fields
func (m *Model) parseTarget(f *ast.File) error {
	obj, ok := f.Scope.Objects[m.Name]
	if !ok {
		return fmt.Errorf("could not find type %s in %s", m.Name, m.FullPath)
	}
	typeSpec, ok := obj.Decl.(*ast.TypeSpec)
	if !ok {
		return fmt.Errorf("%s is not a type", m.Name)
	}
	structType, ok := typeSpec.Type.(*ast.StructType)
	if !ok {
		return fmt.Errorf("%s is not a struct", m.Name)
	}

	for _, field := range structType.Fields.List {
		// We do not handle structs with no name, and un-exported fields
		// Also, I'm not sure in what case we can have more than one name?
		if len(field.Names) > 0 && field.Names[0].IsExported() {
			// Lets be sure the field has a Tag
			if field.Tag == nil {
				continue
			}
			dbName := dbFieldRegex.FindStringSubmatch(field.Tag.Value)
			// for `db:"name"` the func returns [`db:"name"` name], and we want "name"
			if len(dbName) != 2 {
				continue
			}

			newField := &ModelField{
				Name:   field.Names[0].Name,
				DbName: dbName[1],
			}

			m.Fields = append(m.Fields, newField)
		}
	}

	return nil
}

// GenModel is used to generate a new model
func GenModel(c *cli.Context) error {
	model := &Model{
		Name:        c.Args().First(),
		Table:       c.String("table"),
		FileName:    c.String("file"),
		PackageName: c.String("package"),
	}

	return model.Parse()
}
