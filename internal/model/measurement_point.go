package model

import "time"

type MeasurementPoint struct {
	ID          uint64         `json:"id" gorm:"primaryKey"`
	EquipmentID *uint64        `json:"equipment_id,omitempty" gorm:"index"`
	Code        string         `json:"code" gorm:"size:128;not null"`
	Name        string         `json:"name" gorm:"size:128;not null"`
	Description string         `json:"description,omitempty" gorm:"type:text"`
	PointType   string         `json:"point_type" gorm:"size:32;not null"`
	ValueType   uint8          `json:"value_type" gorm:"not null"`
	Unit        string         `json:"unit,omitempty" gorm:"size:32"`
	Precision   uint8          `json:"precision" gorm:"not null"`
	MinValue    *float64       `json:"min_value,omitempty"`
	MaxValue    *float64       `json:"max_value,omitempty"`
	Deadband    float64        `json:"deadband" gorm:"not null"`
	Metadata    map[string]any `json:"metadata,omitempty" gorm:"type:jsonb;serializer:json"`
	Status      uint8          `json:"status" gorm:"not null"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   *time.Time     `json:"deleted_at,omitempty" gorm:"index"`
}

func (MeasurementPoint) TableName() string {
	return "measurement_points"
}
