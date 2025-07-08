package controllers

import (
	"blog-api/config"
	"blog-api/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type PostResponse struct {
	ID       uint   `json:"id"`
	Title    string `json:"title"`
	Content  string `json:"content"`
	AuthorID uint   `json:"author_id"`
}

func toPostResponse(post *models.Post) PostResponse {
	return PostResponse{
		ID:       post.ID,
		Title:    post.Title,
		Content:  post.Content,
		AuthorID: post.AuthorID,
	}
}

func getPostByID(id int) (*models.Post, error) {
	var post models.Post
	if err := config.DB.First(&post, id).Error; err != nil {
		return nil, err
	}
	return &post, nil
}

func GetPosts(c *gin.Context) {
	var posts []models.Post
	config.DB.Find(&posts)
	var resp []PostResponse
	for _, p := range posts {
		resp = append(resp, toPostResponse(&p))
	}
	c.JSON(http.StatusOK, resp)
}

func CreatePost(c *gin.Context) {
	var post models.Post
	if err := c.ShouldBindJSON(&post); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userID, _ := c.Get("user_id")
	post.AuthorID = userID.(uint)
	if err := config.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create post"})
		return
	}
	c.JSON(http.StatusCreated, toPostResponse(&post))
}

func GetPostByID(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	post, err := getPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	c.JSON(http.StatusOK, toPostResponse(post))
}

func UpdatePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid id"})
		return
	}
	post, err := getPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	userID, _ := c.Get("user_id")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only update your own posts"})
		return
	}
	var updateData models.Post
	if err = c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	post.Title = updateData.Title
	post.Content = updateData.Content
	if err := config.DB.Save(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update post"})
		return
	}
	c.JSON(http.StatusOK, toPostResponse(post))
}

func DeletePost(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}
	post, err := getPostByID(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Post not found"})
		return
	}
	userID, _ := c.Get("user_id")
	if post.AuthorID != userID.(uint) {
		c.JSON(http.StatusForbidden, gin.H{"error": "You can only delete your own posts"})
		return
	}
	if err := config.DB.Delete(post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete post"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Post deleted"})
}
