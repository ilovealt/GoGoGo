package model

const (
	RoleAdmin    = "admin"
	RoleOperator = "operator"
	RoleViewer   = "viewer"
)

const (
	StatusDisabled uint8 = 0
	StatusEnabled  uint8 = 1
)

const (
	ValueTypeNumeric uint8 = 1
	ValueTypeText    uint8 = 2
	ValueTypeBool    uint8 = 3
)

const (
	QualityGood         uint8 = 0
	QualityUncertain    uint8 = 1
	QualityBad          uint8 = 2
	QualityCommLost     uint8 = 3
	QualityNotConnected uint8 = 4
	QualityOutOfRange   uint8 = 5
)
