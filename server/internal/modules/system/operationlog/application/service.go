package application

import (
	operationlogdomain "ez-admin-gin/server/internal/modules/system/operationlog/domain"
	operationloginfra "ez-admin-gin/server/internal/modules/system/operationlog/infra"
	"ez-admin-gin/server/internal/pkg/paging"
)

type Service struct {
	repo *operationloginfra.Repository
}

func NewService(repo *operationloginfra.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) List(query operationlogdomain.ListQuery) (operationlogdomain.ListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	success, err := operationlogdomain.NormalizeSuccessFilter(query.Success)
	if err != nil {
		return operationlogdomain.ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, success)
	if err != nil {
		return operationlogdomain.ListResponse{}, err
	}

	result := make([]operationlogdomain.Response, 0, len(items))
	for _, item := range items {
		result = append(result, operationlogdomain.BuildResponse(item))
	}

	return operationlogdomain.ListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}
