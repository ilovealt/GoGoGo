package request

type SourceCreate struct {
	Name     string         `json:"name"`
	Protocol string         `json:"protocol"`
	Endpoint string         `json:"endpoint"`
	Config   map[string]any `json:"config"`
	Status   *uint8         `json:"status"`
}

type SourceUpdate struct {
	Name     *string         `json:"name"`
	Protocol *string         `json:"protocol"`
	Endpoint *string         `json:"endpoint"`
	Config   *map[string]any `json:"config"`
	Status   *uint8          `json:"status"`
}
