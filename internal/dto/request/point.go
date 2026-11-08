package request

type PointCreate struct {
	EquipmentID *uint64        `json:"equipment_id"`
	Code        string         `json:"code"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	PointType   string         `json:"point_type"`
	ValueType   *uint8         `json:"value_type"`
	Unit        string         `json:"unit"`
	Precision   *uint8         `json:"precision"`
	MinValue    *float64       `json:"min_value"`
	MaxValue    *float64       `json:"max_value"`
	Deadband    *float64       `json:"deadband"`
	Metadata    map[string]any `json:"metadata"`
	Status      *uint8         `json:"status"`
}

type PointUpdate struct {
	EquipmentID    *uint64         `json:"equipment_id"`
	ClearEquipment bool            `json:"clear_equipment"`
	Code           *string         `json:"code"`
	Name           *string         `json:"name"`
	Description    *string         `json:"description"`
	PointType      *string         `json:"point_type"`
	ValueType      *uint8          `json:"value_type"`
	Unit           *string         `json:"unit"`
	Precision      *uint8          `json:"precision"`
	MinValue       *float64        `json:"min_value"`
	MaxValue       *float64        `json:"max_value"`
	ClearMinValue  bool            `json:"clear_min_value"`
	ClearMaxValue  bool            `json:"clear_max_value"`
	Deadband       *float64        `json:"deadband"`
	Metadata       *map[string]any `json:"metadata"`
	Status         *uint8          `json:"status"`
}
