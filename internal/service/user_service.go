package service

//go:generate mockery --name=UserRepositoryInterface --output=../../mocks --filename=user_repository_mock.go

import (
	"blog-api/internal/models"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrInvalidCredentials = errors.New("invalid username or password")
	ErrUserNotFound       = errors.New("user not found")
	ErrInvalidPassword    = errors.New("invalid password")
)

type UserRepositoryInterface interface {
	Create(ctx context.Context, user *models.User) error
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByID(ctx context.Context, id uint) (*models.User, error)
	ExistsByUsername(ctx context.Context, username string) (bool, error)
}

type UserService struct {
	userRepo  UserRepositoryInterface
	jwtSecret string
}

func NewUserService(userRepo UserRepositoryInterface, jwtSecret string) *UserService {
	return &UserService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *UserService) Register(ctx context.Context, username, password string) (*models.User, string, error) {
	exists, err := s.userRepo.ExistsByUsername(ctx, username)
	if err != nil {
		return nil, "", err
	}
	if exists {
		return nil, "", ErrUserAlreadyExists
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, "", err
	}

	user := &models.User{
		Username: username,
		Password: string(hash),
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, "", err
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) Login(ctx context.Context, username, password string) (*models.User, string, error) {
	user, err := s.userRepo.FindByUsername(ctx, username)
	if err != nil {
		return nil, "", ErrInvalidCredentials
	}

	if user == nil {
		return nil, "", ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, "", ErrInvalidPassword
	}

	token, err := s.generateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) GetUserByID(ctx context.Context, id uint) (*models.User, error) {
	return s.userRepo.FindByID(ctx, id)
}
