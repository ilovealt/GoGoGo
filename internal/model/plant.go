package model

import "time"

type Plant struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	Code        string    `json:"code" gorm:"size:64;not null"`
	Name        string    `json:"name" gorm:"size:128;not null"`
	Description string    `json:"description,omitempty" gorm:"type:text"`
	Timezone    string    `json:"timezone" gorm:"size:64;not null"`
	Status      uint8     `json:"status" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (Plant) TableName() string {
	return "plants"
}
