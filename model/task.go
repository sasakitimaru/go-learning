package model

import "time"

type Task struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Title     string    `gorm:"size:255;not null;unique" json:"title"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
	User      User      `gorm:"foreignkey:UserID; constraint:OnDelete:CASCADE" json:"user"`
	UserID    uint64    `gorm:"not null" json:"user_id"`
}

type TaskResponse struct {
	ID        uint64    `json:"id" gorm:"primary_key"`
	Title     string    `json:"title" gorm:"size:255;not null;unique"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}
