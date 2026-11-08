package service

import (
	"context"
	"sort"
	"strings"
	"time"

	"gogogo/internal/model"
	"gogogo/internal/pkg/apperrors"
	"gogogo/internal/repository"
)

var allowedIntervals = map[string]time.Duration{
	"1m": time.Minute,
	"5m": 5 * time.Minute,
	"1h": time.Hour,
	"1d": 24 * time.Hour,
}

var allowedAggregations = map[string]bool{
	"avg":   true,
	"min":   true,
	"max":   true,
	"first": true,
	"last":  true,
}

type SampleService struct {
	store repository.Store
}

func NewSampleService(store repository.Store) *SampleService {
	return &SampleService{store: store}
}

func (s *SampleService) Query(ctx context.Context, pointID uint64, start, end time.Time, interval, agg string, limit int) ([]model.SamplePoint, error) {
	if start.IsZero() || end.IsZero() {
		return nil, apperrors.New(apperrors.ErrInvalid, "start and end are required")
	}
	if !end.After(start) {
		return nil, apperrors.New(apperrors.ErrInvalid, "end must be after start")
	}
	if end.Sub(start) > 31*24*time.Hour {
		return nil, apperrors.New(apperrors.ErrInvalid, "time range must not exceed 31 days")
	}

	interval = strings.TrimSpace(interval)
	agg = strings.TrimSpace(agg)
	step, ok := allowedIntervals[interval]
	if interval != "" && !ok {
		return nil, apperrors.New(apperrors.ErrInvalid, "interval must be one of 1m, 5m, 1h, 1d")
	}
	if agg != "" && !allowedAggregations[agg] {
		return nil, apperrors.New(apperrors.ErrInvalid, "agg must be one of avg, min, max, first, last")
	}
	if interval != "" && agg == "" {
		agg = "last"
	}
	if agg != "" && interval == "" {
		interval = "1m"
		step = time.Minute
	}
	if limit == 0 {
		limit = 1000
	}
	if limit < 1 || limit > 10000 {
		return nil, apperrors.New(apperrors.ErrInvalid, "limit must be between 1 and 10000")
	}

	if _, err := s.store.GetPoint(ctx, pointID); err != nil {
		return nil, err
	}
	samples, err := s.store.ListSamples(ctx, pointID, start, end)
	if err != nil {
		return nil, err
	}

	if interval == "" {
		result := make([]model.SamplePoint, 0, len(samples))
		for _, sample := range samples {
			result = append(result, model.SamplePoint{
				Ts:        sample.Ts,
				ValueNum:  sample.ValueNum,
				ValueText: sample.ValueText,
				ValueBool: sample.ValueBool,
				Quality:   sample.Quality,
			})
			if len(result) >= limit {
				break
			}
		}
		return result, nil
	}

	groups := make(map[int64][]model.PointSample)
	keys := make([]int64, 0)
	for _, sample := range samples {
		key := sample.Ts.UTC().Truncate(step).UnixNano()
		if _, exists := groups[key]; !exists {
			keys = append(keys, key)
		}
		groups[key] = append(groups[key], sample)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })

	result := make([]model.SamplePoint, 0, len(keys))
	for _, key := range keys {
		result = append(result, aggregateSamples(groups[key], time.Unix(0, key).UTC(), agg))
		if len(result) >= limit {
			break
		}
	}
	return result, nil
}

func (s *SampleService) Latest(ctx context.Context, pointID uint64) (model.PointLatestValue, error) {
	if _, err := s.store.GetPoint(ctx, pointID); err != nil {
		return model.PointLatestValue{}, err
	}
	return s.store.GetLatest(ctx, pointID)
}

func (s *SampleService) LatestByPlant(ctx context.Context, plantID uint64) ([]model.PointLatestView, error) {
	if _, err := s.store.GetPlant(ctx, plantID); err != nil {
		return nil, err
	}
	return s.store.ListLatestByPlant(ctx, plantID)
}

func aggregateSamples(samples []model.PointSample, ts time.Time, agg string) model.SamplePoint {
	first := samples[0]
	result := model.SamplePoint{Ts: ts, Quality: first.Quality}
	if first.ValueType != model.ValueTypeNumeric || first.ValueNum == nil {
		result.ValueText = first.ValueText
		result.ValueBool = first.ValueBool
		return result
	}

	switch agg {
	case "first":
		value := *first.ValueNum
		result.ValueNum = &value
	case "last":
		value := *samples[len(samples)-1].ValueNum
		result.ValueNum = &value
	case "min":
		value := *first.ValueNum
		for _, sample := range samples[1:] {
			if sample.ValueNum != nil && *sample.ValueNum < value {
				value = *sample.ValueNum
			}
		}
		result.ValueNum = &value
	case "max":
		value := *first.ValueNum
		for _, sample := range samples[1:] {
			if sample.ValueNum != nil && *sample.ValueNum > value {
				value = *sample.ValueNum
			}
		}
		result.ValueNum = &value
	default: // avg
		sum, count := 0.0, 0
		for _, sample := range samples {
			if sample.ValueNum != nil {
				sum += *sample.ValueNum
				count++
			}
		}
		if count == 0 {
			result.ValueNum = first.ValueNum
		} else {
			value := sum / float64(count)
			result.ValueNum = &value
		}
	}
	return result
}
