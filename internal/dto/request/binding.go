package request

type BindingCreate struct {
	PointID          uint64         `json:"point_id"`
	SourceID         uint64         `json:"source_id"`
	AddressConfig    map[string]any `json:"address_config"`
	ConversionConfig map[string]any `json:"conversion_config"`
	Enabled          *bool          `json:"enabled"`
}

type BindingUpdate struct {
	AddressConfig    *map[string]any `json:"address_config"`
	ConversionConfig *map[string]any `json:"conversion_config"`
	Enabled          *bool           `json:"enabled"`
}
