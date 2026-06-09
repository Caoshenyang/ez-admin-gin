// Package application 实现操作日志的业务逻辑：分页列表查询。
package application

import (
	operationlogdomain "ez-admin-gin/server/internal/modules/system/operationlog/domain"
	"ez-admin-gin/server/internal/pkg/paging"
)

// Service 封装操作日志的业务逻辑。
type Service struct {
	repo OperationLogRepository
}

func NewService(repo OperationLogRepository) *Service {
	return &Service{repo: repo}
}

// List 按用户名、方法、路径和成功状态分页查询操作日志列表。
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
