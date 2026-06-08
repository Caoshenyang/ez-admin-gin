// Package application 实现字典的业务逻辑：字典类型和字典项的 CRUD。
package application

import (
	"context"

	dictdomain "ez-admin-gin/server/internal/modules/system/dict/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

// Service 封装字典类型和字典项的业务逻辑。
type Service struct {
	tx   DictTransactor
	repo DictRepository
}

func NewService(tx DictTransactor, repo DictRepository) *Service {
	return &Service{tx: tx, repo: repo}
}

// ListTypes 按关键词和状态分页查询字典类型列表。
func (s *Service) ListTypes(query dictdomain.TypeListQuery) (dictdomain.TypeListResponse, error) {
	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := dictdomain.NormalizeTypeStatusFilter(query.Status)
	if err != nil {
		return dictdomain.TypeListResponse{}, err
	}

	items, total, err := s.repo.ListTypes(query, page, pageSize, status)
	if err != nil {
		return dictdomain.TypeListResponse{}, err
	}

	result := make([]dictdomain.TypeResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dictdomain.BuildTypeResponse(item))
	}

	return dictdomain.TypeListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateType 创建字典类型，校验编码唯一性后写入数据库。
func (s *Service) CreateType(req dictdomain.CreateTypeRequest) (dictdomain.TypeResponse, error) {
	req, err := dictdomain.NormalizeCreateTypeRequest(req)
	if err != nil {
		return dictdomain.TypeResponse{}, err
	}

	created := dictdomain.DictTypeEntity{
		Code:   req.Code,
		Name:   req.Name,
		Sort:   req.Sort,
		Status: req.Status,
		Remark: req.Remark,
	}

	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		exists, err := s.repo.TypeCodeExists(tx, req.Code)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("字典编码已存在")
		}
		return s.repo.CreateType(tx, &created)
	})
	if err != nil {
		return dictdomain.TypeResponse{}, err
	}

	return dictdomain.BuildTypeResponse(created), nil
}

// UpdateType 更新指定字典类型的基本信息。
func (s *Service) UpdateType(typeID uint, req dictdomain.UpdateTypeRequest) (dictdomain.TypeResponse, error) {
	req, err := dictdomain.NormalizeUpdateTypeRequest(req)
	if err != nil {
		return dictdomain.TypeResponse{}, err
	}

	var updated dictdomain.DictTypeEntity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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
		return dictdomain.TypeResponse{}, err
	}

	return dictdomain.BuildTypeResponse(updated), nil
}

// UpdateTypeStatus 切换字典类型的启用/禁用状态。
func (s *Service) UpdateTypeStatus(typeID uint, status model.SystemDictStatus) error {
	status, err := dictdomain.NormalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTypeByID(tx, typeID)
		if err != nil {
			return err
		}
		return s.repo.UpdateTypeStatus(tx, &item, status)
	})
}

// DeleteType 删除字典类型，要求已经清空字典项。
func (s *Service) DeleteType(typeID uint) error {
	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindTypeByID(tx, typeID)
		if err != nil {
			return err
		}

		itemCount, err := s.repo.CountItemsByType(tx, typeID)
		if err != nil {
			return err
		}
		if itemCount > 0 {
			return errorsx.BadRequest("请先删除该字典类型下的字典项")
		}

		return s.repo.DeleteType(tx, &item)
	})
}

// ListItems 按类型 ID、关键词和状态分页查询字典项列表。
func (s *Service) ListItems(query dictdomain.ItemListQuery) (dictdomain.ItemListResponse, error) {
	if query.TypeID == 0 {
		return dictdomain.ItemListResponse{}, errorsx.BadRequest("字典类型 ID 不正确")
	}

	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := dictdomain.NormalizeItemStatusFilter(query.Status)
	if err != nil {
		return dictdomain.ItemListResponse{}, err
	}

	items, total, err := s.repo.ListItems(query, page, pageSize, status)
	if err != nil {
		return dictdomain.ItemListResponse{}, err
	}

	result := make([]dictdomain.ItemResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dictdomain.BuildItemResponse(item))
	}

	return dictdomain.ItemListResponse{Items: result, Total: total, Page: page, PageSize: pageSize}, nil
}

// CreateItem 在指定字典类型下创建字典项，校验类型存在性和键唯一性。
func (s *Service) CreateItem(req dictdomain.CreateItemRequest) (dictdomain.ItemResponse, error) {
	req, err := dictdomain.NormalizeCreateItemRequest(req)
	if err != nil {
		return dictdomain.ItemResponse{}, err
	}

	created := dictdomain.DictItemEntity{
		TypeID:  req.TypeID,
		ItemKey: req.ItemKey,
		Label:   req.Label,
		Value:   req.Value,
		TagType: req.TagType,
		Sort:    req.Sort,
		Status:  req.Status,
		Remark:  req.Remark,
	}

	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		if _, err := s.repo.FindTypeByID(tx, req.TypeID); err != nil {
			return err
		}

		exists, err := s.repo.ItemKeyExists(tx, req.TypeID, req.ItemKey)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("同一字典类型下的字典项编码不能重复")
		}
		return s.repo.CreateItem(tx, &created)
	})
	if err != nil {
		return dictdomain.ItemResponse{}, err
	}

	return dictdomain.BuildItemResponse(created), nil
}

// UpdateItem 更新指定字典项的基本信息。
func (s *Service) UpdateItem(itemID uint, req dictdomain.UpdateItemRequest) (dictdomain.ItemResponse, error) {
	req, err := dictdomain.NormalizeUpdateItemRequest(req)
	if err != nil {
		return dictdomain.ItemResponse{}, err
	}

	var updated dictdomain.DictItemEntity
	err = s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
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
		return dictdomain.ItemResponse{}, err
	}

	return dictdomain.BuildItemResponse(updated), nil
}

// UpdateItemStatus 切换字典项的启用/禁用状态。
func (s *Service) UpdateItemStatus(itemID uint, status model.SystemDictStatus) error {
	status, err := dictdomain.NormalizeStatus(status, false)
	if err != nil {
		return err
	}

	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			return err
		}
		return s.repo.UpdateItemStatus(tx, &item, status)
	})
}

// DeleteItem 删除指定字典项。
func (s *Service) DeleteItem(itemID uint) error {
	return s.tx.WithinTransaction(context.Background(), func(tx *gorm.DB) error {
		item, err := s.repo.FindItemByID(tx, itemID)
		if err != nil {
			return err
		}

		return s.repo.DeleteItem(tx, &item)
	})
}
