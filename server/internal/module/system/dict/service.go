package dict

import (
	"ez-admin-gin/server/internal/apperror"
	"ez-admin-gin/server/internal/model"

	"gorm.io/gorm"
)

// Service 负责字典类型和字典项的业务规则与事务边界。
type Service struct {
	db   *gorm.DB
	repo *Repository
}

// NewService 创建字典服务。
func NewService(db *gorm.DB, repo *Repository) *Service {
	return &Service{db: db, repo: repo}
}

// ListTypes 返回字典类型分页结果。
func (s *Service) ListTypes(query TypeListQuery) (TypeListResponse, error) {
	page, pageSize := NormalizePage(query.Page, query.PageSize)
	status, err := NormalizeTypeStatusFilter(query.Status)
	if err != nil {
		return TypeListResponse{}, err
	}

	items, total, err := s.repo.ListTypes(query, page, pageSize, status)
	if err != nil {
		return TypeListResponse{}, err
	}

	result := make([]TypeResponse, 0, len(items))
	for _, item := range items {
		result = append(result, BuildTypeResponse(item))
	}

	return TypeListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateType 创建字典类型。
func (s *Service) CreateType(req CreateTypeRequest) (TypeResponse, error) {
	req, err := NormalizeCreateTypeRequest(req)
	if err != nil {
		return TypeResponse{}, err
	}

	created := DictTypeEntity{
		Code:   req.Code,
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.TypeCodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("字典编码已存在")
		}
		return s.repo.CreateType(tx, &created)
	})
	if err != nil {
		return TypeResponse{}, err
	}

	return BuildTypeResponse(created), nil
}

// UpdateType 编辑字典类型。
func (s *Service) UpdateType(typeID uint, req UpdateTypeRequest) (TypeResponse, error) {
	req, err := NormalizeUpdateTypeRequest(req)
	if err != nil {
		return TypeResponse{}, err
	}

	var updated DictTypeEntity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindTypeByID(tx, typeID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateTypeBase(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	})
	if err != nil {
		return TypeResponse{}, err
	}

	return BuildTypeResponse(updated), nil
}

// UpdateTypeStatus 单独修改字典类型状态。
func (s *Service) UpdateTypeStatus(typeID uint, status model.SystemDictStatus) error {
	status, err := normalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindTypeByID(tx, typeID)
		if err != nil {
			return err
		}
		return s.repo.UpdateTypeStatus(tx, &item, status)
	})
}

// ListItems 返回字典项分页结果。
func (s *Service) ListItems(query ItemListQuery) (ItemListResponse, error) {
	if query.TypeID == 0 {
		return ItemListResponse{}, apperror.BadRequest("字典类型 ID 不正确")
	}

	page, pageSize := NormalizePage(query.Page, query.PageSize)
	status, err := NormalizeItemStatusFilter(query.Status)
	if err != nil {
		return ItemListResponse{}, err
	}

	if _, err := s.repo.FindTypeByID(s.db, query.TypeID); err != nil {
		return ItemListResponse{}, err
	}

	items, total, err := s.repo.ListItems(query, page, pageSize, status)
	if err != nil {
		return ItemListResponse{}, err
	}

	result := make([]ItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, BuildItemResponse(item))
	}

	return ItemListResponse{
		Items:    result,
		Total:    total,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// CreateItem 创建字典项。
func (s *Service) CreateItem(req CreateItemRequest) (ItemResponse, error) {
	req, err := NormalizeCreateItemRequest(req)
	if err != nil {
		return ItemResponse{}, err
	}

	created := DictItemEntity{
		TypeID:  req.TypeID,
		ItemKey: req.ItemKey,
		Label:   req.Label,
		Value:   req.Value,
		TagType: req.TagType,
		Sort:    req.Sort,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	err = s.db.Transaction(func(tx *gorm.DB) error {
		if _, err := s.repo.FindTypeByID(tx, req.TypeID); err != nil {
			return err
		}

		exists, err := s.repo.ItemKeyExists(tx, req.TypeID, req.ItemKey)
		if err != nil {
			return err
		}
		if exists {
			return apperror.BadRequest("同一字典类型下的字典项编码不能重复")
		}

		return s.repo.CreateItem(tx, &created)
	})
	if err != nil {
		return ItemResponse{}, err
	}

	return BuildItemResponse(created), nil
}

// UpdateItem 编辑字典项。
func (s *Service) UpdateItem(itemID uint, req UpdateItemRequest) (ItemResponse, error) {
	req, err := NormalizeUpdateItemRequest(req)
	if err != nil {
		return ItemResponse{}, err
	}

	var updated DictItemEntity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			return err
		}
		if err := s.repo.UpdateItemBase(tx, &item, req); err != nil {
			return err
		}
		updated = item
		return nil
	})
	if err != nil {
		return ItemResponse{}, err
	}

	return BuildItemResponse(updated), nil
}

// UpdateItemStatus 单独修改字典项状态。
func (s *Service) UpdateItemStatus(itemID uint, status model.SystemDictStatus) error {
	status, err := normalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			return err
		}
		return s.repo.UpdateItemStatus(tx, &item, status)
	})
}
