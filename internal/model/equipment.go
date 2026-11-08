package model

import "time"

type Equipment struct {
	ID            uint64     `json:"id" gorm:"primaryKey"`
	AreaID        uint64     `json:"area_id" gorm:"not null;index"`
	Code          string     `json:"code" gorm:"size:64;not null"`
	Name          string     `json:"name" gorm:"size:128;not null"`
	EquipmentType string     `json:"equipment_type,omitempty" gorm:"size:64"`
	Description   string     `json:"description,omitempty" gorm:"type:text"`
	Status        uint8      `json:"status" gorm:"not null"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	DeletedAt     *time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

func (Equipment) TableName() string {
	return "equipments"
}
