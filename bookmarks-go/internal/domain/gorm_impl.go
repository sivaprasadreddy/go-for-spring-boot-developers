package domain

import (
	"context"
	"errors"
	"github.com/sivaprasadreddy/bookmarks/internal/config"
	"gorm.io/gorm"
)

type gormBookmarkRepo struct {
	db     *gorm.DB
	logger *config.Logger
}

func NewGormBookmarkRepository(db *gorm.DB, logger *config.Logger) BookmarkRepository {
	return gormBookmarkRepo{
		db:     db,
		logger: logger,
	}
}

func (r gormBookmarkRepo) GetAll(ctx context.Context) ([]Bookmark, error) {
	var bookmarks []Bookmark
	result := r.db.WithContext(ctx).Find(&bookmarks)
	return bookmarks, result.Error
}

func (r gormBookmarkRepo) GetByID(ctx context.Context, id int) (*Bookmark, error) {
	var bookmark Bookmark
	err := r.db.WithContext(ctx).First(&bookmark, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &bookmark, err
}

func (r gormBookmarkRepo) Create(ctx context.Context, b Bookmark) (*Bookmark, error) {
	return &b, r.db.WithContext(ctx).Save(&b).Error
}

func (r gormBookmarkRepo) Update(ctx context.Context, b Bookmark) error {
	return r.db.WithContext(ctx).Model(&b).Select("Title", "Url").Updates(b).Error
}

func (r gormBookmarkRepo) Delete(ctx context.Context, id int) error {
	return r.db.WithContext(ctx).Delete(&Bookmark{}, id).Error
}
