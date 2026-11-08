package model

import "time"

// User 系统用户。PasswordHash 和 TokenVersion 不参与 JSON 序列化。
type User struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	Username     string    `json:"username" gorm:"size:64;not null"`
	PasswordHash string    `json:"-" gorm:"size:255;not null"`
	Nickname     string    `json:"nickname" gorm:"size:100;not null"`
	Role         string    `json:"role" gorm:"size:32;not null"`
	Status       uint8     `json:"status" gorm:"not null"`
	TokenVersion uint32    `json:"-" gorm:"not null"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (User) TableName() string {
	return "users"
}
