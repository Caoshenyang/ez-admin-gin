---
title: 系统配置
description: "实现系统配置的管理接口，并为后续模块提供可复用的配置读取与缓存能力。"
---

# 系统配置

前面已经有了用户、角色和菜单管理。现在补上一块后台底座里很常见的能力：系统配置。它的目标不是替代环境变量，而是承载“可在管理台维护、可被业务代码读取”的普通业务配置。

::: tip 🎯 本节目标
完成后，系统会新增 `sys_config` 表；`super_admin` 可以管理配置项；后端可以按配置键读取启用中的配置值，并优先走 Redis 缓存。
:::

## 先说明边界

这一节的系统配置，适合放这类内容：

- 站点标题
- 默认上传目录
- 首页公告开关
- 某个模块的默认分页大小

不适合放这类内容：

- 数据库密码
- JWT 密钥
- 第三方平台 Access Key / Secret
- Redis 连接信息

::: warning ⚠️ 系统配置不是密钥管理
需要跟随部署环境变化、且具备敏感性的内容，仍然应该放在 `.env` 或配置文件中，由运维配置和环境变量管理。系统配置更适合放“业务可调参数”。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
docs/
└─ reference/
   └─ database-ddl.md

server/
├─ internal/
│  ├─ handler/
│  │  └─ system/
│  │     └─ configs.go
│  ├─ model/
│  │  └─ system_config.go
│  └─ router/
│     └─ router.go
└─ migrations/
   ├─ pgsql/
   │  └─ 000002_seed_data.up.sql
   └─ mysql/
      └─ 000002_seed_data.up.sql
```

| 位置 | 用途 |
| --- | --- |
| `docs/reference/database-ddl.md` | 补充 `sys_config` 建表语句 |
| `internal/model/system_config.go` | 定义系统配置模型 |
| `internal/handler/system/configs.go` | 提供配置管理与按键读取接口 |
| `internal/router/router.go` | 注册系统配置路由 |
| `migrations/{pgsql,mysql}/000002_seed_data.up.sql` | 初始化系统配置权限和菜单 |

## 先创建数据表

本节新增 `sys_config`，用于保存后台可维护的普通业务配置。

::: tip 建表 SQL
字段说明、配置值存储方式、索引设计和 PostgreSQL / MySQL 建表语句统一放在参考手册：[数据库建表语句 - `sys_config`](../../reference/database-ddl#sys-config)。
:::

## 接口规划

本节实现 5 个接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/configs` | 配置分页列表 |
| `POST` | `/api/v1/system/configs` | 创建配置 |
| `POST` | `/api/v1/system/configs/:id/update` | 编辑配置 |
| `POST` | `/api/v1/system/configs/:id/status` | 修改配置状态 |
| `GET` | `/api/v1/system/configs/value/:key` | 按配置键读取启用中的配置值 |

其中 `/api/v1/system/configs/value/:key` 有两个作用：

1. 给后续业务模块提供统一读取入口。
2. 方便这一节验证 Redis 缓存是否生效。

## 🛠️ 创建系统配置模型

创建 `server/internal/model/system_config.go`。这是新增文件，直接完整写入即可。

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
	ID         uint               `gorm:"primaryKey" json:"id"`
	GroupCode  string             `gorm:"size:64;not null;index" json:"group_code"`
	ConfigKey  string             `gorm:"column:config_key;size:128;not null;uniqueIndex" json:"key"`
	Name       string             `gorm:"size:64;not null" json:"name"`
	Value      string             `gorm:"type:text;not null" json:"value"`
	Sort       int                `gorm:"not null;default:0" json:"sort"`
	Status     SystemConfigStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark     string             `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt  time.Time          `json:"created_at"`
	UpdatedAt  time.Time          `json:"updated_at"`
	DeletedAt  gorm.DeletedAt     `gorm:"index" json:"-"`
}

// TableName 固定系统配置表名。
func (SystemConfig) TableName() string {
	return "sys_config"
}
```

## 🛠️ 创建系统配置 Handler

创建 `server/internal/handler/system/configs.go`。这是新增文件，直接完整写入即可。

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

::: warning ⚠️ 配置值统一按字符串存储
这一版不在配置表里区分整数、布尔值或 JSON。后续业务模块如果需要布尔值或数字，读取配置后再自行解析。
:::

::: details 为什么缓存异常不让接口直接失败
系统配置的缓存是性能优化，不是唯一数据源。即使 Redis 临时不可用，也应该还能回退到数据库读取；否则一个缓存故障会把整个后台管理能力一起拖下去。
:::

## 🛠️ 注册系统配置路由

修改 `server/internal/router/router.go`。在系统路由里新增系统配置 Handler 和路由：

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log) // [!code ++]

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
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
	system.GET("/configs", configs.List) // [!code ++]
	system.POST("/configs", configs.Create) // [!code ++]
	system.POST("/configs/:id/update", configs.Update) // [!code ++]
	system.POST("/configs/:id/status", configs.UpdateStatus) // [!code ++]
	system.GET("/configs/value/:key", configs.Value) // [!code ++]
}
```

## 🛠️ 初始化系统配置权限和菜单

系统配置的权限和菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建系统配置相关的权限策略和菜单数据。

::: tip 💡 权限和菜单初始化
- 权限策略：在 `migrations/{pgsql,mysql}/000002_seed_data.up.sql` 中插入系统配置接口的 Casbin 规则
- 菜单数据：在同一迁移文件中插入系统配置菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到系统配置菜单
:::

::: details 为什么配置键创建后不允许修改
配置键会被业务代码、缓存键和接口路径共同依赖。创建后保持稳定，后续排查问题会轻松很多；如果真的要改，一般会单独做迁移或改名方案，而不是在后台里随手改。
:::

## ✅ 启动并观察初始化日志

本节没有新增第三方依赖，可以直接启动：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/configs", "method": "GET"}
INFO	default menu created	{"menu_code": "system:config"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 15}
```

## ✅ 验证权限和菜单数据

先确认系统配置接口权限已经写入：

```bash
# 查看系统配置相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/configs%' order by v1, v2;"
```

应该能看到 `GET` 和 `POST` 对应的策略。

再确认系统配置菜单和按钮已经写入：

```bash
# 查看系统配置菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:config%' order by sort, id;"
```

应该能看到 `system:config` 以及几个 `system:config:*` 按钮编码。

## ✅ 验证系统配置接口

先登录拿到 Token：

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "admin"
  password = "EzAdmin@123456"
} | ConvertTo-Json

$login = Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/auth/login `
  -ContentType "application/json" `
  -Body $body

$token = $login.data.access_token
```

```bash [macOS / Linux]
TOKEN=$(curl -s -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"EzAdmin@123456"}' | jq -r '.data.access_token')
```

:::

创建一条站点标题配置：

::: warning ⚠️ Windows PowerShell 发送中文 JSON 时仍然建议显式使用 UTF-8
如果请求体里包含中文，继续使用 UTF-8 字节发送最稳妥。
:::

::: code-group

```powershell [Windows PowerShell]
$body = @{
  group_code = "site"
  key = "site:title"
  name = "站点标题"
  value = "EZ Admin"
  sort = 10
  status = 1
  remark = "站点基础配置"
} | ConvertTo-Json

$utf8Body = [System.Text.Encoding]::UTF8.GetBytes($body)

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/system/configs `
  -ContentType "application/json; charset=utf-8" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $utf8Body
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/system/configs \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"group_code":"site","key":"site:title","name":"站点标题","value":"EZ Admin","sort":10,"status":1,"remark":"站点基础配置"}'
```

:::

创建成功后，查询列表：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/configs?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/configs?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到刚创建的 `site:title`。

## ✅ 验证按键读取与缓存

第一次读取配置值：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri http://localhost:8080/api/v1/system/configs/value/site:title `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl http://localhost:8080/api/v1/system/configs/value/site:title \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

第一次通常会返回：

```json
{
  "code": 0,
  "message": "ok",
  "data": {
    "key": "site:title",
    "value": "EZ Admin",
    "source": "db"
  }
}
```

紧接着再请求一次，同一个接口应该优先命中 Redis，`source` 变成 `cache`。

如果你想直接看 Redis，也可以执行：

```bash
# 直接查看 Redis 中是否已经写入缓存
docker compose -f deploy/compose.local.yml exec redis redis-cli GET sys_config:site:title
```

应该能看到 `EZ Admin`。

## ✅ 验证禁用后缓存失效

把刚创建的配置状态改成禁用。下面示例假设这条配置的 ID 是 `1`，实际请替换成你自己的返回值：

::: code-group

```powershell [Windows PowerShell]
$configId = 1
$body = @{ status = 2 } | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/configs/$configId/status" `
  -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $body
```

```bash [macOS / Linux]
CONFIG_ID=1

curl -X POST "http://localhost:8080/api/v1/system/configs/${CONFIG_ID}/status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"status":2}'
```

:::

然后再次读取配置值：

```bash
curl http://localhost:8080/api/v1/system/configs/value/site:title \
  -H "Authorization: Bearer ${TOKEN}"
```

这时应该返回“配置不存在或已禁用”。同时 Redis 中的缓存键也应该已经被删除：

```bash
docker compose -f deploy/compose.local.yml exec redis redis-cli GET sys_config:site:title
```

这时应该返回空结果。

## 常见问题

::: details 创建配置时提示“配置键已存在”
换一个新的配置键即可，例如从 `site:title` 改成 `site:subtitle`。配置键唯一规则见：[数据库建表语句 - `sys_config`](../../reference/database-ddl#sys-config)。
:::

::: details 为什么配置值接口没有直接返回完整配置对象
按键读取接口更偏向“给业务代码拿值”，所以只返回 `key`、`value` 和当前读取来源。后台管理页如果要看完整信息，直接走列表接口即可。
:::

::: details 为什么禁用配置时要顺便删缓存
如果只改数据库状态，不清缓存，短时间内业务代码还能继续读到旧值。状态和缓存同步，才能避免“后台显示已禁用，但系统还在继续使用”的错觉。
:::

下一节继续补齐文件能力：[文件上传](./file-upload)。
