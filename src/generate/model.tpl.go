package generate

var modelTpl = `package {{.PackageName}}

// Code generated by api-cli; DO NOT EDIT\n

import (
	"errors"
	{{ if .Generate "JoinSQL" -}}
	"fmt"
	{{- end }}

	"github.com/Nivl/sqalx"
	"github.com/melvin-laplanche/ml-api/src/apierror"
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
	return {{.ModelVar}}.SaveTx(db.Con())
}
{{- end }}

{{ if .Generate "SaveTx" -}}
// SaveTx creates or updates the article depending on the value of the id using
// a transaction
func ({{.ModelVar}} *{{.ModelName}}) SaveTx(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return {{.ModelVar}}.CreateTx(tx)
	}

	return {{.ModelVar}}.UpdateTx(tx)
}
{{- end }}

{{ if .Generate "Create" -}}
// Create persists a user in the database
func ({{.ModelVar}} *{{.ModelName}}) Create() error {
	return {{.ModelVar}}.CreateTx(db.Con())
}
{{- end }}

{{ if .Generate "CreateTx" -}}
// Create persists a user in the database
func ({{.ModelVar}} *{{.ModelName}}) CreateTx(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID != "" {
		return apierror.NewServerError("cannot persist a {{.ModelNameLC}} that already has an ID")
	}

	return {{.ModelVar}}.doCreate(tx)
}
{{- end }}

{{ if .Generate "doCreate" -}}
// doCreate persists an object in the database using a Node
func ({{.ModelVar}} *{{.ModelName}}) doCreate(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.CreatedAt = db.Now()
	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.CreateStmt}}"
	_, err := tx.NamedExec(stmt, {{.ModelVar}})

  return err
}
{{- end }}

{{ if .Generate "Update" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}}.
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) Update() error {
	return {{.ModelVar}}.UpdateTx(db.Con())
}
{{- end }}

{{ if .Generate "UpdateTx" -}}
// Update updates most of the fields of a persisted {{.ModelNameLC}} using a transaction
// Excluded fields are id, created_at, deleted_at, etc.
func ({{.ModelVar}} *{{.ModelName}}) UpdateTx(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	return {{.ModelVar}}.doUpdate(tx)
}
{{- end }}

{{ if .Generate "doUpdate" -}}
// doUpdate updates an object in the database using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doUpdate(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot update a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.UpdateStmt}}"
	_, err := tx.NamedExec(stmt, {{.ModelVar}})

	return err
}
{{- end }}

{{ if .Generate "FullyDelete" -}}
// FullyDelete removes an object from the database
func ({{.ModelVar}} *{{.ModelName}}) FullyDelete() error {
	return {{.ModelVar}}.FullyDeleteTx(db.Con())
}
{{- end }}

{{ if .Generate "FullyDeleteTx" -}}
// FullyDeleteTx removes an object from the database using a transaction
func ({{.ModelVar}} *{{.ModelName}}) FullyDeleteTx(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return errors.New("{{.ModelNameLC}} has not been saved")
	}

	stmt := "DELETE FROM {{.TableName}} WHERE id=$1"
	_, err := tx.Exec(stmt, {{.ModelVar}}.ID)

	return err
}
{{- end }}

{{ if .Generate "Delete" -}}
// Delete soft delete an object.
func ({{.ModelVar}} *{{.ModelName}}) Delete() error {
	return {{.ModelVar}}.DeleteTx(db.Con())
}
{{- end }}

{{ if .Generate "DeleteTx" -}}
// DeleteTx soft delete an object using a transaction
func ({{.ModelVar}} *{{.ModelName}}) DeleteTx(tx sqalx.Node) error {
	return {{.ModelVar}}.doDelete(tx)
}
{{- end }}

{{ if .Generate "doDelete" -}}
// doDelete performs a soft delete operation on an object using an optional transaction
func ({{.ModelVar}} *{{.ModelName}}) doDelete(tx sqalx.Node) error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted {{.ModelNameLC}}")
	}

	{{.ModelVar}}.DeletedAt = db.Now()

	stmt := "UPDATE {{.TableName}} SET deleted_at = $2 WHERE id=$1"
	_, err := tx.Exec(stmt, {{.ModelVar}}.ID, {{.ModelVar}}.DeletedAt)
	return err
}
{{- end }}

{{ if .Generate "IsZero" -}}
// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}
{{- end }}`
