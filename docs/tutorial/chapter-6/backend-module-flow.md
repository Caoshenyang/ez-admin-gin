---
title: 后端模块接入流程
description: "按步骤说明新业务模块如何接入后端，从 Model 到 Handler 到路由注册。"
---

# 后端模块接入流程

这一页会把一个新业务模块接入后端的过程拆成三步：定义 Model、实现 Handler、注册路由。完成这三步之后，模块的增删改查接口就能通过统一的响应格式返回数据，并自动经过认证、日志和权限中间件的校验。

::: tip 🎯 本节目标
按固定顺序完成 Model → Handler → Router 三层接入，新增的接口可以直接用 `curl` 或前端页面验证。
:::

## 接入步骤总览

整个后端模块的接入只涉及三个位置：

| 步骤 | 涉及目录 | 做什么 |
| --- | --- | --- |
| 1. 定义 Model | `server/internal/model/` | 新建 GORM 结构体，声明表名和状态常量 |
| 2. 实现 Handler | `server/internal/handler/` | 新建 Handler 结构体，注入 `*gorm.DB` 和 `*zap.Logger`，编写业务方法 |
| 3. 注册路由 | `server/internal/router/router.go` | 构造 Handler 实例，把方法挂到路由分组上 |

本项目没有 Service 和 Repository 层，Handler 直接使用 `h.db` 完成数据库操作。这种结构在个人项目中足够简洁，前提是每个 Handler 内部保持清晰的职责边界。

::: details Go vs Java：为什么没有 Service 和 Repository
在 Java 后台项目里，通常会把数据库访问封装到 Repository，把业务逻辑放在 Service，Controller 只负责参数绑定和响应。这种分层在大团队协作时有明确的职责边界。

当前项目面向个人项目快速上线，Handler 直接操作 GORM 可以减少代码量和文件跳转。如果后续某个模块的业务逻辑变得复杂，可以在 Handler 内部拆分私有方法，而不需要提前引入额外层级。
:::

## Step 1：定义 Model

Model 负责定义数据库表结构和 JSON 序列化规则。每个 Model 需要包含三个部分：状态常量、GORM 结构体和 `TableName()` 方法。

下面是项目中已有的 Model 模式，以用户模型为例：

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// UserStatus 表示用户状态。
type UserStatus int

const (
	// UserStatusEnabled 表示用户可以正常登录。
	UserStatusEnabled UserStatus = 1
	// UserStatusDisabled 表示用户已被禁用。
	UserStatusDisabled UserStatus = 2
)

// User 是后台用户表模型。
type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Username     string         `gorm:"size:64;not null;uniqueIndex" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:64;not null;default:''" json:"nickname"`
	Status       UserStatus     `gorm:"type:smallint;not null;default:1" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

// TableName 固定用户表名，避免后续调整命名策略时影响已有表。
func (User) TableName() string {
	return "sys_user"
}
```

可以注意到几个固定写法：

- **状态常量**用自定义类型 + `const` 块表达，启用值为 `1`，禁用值为 `2`。
- **主键**使用 `gorm:"primaryKey"`，类型为 `uint`。
- **软删除**使用 `gorm.DeletedAt`，并在 JSON 中标记为 `json:"-"`，避免敏感字段暴露。
- **TableName()** 固定表名，防止 GORM 的命名策略变化影响已有表。

再看一个系统配置模型的例子，它的字段比用户模型多，但结构完全一致：

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemConfigStatus 表示系统配置状态。
type SystemConfigStatus int

const (
	// SystemConfigStatusEnabled 表示配置可用。
	SystemConfigStatusEnabled SystemConfigStatus = 1
	// SystemConfigStatusDisabled 表示配置已停用。
	SystemConfigStatusDisabled SystemConfigStatus = 2
)

// SystemConfig 是系统配置表模型。
type SystemConfig struct {
	ID        uint               `gorm:"primaryKey" json:"id"`
	GroupCode string             `gorm:"size:64;not null;index" json:"group_code"`
	ConfigKey string             `gorm:"column:config_key;size:128;not null;uniqueIndex" json:"key"`
	Name      string             `gorm:"size:64;not null" json:"name"`
	Value     string             `gorm:"type:text;not null" json:"value"`
	Sort      int                `gorm:"not null;default:0" json:"sort"`
	Status    SystemConfigStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark    string             `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt time.Time          `json:"created_at"`
	UpdatedAt time.Time          `json:"updated_at"`
	DeletedAt gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 固定系统配置表名。
func (SystemConfig) TableName() string {
	return "sys_config"
}
```

两个 Model 的共同特征：

| 要素 | 写法 |
| --- | --- |
| 状态类型 | `type XxxStatus int`，配合 `Enabled = 1`、`Disabled = 2` |
| 主键 | `ID uint \`gorm:"primaryKey"\` |
| 软删除 | `DeletedAt gorm.DeletedAt \`gorm:"index" json:"-"\` |
| 表名 | `func (Xxx) TableName() string { return "sys_xxx" }` |

::: tip 📌 新表通过迁移文件创建
本项目使用 golang-migrate 管理表结构，不在代码中使用 `AutoMigrate`。定义完 Model 后，在 `server/migrations/` 下编写对应的迁移文件，服务启动时自动执行。建表语句的参考格式可以查看 [数据库建表语句](/reference/database-ddl)。

建表时注意：带 `DeletedAt` 的表要同时创建 `deleted_at` 列和索引，否则软删除查询不会走索引。
:::

## Step 2：实现 Handler

Handler 是后端模块的核心，包含所有业务逻辑。每个 Handler 遵循固定的构造模式：结构体持有 `*gorm.DB` 和 `*zap.Logger`，通过构造函数注入。

### Handler 结构体和构造函数

```go
// XxxHandler 负责某某业务接口。
type XxxHandler struct {
    db  *gorm.DB
    log *zap.Logger
}

// NewXxxHandler 创建 Handler，由路由层调用。
func NewXxxHandler(db *gorm.DB, log *zap.Logger) *XxxHandler {
    return &XxxHandler{db: db, log: log}
}
```

如果模块需要 Redis 缓存（比如系统配置），构造函数会增加 `*goredis.Client` 参数。大多数业务模块只需要 DB 和 Logger 两个依赖。

### 请求和响应结构体

Handler 文件中通常会定义三组结构体，分别用于接收请求参数和构造响应数据：

```go
// 列表查询参数，通过 URL Query 绑定。
type xxxListQuery struct {
    Page     int    `form:"page"`
    PageSize int    `form:"page_size"`
    Keyword  string `form:"keyword"`
    Status   int    `form:"status"`
}

// 创建请求参数，通过 JSON Body 绑定。
type createXxxRequest struct {
    Name   string         `json:"name"`
    Status model.XxxStatus `json:"status"`
    Remark string         `json:"remark"`
}

// 列表响应结构，包含分页信息。
type xxxListResponse struct {
    Items    []xxxResponse `json:"items"`
    Total    int64         `json:"total"`
    Page     int           `json:"page"`
    PageSize int           `json:"page_size"`
}
```

### 常见方法模式

Handler 的核心方法通常包含四种操作：List（分页列表）、Create（创建）、Update（编辑）、UpdateStatus（状态变更）。下面是每个方法的标准流程，以系统配置 Handler 为例：

```go
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
```

这段代码比较长，但阅读时抓住下面几条主线即可：

| 方法 | 核心流程 |
| --- | --- |
| `List` | 绑定查询参数 → 构建条件 → 统计总数 → 分页查询 → 构造响应 |
| `Create` | 绑定 JSON → 校验参数 → 开启事务写入 → 返回新记录 |
| `Update` | 解析路径参数 ID → 绑定 JSON → 事务内先查再改 → 返回更新后记录 |
| `UpdateStatus` | 解析 ID → 校验状态值 → 事务内更新单个字段 → 返回 ID 和新状态 |

### 统一响应格式

所有接口的返回值都经过 `response` 包统一处理，调用方不需要自己构造 JSON：

```go
// 成功时返回数据和 code=0。
response.Success(c, data)

// 错误时自动判断是否为 apperror.Error：
//   - 是：使用其 Code、Message 和 HTTP Status
//   - 否：记录日志，返回 500 和通用错误信息
response.Error(c, err, h.log)
```

对应的响应体结构：

```json
{
  "code": 0,
  "message": "ok",
  "data": { ... }
}
```

### 错误处理模式

Handler 内部的错误处理遵循固定套路：

```go
// 参数校验失败 → 使用 apperror.BadRequest
response.Error(c, apperror.BadRequest("参数不正确"), h.log)

// 记录不存在 → 使用 apperror.NotFound
response.Error(c, apperror.NotFound("资源不存在"), h.log)

// 数据库操作失败 → 使用 apperror.Internal，附带底层错误
response.Error(c, apperror.Internal("查询失败", err), h.log)

// 事务内的混合错误 → 使用统一的错误写入函数
func writeXxxError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
    var appErr *apperror.Error
    if errors.As(err, &appErr) {
        response.Error(c, appErr, log)
        return
    }
    response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

`apperror` 包提供的常用错误构造函数：

| 函数 | HTTP 状态码 | 业务码 | 适用场景 |
| --- | --- | --- | --- |
| `BadRequest(msg)` | 400 | 40000 | 参数校验失败 |
| `Unauthorized(msg)` | 401 | 40100 | 未登录或 Token 过期 |
| `Forbidden(msg)` | 403 | 40300 | 无权限访问 |
| `NotFound(msg)` | 404 | 40400 | 资源不存在 |
| `Internal(msg, err)` | 500 | 50000 | 数据库或其他内部错误 |

## Step 3：注册路由

Handler 写好后，需要在路由文件中完成两件事：构造 Handler 实例、把方法绑定到路由。

打开 `server/internal/router/router.go`，找到对应的路由注册函数：

```go
package router

import (
	"ez-admin-gin/server/internal/config"
	authHandler "ez-admin-gin/server/internal/handler/auth"
	setupHandler "ez-admin-gin/server/internal/handler/setup"
	systemHandler "ez-admin-gin/server/internal/handler/system"
	appLogger "ez-admin-gin/server/internal/logger"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/permission"
	"ez-admin-gin/server/internal/token"

	"github.com/gin-gonic/gin"
	goredis "github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// Options 汇总路由层需要依赖的对象。
type Options struct {
	Config     *config.Config
	Log        *zap.Logger
	DB         *gorm.DB
	Redis      *goredis.Client
	Token      *token.Manager
	Permission *permission.Enforcer
}

// New 创建路由引擎，并统一注册中间件和路由分组。
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	// 配置上传最大内存
	if opts.Config.Upload.MaxSizeMB > 0 {
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20
	}
	// 配置静态文件服务
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)

	registerSystemRoutes(r, opts)
	registerAuthRoutes(r, opts)
	registerSetupRoutes(r, opts)

	return r
}

// registerSetupRoutes 注册系统初始化路由（无需认证）。
func registerSetupRoutes(r *gin.Engine, opts Options) {
	setup := setupHandler.NewSetupHandler(opts.DB, opts.Log)

	api := r.Group("/api/v1")
	setupGroup := api.Group("/setup")
	setupGroup.POST("/init", setup.Init)
}

// registerAuthRoutes 注册认证相关路由。
func registerAuthRoutes(r *gin.Engine, opts Options) {
	login := authHandler.NewLoginHandler(opts.DB, opts.Log, opts.Token)
	me := authHandler.NewMeHandler(opts.Log)
	menus := authHandler.NewMenuHandler(opts.DB, opts.Log)
	dashboard := authHandler.NewDashboardHandler(opts.Config, opts.DB, opts.Redis, opts.Log)

	api := r.Group("/api/v1")
	auth := api.Group("/auth")
	auth.POST("/login", login.Login)

	protectedAuth := auth.Group("")
	protectedAuth.Use(middleware.Auth(opts.Token, opts.Log))
	protectedAuth.GET("/me", me.Me)
	protectedAuth.GET("/menus", menus.Menus)
	protectedAuth.GET("/dashboard", dashboard.Dashboard)
}

// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
	files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log)
	operationLogs := systemHandler.NewOperationLogHandler(opts.DB, opts.Log)
	loginLogs := systemHandler.NewLoginLogHandler(opts.DB, opts.Log)
	notices := systemHandler.NewNoticeHandler(opts.DB, opts.Log)

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.OperationLog(opts.DB, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))

	system.GET("/health", health.Check)
	system.GET("/users", users.List)
	system.POST("/users", users.Create)
	system.POST("/users/:id/update", users.Update)
	system.POST("/users/:id/status", users.UpdateStatus)
	system.POST("/users/:id/roles", users.UpdateRoles)
	system.GET("/roles", roles.List)
	system.POST("/roles", roles.Create)
	system.POST("/roles/:id/update", roles.Update)
	system.POST("/roles/:id/status", roles.UpdateStatus)
	system.POST("/roles/:id/permissions", roles.UpdatePermissions)
	system.POST("/roles/:id/menus", roles.UpdateMenus)
	system.GET("/menus", menus.Tree)
	system.POST("/menus", menus.Create)
	system.POST("/menus/:id/update", menus.Update)
	system.POST("/menus/:id/status", menus.UpdateStatus)
	system.POST("/menus/:id/delete", menus.Delete)
	system.GET("/configs", configs.List)
	system.POST("/configs", configs.Create)
	system.POST("/configs/:id/update", configs.Update)
	system.POST("/configs/:id/status", configs.UpdateStatus)
	system.GET("/configs/value/:key", configs.Value)
	system.GET("/files", files.List)
	system.POST("/files", files.Upload)
	system.GET("/operation-logs", operationLogs.List)
	system.GET("/login-logs", loginLogs.List)
	system.GET("/notices", notices.List)
	system.POST("/notices", notices.Create)
	system.POST("/notices/:id/update", notices.Update)
	system.POST("/notices/:id/status", notices.UpdateStatus)
}
```

阅读这段代码时，关注三件事：

### 构造 Handler 实例

在 `registerSystemRoutes` 函数开头，每一行构造一个 Handler：

```go
configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
```

新增模块时，在已有构造语句下方加一行即可。

### 绑定路由和方法

路由绑定的格式是 `分组.HTTP方法("路径", handler.方法)`：

```go
system.GET("/configs", configs.List)
system.POST("/configs", configs.Create)
system.POST("/configs/:id/update", configs.Update)
system.POST("/configs/:id/status", configs.UpdateStatus)
```

注意当前项目的一个惯例：更新和删除操作使用 `POST` 方法而不是 `PUT` / `DELETE`，路径中用 `/:id/update` 和 `/:id/status` 区分操作类型。

### 中间件链

`system` 路由组挂载了三层中间件，每个请求都会按顺序经过：

```go
system.Use(middleware.Auth(opts.Token, opts.Log))          // 1. 认证：解析 JWT，获取当前用户
system.Use(middleware.OperationLog(opts.DB, opts.Log))      // 2. 操作日志：记录请求信息
system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log)) // 3. 权限：根据角色和接口路径判断是否放行
```

这意味着所有注册到 `system` 分组下的接口默认都需要登录、会被记录操作日志、并且需要对应的权限配置才能访问。

::: warning ⚠️ API 路径必须与权限定义一致
权限中间件使用 `c.FullPath()`（即路由注册时的模板路径）和 `c.Request.Method`（HTTP 方法）作为权限判断的对象。例如注册了 `system.GET("/configs", ...)` 后，权限表中需要配置 `GET` 方法对 `/api/v1/system/configs` 的访问权限。

如果路由路径改了但权限表没更新（或反过来），接口会返回 403。新增模块时，建议先确认路由路径，再在权限初始化数据中补齐对应记录。权限与菜单的具体配置方式在下一节 [权限、菜单与迁移接入](./permission-menu-migration) 中说明。
:::

## 验证清单

完成三步接入后，按下面的顺序检查：

1. **Model 文件存在**：`server/internal/model/` 下有对应的结构体文件，且 `TableName()` 返回的表名与数据库中一致。
2. **数据库表已创建**：通过数据库客户端确认表和索引已存在。
3. **Handler 文件存在**：`server/internal/handler/` 下有对应的 Handler 文件，包含至少 `List` 和 `Create` 方法。
4. **路由已注册**：在 `router.go` 中可以看到对应的构造语句和路由绑定。
5. **接口可访问**：启动后端服务后，使用 `curl` 或 API 工具调用接口，确认返回统一的 JSON 格式。

```bash
# 检查列表接口是否返回分页数据（需要先登录获取 Token）。
curl -s -H "Authorization: Bearer <token>" \
  http://localhost:8080/api/v1/system/configs | jq .
```

期望返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "items": [],
    "total": 0,
    "page": 1,
    "page_size": 10
  }
}
```

## 小结

后端模块接入的核心是固定的三步流程：

- **Model**：定义 GORM 结构体、状态常量和表名，同时编写迁移文件。
- **Handler**：构造函数注入 DB 和 Logger，方法内直接使用 `h.db` 操作数据库，通过 `response.Success` / `response.Error` 统一返回。
- **Router**：在 `router.go` 中构造 Handler 实例，把方法绑定到带中间件保护的路由分组。

这三步完成后，接口就已经具备认证、日志和权限保护。接下来需要为模块补齐权限配置和菜单入口：[权限、菜单与迁移接入](./permission-menu-migration)。
