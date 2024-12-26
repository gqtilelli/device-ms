package controller

import (
	"context"

	"github.com/device-ms/mongo"
)

// ServiceController is the service interface
type ServiceController interface {
	DeviceController() DeviceController
}

// Service represents the service with all controllers and clients inside
type Service struct {
	device DeviceController
}

// New returns a new service
func New(ctx context.Context, deviceDB mongo.DeviceDB) Service {
	return Service{
		device: NewDeviceService(deviceDB),
	}
}

// DeviceController returns the device controller.
func (s Service) DeviceController() DeviceController {
	return s.device
}
