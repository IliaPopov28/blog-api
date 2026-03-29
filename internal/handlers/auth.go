package handlers

import (
	"net/http"

	"blog-api/internal/dto/request"
	"blog-api/internal/dto/response"
	"blog-api/internal/middleware"
	"blog-api/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService *service.UserService
}

func NewAuthHandler(userService *service.UserService) *AuthHandler {
	return &AuthHandler{userService: userService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req request.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	user, token, err := h.userService.Register(ctx, req.Username, req.Password)
	if err != nil {
		if err == service.ErrUserAlreadyExists {
			c.JSON(http.StatusConflict, response.ErrorResponse{
				Error: "Username already taken",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to register user",
		})
		return
	}

	c.JSON(http.StatusCreated, response.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req request.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	ctx := c.Request.Context()
	user, token, err := h.userService.Login(ctx, req.Username, req.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Invalid username or password",
		})
		return
	}

	c.JSON(http.StatusOK, response.AuthResponse{
		Token:    token,
		UserID:   user.ID,
		Username: user.Username,
	})
}

func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	user, err := h.userService.GetUserByID(ctx, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":       user.ID,
		"username": user.Username,
	})
}
