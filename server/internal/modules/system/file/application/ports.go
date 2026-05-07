package application

import (
	"context"
	"mime/multipart"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	"ez-admin-gin/server/internal/platform/database"
	"ez-admin-gin/server/internal/platform/model"

	"gorm.io/gorm"
)

type FileRepository interface {
	List(query filedomain.ListQuery, page int, pageSize int) ([]model.SystemFile, int64, error)
	Create(db *gorm.DB, item *model.SystemFile) error
	DeleteByID(db *gorm.DB, id uint) error
}

type FileStorage interface {
	SaveUploadedFile(fileHeader *multipart.FileHeader) (filedomain.SavedUploadedFile, error)
	Delete(path string) error
}

type FileTransactor = database.Transactor

type AssetService interface {
	UploadEntity(ctx context.Context, uploaderID uint, fileHeader *multipart.FileHeader) (model.SystemFile, error)
	CleanupUploadedFile(item model.SystemFile)
}
