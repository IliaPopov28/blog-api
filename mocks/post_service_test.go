package mocks

import (
	"context"
	"testing"

	"blog-api/internal/models"
	svc "blog-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestPostService_Create_Success(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*models.Post")).Return(nil)

	post, err := service.Create(context.Background(), "Test Title", "Test Content", 1)

	assert.NoError(t, err)
	assert.Equal(t, "Test Title", post.Title)
	assert.Equal(t, "Test Content", post.Content)
	assert.Equal(t, uint(1), post.AuthorID)
}

func TestPostService_GetByID_Success(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	post := &models.Post{ID: 1, Title: "Test", AuthorID: 1}
	mockRepo.On("FindByID", context.Background(), uint(1)).Return(post, nil)

	result, err := service.GetByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, post.Title, result.Title)
}

func TestPostService_GetByID_NotFound(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	mockRepo.On("FindByID", context.Background(), uint(999)).Return(nil, assert.AnError)

	_, err := service.GetByID(context.Background(), 999)

	assert.Error(t, err)
	assert.Equal(t, svc.ErrPostNotFound, err)
}

func TestPostService_Update_Success(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	post := &models.Post{ID: 1, Title: "Old", Content: "Old", AuthorID: 1}
	mockRepo.On("FindByID", context.Background(), uint(1)).Return(post, nil)
	mockRepo.On("Update", context.Background(), mock.AnythingOfType("*models.Post")).Return(nil)

	result, err := service.Update(context.Background(), 1, "New Title", "New Content", 1)

	assert.NoError(t, err)
	assert.Equal(t, "New Title", result.Title)
}

func TestPostService_Update_Forbidden(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	post := &models.Post{ID: 1, Title: "Old", AuthorID: 1}
	mockRepo.On("FindByID", context.Background(), uint(1)).Return(post, nil)

	_, err := service.Update(context.Background(), 1, "New Title", "New Content", 999)

	assert.Error(t, err)
	assert.Equal(t, svc.ErrForbidden, err)
}

func TestPostService_Delete_Success(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	post := &models.Post{ID: 1, AuthorID: 1}
	mockRepo.On("FindByID", context.Background(), uint(1)).Return(post, nil)
	mockRepo.On("Delete", context.Background(), uint(1)).Return(nil)

	err := service.Delete(context.Background(), 1, 1)

	assert.NoError(t, err)
}

func TestPostService_GetAll_WithPagination(t *testing.T) {
	mockRepo := NewMockPostRepositoryInterface(t)
	service := svc.NewPostService(mockRepo)

	posts := []models.Post{
		{ID: 1, Title: "Post 1"},
		{ID: 2, Title: "Post 2"},
	}
	mockRepo.On("FindAll", context.Background(), 1, 20).Return(posts, int64(2), nil)

	result, total, totalPages, err := service.GetAll(context.Background(), 1, 20)

	assert.NoError(t, err)
	assert.Equal(t, 2, len(result))
	assert.Equal(t, int64(2), total)
	assert.Equal(t, 1, totalPages)
}
