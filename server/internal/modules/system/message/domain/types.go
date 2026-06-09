package domain

import (
	"regexp"
	"strconv"
	"strings"
	"time"

	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/model"
)

var codePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

type TemplateListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Type     int    `form:"type"`
	Status   int    `form:"status"`
}

type CreateTemplateRequest struct {
	Code      string                    `json:"code"`
	Name      string                    `json:"name"`
	Title     string                    `json:"title"`
	Content   string                    `json:"content"`
	Type      model.MessageTemplateType `json:"type"`
	Variables string                    `json:"variables"`
	Sort      int                       `json:"sort"`
	Status    model.MessageConfigStatus `json:"status"`
	IsSystem  bool                      `json:"is_system"`
	Remark    string                    `json:"remark"`
}

type UpdateTemplateRequest struct {
	Name      string                    `json:"name"`
	Title     string                    `json:"title"`
	Content   string                    `json:"content"`
	Type      model.MessageTemplateType `json:"type"`
	Variables string                    `json:"variables"`
	Sort      int                       `json:"sort"`
	Status    model.MessageConfigStatus `json:"status"`
	IsSystem  bool                      `json:"is_system"`
	Remark    string                    `json:"remark"`
}

type UpdateTemplateStatusRequest struct {
	Status model.MessageConfigStatus `json:"status"`
}

type TemplateResponse struct {
	ID        uint                      `json:"id"`
	Code      string                    `json:"code"`
	Name      string                    `json:"name"`
	Title     string                    `json:"title"`
	Content   string                    `json:"content"`
	Type      model.MessageTemplateType `json:"type"`
	Variables string                    `json:"variables"`
	Sort      int                       `json:"sort"`
	Status    model.MessageConfigStatus `json:"status"`
	IsSystem  bool                      `json:"is_system"`
	Remark    string                    `json:"remark"`
	CreatedAt time.Time                 `json:"created_at"`
	UpdatedAt time.Time                 `json:"updated_at"`
}

type TemplateListResponse struct {
	Items    []TemplateResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type ReminderListQuery struct {
	Page         int    `form:"page"`
	PageSize     int    `form:"page_size"`
	Keyword      string `form:"keyword"`
	TriggerEvent string `form:"trigger_event"`
	TemplateID   uint   `form:"template_id"`
	ReceiverType int    `form:"receiver_type"`
	Status       int    `form:"status"`
}

type CreateReminderRequest struct {
	Code           string                    `json:"code"`
	Name           string                    `json:"name"`
	TriggerEvent   string                    `json:"trigger_event"`
	TemplateID     uint                      `json:"template_id"`
	Channels       string                    `json:"channels"`
	ReceiverType   model.MessageReceiverType `json:"receiver_type"`
	ReceiverValues string                    `json:"receiver_values"`
	AdvanceMinutes int                       `json:"advance_minutes"`
	LinkURL        string                    `json:"link_url"`
	Sort           int                       `json:"sort"`
	Status         model.MessageConfigStatus `json:"status"`
	IsSystem       bool                      `json:"is_system"`
	Remark         string                    `json:"remark"`
}

type UpdateReminderRequest struct {
	Name           string                    `json:"name"`
	TriggerEvent   string                    `json:"trigger_event"`
	TemplateID     uint                      `json:"template_id"`
	Channels       string                    `json:"channels"`
	ReceiverType   model.MessageReceiverType `json:"receiver_type"`
	ReceiverValues string                    `json:"receiver_values"`
	AdvanceMinutes int                       `json:"advance_minutes"`
	LinkURL        string                    `json:"link_url"`
	Sort           int                       `json:"sort"`
	Status         model.MessageConfigStatus `json:"status"`
	IsSystem       bool                      `json:"is_system"`
	Remark         string                    `json:"remark"`
}

type UpdateReminderStatusRequest struct {
	Status model.MessageConfigStatus `json:"status"`
}

type ReminderResponse struct {
	ID             uint                      `json:"id"`
	Code           string                    `json:"code"`
	Name           string                    `json:"name"`
	TriggerEvent   string                    `json:"trigger_event"`
	TemplateID     uint                      `json:"template_id"`
	TemplateCode   string                    `json:"template_code"`
	TemplateName   string                    `json:"template_name"`
	Channels       string                    `json:"channels"`
	ReceiverType   model.MessageReceiverType `json:"receiver_type"`
	ReceiverValues string                    `json:"receiver_values"`
	AdvanceMinutes int                       `json:"advance_minutes"`
	LinkURL        string                    `json:"link_url"`
	Sort           int                       `json:"sort"`
	Status         model.MessageConfigStatus `json:"status"`
	IsSystem       bool                      `json:"is_system"`
	Remark         string                    `json:"remark"`
	CreatedAt      time.Time                 `json:"created_at"`
	UpdatedAt      time.Time                 `json:"updated_at"`
}

type ReminderListResponse struct {
	Items    []ReminderResponse `json:"items"`
	Total    int64              `json:"total"`
	Page     int                `json:"page"`
	PageSize int                `json:"page_size"`
}

type TemplateEntity = model.MessageTemplate
type ReminderEntity = model.MessageReminder

type ReminderListItem struct {
	Reminder ReminderEntity
	Template TemplateEntity
}

const (
	PermissionTemplateList         = "system:message:template:list"
	PermissionTemplateCreate       = "system:message:template:create"
	PermissionTemplateUpdate       = "system:message:template:update"
	PermissionTemplateUpdateStatus = "system:message:template:status"

	PermissionReminderList         = "system:message:reminder:list"
	PermissionReminderCreate       = "system:message:reminder:create"
	PermissionReminderUpdate       = "system:message:reminder:update"
	PermissionReminderUpdateStatus = "system:message:reminder:status"
)

func NormalizeCreateTemplateRequest(req CreateTemplateRequest) (CreateTemplateRequest, error) {
	code, err := normalizeCode("模板编码", req.Code, 128)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	name, err := normalizeName("模板名称", req.Name)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	title, err := normalizeText("消息标题", req.Title, 128, true)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	content, err := normalizeContent(req.Content)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	templateType, err := NormalizeTemplateType(req.Type, true)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, true)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	variables, err := normalizeLongText("变量说明", req.Variables, 1000)
	if err != nil {
		return CreateTemplateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateTemplateRequest{}, err
	}

	req.Code = code
	req.Name = name
	req.Title = title
	req.Content = content
	req.Type = templateType
	req.Variables = variables
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateTemplateRequest(req UpdateTemplateRequest) (UpdateTemplateRequest, error) {
	name, err := normalizeName("模板名称", req.Name)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	title, err := normalizeText("消息标题", req.Title, 128, true)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	content, err := normalizeContent(req.Content)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	templateType, err := NormalizeTemplateType(req.Type, false)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, false)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	variables, err := normalizeLongText("变量说明", req.Variables, 1000)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateTemplateRequest{}, err
	}

	req.Name = name
	req.Title = title
	req.Content = content
	req.Type = templateType
	req.Variables = variables
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeCreateReminderRequest(req CreateReminderRequest) (CreateReminderRequest, error) {
	code, err := normalizeCode("提醒编码", req.Code, 128)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	name, err := normalizeName("提醒名称", req.Name)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	triggerEvent, err := normalizeCode("触发事件", req.TriggerEvent, 128)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	channels, err := normalizeChannels(req.Channels)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	receiverType, err := NormalizeReceiverType(req.ReceiverType, true)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	receiverValues, err := normalizeReceiverValues(receiverType, req.ReceiverValues)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	advanceMinutes, err := normalizeAdvanceMinutes(req.AdvanceMinutes)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	linkURL, err := normalizeText("跳转链接", req.LinkURL, 255, false)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, true)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return CreateReminderRequest{}, err
	}
	if req.TemplateID == 0 {
		return CreateReminderRequest{}, errorsx.BadRequest("消息模板不能为空")
	}

	req.Code = code
	req.Name = name
	req.TriggerEvent = triggerEvent
	req.Channels = channels
	req.ReceiverType = receiverType
	req.ReceiverValues = receiverValues
	req.AdvanceMinutes = advanceMinutes
	req.LinkURL = linkURL
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeUpdateReminderRequest(req UpdateReminderRequest) (UpdateReminderRequest, error) {
	name, err := normalizeName("提醒名称", req.Name)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	triggerEvent, err := normalizeCode("触发事件", req.TriggerEvent, 128)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	channels, err := normalizeChannels(req.Channels)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	receiverType, err := NormalizeReceiverType(req.ReceiverType, false)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	receiverValues, err := normalizeReceiverValues(receiverType, req.ReceiverValues)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	advanceMinutes, err := normalizeAdvanceMinutes(req.AdvanceMinutes)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	linkURL, err := normalizeText("跳转链接", req.LinkURL, 255, false)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	status, err := NormalizeStatus(req.Status, false)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	remark, err := normalizeRemark(req.Remark)
	if err != nil {
		return UpdateReminderRequest{}, err
	}
	if req.TemplateID == 0 {
		return UpdateReminderRequest{}, errorsx.BadRequest("消息模板不能为空")
	}

	req.Name = name
	req.TriggerEvent = triggerEvent
	req.Channels = channels
	req.ReceiverType = receiverType
	req.ReceiverValues = receiverValues
	req.AdvanceMinutes = advanceMinutes
	req.LinkURL = linkURL
	req.Status = status
	req.Remark = remark
	return req, nil
}

func NormalizeStatus(status model.MessageConfigStatus, allowDefault bool) (model.MessageConfigStatus, error) {
	if status == 0 && allowDefault {
		status = model.MessageConfigStatusEnabled
	}
	if !ValidStatus(status) {
		return 0, errorsx.BadRequest("消息配置状态不正确")
	}
	return status, nil
}

func NormalizeStatusFilter(value int) (*model.MessageConfigStatus, error) {
	if value == 0 {
		return nil, nil
	}
	status := model.MessageConfigStatus(value)
	if !ValidStatus(status) {
		return nil, errorsx.BadRequest("消息配置状态不正确")
	}
	return &status, nil
}

func NormalizeTemplateType(value model.MessageTemplateType, allowDefault bool) (model.MessageTemplateType, error) {
	if value == 0 && allowDefault {
		value = model.MessageTemplateTypeNotification
	}
	if !ValidTemplateType(value) {
		return 0, errorsx.BadRequest("消息模板类型不正确")
	}
	return value, nil
}

func NormalizeTemplateTypeFilter(value int) (*model.MessageTemplateType, error) {
	if value == 0 {
		return nil, nil
	}
	templateType := model.MessageTemplateType(value)
	if !ValidTemplateType(templateType) {
		return nil, errorsx.BadRequest("消息模板类型不正确")
	}
	return &templateType, nil
}

func NormalizeReceiverType(value model.MessageReceiverType, allowDefault bool) (model.MessageReceiverType, error) {
	if value == 0 && allowDefault {
		value = model.MessageReceiverTypeRole
	}
	if !ValidReceiverType(value) {
		return 0, errorsx.BadRequest("接收人类型不正确")
	}
	return value, nil
}

func NormalizeReceiverTypeFilter(value int) (*model.MessageReceiverType, error) {
	if value == 0 {
		return nil, nil
	}
	receiverType := model.MessageReceiverType(value)
	if !ValidReceiverType(receiverType) {
		return nil, errorsx.BadRequest("接收人类型不正确")
	}
	return &receiverType, nil
}

func ValidStatus(status model.MessageConfigStatus) bool {
	return status == model.MessageConfigStatusEnabled || status == model.MessageConfigStatusDisabled
}

func ValidTemplateType(templateType model.MessageTemplateType) bool {
	return templateType == model.MessageTemplateTypeNotification ||
		templateType == model.MessageTemplateTypeTodo ||
		templateType == model.MessageTemplateTypeAlert
}

func ValidReceiverType(receiverType model.MessageReceiverType) bool {
	return receiverType == model.MessageReceiverTypeRole ||
		receiverType == model.MessageReceiverTypeUser ||
		receiverType == model.MessageReceiverTypeDepartment ||
		receiverType == model.MessageReceiverTypeInitiator ||
		receiverType == model.MessageReceiverTypeAssignee
}

func BuildTemplateResponse(item model.MessageTemplate) TemplateResponse {
	return TemplateResponse{
		ID:        item.ID,
		Code:      item.Code,
		Name:      item.Name,
		Title:     item.Title,
		Content:   item.Content,
		Type:      item.Type,
		Variables: item.Variables,
		Sort:      item.Sort,
		Status:    item.Status,
		IsSystem:  item.IsSystem,
		Remark:    item.Remark,
		CreatedAt: item.CreatedAt,
		UpdatedAt: item.UpdatedAt,
	}
}

func BuildReminderResponse(item model.MessageReminder, template model.MessageTemplate) ReminderResponse {
	return ReminderResponse{
		ID:             item.ID,
		Code:           item.Code,
		Name:           item.Name,
		TriggerEvent:   item.TriggerEvent,
		TemplateID:     item.TemplateID,
		TemplateCode:   template.Code,
		TemplateName:   template.Name,
		Channels:       item.Channels,
		ReceiverType:   item.ReceiverType,
		ReceiverValues: item.ReceiverValues,
		AdvanceMinutes: item.AdvanceMinutes,
		LinkURL:        item.LinkURL,
		Sort:           item.Sort,
		Status:         item.Status,
		IsSystem:       item.IsSystem,
		Remark:         item.Remark,
		CreatedAt:      item.CreatedAt,
		UpdatedAt:      item.UpdatedAt,
	}
}

func normalizeCode(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	if !codePattern.MatchString(value) {
		return "", errorsx.BadRequest(fieldName + "只能使用小写字母、数字、冒号、短横线和下划线")
	}
	return value, nil
}

func normalizeName(fieldName string, value string) (string, error) {
	return normalizeText(fieldName, value, 64, true)
}

func normalizeContent(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", errorsx.BadRequest("消息内容不能为空")
	}
	if len(value) > 4000 {
		return "", errorsx.BadRequest("消息内容不能超过 4000 个字符")
	}
	return value, nil
}

func normalizeLongText(fieldName string, value string, maxLen int) (string, error) {
	return normalizeText(fieldName, value, maxLen, false)
}

func normalizeText(fieldName string, value string, maxLen int, required bool) (string, error) {
	value = strings.TrimSpace(value)
	if required && value == "" {
		return "", errorsx.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", errorsx.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	return value, nil
}

func normalizeChannels(value string) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		value = "notification"
	}
	if len(value) > 128 {
		return "", errorsx.BadRequest("提醒渠道不能超过 128 个字符")
	}

	seen := map[string]bool{}
	channels := strings.Split(value, ",")
	result := make([]string, 0, len(channels))
	for _, channel := range channels {
		channel = strings.TrimSpace(channel)
		if channel == "" {
			continue
		}
		if !codePattern.MatchString(channel) {
			return "", errorsx.BadRequest("提醒渠道只能使用小写字母、数字、冒号、短横线和下划线")
		}
		if !seen[channel] {
			result = append(result, channel)
			seen[channel] = true
		}
	}
	if len(result) == 0 {
		return "", errorsx.BadRequest("提醒渠道不能为空")
	}
	return strings.Join(result, ","), nil
}

func normalizeReceiverValues(receiverType model.MessageReceiverType, value string) (string, error) {
	value = strings.TrimSpace(value)
	if receiverType == model.MessageReceiverTypeInitiator || receiverType == model.MessageReceiverTypeAssignee {
		return value, nil
	}
	if value == "" {
		return "", errorsx.BadRequest("接收人不能为空")
	}
	if len(value) > 1000 {
		return "", errorsx.BadRequest("接收人配置不能超过 1000 个字符")
	}
	return value, nil
}

func normalizeAdvanceMinutes(value int) (int, error) {
	if value < 0 {
		return 0, errorsx.BadRequest("提前提醒分钟数不能小于 0")
	}
	if value > 43200 {
		return 0, errorsx.BadRequest("提前提醒分钟数不能超过 43200")
	}
	return value, nil
}

func normalizeRemark(value string) (string, error) {
	return normalizeText("备注", value, 255, false)
}
