package model

// Brand enum
type Brand string

// Enum values
const (
	Bbrand1 Brand = "brand1"
	Bbrand2 Brand = "brand2"
	Bbrand3 Brand = "brand3"
)

var mapDeviceBrand = map[Brand]bool{
	Bbrand1: true,
	Bbrand2: true,
	Bbrand3: true,
}

// IsValid is valid enum value
func (brand Brand) IsValid() bool {
	return mapDeviceBrand[brand]
}
