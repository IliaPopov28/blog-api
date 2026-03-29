package mocks

import (
	"context"
	"testing"

	"blog-api/internal/models"
	svc "blog-api/internal/service"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestUserService_Register_Success(t *testing.T) {
	mockRepo := NewMockUserRepositoryInterface(t)
	service := svc.NewUserService(mockRepo, "test-secret")

	mockRepo.On("ExistsByUsername", context.Background(), "newuser").Return(false, nil)
	mockRepo.On("Create", context.Background(), mock.AnythingOfType("*models.User")).Return(nil)

	user, token, err := service.Register(context.Background(), "newuser", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, "newuser", user.Username)
}

func TestUserService_Register_UserAlreadyExists(t *testing.T) {
	mockRepo := NewMockUserRepositoryInterface(t)
	service := svc.NewUserService(mockRepo, "test-secret")

	mockRepo.On("ExistsByUsername", context.Background(), "existinguser").Return(true, nil)

	user, token, err := service.Register(context.Background(), "existinguser", "password123")

	assert.Error(t, err)
	assert.Equal(t, svc.ErrUserAlreadyExists, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
}

func TestUserService_Login_Success(t *testing.T) {
	mockRepo := NewMockUserRepositoryInterface(t)
	service := svc.NewUserService(mockRepo, "test-secret")

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	user := &models.User{
		ID:       1,
		Username: "testuser",
		Password: string(hashedPassword),
	}

	mockRepo.On("FindByUsername", context.Background(), "testuser").Return(user, nil)

	returnedUser, token, err := service.Login(context.Background(), "testuser", "password123")

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Equal(t, user.Username, returnedUser.Username)
}

func TestUserService_Login_InvalidCredentials(t *testing.T) {
	mockRepo := NewMockUserRepositoryInterface(t)
	service := svc.NewUserService(mockRepo, "test-secret")

	// Return error from repository to simulate invalid user
	mockRepo.On("FindByUsername", context.Background(), "testuser").Return(nil, assert.AnError)

	user, token, err := service.Login(context.Background(), "testuser", "wrongpassword")

	assert.Error(t, err)
	assert.Equal(t, svc.ErrInvalidCredentials, err)
	assert.Empty(t, token)
	assert.Nil(t, user)
}

func TestUserService_GetUserByID(t *testing.T) {
	mockRepo := NewMockUserRepositoryInterface(t)
	service := svc.NewUserService(mockRepo, "test-secret")

	user := &models.User{ID: 1, Username: "testuser"}
	mockRepo.On("FindByID", context.Background(), uint(1)).Return(user, nil)

	result, err := service.GetUserByID(context.Background(), 1)

	assert.NoError(t, err)
	assert.Equal(t, user.Username, result.Username)
}
