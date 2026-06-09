package application

import (
	loginlogdomain "ez-admin-gin/server/internal/modules/system/loginlog/domain"
	"ez-admin-gin/server/internal/platform/model"
)

type LoginLogRepository interface {
	List(query loginlogdomain.ListQuery, page int, pageSize int, status *model.LoginLogStatus) ([]loginlogdomain.Entity, int64, error)
}
