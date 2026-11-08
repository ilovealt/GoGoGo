package request

type EquipmentCreate struct {
	AreaID        uint64 `json:"area_id"`
	Code          string `json:"code"`
	Name          string `json:"name"`
	EquipmentType string `json:"equipment_type"`
	Description   string `json:"description"`
	Status        *uint8 `json:"status"`
}

type EquipmentUpdate struct {
	AreaID        *uint64 `json:"area_id"`
	Code          *string `json:"code"`
	Name          *string `json:"name"`
	EquipmentType *string `json:"equipment_type"`
	Description   *string `json:"description"`
	Status        *uint8  `json:"status"`
}
