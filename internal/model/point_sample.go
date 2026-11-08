package model

import "time"

// PointSample 历史时序数据，主键为 (point_id, ts)。
type PointSample struct {
	PointID     uint64    `json:"point_id" gorm:"primaryKey;autoIncrement:false"`
	Ts          time.Time `json:"ts" gorm:"primaryKey"`
	ValueType   uint8     `json:"value_type" gorm:"not null"`
	ValueNum    *float64  `json:"value_num,omitempty"`
	ValueText   *string   `json:"value_text,omitempty"`
	ValueBool   *bool     `json:"value_bool,omitempty"`
	Quality     uint8     `json:"quality" gorm:"not null"`
	SourceID    *uint64   `json:"source_id,omitempty" gorm:"index"`
	CollectedAt time.Time `json:"collected_at" gorm:"not null"`
}

func (PointSample) TableName() string {
	return "point_samples"
}

// PointLatestValue 每个测点一行，供当前值查询使用。
type PointLatestValue struct {
	PointID   uint64    `json:"point_id" gorm:"primaryKey;autoIncrement:false"`
	Ts        time.Time `json:"ts" gorm:"not null"`
	ValueType uint8     `json:"value_type" gorm:"not null"`
	ValueNum  *float64  `json:"value_num,omitempty"`
	ValueText *string   `json:"value_text,omitempty"`
	ValueBool *bool     `json:"value_bool,omitempty"`
	Quality   uint8     `json:"quality" gorm:"not null"`
	SourceID  *uint64   `json:"source_id,omitempty" gorm:"index"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (PointLatestValue) TableName() string {
	return "point_latest_values"
}

// SamplePoint 是历史查询接口返回的单个样本点。
type SamplePoint struct {
	Ts        time.Time `json:"ts"`
	ValueNum  *float64  `json:"value_num,omitempty"`
	ValueText *string   `json:"value_text,omitempty"`
	ValueBool *bool     `json:"value_bool,omitempty"`
	Quality   uint8     `json:"quality"`
}

// BindingView 是采集器使用的有效绑定视图。
type BindingView struct {
	BindingID        uint64         `json:"binding_id"`
	PointID          uint64         `json:"point_id"`
	PointCode        string         `json:"point_code"`
	PointName        string         `json:"point_name"`
	ValueType        uint8          `json:"value_type"`
	Unit             string         `json:"unit,omitempty"`
	Deadband         float64        `json:"deadband"`
	SourceID         uint64         `json:"source_id"`
	SourceName       string         `json:"source_name"`
	Protocol         string         `json:"protocol"`
	Endpoint         string         `json:"endpoint"`
	AddressConfig    map[string]any `json:"address_config,omitempty" gorm:"serializer:json"`
	ConversionConfig map[string]any `json:"conversion_config,omitempty" gorm:"serializer:json"`
}

// PointLatestView 是工厂当前值接口返回的测点信息加最新值。
type PointLatestView struct {
	PointID   uint64    `json:"point_id"`
	PointCode string    `json:"point_code"`
	PointName string    `json:"point_name"`
	PointType string    `json:"point_type"`
	Unit      string    `json:"unit,omitempty"`
	Ts        time.Time `json:"ts"`
	ValueType uint8     `json:"value_type"`
	ValueNum  *float64  `json:"value_num,omitempty"`
	ValueText *string   `json:"value_text,omitempty"`
	ValueBool *bool     `json:"value_bool,omitempty"`
	Quality   uint8     `json:"quality"`
	SourceID  *uint64   `json:"source_id,omitempty"`
}
