package application

import (
	"context"
	"errors"
	"mime/multipart"
	"testing"

	filedomain "ez-admin-gin/server/internal/modules/system/file/domain"
	errorsx "ez-admin-gin/server/internal/pkg/errorsx"
	platformConfig "ez-admin-gin/server/internal/platform/config"
	"ez-admin-gin/server/internal/platform/model"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type fileTestTransactor struct {
	err error
}

func (t fileTestTransactor) WithinTransaction(_ context.Context, fn func(tx *gorm.DB) error) error {
	if t.err != nil {
		return t.err
	}
	return fn(nil)
}

type fileTestRepo struct{}

func (r *fileTestRepo) List(query filedomain.ListQuery, page int, pageSize int) ([]model.SystemFile, int64, error) {
	return nil, 0, nil
}
func (r *fileTestRepo) Create(db *gorm.DB, item *model.SystemFile) error { return nil }
func (r *fileTestRepo) DeleteByID(db *gorm.DB, id uint) error            { return nil }

type fileTestStorage struct {
	saveErr     error
	deletePaths []string
	savedUpload filedomain.SavedUploadedFile
}

func (s *fileTestStorage) SaveUploadedFile(fileHeader *multipart.FileHeader) (filedomain.SavedUploadedFile, error) {
	if s.saveErr != nil {
		return filedomain.SavedUploadedFile{}, s.saveErr
	}
	return s.savedUpload, nil
}

func (s *fileTestStorage) Delete(path string) error {
	s.deletePaths = append(s.deletePaths, path)
	return nil
}

func TestUploadEntityReturnsInternalErrorWhenStorageFails(t *testing.T) {
	saveErr := errors.New("disk failed")
	storage := &fileTestStorage{saveErr: saveErr}
	service := NewService(
		fileTestTransactor{},
		&fileTestRepo{},
		storage,
		platformConfig.UploadConfig{AllowedExts: []string{".png"}},
		zap.NewNop(),
	)

	_, err := service.UploadEntity(context.Background(), 1, &multipart.FileHeader{Filename: "avatar.png", Size: 16})
	var internalErr *errorsx.Error
	if !errors.As(err, &internalErr) || internalErr.Message != "保存文件失败" || !errors.Is(err, saveErr) {
		t.Fatalf("expected storage failure to be wrapped, got %v", err)
	}
}

func TestUploadEntityDeletesSavedFileWhenTransactionFails(t *testing.T) {
	txErr := errors.New("db down")
	storage := &fileTestStorage{
		savedUpload: filedomain.SavedUploadedFile{
			OriginalName: "avatar.png",
			FileName:     "avatar.png",
			Ext:          ".png",
			MimeType:     "image/png",
			Size:         16,
			Sha256:       "hash",
			Path:         "uploads/20260101/avatar.png",
			URL:          "/uploads/20260101/avatar.png",
			AbsolutePath: "/tmp/avatar.png",
		},
	}
	service := NewService(
		fileTestTransactor{err: txErr},
		&fileTestRepo{},
		storage,
		platformConfig.UploadConfig{AllowedExts: []string{".png"}},
		zap.NewNop(),
	)

	_, err := service.UploadEntity(context.Background(), 1, &multipart.FileHeader{Filename: "avatar.png", Size: 16})
	var internalErr *errorsx.Error
	if !errors.As(err, &internalErr) || internalErr.Message != "保存文件记录失败" || !errors.Is(err, txErr) {
		t.Fatalf("expected transaction failure to be wrapped, got %v", err)
	}
	if len(storage.deletePaths) != 1 || storage.deletePaths[0] != "/tmp/avatar.png" {
		t.Fatalf("expected cleanup of saved file, got %v", storage.deletePaths)
	}
}
