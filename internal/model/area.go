package model

import "time"

type Area struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	PlantID     uint64    `json:"plant_id" gorm:"not null;index"`
	ParentID    *uint64   `json:"parent_id,omitempty" gorm:"index"`
	Code        string    `json:"code" gorm:"size:64;not null"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	Status      uint8     `json:"status" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Area) TableName() string {
	return "areas"
}
