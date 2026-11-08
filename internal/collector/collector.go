package collector

import (
	"context"
	"time"
)

type Binding struct {
	BindingID        uint64
	PointID          uint64
	SourceID         uint64
	ValueType        uint8
	AddressConfig    map[string]any
	ConversionConfig map[string]any
	Deadband         float64
}

type Sample struct {
	PointID   uint64
	Ts        time.Time
	ValueType uint8
	ValueNum  *float64
	ValueText *string
	ValueBool *bool
	Quality   uint8
	SourceID  uint64
}

type Collector interface {
	Collect(ctx context.Context, bindings []Binding) ([]Sample, error)
	Close() error
}
