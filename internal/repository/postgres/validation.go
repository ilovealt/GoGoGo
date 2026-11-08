package postgres

import (
	"math"

	"gogogo/internal/model"
)

var validProtocols = map[string]bool{
	"opcua":      true,
	"modbus_tcp": true,
	"mqtt":       true,
}

func validateStatus(status uint8) error {
	if status != model.StatusDisabled && status != model.StatusEnabled {
		return invalid("status must be 0 or 1")
	}
	return nil
}

// normalizePoint 统一做 trim 和字段校验，和内存实现保持同一套规则。
func normalizePoint(point *model.MeasurementPoint) error {
	point.Code = trim(point.Code)
	point.Name = trim(point.Name)
	point.Description = trim(point.Description)
	point.PointType = trim(point.PointType)
	point.Unit = trim(point.Unit)

	if point.Code == "" || point.Name == "" || point.PointType == "" {
		return invalid("point code, name and point_type are required")
	}
	if point.ValueType < model.ValueTypeNumeric || point.ValueType > model.ValueTypeBool {
		return invalid("value_type must be 1, 2 or 3")
	}
	if point.Precision > 15 {
		return invalid("precision must be between 0 and 15")
	}
	if point.MinValue != nil && point.MaxValue != nil && *point.MinValue > *point.MaxValue {
		return invalid("min_value must be less than or equal to max_value")
	}
	if point.MinValue != nil && (math.IsNaN(*point.MinValue) || math.IsInf(*point.MinValue, 0)) {
		return invalid("min_value must be a finite number")
	}
	if point.MaxValue != nil && (math.IsNaN(*point.MaxValue) || math.IsInf(*point.MaxValue, 0)) {
		return invalid("max_value must be a finite number")
	}
	if math.IsNaN(point.Deadband) || math.IsInf(point.Deadband, 0) || point.Deadband < 0 {
		return invalid("deadband must be a non-negative finite number")
	}
	return validateStatus(point.Status)
}

func normalizeDataSource(source *model.DataSource) error {
	source.Name = trim(source.Name)
	source.Protocol = trim(source.Protocol)
	source.Endpoint = trim(source.Endpoint)

	if source.Name == "" || source.Endpoint == "" {
		return invalid("data source name and endpoint are required")
	}
	if !validProtocols[source.Protocol] {
		return invalid("protocol must be opcua, modbus_tcp or mqtt")
	}
	return validateStatus(source.Status)
}

// validateSample 校验单个样本与测点定义是否匹配。
func validateSample(sample *model.PointSample, point model.MeasurementPoint) error {
	if sample.Ts.IsZero() {
		return invalid("sample ts is required")
	}
	if sample.ValueType != point.ValueType {
		return invalid("sample value_type does not match point definition")
	}
	if sample.Quality > model.QualityOutOfRange {
		return invalid("quality must be between 0 and 5")
	}

	switch sample.ValueType {
	case model.ValueTypeNumeric:
		if sample.ValueNum == nil || sample.ValueText != nil || sample.ValueBool != nil {
			return invalid("numeric sample requires value_num only")
		}
		if math.IsNaN(*sample.ValueNum) || math.IsInf(*sample.ValueNum, 0) {
			return invalid("value_num must be a finite number")
		}
	case model.ValueTypeText:
		if sample.ValueText == nil || sample.ValueNum != nil || sample.ValueBool != nil {
			return invalid("text sample requires value_text only")
		}
	case model.ValueTypeBool:
		if sample.ValueBool == nil || sample.ValueNum != nil || sample.ValueText != nil {
			return invalid("bool sample requires value_bool only")
		}
	default:
		return invalid("unsupported value_type")
	}
	return nil
}
