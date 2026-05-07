package application

import (
	authdomain "ez-admin-gin/server/internal/modules/auth/domain"
	"ez-admin-gin/server/internal/platform/datascope"
)

type MeService struct{}

func NewMeService() *MeService {
	return &MeService{}
}

func (s *MeService) Build(actor datascope.Actor) authdomain.MeResponse {
	return authdomain.BuildMeResponse(actor)
}
