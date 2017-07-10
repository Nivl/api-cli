package generate

var modelTpl = `package {{.PackageName}}

// Code auto-generated; DO NOT EDIT

import (
	"errors"
	{{ if .Generate "JoinSQL" -}}
	"fmt"
	{{- end }}

	"github.com/Nivl/go-rest-tools/network/http/httperr"
	"github.com/Nivl/go-rest-tools/storage/db"
	uuid "github.com/satori/go.uuid"
)

{{ if .Generate "JoinSQL" -}}
// Join{{.OptionalName}}SQL returns a string ready to be embed in a JOIN query
func Join{{.OptionalName}}SQL(prefix string) string {
	fields := []string{ {{.FieldsAsArray}} }
	output := ""

	for i, field := range fields {
		if i != 0 {
			output += ", "
		}

		fullName := fmt.Sprintf("%s.%s", prefix, field)
		output += fmt.Sprintf("%s \"%s\"", fullName, fullName)
	}

	return output
}
{{- end }}

{{ if .Generate "Get" -}}
// Get{{.OptionalName}}ByID finds and returns an active {{.ModelNameLC}} by ID
func Get{{.OptionalName}}ByID(q db.DB, id string) (*{{.ModelName}}, error) {
	{{.ModelVar}} := &{{.ModelName}}{}
	stmt := "SELECT * from {{.TableName}} WHERE id=$1 and deleted_at IS NULL LIMIT 1"
	err := db.Get(q, {{.ModelVar}}, stmt, id)
	// We want to return nil if a {{.ModelNameLC}} is not found
	if {{.ModelVar}}.ID == "" {
		return nil, err
	}
	return {{.ModelVar}}, err
}
{{- end }}

{{ if .Generate "Exists" -}}
// {{.OptionalName}}Exists checks if a {{.ModelNameLC}} exists for a specific ID
func {{.OptionalName}}Exists(q db.DB, id string) (bool, error) {
	exists := false
	stmt := "SELECT exists(SELECT 1 FROM {{.TableName}} WHERE id=$1 and deleted_at IS NULL)"
	err := db.Get(q, &exists, stmt, id)
	return exists, err
}
{{- end }}

{{ if .Generate "Save" -}}
// Save creates or updates the article depending on the value of the id using
// a transaction
func ({{.ModelVar}} *{{.ModelName}}) Save(q db.DB) error {
	if {{.ModelVar}}.ID == "" {
		return {{.ModelVar}}.Create(q)
	}

	return {{.ModelVar}}.Update(q)
}
{{- end }}

{{ if .Generate "Create" -}}
// Create persists a {{.ModelNameLC}} in the database
func ({{.ModelVar}} *{{.ModelName}}) Create(q db.DB) error {
	if {{.ModelVar}}.ID != "" {
		return errors.New("cannot persist a {{.ModelNameLC}} that already has an ID")
	}

	return {{.ModelVar}}.doCreate(q)
}
{{- end }}

{{ if .Generate "doCreate" -}}
// doCreate persists a {{.ModelNameLC}} in the database using a Node
func ({{.ModelVar}} *{{.ModelName}}) doCreate(q db.DB) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.UpdatedAt = db.Now()
	if {{.ModelVar}}.CreatedAt == nil {
		{{.ModelVar}}.CreatedAt = db.Now()
	}

	stmt := "{{.CreateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

  return httperr.NewFromSQL(err)
}
{{- end }}

{{ if .Generate "Update" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}} using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) Update(q db.DB) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("cannot update a non-persisted {{.ModelNameLC}}")
	}

	return {{.ModelVar}}.doUpdate(q)
}
{{- end }}

{{ if .Generate "doUpdate" -}}
// doUpdate updates a {{.ModelNameLC}} in the database using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doUpdate(q db.DB) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("cannot update a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.UpdateStmt}}"
	_, err := q.NamedExec(stmt, {{.ModelVar}})

	return httperr.NewFromSQL(err)
}
{{- end }}

{{ if .Generate "Delete" -}}
// Delete removes a {{.ModelNameLC}} from the database using a transaction
func ({{.ModelVar}} *{{.ModelName}}) Delete(q db.DB) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return errors.New("{{.ModelNameLC}} has not been saved")
	}

	stmt := "DELETE FROM {{.TableName}} WHERE id=$1"
	_, err := q.Exec(stmt, {{.ModelVar}}.ID)

	return err
}
{{- end }}

{{ if .Generate "Trash" -}}
// Trash soft delete a {{.ModelNameLC}} using a transaction
func ({{.ModelVar}} *{{.ModelName}}) Trash(q db.DB) error {
	return {{.ModelVar}}.doTrash(q)
}
{{- end }}

{{ if .Generate "doTrash" -}}
// doTrash performs a soft delete operation on a {{.ModelNameLC}} using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doTrash(q db.DB) error {
	if {{.ModelVar}}.ID == "" {
		return errors.New("cannot trash a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.DeletedAt = db.Now()

	stmt := "UPDATE {{.TableName}} SET deleted_at = $2 WHERE id=$1"
	_, err := q.Exec(stmt, {{.ModelVar}}.ID, {{.ModelVar}}.DeletedAt)
	return err
}
{{- end }}

{{ if .Generate "IsZero" -}}
// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}
{{- end }}`
