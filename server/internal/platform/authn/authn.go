package authn

import (
	"ez-admin-gin/server/internal/config"
	legacyToken "ez-admin-gin/server/internal/token"
)

// Manager 复用现有 token 管理器，实现向 v2 authn 命名空间平滑迁移。
type Manager = legacyToken.Manager

// Claims 复用现有访问令牌声明结构。
type Claims = legacyToken.Claims

// NewManager 创建访问令牌管理器。
func NewManager(cfg config.AuthConfig) (*Manager, error) {
	return legacyToken.NewManager(cfg)
}
