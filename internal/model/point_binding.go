package model

import "time"

type PointBinding struct {
	ID               uint64         `json:"id" gorm:"primaryKey"`
	PointID          uint64         `json:"point_id" gorm:"not null;index"`
	SourceID         uint64         `json:"source_id" gorm:"not null;index"`
	AddressConfig    map[string]any `json:"address_config,omitempty" gorm:"type:jsonb;serializer:json"`
	ConversionConfig map[string]any `json:"conversion_config,omitempty" gorm:"type:jsonb;serializer:json"`
	Enabled          bool           `json:"enabled" gorm:"not null"`
	CreatedAt        time.Time      `json:"created_at"`
	UpdatedAt        time.Time      `json:"updated_at"`
}

func (PointBinding) TableName() string {
	return "point_bindings"
}
