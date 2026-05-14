package model

import (
	"time"

	"gorm.io/gorm"
)

{{- $statusType := "" }}
{{- range .Fields }}
{{- if .IsStatus }}
{{- $statusType = .StatusTypeName }}
{{- end }}
{{- end }}

{{ if $statusType }}
// {{ $statusType }} 表示{{ $.Label }}状态。
type {{ $statusType }} int

const (
	// {{ $statusType }}Enabled 表示{{ $.Label }}可见。
	{{ $statusType }}Enabled {{ $statusType }} = 1
	// {{ $statusType }}Disabled 表示{{ $.Label }}已隐藏。
	{{ $statusType }}Disabled {{ $statusType }} = 2
)
{{ end }}

// {{ $.LabelEn }} 是{{ $.Label }}表模型。
type {{ $.LabelEn }} struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
{{- range .Fields }}
{{- if eq .Name "id" "created_at" "updated_at" "deleted_at" }}
{{- continue }}
{{- end }}
{{- if .IsStatus }}
	{{ title .Name }}    {{ .GoTypeResolved }}   `gorm:"{{ .DbType }}" json:"{{ .Name }}"`
{{- else }}
	{{ title .Name }}    {{ .Type }}              `gorm:"{{ .DbType }}" json:"{{ .Name }}"`
{{- end }}
{{- end }}
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定{{ $.Label }}表名。
func ({{ $.LabelEn }}) TableName() string {
	return "{{ .Table }}"
}
