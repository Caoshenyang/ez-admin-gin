package api

import messagedomain "ez-admin-gin/server/internal/modules/system/message/domain"

type TemplateListQuery = messagedomain.TemplateListQuery
type CreateTemplateRequest = messagedomain.CreateTemplateRequest
type UpdateTemplateRequest = messagedomain.UpdateTemplateRequest
type UpdateTemplateStatusRequest = messagedomain.UpdateTemplateStatusRequest
type TemplateResponse = messagedomain.TemplateResponse
type TemplateListResponse = messagedomain.TemplateListResponse

type ReminderListQuery = messagedomain.ReminderListQuery
type CreateReminderRequest = messagedomain.CreateReminderRequest
type UpdateReminderRequest = messagedomain.UpdateReminderRequest
type UpdateReminderStatusRequest = messagedomain.UpdateReminderStatusRequest
type ReminderResponse = messagedomain.ReminderResponse
type ReminderListResponse = messagedomain.ReminderListResponse
