// Package application 实现接口权限元数据查询逻辑。
package application

import apiresourcedomain "ez-admin-gin/server/internal/modules/iam/apiresource/domain"

// Service 提供接口权限元数据的业务操作服务。
type Service struct {
	repo APIRepository
}

func NewService(repo APIRepository) *Service {
	return &Service{repo: repo}
}

// List 返回全部接口权限元数据。
func (s *Service) List() ([]apiresourcedomain.Response, error) {
	items, err := s.repo.List()
	if err != nil {
		return nil, err
	}

	result := make([]apiresourcedomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, apiresourcedomain.BuildResponse(item))
	}
	return result, nil
}
