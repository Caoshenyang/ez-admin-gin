---
title: 用户管理
description: "实现后台用户的列表、创建、编辑、禁用和角色分配能力。"
---

# 用户管理

第 3 章已经有了默认管理员和登录能力。现在把 `sys_user` 真正变成可管理的系统模块：能查用户列表，能创建用户，能编辑用户状态和昵称，也能给用户分配角色。

::: tip 🎯 本节目标
完成后，`super_admin` 可以访问用户管理接口；系统会初始化用户管理菜单和按钮；通过接口可以完成用户列表、创建用户、编辑用户、禁用用户和分配角色。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
server/
├─ internal/
│  ├─ handler/
│  │  └─ system/
│  │     └─ users.go
│  └─ router/
│     └─ router.go
└─ migrations/
   ├─ postgres/
   │  └─ 000002_seed_data.up.sql
   └─ mysql/
      └─ 000002_seed_data.up.sql
```

| 位置 | 用途 |
| --- | --- |
| `internal/handler/system/users.go` | 用户管理接口 |
| `internal/router/router.go` | 注册用户管理路由 |
| `migrations/{postgres,mysql}/000002_seed_data.up.sql` | 初始化用户管理权限和菜单 |

::: info 本节不新增数据库表
用户管理复用前面已经创建的 `sys_user`、`sys_role`、`sys_user_role`。如果本地还没有这些表，先回到参考手册执行对应建表语句：[数据库建表语句](../../reference/database-ddl)。
:::

## 接口规划

本节先实现后台管理常用的 5 个接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/users` | 用户列表 |
| `POST` | `/api/v1/system/users` | 创建用户 |
| `POST` | `/api/v1/system/users/:id/update` | 编辑用户基础信息 |
| `POST` | `/api/v1/system/users/:id/status` | 修改用户状态 |
| `POST` | `/api/v1/system/users/:id/roles` | 分配用户角色 |

::: warning ⚠️ 用户管理接口必须走权限校验
这些接口都会挂在 `/api/v1/system` 分组下，需要先登录，再通过 Casbin 权限判断。只创建路由但不补 `casbin_rule`，请求会返回 `403`。
:::

## 🛠️ 创建用户管理 Handler

创建 `server/internal/handler/system/users.go`。这是新增文件，直接完整写入即可。

```go
package system

import (
	"errors"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// UserHandler 负责后台用户管理接口。
type UserHandler struct {
	db  *gorm.DB
	log *zap.Logger
}

// NewUserHandler 创建用户管理 Handler。
func NewUserHandler(db *gorm.DB, log *zap.Logger) *UserHandler {
	return &UserHandler{
		db:  db,
		log: log,
	}
}

type userListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Status   int    `form:"status"`
}

type createUserRequest struct {
	Username string           `json:"username"`
	Password string           `json:"password"`
	Nickname string           `json:"nickname"`
	Status   model.UserStatus `json:"status"`
	RoleIDs  []uint           `json:"role_ids"`
}

type updateUserRequest struct {
	Nickname string           `json:"nickname"`
	Status   model.UserStatus `json:"status"`
}

type updateUserStatusRequest struct {
	Status model.UserStatus `json:"status"`
}

type updateUserRolesRequest struct {
	RoleIDs []uint `json:"role_ids"`
}

type userResponse struct {
	ID        uint             `json:"id"`
	Username  string           `json:"username"`
	Nickname  string           `json:"nickname"`
	Status    model.UserStatus `json:"status"`
	RoleIDs   []uint           `json:"role_ids"`
	CreatedAt time.Time        `json:"created_at"`
	UpdatedAt time.Time        `json:"updated_at"`
}

type userListResponse struct {
	Items    []userResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

// List 返回后台用户分页列表。
func (h *UserHandler) List(c *gin.Context) {
	var query userListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.User{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("username LIKE ? OR nickname LIKE ?", like, like)
	}

	if query.Status != 0 {
		status := model.UserStatus(query.Status)
		if !validUserStatus(status) {
			response.Error(c, apperror.BadRequest("用户状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询用户总数失败", err), h.log)
		return
	}

	var users []model.User
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error; err != nil {
		response.Error(c, apperror.Internal("查询用户列表失败", err), h.log)
		return
	}

	roleIDs, err := h.userRoleIDs(users)
	if err != nil {
		response.Error(c, apperror.Internal("查询用户角色失败", err), h.log)
		return
	}

	items := make([]userResponse, 0, len(users))
	for _, user := range users {
		items = append(items, buildUserResponse(user, roleIDs[user.ID]))
	}

	response.Success(c, userListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Create 创建后台用户。
func (h *UserHandler) Create(c *gin.Context) {
	var req createUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	username, password, nickname, status, roleIDs, err := normalizeCreateUserRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		response.Error(c, apperror.Internal("生成密码哈希失败", err), h.log)
		return
	}

	var created model.User
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := ensureUsernameAvailable(tx, username); err != nil {
			return err
		}

		if err := ensureRolesUsable(tx, roleIDs); err != nil {
			return err
		}

		user := model.User{
			Username:     username,
			PasswordHash: string(passwordHash),
			Nickname:     nickname,
			Status:       status,
		}

		if err := tx.Create(&user).Error; err != nil {
			return err
		}

		if err := replaceUserRoles(tx, user.ID, roleIDs); err != nil {
			return err
		}

		created = user
		return nil
	})
	if err != nil {
		writeError(c, err, "创建用户失败", h.log)
		return
	}

	response.Success(c, buildUserResponse(created, roleIDs))
}

// Update 编辑用户基础信息。
func (h *UserHandler) Update(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	nickname, status, err := normalizeUpdateUserRequest(req)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID && status == model.UserStatusDisabled {
		response.Error(c, apperror.BadRequest("不能禁用当前登录用户"), h.log)
		return
	}

	var user model.User
	err = h.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		if err := tx.Model(&user).Updates(map[string]any{
			"nickname": nickname,
			"status":   status,
		}).Error; err != nil {
			return err
		}

		user.Nickname = nickname
		user.Status = status
		return nil
	})
	if err != nil {
		writeError(c, err, "更新用户失败", h.log)
		return
	}

	roleIDs, err := h.userRoleIDs([]model.User{user})
	if err != nil {
		response.Error(c, apperror.Internal("查询用户角色失败", err), h.log)
		return
	}

	response.Success(c, buildUserResponse(user, roleIDs[user.ID]))
}

// UpdateStatus 修改用户启用状态。
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	var req updateUserStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	if !validUserStatus(req.Status) {
		response.Error(c, apperror.BadRequest("用户状态不正确"), h.log)
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID && req.Status == model.UserStatusDisabled {
		response.Error(c, apperror.BadRequest("不能禁用当前登录用户"), h.log)
		return
	}

	err := h.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		return tx.Model(&user).Update("status", req.Status).Error
	})
	if err != nil {
		writeError(c, err, "更新用户状态失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":     userID,
		"status": req.Status,
	})
}

// UpdateRoles 更新用户绑定的角色。
func (h *UserHandler) UpdateRoles(c *gin.Context) {
	userID, ok := userIDParam(c, h.log)
	if !ok {
		return
	}

	if currentUserID, ok := middleware.CurrentUserID(c); ok && currentUserID == userID {
		response.Error(c, apperror.BadRequest("不能修改当前登录用户的角色"), h.log)
		return
	}

	var req updateUserRolesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, apperror.BadRequest("请求参数不正确"), h.log)
		return
	}

	roleIDs, err := normalizeRoleIDs(req.RoleIDs)
	if err != nil {
		response.Error(c, err, h.log)
		return
	}

	err = h.db.Transaction(func(tx *gorm.DB) error {
		var user model.User
		if err := tx.First(&user, userID).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return apperror.NotFound("用户不存在")
			}
			return err
		}

		if err := ensureRolesUsable(tx, roleIDs); err != nil {
			return err
		}

		return replaceUserRoles(tx, userID, roleIDs)
	})
	if err != nil {
		writeError(c, err, "更新用户角色失败", h.log)
		return
	}

	response.Success(c, gin.H{
		"id":       userID,
		"role_ids": roleIDs,
	})
}

func (h *UserHandler) userRoleIDs(users []model.User) (map[uint][]uint, error) {
	result := make(map[uint][]uint, len(users))
	if len(users) == 0 {
		return result, nil
	}

	userIDs := make([]uint, 0, len(users))
	for _, user := range users {
		userIDs = append(userIDs, user.ID)
	}

	var rows []model.UserRole
	if err := h.db.Where("user_id IN ?", userIDs).Order("role_id ASC").Find(&rows).Error; err != nil {
		return nil, err
	}

	for _, row := range rows {
		result[row.UserID] = append(result[row.UserID], row.RoleID)
	}

	return result, nil
}

func normalizeCreateUserRequest(req createUserRequest) (string, string, string, model.UserStatus, []uint, error) {
	username := strings.TrimSpace(req.Username)
	if username == "" {
		return "", "", "", 0, nil, apperror.BadRequest("用户名不能为空")
	}
	if len(username) > 64 {
		return "", "", "", 0, nil, apperror.BadRequest("用户名不能超过 64 个字符")
	}

	if len(req.Password) < 8 || len(req.Password) > 72 {
		return "", "", "", 0, nil, apperror.BadRequest("密码长度需要在 8 到 72 个字符之间")
	}

	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		nickname = username
	}
	if len(nickname) > 64 {
		return "", "", "", 0, nil, apperror.BadRequest("昵称不能超过 64 个字符")
	}

	status := req.Status
	if status == 0 {
		status = model.UserStatusEnabled
	}
	if !validUserStatus(status) {
		return "", "", "", 0, nil, apperror.BadRequest("用户状态不正确")
	}

	roleIDs, err := normalizeRoleIDs(req.RoleIDs)
	if err != nil {
		return "", "", "", 0, nil, err
	}

	return username, req.Password, nickname, status, roleIDs, nil
}

func normalizeUpdateUserRequest(req updateUserRequest) (string, model.UserStatus, error) {
	nickname := strings.TrimSpace(req.Nickname)
	if nickname == "" {
		return "", 0, apperror.BadRequest("昵称不能为空")
	}
	if len(nickname) > 64 {
		return "", 0, apperror.BadRequest("昵称不能超过 64 个字符")
	}

	if !validUserStatus(req.Status) {
		return "", 0, apperror.BadRequest("用户状态不正确")
	}

	return nickname, req.Status, nil
}

func normalizeRoleIDs(roleIDs []uint) ([]uint, error) {
	unique := make([]uint, 0, len(roleIDs))
	seen := make(map[uint]struct{}, len(roleIDs))

	for _, roleID := range roleIDs {
		if roleID == 0 {
			return nil, apperror.BadRequest("角色 ID 不正确")
		}
		if _, ok := seen[roleID]; ok {
			continue
		}

		seen[roleID] = struct{}{}
		unique = append(unique, roleID)
	}

	return unique, nil
}

func normalizePage(page int, pageSize int) (int, int) {
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

func validUserStatus(status model.UserStatus) bool {
	return status == model.UserStatusEnabled || status == model.UserStatusDisabled
}

func userIDParam(c *gin.Context, log *zap.Logger) (uint, bool) {
	rawID := c.Param("id")
	id, err := strconv.ParseUint(rawID, 10, 64)
	if err != nil || id == 0 {
		response.Error(c, apperror.BadRequest("用户 ID 不正确"), log)
		return 0, false
	}

	return uint(id), true
}

func ensureUsernameAvailable(db *gorm.DB, username string) error {
	var user model.User
	err := db.Unscoped().Where("username = ?", username).First(&user).Error
	if err == nil {
		return apperror.BadRequest("用户名已存在")
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil
	}

	return err
}

func ensureRolesUsable(db *gorm.DB, roleIDs []uint) error {
	if len(roleIDs) == 0 {
		return nil
	}

	var count int64
	err := db.Model(&model.Role{}).
		Where("id IN ?", roleIDs).
		Where("status = ?", model.RoleStatusEnabled).
		Count(&count).Error
	if err != nil {
		return err
	}

	if count != int64(len(roleIDs)) {
		return apperror.BadRequest("角色不存在或已禁用")
	}

	return nil
}

func replaceUserRoles(db *gorm.DB, userID uint, roleIDs []uint) error {
	if err := db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	if len(roleIDs) == 0 {
		return nil
	}

	rows := make([]model.UserRole, 0, len(roleIDs))
	for _, roleID := range roleIDs {
		rows = append(rows, model.UserRole{
			UserID: userID,
			RoleID: roleID,
		})
	}

	return db.Create(&rows).Error
}

func buildUserResponse(user model.User, roleIDs []uint) userResponse {
	return userResponse{
		ID:        user.ID,
		Username:  user.Username,
		Nickname:  user.Nickname,
		Status:    user.Status,
		RoleIDs:   roleIDs,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func writeError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

::: details 为什么编辑用户时不允许修改 `username`
`username` 是登录身份标识，已经参与唯一索引和历史审计。后台管理里可以先允许修改昵称、状态和角色；如果后续确实需要改用户名，建议单独做接口，并记录操作日志。
:::

::: details 为什么不允许禁用当前登录用户
这是一条基础保护：避免管理员误把自己禁用，导致当前会话后续无法继续管理系统。真实项目里还可以进一步禁止移除自己的管理员角色。
:::

## 🛠️ 注册用户管理路由

修改 `server/internal/router/router.go`。这一处只需要在系统路由里增加用户 Handler 和路由。

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log) // [!code ++]

	// /health 通常给部署探针和本地快速验证使用。
	r.GET("/health", health.Check)

	// /api/v1/system/health 放在接口版本分组下，方便统一管理后台接口。
	api := r.Group("/api/v1")
	system := api.Group("/system")
	system.Use(middleware.Auth(opts.Token, opts.Log))
	system.Use(middleware.Permission(opts.DB, opts.Permission, opts.Log))
	system.GET("/health", health.Check)
	system.GET("/users", users.List) // [!code ++]
	system.POST("/users", users.Create) // [!code ++]
	system.POST("/users/:id/update", users.Update) // [!code ++]
	system.POST("/users/:id/status", users.UpdateStatus) // [!code ++]
	system.POST("/users/:id/roles", users.UpdateRoles) // [!code ++]
}
```

## 🛠️ 初始化用户管理权限和菜单

用户管理的权限和菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建用户管理相关的权限策略和菜单数据。

::: tip 💡 权限和菜单初始化
- 权限策略：在 `migrations/{postgres,mysql}/000002_seed_data.up.sql` 中插入用户管理接口的 Casbin 规则
- 菜单数据：在同一迁移文件中插入用户管理菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到用户管理菜单
:::

## ✅ 整理依赖并启动

本节没有新增第三方依赖，但修改了后端文件，仍然可以整理一次：

```bash
# 在 server/ 目录执行
go mod tidy
```

确认数据库和 Redis 正在运行：

```bash
# 在项目根目录执行，确认本地依赖服务处于运行状态
docker compose -f deploy/compose.local.yml ps
```

回到 `server/` 目录启动服务：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	database migrations applied
INFO	server started	{"addr": ":8080", "env": "dev"}
```

### 创建管理员账号

服务启动后，先通过初始化接口创建管理员账号：

```bash
# 创建管理员账号
curl -X POST http://localhost:8080/api/v1/setup/init \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"YourPassword123","nickname":"管理员"}'
```

## ✅ 验证权限和菜单数据

先确认用户管理接口权限已经写入：

```bash
# 查看用户管理相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/users%' order by v1, v2;"
```

应该能看到 `GET` 和 `POST` 用户管理权限。

再确认用户管理菜单已经写入：

```bash
# 查看用户管理相关菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:user%' order by sort, id;"
```

应该能看到 `system:user` 以及几个 `system:user:*` 按钮编码。

## ✅ 验证用户管理接口

先登录拿到 Token：

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "admin"
  password = "YourPassword123"
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
  -d '{"username":"admin","password":"YourPassword123"}' | jq -r '.data.access_token')
```

:::

查看用户列表：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/users?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/users?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到包含 `admin` 的分页结果。

创建一个测试用户。下面示例假设 `super_admin` 的角色 ID 是 `1`，如果不确定，可以先查询 `sys_role`：

::: warning ⚠️ Windows PowerShell 发送中文 JSON 时要显式使用 UTF-8
如果直接把 `$body` 传给 `Invoke-RestMethod -Body`，Windows PowerShell 5.1 可能把中文发送成 `????`。下面示例会先把 JSON 转成 UTF-8 字节再发送。
:::

::: code-group

```powershell [Windows PowerShell]
$body = @{
  username = "demo"
  password = "Demo@123456"
  nickname = "演示用户"
  status = 1
  role_ids = @(1)
} | ConvertTo-Json

$utf8Body = [System.Text.Encoding]::UTF8.GetBytes($body)

Invoke-RestMethod `
  -Method Post `
  -Uri http://localhost:8080/api/v1/system/users `
  -ContentType "application/json; charset=utf-8" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $utf8Body
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/system/users \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"username":"demo","password":"Demo@123456","nickname":"演示用户","status":1,"role_ids":[1]}'
```

:::

创建成功后，再查询数据库确认用户和角色关系：

```bash
# 查看 demo 用户及其绑定角色
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select u.id, u.username, u.nickname, u.status, r.code from sys_user u left join sys_user_role ur on ur.user_id = u.id left join sys_role r on r.id = ur.role_id where u.username = 'demo';"
```

修改用户状态时，把上一步返回的用户 ID 替换到路径里：

::: code-group

```powershell [Windows PowerShell]
$userId = 2
$body = @{ status = 2 } | ConvertTo-Json

Invoke-RestMethod `
  -Method Post `
  -Uri "http://localhost:8080/api/v1/system/users/$userId/status" `
  -ContentType "application/json" `
  -Headers @{ Authorization = "Bearer $token" } `
  -Body $body
```

```bash [macOS / Linux]
USER_ID=2

curl -X POST "http://localhost:8080/api/v1/system/users/${USER_ID}/status" \
  -H "Authorization: Bearer ${TOKEN}" \
  -H "Content-Type: application/json" \
  -d '{"status":2}'
```

:::

`status = 2` 表示禁用。禁用后，这个用户不能再登录。

::: warning ⚠️ 不要拿当前登录用户做禁用验证
本节代码会阻止禁用当前登录用户。验证禁用逻辑时，使用新创建的测试用户，例如 `demo`。
:::

## 常见问题

::: details 创建用户时提示“用户名已存在”
换一个用户名即可。账号唯一规则见：[数据库建表语句 - `sys_user`](../../reference/database-ddl#sys-user)。
:::

::: details 创建用户时提示“角色不存在或已禁用”
请求里的 `role_ids` 必须对应已经存在且启用的角色。可以先执行下面的 SQL 查看角色：

```sql
select id, code, name, status from sys_role order by id;
```
:::

::: details 请求用户管理接口返回 `403`
优先检查两件事：

- `casbin_rule` 中是否已经有 `/api/v1/system/users` 相关策略。
- 新增策略后是否已经重启服务，让 Enforcer 重新加载策略。
:::

::: details 为什么角色分配用“整体替换”
用户角色通常来自多选框提交。后端收到完整的 `role_ids` 后，先删除旧关系，再写入新关系，逻辑更简单，也更容易验证最终结果。
:::

下一节会继续补齐角色自身的管理能力：[角色管理](./role-management)。
