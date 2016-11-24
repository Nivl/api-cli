package generate

var modelTpl = `package {{.PackageName}}

// Code generated by api-cli; DO NOT EDIT\n

import (
	"errors"
	{{ if .Generate "JoinSQL" -}}
	"fmt"
	{{- end }}

	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/db"
	uuid "github.com/satori/go.uuid"
)

{{ if .Generate "JoinSQL" -}}
// JoinSQL returns a string ready to be embed in a JOIN query
func JoinSQL(prefix string) string {
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

{{ if .Generate "Save" -}}
// Save creates or updates the {{.ModelNameLC}} depending on the value of the id
func ({{.ModelVar}} *{{.ModelName}}) Save() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return {{.ModelVar}}.Create()
	}

	return {{.ModelVar}}.Update()
}
{{- end }}

{{ if .Generate "Create" -}}
// Create persists a user in the database
func ({{.ModelVar}} *{{.ModelName}}) Create() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID != "" {
		return apierror.NewServerError("cannot persist a {{.ModelNameLC}} that already has a ID")
	}

	return {{.ModelVar}}.doCreate()
}
{{- end }}

{{ if .Generate "doCreate" -}}
// doCreate persists an object in the database
func ({{.ModelVar}} *{{.ModelName}}) doCreate() error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.CreatedAt = db.Now()
	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.CreateStmt}}"
	_, err := app.GetContext().SQL.NamedExec(stmt, {{.ModelVar}})
  return err
}
{{- end }}

{{ if .Generate "Update" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}}.
// Excluded fields are id, created_at, deleted_at
func ({{.ModelVar}} *{{.ModelName}}) Update() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	return {{.ModelVar}}.doUpdate()
}
{{- end }}

{{ if .Generate "doUpdate" -}}
// doUpdate updates an object in the database
func ({{.ModelVar}} *{{.ModelName}}) doUpdate() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.UpdateStmt}}"
	_, err := app.GetContext().SQL.Exec(stmt, {{.UpdateStmtArgs}})
	return err
}
{{- end }}

{{ if .Generate "FullyDelete" -}}
// FullyDelete removes an object from the database
func ({{.ModelVar}} *{{.ModelName}}) FullyDelete() error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return errors.New("{{.ModelNameLC}} has not been saved")
	}

	_, err := sql().Exec("DELETE FROM {{.TableName}} WHERE id=$1", {{.ModelVar}}.ID)
	return err
}
{{- end }}

{{ if .Generate "Delete" -}}
// Delete soft delete an object.
func ({{.ModelVar}} *{{.ModelName}}) Delete() error {
	return {{.ModelVar}}.doDelete()
}
{{- end }}

{{ if .Generate "doDelete" -}}
// doDelete performs a soft delete operation on an object
func ({{.ModelVar}} *{{.ModelName}}) doDelete() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.DeletedAt = db.Now()

	stmt := "UPDATE {{.TableName}} SET deleted_at = $2 WHERE id=$1"
	_, err := sql().Exec(stmt, {{.ModelVar}}.ID, *{{.ModelVar}}.DeletedAt)
	return err
}
{{- end }}

{{ if .Generate "IsZero" -}}
// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}
{{- end }}`
