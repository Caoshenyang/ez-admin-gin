package operationlog

// Service 负责操作日志的查询规则编排。
type Service struct {
	repo *Repository
}

// NewService 创建操作日志服务。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List 返回操作日志分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)
	success, err := NormalizeSuccessFilter(query.Success)
	if err != nil {
		return ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, success)
	if err != nil {
		return ListResponse{}, err
	}

	result := make([]Response, 0, len(items))
	for _, item := range items {
		result = append(result, BuildResponse(item))
	}

	return ListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}
