package repository

import (
	"blog-api/internal/models"
	"context"

	"gorm.io/gorm"
)

type PostRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) *PostRepository {
	return &PostRepository{db: db}
}

func (r *PostRepository) Create(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Create(post).Error
}

func (r *PostRepository) FindByID(ctx context.Context, id uint) (*models.Post, error) {
	var post models.Post
	err := r.db.WithContext(ctx).First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *PostRepository) FindAll(ctx context.Context, page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	if err := r.db.WithContext(ctx).Model(&models.Post{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * limit
	err := r.db.WithContext(ctx).Preload("Author").Order("created_at DESC").Offset(offset).Limit(limit).Find(&posts).Error

	return posts, total, err
}

func (r *PostRepository) Update(ctx context.Context, post *models.Post) error {
	return r.db.WithContext(ctx).Save(post).Error
}

func (r *PostRepository) Delete(ctx context.Context, id uint) error {
	return r.db.WithContext(ctx).Delete(&models.Post{}, id).Error
}

func (r *PostRepository) FindByAuthorID(ctx context.Context, authorID uint) ([]models.Post, error) {
	var posts []models.Post
	err := r.db.WithContext(ctx).Where("author_id = ?", authorID).Find(&posts).Error
	return posts, err
}
