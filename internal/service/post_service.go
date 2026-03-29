package service

//go:generate mockery --name=PostRepositoryInterface --output=../../mocks --filename=post_repository_mock.go

import (
	"blog-api/internal/models"
	"context"
	"errors"
	"math"
)

var (
	ErrPostNotFound = errors.New("post not found")
	ErrForbidden    = errors.New("you can only modify your own posts")
)

type PostRepositoryInterface interface {
	Create(ctx context.Context, post *models.Post) error
	FindByID(ctx context.Context, id uint) (*models.Post, error)
	FindAll(ctx context.Context, page, limit int) ([]models.Post, int64, error)
	Update(ctx context.Context, post *models.Post) error
	Delete(ctx context.Context, id uint) error
	FindByAuthorID(ctx context.Context, authorID uint) ([]models.Post, error)
}

type PostService struct {
	postRepo PostRepositoryInterface
}

func NewPostService(postRepo PostRepositoryInterface) *PostService {
	return &PostService{postRepo: postRepo}
}

func (s *PostService) Create(ctx context.Context, title, content string, authorID uint) (*models.Post, error) {
	post := &models.Post{
		Title:    title,
		Content:  content,
		AuthorID: authorID,
	}

	if err := s.postRepo.Create(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) GetByID(ctx context.Context, id uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}
	return post, nil
}

func (s *PostService) GetAll(ctx context.Context, page, limit int) ([]models.Post, int64, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.postRepo.FindAll(ctx, page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := int(math.Ceil(float64(total) / float64(limit)))

	return posts, total, totalPages, nil
}

func (s *PostService) Update(ctx context.Context, id uint, title, content string, userID uint) (*models.Post, error) {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return nil, ErrPostNotFound
	}

	if post.AuthorID != userID {
		return nil, ErrForbidden
	}

	post.Title = title
	post.Content = content

	if err := s.postRepo.Update(ctx, post); err != nil {
		return nil, err
	}

	return post, nil
}

func (s *PostService) Delete(ctx context.Context, id uint, userID uint) error {
	post, err := s.postRepo.FindByID(ctx, id)
	if err != nil {
		return ErrPostNotFound
	}

	if post.AuthorID != userID {
		return ErrForbidden
	}

	return s.postRepo.Delete(ctx, id)
}

func (s *PostService) GetByAuthor(ctx context.Context, authorID uint) ([]models.Post, error) {
	return s.postRepo.FindByAuthorID(ctx, authorID)
}
