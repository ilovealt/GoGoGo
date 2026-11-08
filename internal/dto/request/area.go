package request

type AreaCreate struct {
	PlantID     uint64  `json:"plant_id"`
	ParentID    *uint64 `json:"parent_id"`
	Code        string  `json:"code"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Status      *uint8  `json:"status"`
}

type AreaUpdate struct {
	ParentID    *uint64 `json:"parent_id"`
	ClearParent bool    `json:"clear_parent"`
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Status      *uint8  `json:"status"`
}
