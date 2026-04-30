package auth

import "ez-admin-gin/server/internal/platform/datascope"

// MeService 负责当前登录用户摘要组装。
type MeService struct{}

// NewMeService 创建当前用户服务。
func NewMeService() *MeService {
	return &MeService{}
}

// Build 把 Actor 组装成 /auth/me 响应。
func (s *MeService) Build(actor datascope.Actor) MeResponse {
	return BuildMeResponse(actor)
}
