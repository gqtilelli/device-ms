package controller

import (
	"context"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/model"
	"github.com/device-ms/mongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodrv "go.mongodb.org/mongo-driver/mongo"
)

// DeviceController service
type DeviceController interface {
	Create(ctx context.Context, dv *model.Device) error
	GetDevice(ctx context.Context, deviceID primitive.ObjectID) (*model.Device, error)
	GetDevices(ctx context.Context) ([]dto.DeviceDTO, error)
	Update(ctx context.Context, dv *model.Device) error
	UpdateName(ctx context.Context, deviceID primitive.ObjectID, name string) error
	UpdateBrand(ctx context.Context, deviceID primitive.ObjectID, brand model.Brand) error
	Delete(ctx context.Context, deviceID primitive.ObjectID) error
	GetDevicesByBrand(ctx context.Context, brand model.Brand) ([]dto.DeviceDTO, error)
}

// DeviceService service
type DeviceService struct {
	deviceDB mongo.DeviceDB
}

// NewDeviceService DeviceService constructor
func NewDeviceService(deviceDB mongo.DeviceDB) DeviceController {
	return DeviceService{
		deviceDB: deviceDB,
	}
}

// Create creates a device in the database
func (dvs DeviceService) Create(ctx context.Context, device *model.Device) error {
	err := dvs.deviceDB.Create(ctx, device)
	if err != nil {
		return err
	}
	return nil
}

// Delete deletes a device from the database
func (dvs DeviceService) Delete(ctx context.Context, deviceID primitive.ObjectID) error {
	err := dvs.deviceDB.Delete(ctx, deviceID)
	if err != nil {
		return err
	}
	return nil
}

// Gets a device from the database by ID
func (dvs DeviceService) GetDevice(ctx context.Context, deviceID primitive.ObjectID) (*model.Device, error) {
	dv, err := dvs.deviceDB.ByID(ctx, deviceID)
	if err != nil {
		if err == mongodrv.ErrNoDocuments {
			return nil, errors.CouldNotFindObjectError("device", deviceID.Hex(), err)
		} else {
			return nil, err
		}
	}
	return dv, nil
}

// GetDevices gets all devices
func (dvs DeviceService) GetDevices(ctx context.Context) ([]dto.DeviceDTO, error) {
	models, err := dvs.deviceDB.List(ctx)
	if err != nil {
		return nil, err
	}
	dtos := make([]dto.DeviceDTO, len(models))
	for i := range models {
		dtos[i] = *dto.ToDeviceDTO(&models[i])
	}
	return dtos, nil
}

// GetDevicesByBrand gets all devices of a certain brand
func (dvs DeviceService) GetDevicesByBrand(ctx context.Context, brand model.Brand) ([]dto.DeviceDTO, error) {
	models, err := dvs.deviceDB.ListByBrand(ctx, brand)
	if err != nil {
		return nil, err
	}
	dtos := make([]dto.DeviceDTO, len(models))
	for i := range models {
		dtos[i] = *dto.ToDeviceDTO(&models[i])
	}
	return dtos, nil
}

// Update updates the information of a device, except infra fields like CreatedAt and UpdatedAt
func (dvs DeviceService) Update(ctx context.Context, dv *model.Device) error {
	_, err := dvs.deviceDB.Update(ctx, dv)
	if err != nil {
		return err
	}
	return nil
}

// UpdateBrand updates the brand of a device
func (dvs DeviceService) UpdateBrand(ctx context.Context, deviceID primitive.ObjectID, brand model.Brand) error {
	err := dvs.deviceDB.UpdateBrand(ctx, deviceID, brand)
	if err != nil {
		return err
	}
	return nil
}

// UpdateName updates the name of a device
func (dvs DeviceService) UpdateName(ctx context.Context, deviceID primitive.ObjectID, name string) error {
	err := dvs.deviceDB.UpdateName(ctx, deviceID, name)
	if err != nil {
		return err
	}
	return nil
}
