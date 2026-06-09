package api

import maildomain "ez-admin-gin/server/internal/modules/system/mail/domain"

type AccountListQuery = maildomain.AccountListQuery
type CreateAccountRequest = maildomain.CreateAccountRequest
type UpdateAccountRequest = maildomain.UpdateAccountRequest
type UpdateAccountStatusRequest = maildomain.UpdateAccountStatusRequest
type AccountResponse = maildomain.AccountResponse
type AccountListResponse = maildomain.AccountListResponse

type TemplateListQuery = maildomain.TemplateListQuery
type CreateTemplateRequest = maildomain.CreateTemplateRequest
type UpdateTemplateRequest = maildomain.UpdateTemplateRequest
type UpdateTemplateStatusRequest = maildomain.UpdateTemplateStatusRequest
type TemplateResponse = maildomain.TemplateResponse
type TemplateListResponse = maildomain.TemplateListResponse
type RenderTemplateRequest = maildomain.RenderTemplateRequest
type RenderTemplateResponse = maildomain.RenderTemplateResponse

type SendMailRequest = maildomain.SendMailRequest
type SendMailResponse = maildomain.SendMailResponse
type TestAccountRequest = maildomain.TestAccountRequest

type LogListQuery = maildomain.LogListQuery
type LogResponse = maildomain.LogResponse
type LogListResponse = maildomain.LogListResponse
