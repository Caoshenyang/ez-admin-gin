package system

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	systemConfigCachePrefix = "sys_config:"
	systemConfigCacheTTL    = time.Hour
)

var systemConfigCodePattern = regexp.MustCompile(`^[a-z0-9:_-]+$`)

// SystemConfigHandler 负责系统配置管理接口。
type SystemConfigHandler struct {
	db    *gorm.DB
	redis *goredis.Client
	log   *zap.Logger
}

// NewSystemConfigHandler 创建系统配置 Handler。
func NewSystemConfigHandler(db *gorm.DB, redis *goredis.Client, log *zap.Logger) *SystemConfigHandler {
	return &SystemConfigHandler{
		db:    db,
		redis: redis,
		log:   log,
	}
}

type systemConfigListQuery struct {
	Page      int    `form:"page"`
	PageSize  int    `form:"page_size"`
	Keyword   string `form:"keyword"`
	GroupCode string `form:"group_code"`
	Status    int    `form:"status"`
}

type createSystemConfigRequest struct {
	GroupCode string                   `json:"group_code"`
	Key       string                   `json:"key"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type updateSystemConfigRequest struct {
	GroupCode string                   `json:"group_code"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
}

type updateSystemConfigStatusRequest struct {
	Status model.SystemConfigStatus `json:"status"`
}

type systemConfigResponse struct {
	ID        uint                     `json:"id"`
	GroupCode string                   `json:"group_code"`
	Key       string                   `json:"key"`
	Name      string                   `json:"name"`
	Value     string                   `json:"value"`
	Sort      int                      `json:"sort"`
	Status    model.SystemConfigStatus `json:"status"`
	Remark    string                   `json:"remark"`
	CreatedAt time.Time                `json:"created_at"`
	UpdatedAt time.Time                `json:"updated_at"`
}

type systemConfigListResponse struct {
	Items    []systemConfigResponse `json:"items"`
	Total    int64                  `json:"total"`
	Page     int                    `json:"page"`
	PageSize int                    `json:"page_size"`
}

type systemConfigValueResponse struct {
	Key    string `json:"key"`
	Value  string `json:"value"`
	Source string `json:"source"`
}

// List 返回系统配置分页列表。
func (h *SystemConfigHandler) List(c *gin.Context) {
	var query systemConfigListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeSystemConfigPage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.SystemConfig{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("config_key LIKE ? OR name LIKE ?", like, like)
	}

	groupCode := strings.TrimSpace(query.GroupCode)
	if groupCode != "" {
		queryDB = queryDB.Where("group_code = ?", groupCode)
	}

	if query.Status != 0 {
		status := model.SystemConfigStatus(query.Status)
		if !validSystemConfigStatus(status) {
			response.Error(c, apperror.BadRequest("配置状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询配置总数失败", err), h.log)
		return
	}

	var configs []model.SystemConfig
	if err := queryDB.
		Order("group_code ASC, sort ASC, id ASC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&configs).Error; err != nil {
		response.Error(c, apperror.Internal("查询配置列表失败", err), h.log)
		return
	}

	items := make([]systemConfigResponse, 0, len(configs))
	for _, config := range configs {
		items = append(items, buildSystemConfigResponse(config))
	}

	response.Success(c, systemConfigListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建系统配置。
func (h *SystemConfigHandler) Create(c *gin.Context) {
	var req createSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	groupCode, key, name, status, remark, err := normalizeCreateSystemConfigRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	config := model.SystemConfig{
		GroupCode: groupCode,
		ConfigKey: key,
		Name:      name,
		Value:     req.Value,
		Sort:      req.Sort,
		Status:    status,
		Remark:    remark,
	}

	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureSystemConfigKeyAvailable(tx, config.ConfigKey); err != nil {
			return err
		}

		return tx.Create(&config).Error
	}); err != nil {
		writeSystemConfigError(c, err, "创建系统配置失败", h.log)
		return
	}

	h.syncSystemConfigCache(c, config)
	response.Success(c, buildSystemConfigResponse(config))
}

// Update 编辑系统配置。
func (h *SystemConfigHandler) Update(c *gin.Context) {
	configID, ok := systemConfigIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateSystemConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	groupCode, name, status, remark, err := normalizeUpdateSystemConfigRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	var config model.SystemConfig
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&config, configID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("配置不存在")
			}
			return err
		}

		if err := tx.Model(&config).Updates(map[string]any{
			"group_code": groupCode,
			"name":       name,
			"value":      req.Value,
			"sort":       req.Sort,
			"status":     status,
			"remark":     remark,
		}).Error; err != nil {
			return err
		}

		config.GroupCode = groupCode
		config.Name = name
		config.Value = req.Value
		config.Sort = req.Sort
		config.Status = status
		config.Remark = remark
		return nil
	}); err != nil {
		writeSystemConfigError(c, err, "更新系统配置失败", h.log)
		return
	}

	h.syncSystemConfigCache(c, config)
	response.Success(c, buildSystemConfigResponse(config))
}

// UpdateStatus 修改系统配置状态。
func (h *SystemConfigHandler) UpdateStatus(c *gin.Context) {
	configID, ok := systemConfigIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateSystemConfigStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validSystemConfigStatus(req.Status) {
		response.Error(c, apperror.BadRequest("配置状态不正确"), h.log)
		return
	}

	var config model.SystemConfig
	if err := h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&config, configID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("配置不存在")
			}
			return err
		}

		if err := tx.Model(&config).Update("status", req.Status).Error; err != nil {
			return err
		}

		config.Status = req.Status
		return nil
	}); err != nil {
		writeSystemConfigError(c, err, "更新配置状态失败", h.log)
		return
	}

	h.syncSystemConfigCache(c, config)
	response.Success(c, gin.H{
		"id":     configID,
		"status": req.Status,
	})
}

// Value 按配置键读取启用中的配置值，优先走 Redis 缓存。
func (h *SystemConfigHandler) Value(c *gin.Context) {
	key := strings.TrimSpace(c.Param("key"))
	if err := validateSystemConfigCode("配置键", key, 128); err != nil {
		response.Error(c, err, h.log)
		return
	}

	if h.redis != nil {
		value, err := h.redis.Get(c.Request.Context(), h.systemConfigCacheKey(key)).Result()
		if err == nil {
			response.Success(c, systemConfigValueResponse{
				Key:    key,
				Value:  value,
				Source: "cache",
			})
			return
		}

		if !errors.Is(err, goredis.Nil) {
			h.log.Warn("get system config cache failed", zap.String("key", key), zap.Error(err))
		}
	}

	var config model.SystemConfig
	if err := h.db.
		Where("config_key = ?", key).
		Where("status = ?", model.SystemConfigStatusEnabled).
		First(&config).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Error(c, apperror.NotFound("配置不存在或已禁用"), h.log)
			return
		}

		response.Error(c, apperror.Internal("读取系统配置失败", err), h.log)
		return
	}

	h.writeSystemConfigCache(c, config)
	response.Success(c, systemConfigValueResponse{
		Key:    config.ConfigKey,
		Value:  config.Value,
		Source: "db",
	})
}

func normalizeCreateSystemConfigRequest(req createSystemConfigRequest) (string, string, string, model.SystemConfigStatus, string, error) {
	groupCode, err := normalizeSystemConfigCode("配置分组", req.GroupCode, 64)
	if err != nil {
		return "", "", "", 0, "", err
	}

	key, err := normalizeSystemConfigCode("配置键", req.Key, 128)
	if err != nil {
		return "", "", "", 0, "", err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", "", "", 0, "", apperror.BadRequest("配置名称不能为空")
	}
	if len(name) > 64 {
		return "", "", "", 0, "", apperror.BadRequest("配置名称不能超过 64 个字符")
	}

	status := req.Status
	if status == 0 {
		status = model.SystemConfigStatusEnabled
	}
	if !validSystemConfigStatus(status) {
		return "", "", "", 0, "", apperror.BadRequest("配置状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return groupCode, key, name, status, remark, nil
}

func normalizeUpdateSystemConfigRequest(req updateSystemConfigRequest) (string, string, model.SystemConfigStatus, string, error) {
	groupCode, err := normalizeSystemConfigCode("配置分组", req.GroupCode, 64)
	if err != nil {
		return "", "", 0, "", err
	}

	name := strings.TrimSpace(req.Name)
	if name == "" {
		return "", "", 0, "", apperror.BadRequest("配置名称不能为空")
	}
	if len(name) > 64 {
		return "", "", 0, "", apperror.BadRequest("配置名称不能超过 64 个字符")
	}

	if !validSystemConfigStatus(req.Status) {
		return "", "", 0, "", apperror.BadRequest("配置状态不正确")
	}

	remark := strings.TrimSpace(req.Remark)
	if len(remark) > 255 {
		return "", "", 0, "", apperror.BadRequest("备注不能超过 255 个字符")
	}

	return groupCode, name, req.Status, remark, nil
}

func normalizeSystemConfigCode(fieldName string, value string, maxLen int) (string, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return "", apperror.BadRequest(fieldName + "不能为空")
	}
	if len(value) > maxLen {
		return "", apperror.BadRequest(fieldName + "长度不能超过 " + strconv.Itoa(maxLen) + " 个字符")
	}
	if !systemConfigCodePattern.MatchString(value) {
		return "", apperror.BadRequest(fieldName + "只能使用小写字母、数字、冒号、短横线和下划线")
	}

	return value, nil
}

func validateSystemConfigCode(fieldName string, value string, maxLen int) error {
	_, err := normalizeSystemConfigCode(fieldName, value, maxLen)
	return err
}

func normalizeSystemConfigPage(page int, pageSize int) (int, int) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	return page, pageSize
}

func validSystemConfigStatus(status model.SystemConfigStatus) bool {
	return status == model.SystemConfigStatusEnabled || status == model.SystemConfigStatusDisabled
}

func systemConfigIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("配置 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func ensureSystemConfigKeyAvailable(db *gorm.DB, key string) error {
	var config model.SystemConfig
	err := db.Unscoped().Where("config_key = ?", key).First(&config).Error
	if err == nil {
		return apperror.BadRequest("配置键已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func buildSystemConfigResponse(config model.SystemConfig) systemConfigResponse {
	return systemConfigResponse{
		ID:        config.ID,
		GroupCode: config.GroupCode,
		Key:       config.ConfigKey,
		Name:      config.Name,
		Value:     config.Value,
		Sort:      config.Sort,
		Status:    config.Status,
		Remark:    config.Remark,
		CreatedAt: config.CreatedAt,
		UpdatedAt: config.UpdatedAt,
	}
}

func (h *SystemConfigHandler) systemConfigCacheKey(key string) string {
	return systemConfigCachePrefix + key
}

func (h *SystemConfigHandler) writeSystemConfigCache(c *gin.Context, config model.SystemConfig) {
	if h.redis == nil {
		return
	}

	if err := h.redis.Set(
		c.Request.Context(),
		h.systemConfigCacheKey(config.ConfigKey),
		config.Value,
		systemConfigCacheTTL,
	).Err(); err != nil {
		h.log.Warn("set system config cache failed", zap.String("key", config.ConfigKey), zap.Error(err))
	}
}

func (h *SystemConfigHandler) deleteSystemConfigCache(c *gin.Context, key string) {
	if h.redis == nil {
		return
	}

	if err := h.redis.Del(c.Request.Context(), h.systemConfigCacheKey(key)).Err(); err != nil {
		h.log.Warn("delete system config cache failed", zap.String("key", key), zap.Error(err))
	}
}

func (h *SystemConfigHandler) syncSystemConfigCache(c *gin.Context, config model.SystemConfig) {
	if config.Status == model.SystemConfigStatusEnabled {
		h.writeSystemConfigCache(c, config)
		return
	}

	h.deleteSystemConfigCache(c, config.ConfigKey)
}

func writeSystemConfigError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
