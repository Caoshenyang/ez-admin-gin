package application

import operationlogdomain "ez-admin-gin/server/internal/modules/system/operationlog/domain"

type OperationLogRepository interface {
	List(query operationlogdomain.ListQuery, page int, pageSize int, success *bool) ([]operationlogdomain.Entity, int64, error)
}
