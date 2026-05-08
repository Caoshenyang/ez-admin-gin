// Package application 实现登录日志的业务逻辑：分页列表查询。
package application

import (
	loginlogdomain "ez-admin-gin/server/internal/modules/system/loginlog/domain"
	"ez-admin-gin/server/internal/pkg/paging"
)

// Service 封装登录日志的业务逻辑。
type Service struct {
	repo LoginLogRepository
}

func NewService(repo LoginLogRepository) *Service {
	return &Service{repo: repo}
}

// List 按用户名、IP 和状态分页查询登录日志列表。
func (s *Service) List(query loginlogdomain.ListQuery) (loginlogdomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := loginlogdomain.NormalizeStatusFilter(query.Status)
	if err != nil {
		return loginlogdomain.ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, status)
	if err != nil {
		return loginlogdomain.ListResponse{}, err
	}

	result := make([]loginlogdomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, loginlogdomain.BuildResponse(item))
	}

	return loginlogdomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}
