package model

import "time"

// RefreshToken 只保存令牌哈希，不保存明文。
type RefreshToken struct {
	ID           uint64     `json:"id" gorm:"primaryKey"`
	UserID       uint64     `json:"user_id" gorm:"not null;index"`
	TokenHash    string     `json:"-" gorm:"size:64;not null;uniqueIndex"`
	FamilyID     string     `json:"family_id" gorm:"size:64;not null;index"`
	ReplacedByID *uint64    `json:"replaced_by_id,omitempty" gorm:"index"`
	ExpiresAt    time.Time  `json:"expires_at" gorm:"not null;index"`
	RevokedAt    *time.Time `json:"revoked_at,omitempty"`
	UserAgent    string     `json:"user_agent,omitempty" gorm:"type:text"`
	IP           string     `json:"ip,omitempty" gorm:"size:64"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (RefreshToken) TableName() string {
	return "refresh_tokens"
}
