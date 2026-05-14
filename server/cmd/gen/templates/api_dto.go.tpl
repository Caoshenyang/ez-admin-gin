package api

import {{ .Group }}domain "ez-admin-gin/server/internal/modules/{{ .Group }}/{{ .Module }}/domain"

type ListQuery = {{ .Group }}domain.ListQuery
type CreateRequest = {{ .Group }}domain.CreateRequest
type UpdateRequest = {{ .Group }}domain.UpdateRequest
{{- if .HasStatus }}
type UpdateStatusRequest = {{ .Group }}domain.UpdateStatusRequest
{{- end }}
type Response = {{ .Group }}domain.Response
type ListResponse = {{ .Group }}domain.ListResponse
