package generate

var modelTestTpl = `package {{.PackageName}}

// Code generated; DO NOT EDIT.

import (
	"testing"

		"github.com/stretchr/testify/assert"

		"github.com/satori/go.uuid"

		"github.com/Nivl/go-rest-tools/storage/db/mockdb"

	{{ if or (.Generate "doCreate") (.Generate "doUpdate") }}"github.com/Nivl/go-rest-tools/types/datetime"{{ end }}
)


{{ if .Generate "Save" -}}
func Test{{.ModelName}}SaveNew(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
}

func Test{{.ModelName}}SaveExisting(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	id := uuid.NewV4().String()
	{{.ModelVar}}.ID = id
	err := {{.ModelVar}}.Save(mockDB)

	assert.NoError(t, err, "Save() should not have fail")
	mockDB.AssertExpectations(t)
	assert.Equal(t, id, {{.ModelVar}}.ID, "ID should not have changed")
}
{{- end }}

{{ if .Generate "Create" -}}
func Test{{.ModelName}}Create(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.Create(mockDB)

	assert.NoError(t, err, "Create() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
	assert.NotNil(t, {{.ModelVar}}.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, {{.ModelVar}}.UpdatedAt, "UpdatedAt should have been set")
}

func Test{{.ModelName}}CreateWithID(t *testing.T) {
	mockDB := &mockdb.Queryable{}

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()

	err := {{.ModelVar}}.Create(mockDB)
	assert.Error(t, err, "Create() should have fail")
	mockDB.AssertExpectations(t)
}
{{- end }}

{{ if .Generate "doCreate" -}}
func Test{{.ModelName}}DoCreate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
	assert.NotNil(t, {{.ModelVar}}.CreatedAt, "CreatedAt should have been set")
	assert.NotNil(t, {{.ModelVar}}.UpdatedAt, "UpdatedAt should have been set")
}

func Test{{.ModelName}}DoCreateWithDate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsert("*{{.PackageName}}.{{.ModelName}}")

	createdAt := datetime.Now().AddDate(0, 0, 1)
	{{.ModelVar}} := &{{.ModelName}}{CreatedAt: createdAt}
	err := {{.ModelVar}}.doCreate(mockDB)

	assert.NoError(t, err, "doCreate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
	assert.True(t, {{.ModelVar}}.CreatedAt.Equal(createdAt), "CreatedAt should not have been updated")
	assert.NotNil(t, {{.ModelVar}}.UpdatedAt, "UpdatedAt should have been set")
}

func Test{{.ModelName}}DoCreateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectInsertError("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.doCreate(mockDB)

	assert.Error(t, err, "doCreate() should have fail")
	mockDB.AssertExpectations(t)
}
{{- end }}


{{ if .Generate "Update" -}}
func Test{{.ModelName}}Update(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	err := {{.ModelVar}}.Update(mockDB)

	assert.NoError(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
	assert.NotNil(t, {{.ModelVar}}.UpdatedAt, "UpdatedAt should have been set")
}

func Test{{.ModelName}}UpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.Update(mockDB)

	assert.Error(t, err, "Update() should not have fail")
	mockDB.AssertExpectations(t)
}
{{- end }}


{{ if .Generate "doUpdate" -}}
func Test{{.ModelName}}DoUpdate(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdate("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	err := {{.ModelVar}}.doUpdate(mockDB)

	assert.NoError(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
	assert.NotEmpty(t, {{.ModelVar}}.ID, "ID should have been set")
	assert.NotNil(t, {{.ModelVar}}.UpdatedAt, "UpdatedAt should have been set")
}

func Test{{.ModelName}}DoUpdateWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should not have fail")
	mockDB.AssertExpectations(t)
}

func Test{{.ModelName}}DoUpdateFail(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectUpdateError("*{{.PackageName}}.{{.ModelName}}")

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	err := {{.ModelVar}}.doUpdate(mockDB)

	assert.Error(t, err, "doUpdate() should have fail")
	mockDB.AssertExpectations(t)
}
{{- end }}

{{ if .Generate "Delete" -}}
func Test{{.ModelName}}Delete(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletion()

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	err := {{.ModelVar}}.Delete(mockDB)

	assert.NoError(t, err, "Delete() should not have fail")
	mockDB.AssertExpectations(t)
}

func Test{{.ModelName}}DeleteWithoutID(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	{{.ModelVar}} := &{{.ModelName}}{}
	err := {{.ModelVar}}.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}

func Test{{.ModelName}}DeleteError(t *testing.T) {
	mockDB := &mockdb.Queryable{}
	mockDB.ExpectDeletionError()

	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	err := {{.ModelVar}}.Delete(mockDB)

	assert.Error(t, err, "Delete() should have fail")
	mockDB.AssertExpectations(t)
}
{{- end }}

{{ if .Generate "GetID" -}}
func Test{{.ModelName}}GetID(t *testing.T) {
	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.ID = uuid.NewV4().String()
	assert.Equal(t, {{.ModelVar}}.ID, {{.ModelVar}}.GetID(), "GetID() did not return the right ID")
}
{{- end }}

{{ if .Generate "SetID" -}}
func Test{{.ModelName}}SetID(t *testing.T) {
	{{.ModelVar}} := &{{.ModelName}}{}
	{{.ModelVar}}.SetID(uuid.NewV4().String())
	assert.NotEmpty(t, {{.ModelVar}}.ID, "SetID() did not set the ID")
}
{{- end }}

{{ if .Generate "IsZero" -}}
func Test{{.ModelName}}IsZero(t *testing.T) {
	empty := &{{.ModelName}}{}
	assert.True(t, empty.IsZero(), "IsZero() should return true for empty struct")

	var nilStruct *{{.ModelName}}
	assert.True(t, nilStruct.IsZero(), "IsZero() should return true for nil struct")

	valid := &{{.ModelName}}{ID: uuid.NewV4().String()}
	assert.False(t, valid.IsZero(), "IsZero() should return false for valid struct")
}
{{- end }}`
