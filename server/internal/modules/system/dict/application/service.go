package application

import (
	dictdomain "ez-admin-gin/server/internal/modules/system/dict/domain"
	dictinfra "ez-admin-gin/server/internal/modules/system/dict/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/pkg/paging"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	db   *gorm.DB
	repo *dictinfra.Repository
}

func NewService(db *gorm.DB, repo *dictinfra.Repository) *Service {
	return &Service{db: db, repo: repo}
}

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

	err = s.db.Transaction(func(tx *gorm.DB) error {
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

func (s *Service) UpdateType(typeID uint, req dictdomain.UpdateTypeRequest) (dictdomain.TypeResponse, error) {
	req, err := dictdomain.NormalizeUpdateTypeRequest(req)
	if err != nil {
		return dictdomain.TypeResponse{}, err
	}

	var updated dictdomain.DictTypeEntity
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
		return dictdomain.TypeResponse{}, err
	}

	return dictdomain.BuildTypeResponse(updated), nil
}

func (s *Service) UpdateTypeStatus(typeID uint, status model.SystemDictStatus) error {
	status, err := dictdomain.NormalizeStatus(status, false)
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

func (s *Service) ListItems(query dictdomain.ItemListQuery) (dictdomain.ItemListResponse, error) {
	if query.TypeID == 0 {
		return dictdomain.ItemListResponse{}, errorsx.BadRequest("字典类型 ID 不正确")
	}

	page, pageSize := paging.NormalizePage(query.Page, query.PageSize)
	status, err := dictdomain.NormalizeItemStatusFilter(query.Status)
	if err != nil {
		return dictdomain.ItemListResponse{}, err
	}

	if _, err := s.repo.FindTypeByID(s.db, query.TypeID); err != nil {
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

	err = s.db.Transaction(func(tx *gorm.DB) error {
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

func (s *Service) UpdateItem(itemID uint, req dictdomain.UpdateItemRequest) (dictdomain.ItemResponse, error) {
	req, err := dictdomain.NormalizeUpdateItemRequest(req)
	if err != nil {
		return dictdomain.ItemResponse{}, err
	}

	var updated dictdomain.DictItemEntity
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
		return dictdomain.ItemResponse{}, err
	}

	return dictdomain.BuildItemResponse(updated), nil
}

func (s *Service) UpdateItemStatus(itemID uint, status model.SystemDictStatus) error {
	status, err := dictdomain.NormalizeStatus(status, false)
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
