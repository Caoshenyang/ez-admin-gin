package database

import (
	"context"

	"gorm.io/gorm"
)

type Transactor interface {
	WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error
}

type GormTransactor struct {
	db *gorm.DB
}

func NewTransactor(db *gorm.DB) *GormTransactor {
	return &GormTransactor{db: db}
}

func (t *GormTransactor) WithinTransaction(ctx context.Context, fn func(tx *gorm.DB) error) error {
	query := t.db
	if ctx != nil {
		query = query.WithContext(ctx)
	}

	return query.Transaction(fn)
}
