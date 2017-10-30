package generate

var modelTpl = `package {{.PackageName}}

// Code generated; DO NOT EDIT.

import (
	"errors"
	{{ if .Generate "JoinSQL"}}
	"fmt"
	"strings"
	{{ end }}

	"github.com/Nivl/go-rest-tools/types/apierror"
	{{ if or (.Generate "doCreate") (.Generate "doUpdate") }}"github.com/Nivl/go-types/datetime"{{ end }}
	"github.com/Nivl/go-sqldb"
	uuid "github.com/satori/go.uuid"
)

{{ if .Generate "JoinSQL" -}}
// Join{{.OptionalName}}SQL returns a string ready to be embed in a JOIN query
func Join{{.OptionalName}}SQL(prefix string) string {
	fields := []string{ {{.FieldsAsArray}} }
	output := ""

	for _, field := range fields {
		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\", ", fullName, fullName)
	}
	return strings.TrimSuffix(output, ", ")
}
{{- end }}

{{ if .Generate "Get" -}}
// Get{{.OptionalName}}ByID finds and returns an active {{.ModelNameLC}} by ID
// Deleted object are not returned
func Get{{.OptionalName}}ByID(q sqldb.Queryable, id string) (*{{.ModelName}}, error) {
	{{.ModelVar}} := &{{.ModelName}}{}
	stmt := "SELECT * from {{.TableName}} WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := q.Get({{.ModelVar}}, stmt, id)
	return {{.ModelVar}}, apierror.NewFromSQL(err)
}
{{- end }}

{{ if .Generate "GetAny" -}}
// GetAny{{.OptionalName}}ByID finds and returns an {{.ModelNameLC}} by ID.
// Deleted object are returned
func GetAny{{.OptionalName}}ByID(q sqldb.Queryable, id string) (*{{.ModelName}}, error) {
	{{.ModelVar}} := &{{.ModelName}}{}
	stmt := "SELECT * from {{.TableName}} WHERE id=$1 LIMIT 1"
	err := q.Get({{.ModelVar}}, stmt, id)
	return {{.ModelVar}}, apierror.NewFromSQL(err)
}
{{- end }}


{{ if .Generate "Save" -}}
// Save creates or updates the article depending on the value of the id using
// a transaction
func ({{.ModelVar}} *{{.ModelName}}) Save(q sqldb.Queryable) error {
	if {{.ModelVar}}.ID == "" {
		return {{.ModelVar}}.Create(q)
	}

	return {{.ModelVar}}.Update(q)
}
{{- end }}

{{ if .Generate "Create" -}}
// Create persists a {{.ModelNameLC}} in the database
func ({{.ModelVar}} *{{.ModelName}}) Create(q sqldb.Queryable) error {
	if {{.ModelVar}}.ID != "" {
		return errors.New("cannot persist a {{.ModelNameLC}} that already has an ID")
	}

	return {{.ModelVar}}.doCreate(q)
}
{{- end }}

{{ if .Generate "doCreate" -}}
// doCreate persists a {{.ModelNameLC}} in the database using a Node
func ({{.ModelVar}} *{{.ModelName}}) doCreate(q sqldb.Queryable) error {
	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.UpdatedAt = datetime.Now()
	if {{.ModelVar}}.CreatedAt == nil {
		{{.ModelVar}}.CreatedAt = datetime.Now()
	}

	stmt := "{{.CreateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

  return apierror.NewFromSQL(err)
}
{{- end }}

{{ if .Generate "Update" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}}
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) Update(q sqldb.Queryable) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("cannot update a non-persisted {{.ModelNameLC}}")
	}

	return {{.ModelVar}}.doUpdate(q)
}
{{- end }}

{{ if .Generate "doUpdate" -}}
// doUpdate updates a {{.ModelNameLC}} in the database
func ({{.ModelVar}} *{{.ModelName}}) doUpdate(q sqldb.Queryable) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("cannot update a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.UpdatedAt = datetime.Now()

	stmt := "{{.UpdateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

	return apierror.NewFromSQL(err)
}
{{- end }}

{{ if .Generate "Delete" -}}
// Delete removes a {{.ModelNameLC}} from the database
func ({{.ModelVar}} *{{.ModelName}}) Delete(q sqldb.Queryable) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("{{.ModelNameLC}} has not been saved")
	}

	stmt := "DELETE FROM {{.TableName}} WHERE id=$1"
	_, err := q.Exec(stmt, {{.ModelVar}}.ID)

	return err
}
{{- end }}

{{ if .Generate "GetID" -}}
// GetID returns the ID field
func ({{.ModelVar}} *{{.ModelName}}) GetID() string {
	return {{.ModelVar}}.ID
}
{{- end }}

{{ if .Generate "SetID" -}}
// SetID sets the ID field
func ({{.ModelVar}} *{{.ModelName}}) SetID(id string) {
	{{.ModelVar}}.ID = id
}
{{- end }}

{{ if .Generate "IsZero" -}}
// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}
{{- end }}`
