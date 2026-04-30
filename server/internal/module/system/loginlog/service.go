package loginlog

// Service 负责登录日志的查询规则编排。
type Service struct {
	repo *Repository
}

// NewService 创建登录日志服务。
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// List 返回登录日志分页结果。
func (s *Service) List(query ListQuery) (ListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)
	status, err := NormalizeStatusFilter(query.Status)
	if err != nil {
		return ListResponse{}, err
	}

	items, total, err := s.repo.List(query, page, pageSize, status)
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
