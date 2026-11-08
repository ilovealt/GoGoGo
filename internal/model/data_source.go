package model

import "time"

type DataSource struct {
	ID        uint64         `json:"id" gorm:"primaryKey"`
	Name      string         `json:"name" gorm:"size:128;not null;uniqueIndex"`
	Protocol  string         `json:"protocol" gorm:"size:32;not null"`
	Endpoint  string         `json:"endpoint" gorm:"type:text;not null"`
	Config    map[string]any `json:"config,omitempty" gorm:"type:jsonb;serializer:json"`
	Status    uint8          `json:"status" gorm:"not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
}

func (DataSource) TableName() string {
	return "data_sources"
}
