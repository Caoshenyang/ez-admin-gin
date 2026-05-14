package application

import (
	{{ .Module }}domain "ez-admin-gin/server/internal/modules/{{ .Group }}/{{ .Module }}/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type {{ .LabelEn }}Transactor = database.Transactor

type {{ .LabelEn }}Repository interface {
	List(query {{ .Module }}domain.ListQuery, page int, pageSize int{{- range .FilterFields }}{{- if eq .FilterType "select" }}, {{ .Name }} *model.{{ .StatusTypeName }}{{- end }}{{- end }}) ([]{{ .Module }}domain.Entity, int64, error)
	FindByID(db *gorm.DB, id uint) ({{ .Module }}domain.Entity, error)
	Create(db *gorm.DB, item *{{ .Module }}domain.Entity) error
	Update(db *gorm.DB, item *{{ .Module }}domain.Entity, req {{ .Module }}domain.UpdateRequest) error
{{- if .HasStatus }}
	UpdateStatus(db *gorm.DB, item *{{ .Module }}domain.Entity, status model.{{ .StatusField.StatusTypeName }}) error
{{- end }}
}
