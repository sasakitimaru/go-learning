package model

import "time"

type User struct {
	ID        uint64    `gorm:"primary_key" json:"id"`
	Email     string    `gorm:"size:255;not null;unique" json:"email"`
	Password  string    `gorm:"size:100;not null;" json:"password"`
	CreatedAt time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"created_at"`
	UpdateAt  time.Time `gorm:"default:CURRENT_TIMESTAMP" json:"updated_at"`
}

type UserResponse struct {
	ID    uint64 `json:"id" gorm:"primary_key"`
	Email string `json:"email" gorm:"size:255;not null;unique"`
}
