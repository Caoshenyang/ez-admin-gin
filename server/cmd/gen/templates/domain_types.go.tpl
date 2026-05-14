package domain

import (
	"strings"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

type ListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
{{- range .FilterFields }}
{{- if eq .FilterType "keyword" }}
	Keyword  string `form:"keyword"`
{{- else if eq .FilterType "select" }}
	{{ title .Name }} int    `form:"{{ .Name }}"`
{{- end }}
{{- end }}
}

type CreateRequest struct {
{{- range .FormFields }}
	{{ title .Name }} {{ if .IsStatus }}{{ .GoTypeResolved }}{{ else }}{{ .Type }}{{ end }} `json:"{{ .Name }}"`
{{- end }}
}

type UpdateRequest struct {
{{- range .FormFields }}
	{{ title .Name }} {{ if .IsStatus }}{{ .GoTypeResolved }}{{ else }}{{ .Type }}{{ end }} `json:"{{ .Name }}"`
{{- end }}
}

{{ if .HasStatus }}
type UpdateStatusRequest struct {
{{- $status := .StatusField }}
{{- if $status }}
	{{ title $status.Name }} {{ $status.GoTypeResolved }} `json:"{{ $status.Name }}"`
{{- end }}
}
{{ end }}

type Response struct {
{{- range .Fields }}
{{- if eq .Name "id" }}
	ID        uint               `json:"id"`
{{- else if eq .Name "created_at" }}
	CreatedAt string             `json:"created_at"`
{{- else if eq .Name "updated_at" }}
	UpdatedAt string             `json:"updated_at"`
{{- else if .IsStatus }}
	{{ title .Name }}    {{ .GoTypeResolved }}           `json:"{{ .Name }}"`
{{- else }}
	{{ title .Name }}    {{ .Type }}                      `json:"{{ .Name }}"`
{{- end }}
{{- end }}
}

type ListResponse struct {
	Items    []Response `json:"items"`
	Total    int64      `json:"total"`
	Page     int        `json:"page"`
	PageSize int        `json:"page_size"`
}

type Entity = model.{{ .LabelEn }}

const (
	PermissionList         = "{{ .Group }}:{{ .Module }}:list"
	PermissionCreate       = "{{ .Group }}:{{ .Module }}:create"
	PermissionUpdate       = "{{ .Group }}:{{ .Module }}:update"
	PermissionDelete       = "{{ .Group }}:{{ .Module }}:delete"
{{- if .HasStatus }}
	PermissionUpdateStatus = "{{ .Group }}:{{ .Module }}:update_status"
{{- end }}
)

func NormalizeCreateRequest(req CreateRequest) (CreateRequest, error) {
{{- range .FormFields }}
{{- if .Required }}
	if strings.TrimSpace(req.{{ title .Name }}) == "" {
		return CreateRequest{}, errorsx.BadRequest("{{ .Label }}不能为空")
	}
{{- end }}
{{- end }}
	return req, nil
}

func NormalizeUpdateRequest(req UpdateRequest) (UpdateRequest, error) {
{{- range .FormFields }}
{{- if .Required }}
	if strings.TrimSpace(req.{{ title .Name }}) == "" {
		return UpdateRequest{}, errorsx.BadRequest("{{ .Label }}不能为空")
	}
{{- end }}
{{- end }}
	return req, nil
}

func BuildResponse(item model.{{ $.LabelEn }}) Response {
	return Response{
{{- range .Fields }}
{{- if eq .Name "id" }}
		ID: item.ID,
{{- else if eq .Name "created_at" }}
		CreatedAt: item.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
{{- else if eq .Name "updated_at" }}
		UpdatedAt: item.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
{{- else }}
		{{ title .Name }}: item.{{ title .Name }},
{{- end }}
{{- end }}
	}
}
