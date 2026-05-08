package application

import (
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/datascope"
)

// MeService 将已加载的 Actor 上下文转换为 /auth/me 响应。
type MeService struct{}

func NewMeService() *MeService {
	return &MeService{}
}

// Build 将已加载的 Actor 上下文转换为 /auth/me 响应。
func (s *MeService) Build(actor datascope.Actor) authdomain.MeResponse {
	return authdomain.BuildMeResponse(actor)
}
