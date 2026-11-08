package collector

import (
	"context"
	"log/slog"
	"math"
	"time"

	"gogogo/internal/model"
	"gogogo/internal/repository"
)

// DemoCollector 周期性为所有有效数值测点生成模拟数据，用来演示采集链路。
// 它不连接真实 OPC UA/Modbus/MQTT，真实协议实现可以替换 Collect 方法。
type DemoCollector struct {
	store    repository.Store
	interval time.Duration
	logger   *slog.Logger
}

func NewDemo(store repository.Store, interval time.Duration, logger *slog.Logger) *DemoCollector {
	return &DemoCollector{store: store, interval: interval, logger: logger}
}

func (d *DemoCollector) Run(ctx context.Context) {
	ticker := time.NewTicker(d.interval)
	defer ticker.Stop()

	d.collectOnce(ctx)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			d.collectOnce(ctx)
		}
	}
}

func (d *DemoCollector) collectOnce(ctx context.Context) {
	bindings, err := d.store.ListActiveBindings(ctx)
	if err != nil {
		d.logger.Error("collector list active bindings", "error", err)
		return
	}
	if len(bindings) == 0 {
		return
	}

	now := time.Now().UTC()
	samples := make([]model.PointSample, 0, len(bindings))
	for index, binding := range bindings {
		if binding.ValueType != model.ValueTypeNumeric {
			continue
		}
		value := 20 + 5*math.Sin(float64(now.Unix())/30+float64(index))
		sourceID := binding.SourceID
		samples = append(samples, model.PointSample{
			PointID:     binding.PointID,
			Ts:          now,
			ValueType:   model.ValueTypeNumeric,
			ValueNum:    &value,
			Quality:     model.QualityGood,
			SourceID:    &sourceID,
			CollectedAt: now,
		})
	}
	if len(samples) == 0 {
		return
	}
	if err := d.store.AddSamples(ctx, samples); err != nil {
		d.logger.Error("collector write samples", "error", err)
		return
	}
	d.logger.Debug("collector wrote samples", "count", len(samples))
}

// Collect 是 Collector 接口的实现。当前 DemoCollector 只按绑定生成数值。
func (d *DemoCollector) Collect(_ context.Context, bindings []Binding) ([]Sample, error) {
	now := time.Now().UTC()
	samples := make([]Sample, 0, len(bindings))
	for index, binding := range bindings {
		if binding.ValueType != model.ValueTypeNumeric {
			continue
		}
		value := 20 + 5*math.Sin(float64(now.Unix())/30+float64(index))
		samples = append(samples, Sample{
			PointID:   binding.PointID,
			Ts:        now,
			ValueType: model.ValueTypeNumeric,
			ValueNum:  &value,
			Quality:   model.QualityGood,
			SourceID:  binding.SourceID,
		})
	}
	return samples, nil
}

func (d *DemoCollector) Close() error {
	return nil
}
