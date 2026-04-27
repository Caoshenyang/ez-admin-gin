---
title: 文件上传
description: "实现后台本地文件上传能力，记录文件元数据，并提供可访问的文件 URL。"
---

# 文件上传

系统配置完成后，继续补齐后台底座里很常见的一项能力：文件上传。本节先实现本地上传，把文件保存到 `server/uploads/`，同时把文件名、大小、后缀、访问地址等元数据写入数据库。

::: tip 🎯 本节目标
完成后，可以通过接口上传文件；后端会校验文件大小和后缀，把文件保存到本地目录，并在 `sys_file` 中记录文件信息。
:::

## 本节会改什么

本节会新增或修改下面这些文件：

```text
docs/
└─ reference/
   └─ database-ddl.md

server/
├─ configs/
│  └─ config.yaml
├─ internal/
│  ├─ config/
│  │  └─ config.go
│  ├─ handler/
│  │  └─ system/
│  │     └─ files.go
│  ├─ model/
│  │  └─ system_file.go
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
| `docs/reference/database-ddl.md` | 补充 `sys_file` 建表语句 |
| `configs/config.yaml` | 增加上传目录、访问前缀、大小限制、后缀白名单 |
| `internal/config/config.go` | 增加上传配置结构 |
| `internal/model/system_file.go` | 定义文件记录模型 |
| `internal/handler/system/files.go` | 实现文件上传和文件列表接口 |
| `internal/router/router.go` | 注册静态文件访问和文件接口 |
| `migrations/{postgres,mysql}/000002_seed_data.up.sql` | 初始化文件上传权限和菜单 |

## 先创建数据表

本节新增 `sys_file`，用于保存文件上传后的元数据，文件内容仍然保存在本地上传目录中。

::: tip 建表 SQL
字段说明、文件元数据设计、索引设计和 PostgreSQL / MySQL 建表语句统一放在参考手册：[数据库建表语句 - `sys_file`](../../reference/database-ddl#sys-file)。
:::

## 文件如何保存

本节采用下面的保存规则：

```text
server/uploads/
└─ 20260423/
   └─ 20260423153000_a1b2c3d4e5f6a7b8.txt
```

| 数据 | 示例 | 说明 |
| --- | --- | --- |
| 原始文件名 | `测试文件.txt` | 用户上传时的文件名 |
| 保存文件名 | `20260423153000_a1b2c3d4e5f6a7b8.txt` | 后端生成，避免重名 |
| 磁盘路径 | `uploads/20260423/...txt` | 写入 `sys_file.path` |
| 访问地址 | `/uploads/20260423/...txt` | 写入 `sys_file.url` |

::: warning ⚠️ 不要直接使用用户上传的文件名作为保存文件名
用户上传的文件名可能重复，也可能包含不适合直接拼接路径的字符。后端应该保留原始文件名用于展示，同时生成新的安全文件名用于保存。
:::

## 接口规划

本节实现 2 个接口：

| 方法 | 路径 | 用途 |
| --- | --- | --- |
| `GET` | `/api/v1/system/files` | 文件分页列表 |
| `POST` | `/api/v1/system/files` | 上传文件 |

上传成功后返回文件记录，其中 `url` 可以直接访问：

```json
{
  "id": 1,
  "original_name": "test-upload.txt",
  "file_name": "20260423153000_a1b2c3d4e5f6a7b8.txt",
  "url": "/uploads/20260423/20260423153000_a1b2c3d4e5f6a7b8.txt"
}
```

## 🛠️ 创建文件记录模型

创建 `server/internal/model/system_file.go`。这是新增文件，直接完整写入即可。

```go
package model

import (
	"time"

	"gorm.io/gorm"
)

// SystemFileStatus 表示文件记录状态。
type SystemFileStatus int

const (
	// SystemFileStatusEnabled 表示文件可正常使用。
	SystemFileStatusEnabled SystemFileStatus = 1
	// SystemFileStatusDisabled 表示文件已停用。
	SystemFileStatusDisabled SystemFileStatus = 2
)

// SystemFile 是文件上传记录模型。
type SystemFile struct {
	ID           uint             `gorm:"primaryKey" json:"id"`
	Storage      string           `gorm:"size:32;not null;default:'local'" json:"storage"`
	OriginalName string           `gorm:"size:255;not null" json:"original_name"`
	FileName     string           `gorm:"size:255;not null" json:"file_name"`
	Ext          string           `gorm:"size:32;not null;default:'';index" json:"ext"`
	MimeType     string           `gorm:"size:128;not null;default:''" json:"mime_type"`
	Size         int64            `gorm:"not null;default:0" json:"size"`
	Sha256       string           `gorm:"size:64;not null;default:'';index" json:"sha256"`
	Path         string           `gorm:"size:500;not null" json:"path"`
	URL          string           `gorm:"column:url;size:500;not null" json:"url"`
	UploaderID   uint             `gorm:"not null;default:0;index" json:"uploader_id"`
	Status       SystemFileStatus `gorm:"type:smallint;not null;default:1" json:"status"`
	Remark       string           `gorm:"size:255;not null;default:''" json:"remark"`
	CreatedAt    time.Time        `json:"created_at"`
	UpdatedAt    time.Time        `json:"updated_at"`
	DeletedAt    gorm.DeletedAt   `gorm:"index" json:"-"`
}

// TableName 固定文件记录表名。
func (SystemFile) TableName() string {
	return "sys_file"
}
```

## 🛠️ 增加上传配置

修改 `server/internal/config/config.go`。本次要改 4 个位置：

- `Config` 增加 `Upload UploadConfig`
- 新增 `UploadConfig` 结构体
- `setDefaults` 增加上传默认值
- `bindEnvs` 增加上传相关环境变量

先在 `Config` 中增加 `Upload`：

```go
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Server   ServerConfig   `mapstructure:"server"`
	Database DatabaseConfig `mapstructure:"database"`
	Redis    RedisConfig    `mapstructure:"redis"`
	Auth     AuthConfig     `mapstructure:"auth"`
	Upload   UploadConfig   `mapstructure:"upload"` // [!code ++]
	Log      LogConfig      `mapstructure:"log"`
}
```

在 `AuthConfig` 后面增加上传配置结构：

```go
// UploadConfig 保存本地文件上传配置。
type UploadConfig struct { // [!code ++]
	// Dir 是文件保存目录，相对于 server/ 目录。 // [!code ++]
	Dir string `mapstructure:"dir"` // [!code ++]
	// PublicPath 是文件公开访问前缀。 // [!code ++]
	PublicPath string `mapstructure:"public_path"` // [!code ++]
	// MaxSizeMB 是单个文件最大大小，单位 MB。 // [!code ++]
	MaxSizeMB int64 `mapstructure:"max_size_mb"` // [!code ++]
	// AllowedExts 是允许上传的文件后缀白名单。 // [!code ++]
	AllowedExts []string `mapstructure:"allowed_exts"` // [!code ++]
} // [!code ++]
```

在 `setDefaults` 中增加默认值：

```go
func setDefaults(v *viper.Viper) {
	v.SetDefault("auth.jwt_secret", "ez-admin-dev-secret-change-me-please-32")
	v.SetDefault("auth.access_token_ttl", 7200)
	v.SetDefault("auth.issuer", "ez-admin")
	v.SetDefault("upload.dir", "uploads") // [!code ++]
	v.SetDefault("upload.public_path", "/uploads") // [!code ++]
	v.SetDefault("upload.max_size_mb", 10) // [!code ++]
	v.SetDefault("upload.allowed_exts", []string{".jpg", ".jpeg", ".png", ".gif", ".webp", ".pdf", ".txt", ".docx", ".xlsx"}) // [!code ++]
}
```

在 `bindEnvs` 的 `keys` 中追加上传配置：

```go
keys := []string{
	"auth.jwt_secret",
	"auth.access_token_ttl",
	"auth.issuer",
	"upload.dir", // [!code ++]
	"upload.public_path", // [!code ++]
	"upload.max_size_mb", // [!code ++]
	"upload.allowed_exts", // [!code ++]
}
```

继续修改 `server/configs/config.yaml`，在 `auth` 后面增加：

```yaml
upload:
  # 文件保存目录，相对于 server/ 目录。
  dir: uploads
  # 文件公开访问前缀，后续会注册为静态资源路径。
  public_path: /uploads
  # 单个文件最大大小，单位 MB。
  max_size_mb: 10
  # 允许上传的文件后缀。先用白名单，避免任意文件都能上传。
  allowed_exts:
    - .jpg
    - .jpeg
    - .png
    - .gif
    - .webp
    - .pdf
    - .txt
    - .docx
    - .xlsx
```

::: warning ⚠️ 上传目录不要提交到 Git
上传目录里保存的是运行时文件，不是项目源码。建议在 `.gitignore` 中追加：

```text
server/uploads/
```
:::

## 🛠️ 创建文件上传 Handler

创建 `server/internal/handler/system/files.go`。这个文件比较长，下面分两段展示；两段代码要放在同一个文件里，第二段紧接在第一段后面。

```go
package system

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/config"
	"ez-admin-gin/server/internal/middleware"
	"ez-admin-gin/server/internal/model"
	"ez-admin-gin/server/internal/response"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

const (
	defaultUploadDir        = "uploads"
	defaultUploadPublicPath = "/uploads"
	defaultUploadMaxSizeMB  = 10
	localFileStorage        = "local"
)

// FileHandler 负责后台文件上传接口。
type FileHandler struct {
	db  *gorm.DB
	cfg config.UploadConfig
	log *zap.Logger
}

// NewFileHandler 创建文件上传 Handler。
func NewFileHandler(db *gorm.DB, cfg config.UploadConfig, log *zap.Logger) *FileHandler {
	return &FileHandler{
		db:  db,
		cfg: normalizeUploadConfig(cfg),
		log: log,
	}
}

type fileListQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Keyword  string `form:"keyword"`
	Ext      string `form:"ext"`
	Status   int    `form:"status"`
}

type fileResponse struct {
	ID           uint                   `json:"id"`
	Storage      string                 `json:"storage"`
	OriginalName string                 `json:"original_name"`
	FileName     string                 `json:"file_name"`
	Ext          string                 `json:"ext"`
	MimeType     string                 `json:"mime_type"`
	Size         int64                  `json:"size"`
	Sha256       string                 `json:"sha256"`
	Path         string                 `json:"path"`
	URL          string                 `json:"url"`
	UploaderID   uint                   `json:"uploader_id"`
	Status       model.SystemFileStatus `json:"status"`
	Remark       string                 `json:"remark"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

type fileListResponse struct {
	Items    []fileResponse `json:"items"`
	Total    int64          `json:"total"`
	Page     int            `json:"page"`
	PageSize int            `json:"page_size"`
}

type savedUploadedFile struct {
	OriginalName string
	FileName     string
	Ext          string
	MimeType     string
	Size         int64
	Sha256       string
	Path         string
	URL          string
	AbsolutePath string
}

// List 返回文件分页列表。
func (h *FileHandler) List(c *gin.Context) {
	var query fileListQuery
	if err := c.ShouldBindQuery(&query); err != nil {
		response.Error(c, apperror.BadRequest("查询参数不正确"), h.log)
		return
	}

	page, pageSize := normalizeFilePage(query.Page, query.PageSize)
	queryDB := h.db.Model(&model.SystemFile{})

	keyword := strings.TrimSpace(query.Keyword)
	if keyword != "" {
		like := "%" + keyword + "%"
		queryDB = queryDB.Where("original_name LIKE ? OR file_name LIKE ?", like, like)
	}

	ext := normalizeFileExt(query.Ext)
	if ext != "" {
		queryDB = queryDB.Where("ext = ?", ext)
	}

	if query.Status != 0 {
		status := model.SystemFileStatus(query.Status)
		if !validFileStatus(status) {
			response.Error(c, apperror.BadRequest("文件状态不正确"), h.log)
			return
		}
		queryDB = queryDB.Where("status = ?", status)
	}

	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		response.Error(c, apperror.Internal("查询文件总数失败", err), h.log)
		return
	}

	var files []model.SystemFile
	if err := queryDB.
		Order("id DESC").
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&files).Error; err != nil {
		response.Error(c, apperror.Internal("查询文件列表失败", err), h.log)
		return
	}

	items := make([]fileResponse, 0, len(files))
	for _, file := range files {
		items = append(items, buildFileResponse(file))
	}

	response.Success(c, fileListResponse{
		Items:    items,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	})
}

// Upload 上传文件到本地目录，并写入文件记录。
func (h *FileHandler) Upload(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		response.Error(c, apperror.BadRequest("请选择要上传的文件"), h.log)
		return
	}

	if err := h.validateUploadFile(fileHeader); err != nil {
		response.Error(c, err, h.log)
		return
	}

	saved, err := h.saveUploadedFile(fileHeader)
	if err != nil {
		response.Error(c, apperror.Internal("保存文件失败", err), h.log)
		return
	}

	uploaderID, _ := middleware.CurrentUserID(c)
	file := model.SystemFile{
		Storage:      localFileStorage,
		OriginalName: saved.OriginalName,
		FileName:     saved.FileName,
		Ext:          saved.Ext,
		MimeType:     saved.MimeType,
		Size:         saved.Size,
		Sha256:       saved.Sha256,
		Path:         saved.Path,
		URL:          saved.URL,
		UploaderID:   uploaderID,
		Status:       model.SystemFileStatusEnabled,
		Remark:       "",
	}

	if err := h.db.Create(&file).Error; err != nil {
		_ = os.Remove(saved.AbsolutePath)
		response.Error(c, apperror.Internal("保存文件记录失败", err), h.log)
		return
	}

response.Success(c, buildFileResponse(file))
}
```

继续在同一个 `files.go` 中追加下面的辅助函数：

```go
func (h *FileHandler) validateUploadFile(fileHeader *multipart.FileHeader) error {
	if fileHeader.Size <= 0 {
		return apperror.BadRequest("不能上传空文件")
	}

	maxBytes := uploadMaxBytes(h.cfg.MaxSizeMB)
	if fileHeader.Size > maxBytes {
		return apperror.BadRequest("文件大小不能超过 " + strconv.FormatInt(h.cfg.MaxSizeMB, 10) + " MB")
	}

	ext := normalizeFileExt(filepath.Ext(fileHeader.Filename))
	if ext == "" {
		return apperror.BadRequest("文件后缀不能为空")
	}

	if !isAllowedFileExt(ext, h.cfg.AllowedExts) {
		return apperror.BadRequest("不支持上传该文件类型")
	}

	return nil
}

func (h *FileHandler) saveUploadedFile(fileHeader *multipart.FileHeader) (savedUploadedFile, error) {
	src, err := fileHeader.Open()
	if err != nil {
		return savedUploadedFile{}, err
	}
	defer src.Close()

	now := time.Now()
	dateDir := now.Format("20060102")
	ext := normalizeFileExt(filepath.Ext(fileHeader.Filename))
	randomPart, err := randomHex(8)
	if err != nil {
		return savedUploadedFile{}, err
	}

	// 保存文件名由后端生成，避免重名和路径拼接风险。
	fileName := fmt.Sprintf("%s_%s%s", now.Format("20060102150405"), randomPart, ext)
	targetDir := filepath.Join(h.cfg.Dir, dateDir)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return savedUploadedFile{}, err
	}

	absolutePath := filepath.Join(targetDir, fileName)
	dst, err := os.OpenFile(absolutePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return savedUploadedFile{}, err
	}
	defer dst.Close()

	hasher := sha256.New()
	written, err := io.Copy(dst, io.TeeReader(src, hasher))
	if err != nil {
		_ = os.Remove(absolutePath)
		return savedUploadedFile{}, err
	}

	publicPath := normalizeUploadPublicPath(h.cfg.PublicPath)
	relativePath := filepath.ToSlash(filepath.Join(h.cfg.Dir, dateDir, fileName))
	url := publicPath + "/" + dateDir + "/" + fileName
	mimeType := strings.TrimSpace(fileHeader.Header.Get("Content-Type"))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return savedUploadedFile{
		OriginalName: filepath.Base(fileHeader.Filename),
		FileName:     fileName,
		Ext:          ext,
		MimeType:     mimeType,
		Size:         written,
		Sha256:       hex.EncodeToString(hasher.Sum(nil)),
		Path:         relativePath,
		URL:          url,
		AbsolutePath: absolutePath,
	}, nil
}

func normalizeUploadConfig(cfg config.UploadConfig) config.UploadConfig {
	cfg.Dir = strings.TrimSpace(cfg.Dir)
	if cfg.Dir == "" {
		cfg.Dir = defaultUploadDir
	}

	cfg.PublicPath = normalizeUploadPublicPath(cfg.PublicPath)
	if cfg.MaxSizeMB <= 0 {
		cfg.MaxSizeMB = defaultUploadMaxSizeMB
	}

	cfg.AllowedExts = normalizeAllowedExts(cfg.AllowedExts)
	return cfg
}

func normalizeUploadPublicPath(publicPath string) string {
	publicPath = strings.TrimSpace(publicPath)
	if publicPath == "" {
		return defaultUploadPublicPath
	}

	if !strings.HasPrefix(publicPath, "/") {
		publicPath = "/" + publicPath
	}

	return strings.TrimRight(publicPath, "/")
}

func normalizeAllowedExts(exts []string) []string {
	result := make([]string, 0, len(exts))
	seen := make(map[string]struct{}, len(exts))

	for _, ext := range exts {
		ext = normalizeFileExt(ext)
		if ext == "" {
			continue
		}
		if _, ok := seen[ext]; ok {
			continue
		}

		seen[ext] = struct{}{}
		result = append(result, ext)
	}

	return result
}

func normalizeFileExt(ext string) string {
	ext = strings.ToLower(strings.TrimSpace(ext))
	if ext == "" {
		return ""
	}
	if !strings.HasPrefix(ext, ".") {
		ext = "." + ext
	}

	return ext
}

func isAllowedFileExt(ext string, allowedExts []string) bool {
	ext = normalizeFileExt(ext)
	if ext == "" || len(allowedExts) == 0 {
		return false
	}

	for _, allowed := range allowedExts {
		if ext == normalizeFileExt(allowed) {
			return true
		}
	}

	return false
}

func uploadMaxBytes(maxSizeMB int64) int64 {
	if maxSizeMB <= 0 {
		maxSizeMB = defaultUploadMaxSizeMB
	}

	return maxSizeMB * 1024 * 1024
}

func randomHex(size int) (string, error) {
	bytes := make([]byte, size)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func normalizeFilePage(page int, pageSize int) (int, int) {
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

func validFileStatus(status model.SystemFileStatus) bool {
	return status == model.SystemFileStatusEnabled || status == model.SystemFileStatusDisabled
}

func buildFileResponse(file model.SystemFile) fileResponse {
	return fileResponse{
		ID:           file.ID,
		Storage:      file.Storage,
		OriginalName: file.OriginalName,
		FileName:     file.FileName,
		Ext:          file.Ext,
		MimeType:     file.MimeType,
		Size:         file.Size,
		Sha256:       file.Sha256,
		Path:         file.Path,
		URL:          file.URL,
		UploaderID:   file.UploaderID,
		Status:       file.Status,
		Remark:       file.Remark,
		CreatedAt:    file.CreatedAt,
		UpdatedAt:    file.UpdatedAt,
	}
}

func writeFileError(c *gin.Context, err error, fallbackMessage string, log *zap.Logger) {
	var appErr *apperror.Error
	if errors.As(err, &appErr) {
		response.Error(c, appErr, log)
		return
	}

	response.Error(c, apperror.Internal(fallbackMessage, err), log)
}
```

::: warning ⚠️ 后缀白名单只是第一层限制
这一节先用后缀白名单和大小限制完成基础能力。真实项目里如果允许上传可执行风险更高的文件，还应该结合 MIME 检测、病毒扫描、私有访问策略等继续加强。
:::

## 🛠️ 注册静态文件和上传路由

修改 `server/internal/router/router.go`。本次要改两处：

- 在 `New` 中注册静态文件访问路径
- 在 `registerSystemRoutes` 中注册文件接口

先修改 `New`：

```go
// New 创建路由引擎，并统一注册中间件和路由分组。
func New(opts Options) *gin.Engine {
	r := gin.New()
	r.Use(appLogger.GinLogger(opts.Log), appLogger.GinRecovery(opts.Log))

	if opts.Config.Upload.MaxSizeMB > 0 { // [!code ++]
		r.MaxMultipartMemory = opts.Config.Upload.MaxSizeMB << 20 // [!code ++]
	} // [!code ++]
	r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir) // [!code ++]

	registerSystemRoutes(r, opts)
	registerAuthRoutes(r, opts)

	return r
}
```

再修改 `registerSystemRoutes`：

```go
// registerSystemRoutes 注册系统级路由。
func registerSystemRoutes(r *gin.Engine, opts Options) {
	health := systemHandler.NewHealthHandler(opts.Config, opts.DB, opts.Redis, opts.Log)
	users := systemHandler.NewUserHandler(opts.DB, opts.Log)
	roles := systemHandler.NewRoleHandler(opts.DB, opts.Log)
	menus := systemHandler.NewMenuAdminHandler(opts.DB, opts.Log)
	configs := systemHandler.NewSystemConfigHandler(opts.DB, opts.Redis, opts.Log)
	files := systemHandler.NewFileHandler(opts.DB, opts.Config.Upload, opts.Log) // [!code ++]

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
	system.GET("/configs", configs.List)
	system.POST("/configs", configs.Create)
	system.POST("/configs/:id/update", configs.Update)
	system.POST("/configs/:id/status", configs.UpdateStatus)
	system.GET("/configs/value/:key", configs.Value)
	system.GET("/files", files.List) // [!code ++]
	system.POST("/files", files.Upload) // [!code ++]
}
```

::: details 为什么静态文件路径不放到 `/api/v1`
`/api/v1` 是接口路径；上传后的文件 URL 更像静态资源路径。把文件访问放到 `/uploads/...`，前端展示图片、下载文件都会更直接。
:::

## 🛠️ 初始化文件上传权限和菜单

文件上传的权限和菜单已经在数据库迁移文件中初始化。迁移文件会在服务启动时自动执行，创建文件上传相关的权限策略和菜单数据。

::: tip 💡 权限和菜单初始化
- 权限策略：在 `migrations/{postgres,mysql}/000002_seed_data.up.sql` 中插入文件上传接口的 Casbin 规则
- 菜单数据：在同一迁移文件中插入文件上传菜单和按钮
- 角色菜单绑定：在同一迁移文件中绑定 `super_admin` 角色到文件上传菜单
:::

## ✅ 启动并观察初始化日志

本节没有新增第三方依赖，可以直接启动：

```bash
# 在 server/ 目录启动服务
go run .
```

第一次启动后，控制台应该能看到类似日志：

```text
INFO	default permission created	{"role_code": "super_admin", "path": "/api/v1/system/files", "method": "POST"}
INFO	default menu created	{"menu_code": "system:file"}
INFO	default role menu bound	{"role_id": 1, "menu_id": 20}
```

## ✅ 验证权限和菜单数据

先确认文件上传接口权限已经写入：

```bash
# 查看文件管理相关接口权限
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select ptype, v0, v1, v2 from casbin_rule where v1 like '/api/v1/system/files%' order by v1, v2;"
```

应该能看到 `GET` 和 `POST` 两条策略。

再确认文件管理菜单和按钮已经写入：

```bash
# 查看文件管理菜单和按钮
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, parent_id, type, code, title from sys_menu where code like 'system:file%' order by sort, id;"
```

应该能看到 `system:file`、`system:file:list`、`system:file:upload`。

## ✅ 验证文件上传接口

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

创建一个测试文件：

::: code-group

```powershell [Windows PowerShell]
Set-Content -Path .\test-upload.txt -Value "hello ez admin" -Encoding UTF8
```

```bash [macOS / Linux]
echo "hello ez admin" > test-upload.txt
```

:::

上传文件：

::: code-group

```powershell [Windows PowerShell]
curl.exe -X POST http://localhost:8080/api/v1/system/files `
  -H "Authorization: Bearer $token" `
  -F "file=@.\test-upload.txt"
```

```bash [macOS / Linux]
curl -X POST http://localhost:8080/api/v1/system/files \
  -H "Authorization: Bearer ${TOKEN}" \
  -F "file=@./test-upload.txt"
```

:::

上传成功后，响应里应该能看到：

- `original_name` 是 `test-upload.txt`
- `ext` 是 `.txt`
- `size` 大于 `0`
- `url` 以 `/uploads/` 开头
- `sha256` 是 64 位哈希字符串

::: warning ⚠️ Windows 下这里使用 `curl.exe`
PowerShell 里直接写 `curl` 可能会调用到别名，不一定是真正的 curl。上传 multipart 文件时，建议明确使用 `curl.exe`。
:::

## ✅ 验证数据库和本地文件

先查数据库记录：

```bash
# 查看最近上传的文件记录
docker compose -f deploy/compose.local.yml exec postgres psql -U ez_admin -d ez_admin -c "select id, original_name, ext, size, path, url, uploader_id from sys_file order by id desc limit 5;"
```

应该能看到刚上传的 `test-upload.txt`，并且 `path` 指向 `uploads/日期目录/文件名`。

再检查本地文件是否存在。下面路径里的日期目录和文件名，以接口返回的 `path` 为准：

```powershell
# 在 server/ 目录下查看上传目录
Get-ChildItem .\uploads -Recurse
```

应该能看到新生成的文件。

最后验证公开访问地址。把 `$fileUrl` 换成接口返回的 `url`：

::: code-group

```powershell [Windows PowerShell]
$fileUrl = "/uploads/20260423/20260423153000_a1b2c3d4e5f6a7b8.txt"
Invoke-RestMethod -Method Get -Uri "http://localhost:8080$fileUrl"
```

```bash [macOS / Linux]
FILE_URL="/uploads/20260423/20260423153000_a1b2c3d4e5f6a7b8.txt"
curl "http://localhost:8080${FILE_URL}"
```

:::

如果上传的是文本文件，应该能看到文件内容 `hello ez admin`。

## ✅ 验证文件列表

调用文件列表接口：

::: code-group

```powershell [Windows PowerShell]
Invoke-RestMethod `
  -Method Get `
  -Uri "http://localhost:8080/api/v1/system/files?page=1&page_size=10" `
  -Headers @{ Authorization = "Bearer $token" }
```

```bash [macOS / Linux]
curl "http://localhost:8080/api/v1/system/files?page=1&page_size=10" \
  -H "Authorization: Bearer ${TOKEN}"
```

:::

应该能看到包含 `test-upload.txt` 的分页结果。

## 常见问题

::: details 上传接口返回“请选择要上传的文件”
确认表单字段名是 `file`。本节代码固定读取：

```text
file
```

如果前端或 curl 里写成 `image`、`upload`、`files`，后端就取不到文件。
:::

::: details 上传接口返回“不支持上传该文件类型”
检查 `server/configs/config.yaml` 里的 `upload.allowed_exts`。本节默认允许 `.jpg`、`.jpeg`、`.png`、`.gif`、`.webp`、`.pdf`、`.txt`、`.docx`、`.xlsx`。
:::

::: details 能上传成功，但访问 `/uploads/...` 返回 404
优先检查三件事：

- `router.New` 中是否已经注册 `r.Static(opts.Config.Upload.PublicPath, opts.Config.Upload.Dir)`。
- `config.yaml` 中的 `upload.public_path` 是否是 `/uploads`。
- 文件是否真的保存在 `server/uploads/日期目录/` 下。
:::

::: details 为什么数据库保存 `path`，还保存 `url`
`path` 表示服务端磁盘位置，主要给后端删除、迁移或排查使用；`url` 表示前端可访问地址，主要给页面展示和下载使用。两者职责不同，分开保存更直观。
:::

下一节继续补齐操作审计能力：[操作日志](./operation-logs)。
