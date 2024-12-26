package dto

import (
	"time"

	"github.com/device-ms/model"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// DeviceDTO is a device DTO
type DeviceDTO struct {
	ID        string      `json:"id"`
	Name      string      `json:"name,omitempty"`
	Brand     model.Brand `json:"brand"`
	CreatedAt *time.Time  `json:"createdAt"`
}

// ToDeviceDTO maps a device model to a device dto response
func ToDeviceDTO(m *model.Device) *DeviceDTO {
	dto := DeviceDTO{
		ID:        m.ID.Hex(),
		Name:      m.Name,
		Brand:     m.Brand,
		CreatedAt: &m.CreatedAt,
	}

	return &dto
}

// UpdateDeviceRequestDTO request when updating a device
type UpdateDeviceRequestDTO struct {
	DeviceID primitive.ObjectID
	Name     string      `json:"name"`
	Brand    model.Brand `json:"brand"`
}

// ToModel maps a device update request dto to a device model
func (req UpdateDeviceRequestDTO) ToModel() (*model.Device, error) {
	return &model.Device{
		ID:    req.DeviceID,
		Name:  req.Name,
		Brand: req.Brand,
	}, nil
}

// UpdateDeviceNameRequestDTO request when updating the name of a device
type UpdateDeviceNameRequestDTO struct {
	DeviceID primitive.ObjectID
	Name     string `json:"name,omitempty"`
}

// ToModel maps a device update name request dto to a device model
func (req UpdateDeviceNameRequestDTO) ToModel() (*model.Device, error) {
	return &model.Device{
		ID:   req.DeviceID,
		Name: req.Name,
	}, nil
}

// UpdateDeviceBrandRequestDTO request when updating the brand of a device
type UpdateDeviceBrandRequestDTO struct {
	DeviceID primitive.ObjectID
	Brand    model.Brand `json:"brand,omitempty"`
}

// ToModel maps a device update name request dto to a device model
func (req UpdateDeviceBrandRequestDTO) ToModel() (*model.Device, error) {
	return &model.Device{
		ID:    req.DeviceID,
		Brand: req.Brand,
	}, nil
}

// CreateDeviceRequestDTO represents the body information to create a new device
type CreateDeviceRequestDTO struct {
	Name  string      `json:"name"`
	Brand model.Brand `json:"brand"`
}

// ToModel maps a device creation dto to a device model
func (req CreateDeviceRequestDTO) ToModel() *model.Device {
	return &model.Device{
		Name:  req.Name,
		Brand: req.Brand,
	}
}

// CreatedDeviceResponseDTO is a device DTO
type CreatedDeviceResponseDTO struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ToCreatedDeviceResponseDTO maps a device model to a created device response
func ToCreatedDeviceResponseDTO(m *model.Device) CreatedDeviceResponseDTO {
	return CreatedDeviceResponseDTO{
		ID:   m.ID.Hex(),
		Name: m.Name,
	}
}

// SearchDevicesRequestDTO is the request information used to search devices
type SearchDevicesRequestDTO struct {
	Brand model.Brand `json:"brand"`
}
