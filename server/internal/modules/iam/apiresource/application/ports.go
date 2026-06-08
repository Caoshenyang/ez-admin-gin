package application

import "ez-admin-gin/server/internal/platform/model"

// APIRepository 定义接口权限元数据的数据访问接口。
type APIRepository interface {
	List() ([]model.API, error)
}
