package application

import (
	loginlogdomain "ez-admin-gin/server/internal/modules/system/loginlog/domain"
	loginloginfra "ez-admin-gin/server/internal/modules/system/loginlog/infra"
	"ez-admin-gin/server/internal/pkg/paging"
)

type Service struct {
	repo *loginloginfra.Repository
}

func NewService(repo *loginloginfra.Repository) *Service {
	return &Service{repo: repo}
}

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
