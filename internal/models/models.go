package models

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Username  string    `json:"username" gorm:"unique;size:50;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Posts     []Post    `json:"-" gorm:"foreignKey:AuthorID"`
}

type Post struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"size:200;not null"`
	Content   string    `json:"content" gorm:"type:text;not null"`
	AuthorID  uint      `json:"author_id" gorm:"not null;index"`
	Author    User      `json:"author,omitempty" gorm:"foreignKey:AuthorID"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
