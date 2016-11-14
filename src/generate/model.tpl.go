package generate

var modelTpl = `package {{.PackageName}}

// Code generated by api-cli; DO NOT EDIT\n

import (
	"errors"

	"github.com/melvin-laplanche/ml-api/src/apierror"
	"github.com/melvin-laplanche/ml-api/src/app"
	"github.com/melvin-laplanche/ml-api/src/db"
	uuid "github.com/satori/go.uuid"
)

// doCreate persists an object in the database
func ({{.ModelVar}} *{{.ModelName}}) doCreate() error {
	if {{.ModelVar}} == nil {
		return errors.New("{{.ModelNameLC}} not instanced")
	}

	{{.ModelVar}}.ID = uuid.NewV4().String()
	{{.ModelVar}}.CreatedAt = db.Now()
	{{.ModelVar}}.UpdatedAt = db.Now()

	stmt := "{{.CreateStmt}}"
	_, err := app.GetContext().SQL.Exec(stmt, {{.CreateStmtArgs}})
  return err
}

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

// doDelete performs a soft delete operation on an object
func ({{.ModelVar}} *{{.ModelName}}) doDelete() error {
	if {{.ModelVar}} == nil {
		return apierror.NewServerError("{{.ModelNameLC}} is not instanced")
	}

	if {{.ModelVar}}.ID == "" {
		return apierror.NewServerError("cannot delete a non-persisted {{.ModelNameLC}}")
	}

	now := db.Now()
	{{.ModelVar}}.DeletedAt = &now

	stmt := "UPDATE {{.TableName}} SET deleted_at = $2 WHERE id=$1"
	_, err := sql().Exec(stmt, {{.ModelVar}}.ID, *{{.ModelVar}}.DeletedAt)
	return err
}

// IsZero checks if the object is either nil or don't have an ID
func ({{.ModelVar}} *{{.ModelName}}) IsZero() bool {
	return {{.ModelVar}} == nil || {{.ModelVar}}.ID == ""
}`
