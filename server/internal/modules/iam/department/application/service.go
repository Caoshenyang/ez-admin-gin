package application

import (
	"fmt"
	"sort"
	"strings"

	departmentdomain "ez-admin-gin/server/internal/modules/iam/department/domain"
	departmentinfra "ez-admin-gin/server/internal/modules/iam/department/infra"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	"ez-admin-gin/server/internal/platform/datascope"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type Service struct {
	repo *departmentinfra.Repository
	db   *gorm.DB
}

func NewService(db *gorm.DB, repo *departmentinfra.Repository) *Service {
	return &Service{db: db, repo: repo}
}

func (s *Service) List(actor datascope.Actor, query departmentdomain.ListQuery) ([]departmentdomain.Response, error) {
	items, err := s.repo.List(actor, query)
	if err != nil {
		return nil, err
	}

	return buildTree(items), nil
}

func (s *Service) Create(actor datascope.Actor, req departmentdomain.CreateRequest) (departmentdomain.Response, error) {
	parentID, name, code, leaderUserID, sortValue, status, remark, err := departmentdomain.NormalizeDepartmentInput(
		req.ParentID, req.Name, req.Code, req.LeaderUserID, req.Sort, req.Status, req.Remark,
	)
	if err != nil {
		return departmentdomain.Response{}, err
	}

	var created departmentdomain.Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		exists, err := s.repo.CodeExists(tx, code, 0)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("部门编码已存在")
		}
		if err := s.repo.LeaderUsable(tx, leaderUserID); err != nil {
			return err
		}

		var parent model.Department
		if parentID != 0 {
			parent, err = s.repo.FindByIDInScope(tx, actor, parentID)
			if err != nil {
				return err
			}
		}

		created = departmentdomain.Entity{
			ParentID:     parentID,
			Ancestors:    departmentinfra.BuildAncestors(parent),
			Name:         name,
			Code:         code,
			LeaderUserID: leaderUserID,
			Sort:         sortValue,
			Status:       status,
			Remark:       remark,
		}
		return s.repo.Create(tx, &created)
	})
	if err != nil {
		return departmentdomain.Response{}, err
	}

	return departmentdomain.BuildResponse(created), nil
}

func (s *Service) Update(actor datascope.Actor, departmentID uint, req departmentdomain.UpdateRequest) (departmentdomain.Response, error) {
	parentID, name, code, leaderUserID, sortValue, status, remark, err := departmentdomain.NormalizeDepartmentInput(
		req.ParentID, req.Name, req.Code, req.LeaderUserID, req.Sort, req.Status, req.Remark,
	)
	if err != nil {
		return departmentdomain.Response{}, err
	}

	var updated departmentdomain.Entity
	err = s.db.Transaction(func(tx *gorm.DB) error {
		current, err := s.repo.FindByIDInScope(tx, actor, departmentID)
		if err != nil {
			return err
		}

		exists, err := s.repo.CodeExists(tx, code, departmentID)
		if err != nil {
			return err
		}
		if exists {
			return errorsx.BadRequest("部门编码已存在")
		}
		if err := s.repo.LeaderUsable(tx, leaderUserID); err != nil {
			return err
		}

		var parent model.Department
		if parentID != 0 {
			parent, err = s.repo.FindByIDInScope(tx, actor, parentID)
			if err != nil {
				return err
			}
		}
		if parentID == current.ID {
			return errorsx.BadRequest("不能把部门挂到自己下面")
		}

		oldFullPath := departmentinfra.FullPath(current)
		oldParentID := current.ParentID
		oldAncestors := current.Ancestors
		if parent.ID != 0 {
			parentFullPath := departmentinfra.FullPath(parent)
			if departmentinfra.IsDescendantPath(parentFullPath, oldFullPath) {
				return errorsx.BadRequest("不能把部门挂到自己的子部门下面")
			}
		}

		newAncestors := departmentinfra.BuildAncestors(parent)
		if err := s.repo.Update(tx, &current, parentID, newAncestors, name, code, leaderUserID, sortValue, status, remark); err != nil {
			return err
		}

		if oldParentID != parentID || oldAncestors != newAncestors {
			newFullPath := fmt.Sprintf("%s,%d", newAncestors, current.ID)
			children, err := s.repo.Subtree(tx, current.ID, oldFullPath)
			if err != nil {
				return err
			}
			for _, child := range children {
				newChildAncestors := newFullPath + strings.TrimPrefix(child.Ancestors, oldFullPath)
				if err := s.repo.UpdateAncestors(tx, child.ID, newChildAncestors); err != nil {
					return err
				}
			}
		}

		updated = current
		return nil
	})
	if err != nil {
		return departmentdomain.Response{}, err
	}

	return departmentdomain.BuildResponse(updated), nil
}

func (s *Service) UpdateStatus(actor datascope.Actor, departmentID uint, status model.DepartmentStatus) error {
	if !departmentdomain.ValidStatus(status) {
		return errorsx.BadRequest("部门状态不正确")
	}

	return s.db.Transaction(func(tx *gorm.DB) error {
		current, err := s.repo.FindByIDInScope(tx, actor, departmentID)
		if err != nil {
			return err
		}

		return s.repo.UpdateStatus(tx, &current, status)
	})
}

func buildTree(items []model.Department) []departmentdomain.Response {
	nodes := make(map[uint]*departmentdomain.Response, len(items))
	roots := make([]*departmentdomain.Response, 0)

	for _, item := range items {
		response := departmentdomain.BuildResponse(item)
		nodes[item.ID] = &response
	}

	for _, item := range items {
		node := nodes[item.ID]
		if parent, ok := nodes[item.ParentID]; ok {
			parent.Children = append(parent.Children, *node)
			continue
		}
		roots = append(roots, node)
	}

	var sortTree func(items []*departmentdomain.Response)
	sortTree = func(items []*departmentdomain.Response) {
		sort.Slice(items, func(i, j int) bool {
			if items[i].Sort == items[j].Sort {
				return items[i].ID < items[j].ID
			}
			return items[i].Sort < items[j].Sort
		})
		for _, item := range items {
			if len(item.Children) == 0 {
				continue
			}

			children := make([]*departmentdomain.Response, 0, len(item.Children))
			for idx := range item.Children {
				children = append(children, &item.Children[idx])
			}
			sortTree(children)
			sort.Slice(item.Children, func(i, j int) bool {
				if item.Children[i].Sort == item.Children[j].Sort {
					return item.Children[i].ID < item.Children[j].ID
				}
				return item.Children[i].Sort < item.Children[j].Sort
			})
		}
	}

	sortTree(roots)

	result := make([]departmentdomain.Response, 0, len(roots))
	for _, item := range roots {
		result = append(result, *item)
	}

	return result
}
