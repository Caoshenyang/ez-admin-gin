package domain

import (
	"encoding/json"
	"net/mail"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

var templateCodePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

type AccountListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type CreateAccountRequest struct {
	Name       string                  `json:"name"`
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Username   string                  `json:"username"`
	Password   string                  `json:"password"`
	FromEmail  string                  `json:"from_email"`
	FromName   string                  `json:"from_name"`
	Encryption model.MailEncryption    `json:"encryption"`
	IsDefault  bool                    `json:"is_default"`
	Status     model.MailAccountStatus `json:"status"`
	Remark     string                  `json:"remark"`
}

type UpdateAccountRequest struct {
	Name       string                  `json:"name"`
	Host       string                  `json:"host"`
	Port       int                     `json:"port"`
	Username   string                  `json:"username"`
	Password   string                  `json:"password"`
	FromEmail  string                  `json:"from_email"`
	FromName   string                  `json:"from_name"`
	Encryption model.MailEncryption    `json:"encryption"`
	IsDefault  bool                    `json:"is_default"`
	Status     model.MailAccountStatus `json:"status"`
	Remark     string                  `json:"remark"`
}

type UpdateAccountStatusRequest struct {
	Status model.MailAccountStatus `json:"status"`
}

type AccountResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Host        string                  `json:"host"`
	Port        int                     `json:"port"`
	Username    string                  `json:"username"`
	FromEmail   string                  `json:"from_email"`
	FromName    string                  `json:"from_name"`
	Encryption  model.MailEncryption    `json:"encryption"`
	IsDefault   bool                    `json:"is_default"`
	Status      model.MailAccountStatus `json:"status"`
	LastTestAt  *time.Time              `json:"last_test_at"`
	LastTestMsg string                  `json:"last_test_msg"`
	Remark      string                  `json:"remark"`
	CreatedAt   time.Time               `json:"created_at"`
	UpdatedAt   time.Time               `json:"updated_at"`
}

type AccountListResponse struct {
	Items    []AccountResponse `json:"items"`
	Total    int64             `json:"total"`
	Page     int               `json:"page"`
	PageSize int               `json:"page_size"`
}

type TemplateListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type CreateTemplateRequest struct {
	Code      string                   `json:"code"`
	Name      string                   `json:"name"`
	Subject   string                   `json:"subject"`
	Content   string                   `json:"content"`
	IsHTML    bool                     `json:"is_html"`
	Variables []string                 `json:"variables"`
	Sort      int                      `json:"sort"`
	Status    model.MailTemplateStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type UpdateTemplateRequest struct {
	Name      string                   `json:"name"`
	Subject   string                   `json:"subject"`
	Content   string                   `json:"content"`
	IsHTML    bool                     `json:"is_html"`
	Variables []string                 `json:"variables"`
	Sort      int                      `json:"sort"`
	Status    model.MailTemplateStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type UpdateTemplateStatusRequest struct {
	Status model.MailTemplateStatus `json:"status"`
}

type TemplateResponse struct {
	ID        uint                     `json:"id"`
	Code      string                   `json:"code"`
	Name      string                   `json:"name"`
	Subject   string                   `json:"subject"`
	Content   string                   `json:"content"`
	IsHTML    bool                     `json:"is_html"`
	Variables []string                 `json:"variables"`
	Sort      int                      `json:"sort"`
	Status    model.MailTemplateStatus `json:"status"`
	Remark    string                   `json:"remark"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type TemplateListResponse struct {
	Items    []TemplateResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type RenderTemplateRequest struct {
	Variables map[string]string `json:"variables"`
}

type RenderTemplateResponse struct {
	Subject string `json:"subject"`
	Content string `json:"content"`
}

type SendMailRequest struct {
	AccountID    uint              `json:"account_id"`
	TemplateCode string            `json:"template_code"`
	To           []string          `json:"to"`
	Cc           []string          `json:"cc"`
	Bcc          []string          `json:"bcc"`
	Subject      string            `json:"subject"`
	Content      string            `json:"content"`
	IsHTML       bool              `json:"is_html"`
	Variables    map[string]string `json:"variables"`
}

type SendMailResponse struct {
	LogID  uint                `json:"log_id"`
	Status model.MailLogStatus `json:"status"`
}

type TestAccountRequest struct {
	To      []string `json:"to"`
	Subject string   `json:"subject"`
	Content string   `json:"content"`
}

type LogListQuery struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Keyword      string `form:"keyword"`
	Status       int    `form:"status"`
	AccountID    uint   `form:"account_id"`
	TemplateCode string `form:"template_code"`
}

type LogResponse struct {
	ID           uint                `json:"id"`
	AccountID    uint                `json:"account_id"`
	AccountName  string              `json:"account_name"`
	TemplateID   uint                `json:"template_id"`
	TemplateCode string              `json:"template_code"`
	Subject      string              `json:"subject"`
	FromEmail    string              `json:"from_email"`
	ToEmails     []string            `json:"to_emails"`
	CcEmails     []string            `json:"cc_emails"`
	BccEmails    []string            `json:"bcc_emails"`
	Status       model.MailLogStatus `json:"status"`
	ErrorMessage string              `json:"error_message"`
	CreatedAt    time.Time           `json:"created_at"`
}

type LogListResponse struct {
	Items    []LogResponse `json:"items"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type AccountEntity = model.MailAccount
type TemplateEntity = model.MailTemplate
type LogEntity = model.MailLog

const (
	PermissionAccountList         = "system:mail:account:list"
	PermissionAccountCreate       = "system:mail:account:create"
	PermissionAccountUpdate       = "system:mail:account:update"
	PermissionAccountUpdateStatus = "system:mail:account:status"
	PermissionAccountDelete       = "system:mail:account:delete"
	PermissionAccountTest         = "system:mail:account:test"

	PermissionTemplateList         = "system:mail:template:list"
	PermissionTemplateCreate       = "system:mail:template:create"
	PermissionTemplateUpdate       = "system:mail:template:update"
	PermissionTemplateUpdateStatus = "system:mail:template:status"
	PermissionTemplateDelete       = "system:mail:template:delete"
	PermissionTemplateRender       = "system:mail:template:render"

	PermissionSend = "system:mail:send"
	PermissionLog  = "system:mail:log"
)

func NormalizeCreateAccountRequest(req CreateAccountRequest) (CreateAccountRequest, error) {
	name, err := normalizeRequiredText("邮箱名称", req.Name, 64)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	host, err := normalizeRequiredText("SMTP 主机", req.Host, 128)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	port, err := normalizePort(req.Port)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	fromEmail, err := normalizeEmail("发件邮箱", req.FromEmail)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	fromName, err := normalizeOptionalText("发件人名称", req.FromName, 64)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	encryption, err := NormalizeEncryption(req.Encryption)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	status, err := NormalizeAccountStatus(req.Status, true)
	if err != nil {
		return CreateAccountRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateAccountRequest{}, err
	}

	req.Name = name
	req.Host = host
	req.Port = port
	req.Username = strings.TrimSpace(req.Username)
	req.Password = strings.TrimSpace(req.Password)
	req.FromEmail = fromEmail
	req.FromName = fromName
	req.Encryption = encryption
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateAccountRequest(req UpdateAccountRequest) (UpdateAccountRequest, error) {
	createLike := CreateAccountRequest{
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
	normalized, err := NormalizeCreateAccountRequest(createLike)
	if err != nil {
		return UpdateAccountRequest{}, err
	}
	return UpdateAccountRequest(normalized), nil
}

func NormalizeAccountStatus(status model.MailAccountStatus, allowDefault bool) (model.MailAccountStatus, error) {
	if status == 0 && allowDefault {
		status = model.MailAccountStatusEnabled
	}
	if status != model.MailAccountStatusEnabled && status != model.MailAccountStatusDisabled {
		return 0, errorsx.BadRequest("邮箱状态不正确")
	}
	return status, nil
}

func NormalizeAccountStatusFilter(value int) (*model.MailAccountStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status, err := NormalizeAccountStatus(model.MailAccountStatus(value), false)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func NormalizeEncryption(value model.MailEncryption) (model.MailEncryption, error) {
	value = model.MailEncryption(strings.ToLower(strings.TrimSpace(string(value))))
	if value == "" {
		return model.MailEncryptionNone, nil
	}
	switch value {
	case model.MailEncryptionNone, model.MailEncryptionSSL, model.MailEncryptionSTARTTLS:
		return value, nil
	default:
		return "", errorsx.BadRequest("邮箱加密方式不正确")
	}
}

func NormalizeCreateTemplateRequest(req CreateTemplateRequest) (CreateTemplateRequest, error) {
	code, err := normalizeCode(req.Code)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	name, err := normalizeRequiredText("模板名称", req.Name, 64)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	subject, err := normalizeRequiredText("邮件主题", req.Subject, 200)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return CreateTemplateRequest{}, errorsx.BadRequest("邮件正文不能为空")
	}
	status, err := NormalizeTemplateStatus(req.Status, true)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	variables, err := NormalizeVariables(req.Variables)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateTemplateRequest{}, err
	}

	req.Code = code
	req.Name = name
	req.Subject = subject
	req.Content = content
	req.Variables = variables
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateTemplateRequest(req UpdateTemplateRequest) (UpdateTemplateRequest, error) {
	name, err := normalizeRequiredText("模板名称", req.Name, 64)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	subject, err := normalizeRequiredText("邮件主题", req.Subject, 200)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		return UpdateTemplateRequest{}, errorsx.BadRequest("邮件正文不能为空")
	}
	status, err := NormalizeTemplateStatus(req.Status, false)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	variables, err := NormalizeVariables(req.Variables)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}

	req.Name = name
	req.Subject = subject
	req.Content = content
	req.Variables = variables
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeTemplateStatus(status model.MailTemplateStatus, allowDefault bool) (model.MailTemplateStatus, error) {
	if status == 0 && allowDefault {
		status = model.MailTemplateStatusEnabled
	}
	if status != model.MailTemplateStatusEnabled && status != model.MailTemplateStatusDisabled {
		return 0, errorsx.BadRequest("模板状态不正确")
	}
	return status, nil
}

func NormalizeTemplateStatusFilter(value int) (*model.MailTemplateStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status, err := NormalizeTemplateStatus(model.MailTemplateStatus(value), false)
	if err != nil {
		return nil, err
	}
	return &status, nil
}

func NormalizeVariables(values []string) ([]string, error) {
	seen := make(map[string]struct{}, len(values))
	result := make([]string, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}
		if len(value) > 64 {
			return nil, errorsx.BadRequest("模板变量名不能超过 64 个字符")
		}
		if !templateCodePattern.MatchString(value) {
			return nil, errorsx.BadRequest("模板变量名只能使用小写字母、数字、冒号、短横线和下划线")
		}
		if _, ok := seen[value]; ok {
			continue
		}
		seen[value] = struct{}{}
		result = append(result, value)
	}
	sort.Strings(result)
	return result, nil
}

func EncodeVariables(values []string) string {
	if len(values) == 0 {
		return "[]"
	}
	data, err := json.Marshal(values)
	if err != nil {
		return "[]"
	}
	return string(data)
}

func DecodeVariables(value string) []string {
	var result []string
	if err := json.Unmarshal([]byte(value), &result); err != nil {
		return []string{}
	}
	normalized, err := NormalizeVariables(result)
	if err != nil {
		return []string{}
	}
	return normalized
}

func NormalizeSendMailRequest(req SendMailRequest) (SendMailRequest, error) {
	to, err := NormalizeEmailList("收件人", req.To, true)
	if err != nil {
		return SendMailRequest{}, err
	}
	cc, err := NormalizeEmailList("抄送人", req.Cc, false)
	if err != nil {
		return SendMailRequest{}, err
	}
	bcc, err := NormalizeEmailList("密送人", req.Bcc, false)
	if err != nil {
		return SendMailRequest{}, err
	}

	req.TemplateCode = strings.TrimSpace(req.TemplateCode)
	if req.TemplateCode == "" {
		subject, err := normalizeRequiredText("邮件主题", req.Subject, 200)
		if err != nil {
			return SendMailRequest{}, err
		}
		content := strings.TrimSpace(req.Content)
		if content == "" {
			return SendMailRequest{}, errorsx.BadRequest("邮件正文不能为空")
		}
		req.Subject = subject
		req.Content = content
	} else {
		code, err := normalizeCode(req.TemplateCode)
		if err != nil {
			return SendMailRequest{}, err
		}
		req.TemplateCode = code
	}

	req.To = to
	req.Cc = cc
	req.Bcc = bcc
	if req.Variables == nil {
		req.Variables = map[string]string{}
	}
	return req, nil
}

func NormalizeTestAccountRequest(req TestAccountRequest) (TestAccountRequest, error) {
	to, err := NormalizeEmailList("收件人", req.To, true)
	if err != nil {
		return TestAccountRequest{}, err
	}
	subject := strings.TrimSpace(req.Subject)
	if subject == "" {
		subject = "EZ Admin 邮箱测试"
	}
	if len(subject) > 200 {
		return TestAccountRequest{}, errorsx.BadRequest("邮件主题不能超过 200 个字符")
	}
	content := strings.TrimSpace(req.Content)
	if content == "" {
		content = "这是一封来自 EZ Admin 的邮箱配置测试邮件。"
	}
	req.To = to
	req.Subject = subject
	req.Content = content
	return req, nil
}

func NormalizeLogStatusFilter(value int) (*model.MailLogStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status := model.MailLogStatus(value)
	if status != model.MailLogStatusSuccess && status != model.MailLogStatusFailed {
		return nil, errorsx.BadRequest("发送状态不正确")
	}
	return &status, nil
}

func NormalizeEmailList(label string, values []string, required bool) ([]string, error) {
	result := make([]string, 0, len(values))
	seen := make(map[string]struct{}, len(values))
	for _, value := range values {
		parts := strings.Split(value, ",")
		for _, part := range parts {
			email, err := normalizeEmail(label, part)
			if err != nil {
				if strings.TrimSpace(part) == "" {
					continue
				}
				return nil, err
			}
			key := strings.ToLower(email)
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			result = append(result, email)
		}
	}
	if required && len(result) == 0 {
		return nil, errorsx.BadRequest(label + "不能为空")
	}
	return result, nil
}

func JoinEmails(values []string) string {
	if len(values) == 0 {
		return ""
	}
	return strings.Join(values, ",")
}

func SplitEmails(value string) []string {
	if strings.TrimSpace(value) == "" {
		return []string{}
	}
	items := strings.Split(value, ",")
	result := make([]string, 0, len(items))
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	return result
}

func normalizeCode(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest("模板编码不能为空")
	}
	if len(value) > 64 {
		return "", errorsx.BadRequest("模板编码长度不能超过 64 个字符")
	}
	if !templateCodePattern.MatchString(value) {
		return "", errorsx.BadRequest("模板编码只能使用小写字母、数字、冒号、短横线和下划线")
	}
	return value, nil
}

func normalizeRequiredText(label string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(label + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(label + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizeOptionalText(label string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > maxLen {
		return "", errorsx.BadRequest(label + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizePort(port int) (int, error) {
	if port <= 0 || port > 65535 {
		return 0, errorsx.BadRequest("SMTP 端口不正确")
	}
	return port, nil
}

func normalizeEmail(label string, value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(label + "不能为空")
	}
	address, err := mail.ParseAddress(value)
	if err != nil || address.Address == "" {
		return "", errorsx.BadRequest(label + "格式不正确")
	}
	if len(address.Address) > 128 {
		return "", errorsx.BadRequest(label + "长度不能超过 128 个字符")
	}
	return address.Address, nil
}

func normalizeRemark(value string) (string, error) {
	value = strings.TrimSpace(value)
	if len(value) > 255 {
		return "", errorsx.BadRequest("备注不能超过 255 个字符")
	}
	return value, nil
}

func BuildAccountResponse(item model.MailAccount) AccountResponse {
	return AccountResponse{
		ID:          item.ID,
		Name:        item.Name,
		Host:        item.Host,
		Port:        item.Port,
		Username:    item.Username,
		FromEmail:   item.FromEmail,
		FromName:    item.FromName,
		Encryption:  item.Encryption,
		IsDefault:   item.IsDefault,
		Status:      item.Status,
		LastTestAt:  item.LastTestAt,
		LastTestMsg: item.LastTestMsg,
		Remark:      item.Remark,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}
}

func BuildTemplateResponse(item model.MailTemplate) TemplateResponse {
	return TemplateResponse{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Subject:   item.Subject,
		Content:   item.Content,
		IsHTML:    item.IsHTML,
		Variables: DecodeVariables(item.Variables),
		Sort:      item.Sort,
		Status:    item.Status,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func BuildLogResponse(item model.MailLog) LogResponse {
	return LogResponse{
		ID:           item.ID,
		AccountID:    item.AccountID,
		AccountName:  item.AccountName,
		TemplateID:   item.TemplateID,
		TemplateCode: item.TemplateCode,
		Subject:      item.Subject,
		FromEmail:    item.FromEmail,
		ToEmails:     SplitEmails(item.ToEmails),
		CcEmails:     SplitEmails(item.CcEmails),
		BccEmails:    SplitEmails(item.BccEmails),
		Status:       item.Status,
		ErrorMessage: item.ErrorMessage,
		CreatedAt:    item.CreatedAt,
	}
}
