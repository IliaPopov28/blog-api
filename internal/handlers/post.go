package handlers

import (
	"net/http"
	"strconv"

	"blog-api/internal/dto/request"
	"blog-api/internal/dto/response"
	"blog-api/internal/middleware"
	"blog-api/internal/service"

	"github.com/gin-gonic/gin"
)

type PostHandler struct {
	postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
	return &PostHandler{postService: postService}
}

func (h *PostHandler) Create(c *gin.Context) {
	var req request.CreatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	post, err := h.postService.Create(ctx, req.Title, req.Content, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to create post",
		})
		return
	}

	c.JSON(http.StatusCreated, response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (h *PostHandler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))

	ctx := c.Request.Context()
	posts, total, totalPages, err := h.postService.GetAll(ctx, page, limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to fetch posts",
		})
		return
	}

	var postResponses []response.PostResponse
	for _, p := range posts {
		postResponses = append(postResponses, response.PostResponse{
			ID:        p.ID,
			Title:     p.Title,
			Content:   p.Content,
			AuthorID:  p.AuthorID,
			CreatedAt: p.CreatedAt,
			UpdatedAt: p.UpdatedAt,
		})
	}

	c.JSON(http.StatusOK, response.PostListResponse{
		Posts:      postResponses,
		Page:       page,
		Limit:      limit,
		TotalCount: total,
		TotalPages: totalPages,
	})
}

func (h *PostHandler) GetByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid post ID",
		})
		return
	}

	ctx := c.Request.Context()
	post, err := h.postService.GetByID(ctx, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, response.ErrorResponse{
			Error: "Post not found",
		})
		return
	}

	c.JSON(http.StatusOK, response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (h *PostHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid post ID",
		})
		return
	}

	var req request.UpdatePostRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error:   "Validation failed",
			Details: err.Error(),
		})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	post, err := h.postService.Update(ctx, uint(id), req.Title, req.Content, userID)
	if err != nil {
		if err == service.ErrPostNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				Error: "Post not found",
			})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, response.ErrorResponse{
				Error: "You can only update your own posts",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to update post",
		})
		return
	}

	c.JSON(http.StatusOK, response.PostResponse{
		ID:        post.ID,
		Title:     post.Title,
		Content:   post.Content,
		AuthorID:  post.AuthorID,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	})
}

func (h *PostHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, response.ErrorResponse{
			Error: "Invalid post ID",
		})
		return
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, response.ErrorResponse{
			Error: "Unauthorized",
		})
		return
	}

	ctx := c.Request.Context()
	err = h.postService.Delete(ctx, uint(id), userID)
	if err != nil {
		if err == service.ErrPostNotFound {
			c.JSON(http.StatusNotFound, response.ErrorResponse{
				Error: "Post not found",
			})
			return
		}
		if err == service.ErrForbidden {
			c.JSON(http.StatusForbidden, response.ErrorResponse{
				Error: "You can only delete your own posts",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, response.ErrorResponse{
			Error: "Failed to delete post",
		})
		return
	}

	c.JSON(http.StatusOK, response.MessageResponse{
		Message: "Post deleted successfully",
	})
}
