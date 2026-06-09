// Package application 实现邮件模块的业务逻辑：邮箱配置、模板、发送和日志查询。
package application

import (
	"context"
	"regexp"
	"strings"

	maildomain "ez-admin-gin/server/internal/modules/system/mail/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

var placeholderPattern = regexp.MustCompile(`\{\{\s*([a-z0-9:_-]+)\s*\}\}`)

// Service 封装邮件模块的业务逻辑。
type Service struct {
	tx     MailTransactor
	repo   MailRepository
	sender Sender
}

func NewService(tx MailTransactor, repo MailRepository, sender Sender) *Service {
	return &Service{tx: tx, repo: repo, sender: sender}
}

func (s *Service) ListAccounts(query maildomain.AccountListQuery) (maildomain.AccountListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := maildomain.NormalizeAccountStatusFilter(query.Status)
	if err != nil {
		return maildomain.AccountListResponse{}, err
	}

	items, total, err := s.repo.ListAccounts(query, page, pageSize, status)
	if err != nil {
		return maildomain.AccountListResponse{}, err
	}

	result := make([]maildomain.AccountResponse, 0, len(items))
	for _, item := range items {
		result = append(result, maildomain.BuildAccountResponse(item))
	}
	return maildomain.AccountListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) CreateAccount(req maildomain.CreateAccountRequest) (maildomain.AccountResponse, error) {
	req, err := maildomain.NormalizeCreateAccountRequest(req)
	if err != nil {
		return maildomain.AccountResponse{}, err
	}

	created := maildomain.AccountEntity{
		Name:       req.Name,
		Host:       req.Host,
		Port:       req.Port,
		Username:   req.Username,
		Password:   req.Password,
		FromEmail:  req.FromEmail,
		FromName:   req.FromName,
		Encryption: req.Encryption,
		IsDefault:  req.IsDefault,
		Status:     req.Status,
		Remark:     req.Remark,
	}

	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.AccountNameExists(tx, req.Name, 0)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("邮箱名称已存在")
		}
		if req.IsDefault {
			if err := s.repo.ClearDefaultAccounts(tx, 0); err != nil {
				return err
			}
		}
		return s.repo.CreateAccount(tx, &created)
	}); err != nil {
		return maildomain.AccountResponse{}, err
	}

	return maildomain.BuildAccountResponse(created), nil
}

func (s *Service) UpdateAccount(accountID uint, req maildomain.UpdateAccountRequest) (maildomain.AccountResponse, error) {
	req, err := maildomain.NormalizeUpdateAccountRequest(req)
	if err != nil {
		return maildomain.AccountResponse{}, err
	}

	var updated maildomain.AccountEntity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindAccountByID(tx, accountID)
		if err != nil {
			return err
		}
		exists, err := s.repo.AccountNameExists(tx, req.Name, accountID)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("邮箱名称已存在")
		}
		if req.IsDefault {
			if err := s.repo.ClearDefaultAccounts(tx, accountID); err != nil {
				return err
			}
		}
		if err := s.repo.UpdateAccount(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	})
	if err != nil {
		return maildomain.AccountResponse{}, err
	}

	return maildomain.BuildAccountResponse(updated), nil
}

func (s *Service) UpdateAccountStatus(accountID uint, status model.MailAccountStatus) error {
	status, err := maildomain.NormalizeAccountStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindAccountByID(tx, accountID)
		if err != nil {
			return err
		}
		if item.IsDefault && status == model.MailAccountStatusDisabled {
			return errorsx.BadRequest("默认邮箱不能直接禁用，请先切换默认邮箱")
		}
		return s.repo.UpdateAccountStatus(tx, &item, status)
	})
}

func (s *Service) DeleteAccount(accountID uint) error {
	if accountID == 0 {
		return errorsx.BadRequest("邮箱账号 ID 不正确")
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindAccountByID(tx, accountID)
		if err != nil {
			return err
		}
		if item.IsDefault {
			return errorsx.BadRequest("默认邮箱不能删除，请先切换默认邮箱")
		}
		return s.repo.DeleteAccount(tx, &item)
	})
}

func (s *Service) TestAccount(accountID uint, req maildomain.TestAccountRequest) (maildomain.SendMailResponse, error) {
	req, err := maildomain.NormalizeTestAccountRequest(req)
	if err != nil {
		return maildomain.SendMailResponse{}, err
	}

	account, err := s.repo.FindEnabledAccountByID(accountID)
	if err != nil {
		return maildomain.SendMailResponse{}, err
	}

	message := Message{To: req.To, Subject: req.Subject, Content: req.Content, IsHTML: false}
	logItem, sendErr, logErr := s.sendAndLog(account, maildomain.TemplateEntity{}, message)
	testMessage := "测试发送成功"
	if sendErr != nil {
		testMessage = trimForMessage(sendErr.Error(), 255)
	}
	_ = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item := account
		return s.repo.UpdateAccountTestResult(tx, &item, testMessage)
	})
	if logErr != nil {
		return maildomain.SendMailResponse{LogID: logItem.ID, Status: logItem.Status}, errorsx.Internal("测试邮件发送记录保存失败", logErr)
	}
	if sendErr != nil {
		result := maildomain.SendMailResponse{LogID: logItem.ID, Status: logItem.Status}
		if isLoggedSendFailure(logItem) {
			return result, nil
		}
		return result, errorsx.ServiceUnavailable("测试邮件发送失败", sendErr)
	}
	return maildomain.SendMailResponse{LogID: logItem.ID, Status: model.MailLogStatusSuccess}, nil
}

func (s *Service) ListTemplates(query maildomain.TemplateListQuery) (maildomain.TemplateListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := maildomain.NormalizeTemplateStatusFilter(query.Status)
	if err != nil {
		return maildomain.TemplateListResponse{}, err
	}

	items, total, err := s.repo.ListTemplates(query, page, pageSize, status)
	if err != nil {
		return maildomain.TemplateListResponse{}, err
	}

	result := make([]maildomain.TemplateResponse, 0, len(items))
	for _, item := range items {
		result = append(result, maildomain.BuildTemplateResponse(item))
	}
	return maildomain.TemplateListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) CreateTemplate(req maildomain.CreateTemplateRequest) (maildomain.TemplateResponse, error) {
	req, err := maildomain.NormalizeCreateTemplateRequest(req)
	if err != nil {
		return maildomain.TemplateResponse{}, err
	}

	created := maildomain.TemplateEntity{
		Code:      req.Code,
		Name:      req.Name,
		Subject:   req.Subject,
		Content:   req.Content,
		IsHTML:    req.IsHTML,
		Variables: maildomain.EncodeVariables(req.Variables),
		Sort:      req.Sort,
		Status:    req.Status,
		Remark:    req.Remark,
	}

	if err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.TemplateCodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("模板编码已存在")
		}
		return s.repo.CreateTemplate(tx, &created)
	}); err != nil {
		return maildomain.TemplateResponse{}, err
	}

	return maildomain.BuildTemplateResponse(created), nil
}

func (s *Service) UpdateTemplate(templateID uint, req maildomain.UpdateTemplateRequest) (maildomain.TemplateResponse, error) {
	req, err := maildomain.NormalizeUpdateTemplateRequest(req)
	if err != nil {
		return maildomain.TemplateResponse{}, err
	}

	var updated maildomain.TemplateEntity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateTemplate(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	})
	if err != nil {
		return maildomain.TemplateResponse{}, err
	}

	return maildomain.BuildTemplateResponse(updated), nil
}

func (s *Service) UpdateTemplateStatus(templateID uint, status model.MailTemplateStatus) error {
	status, err := maildomain.NormalizeTemplateStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		return s.repo.UpdateTemplateStatus(tx, &item, status)
	})
}

func (s *Service) DeleteTemplate(templateID uint) error {
	if templateID == 0 {
		return errorsx.BadRequest("邮件模板 ID 不正确")
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		return s.repo.DeleteTemplate(tx, &item)
	})
}

func (s *Service) RenderTemplate(templateID uint, req maildomain.RenderTemplateRequest) (maildomain.RenderTemplateResponse, error) {
	var item maildomain.TemplateEntity
	err := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		found, err := s.repo.FindTemplateByID(tx, templateID)
		if err != nil {
			return err
		}
		item = found
		return nil
	})
	if err != nil {
		return maildomain.RenderTemplateResponse{}, err
	}

	subject, content := renderTemplate(item, req.Variables)
	return maildomain.RenderTemplateResponse{Subject: subject, Content: content}, nil
}

func (s *Service) Send(req maildomain.SendMailRequest) (maildomain.SendMailResponse, error) {
	req, err := maildomain.NormalizeSendMailRequest(req)
	if err != nil {
		return maildomain.SendMailResponse{}, err
	}

	account, err := s.resolveAccount(req.AccountID)
	if err != nil {
		return maildomain.SendMailResponse{}, err
	}

	var template maildomain.TemplateEntity
	message := Message{To: req.To, Cc: req.Cc, Bcc: req.Bcc, Subject: req.Subject, Content: req.Content, IsHTML: req.IsHTML}
	if req.TemplateCode != "" {
		template, err = s.repo.FindEnabledTemplateByCode(req.TemplateCode)
		if err != nil {
			return maildomain.SendMailResponse{}, err
		}
		subject, content := renderTemplate(template, req.Variables)
		message.Subject = subject
		message.Content = content
		message.IsHTML = template.IsHTML
	}

	logItem, sendErr, logErr := s.sendAndLog(account, template, message)
	if logErr != nil {
		return maildomain.SendMailResponse{LogID: logItem.ID, Status: logItem.Status}, errorsx.Internal("邮件发送记录保存失败", logErr)
	}
	if sendErr != nil {
		result := maildomain.SendMailResponse{LogID: logItem.ID, Status: logItem.Status}
		if isLoggedSendFailure(logItem) {
			return result, nil
		}
		return result, errorsx.ServiceUnavailable("邮件发送失败", sendErr)
	}
	return maildomain.SendMailResponse{LogID: logItem.ID, Status: model.MailLogStatusSuccess}, nil
}

func (s *Service) ListLogs(query maildomain.LogListQuery) (maildomain.LogListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := maildomain.NormalizeLogStatusFilter(query.Status)
	if err != nil {
		return maildomain.LogListResponse{}, err
	}

	items, total, err := s.repo.ListLogs(query, page, pageSize, status)
	if err != nil {
		return maildomain.LogListResponse{}, err
	}

	result := make([]maildomain.LogResponse, 0, len(items))
	for _, item := range items {
		result = append(result, maildomain.BuildLogResponse(item))
	}
	return maildomain.LogListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

func (s *Service) resolveAccount(accountID uint) (maildomain.AccountEntity, error) {
	if accountID > 0 {
		return s.repo.FindEnabledAccountByID(accountID)
	}
	return s.repo.FindDefaultAccount()
}

func (s *Service) sendAndLog(account maildomain.AccountEntity, template maildomain.TemplateEntity, message Message) (maildomain.LogEntity, error, error) {
	logItem := maildomain.LogEntity{
		AccountID:    account.ID,
		AccountName:  account.Name,
		TemplateID:   template.ID,
		TemplateCode: template.Code,
		Subject:      message.Subject,
		FromEmail:    account.FromEmail,
		ToEmails:     maildomain.JoinEmails(message.To),
		CcEmails:     maildomain.JoinEmails(message.Cc),
		BccEmails:    maildomain.JoinEmails(message.Bcc),
		Status:       model.MailLogStatusSuccess,
	}

	err := s.sender.Send(account, message)
	if err != nil {
		logItem.Status = model.MailLogStatusFailed
		logItem.ErrorMessage = trimForMessage(err.Error(), 500)
	}

	if logErr := s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		return s.repo.CreateLog(tx, &logItem)
	}); logErr != nil {
		return logItem, err, logErr
	}
	return logItem, err, nil
}

func renderTemplate(template maildomain.TemplateEntity, variables map[string]string) (string, string) {
	subject := replacePlaceholders(template.Subject, variables)
	content := replacePlaceholders(template.Content, variables)
	return subject, content
}

func replacePlaceholders(value string, variables map[string]string) string {
	return placeholderPattern.ReplaceAllStringFunc(value, func(token string) string {
		matches := placeholderPattern.FindStringSubmatch(token)
		if len(matches) < 2 {
			return token
		}
		key := matches[1]
		if variables == nil {
			return ""
		}
		return strings.TrimSpace(variables[key])
	})
}

func trimForMessage(value string, maxLen int) string {
	value = strings.TrimSpace(value)
	if len(value) <= maxLen {
		return value
	}
	return value[:maxLen]
}

func isLoggedSendFailure(logItem maildomain.LogEntity) bool {
	return logItem.ID > 0 && logItem.Status == model.MailLogStatusFailed
}
