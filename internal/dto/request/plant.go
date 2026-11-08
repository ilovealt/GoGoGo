package request

type PlantCreate struct {
	Code        string `json:"code"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Timezone    string `json:"timezone"`
	Status      *uint8 `json:"status"`
}

type PlantUpdate struct {
	Code        *string `json:"code"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	Timezone    *string `json:"timezone"`
	Status      *uint8  `json:"status"`
}
