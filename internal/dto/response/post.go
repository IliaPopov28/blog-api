package response

import "time"

type PostResponse struct {
	ID        uint      `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	AuthorID  uint      `json:"author_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type PostListResponse struct {
	Posts      []PostResponse `json:"posts"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalCount int64          `json:"total_count"`
	TotalPages int            `json:"total_pages"`
}
