package controller

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/device-ms/dto"
	"github.com/device-ms/errors"
	"github.com/device-ms/model"
	mongoMocks "github.com/device-ms/mongo/mocks"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func TestDeviceController_GetDevice(t *testing.T) {
	ctx := context.Background()

	deviceDB := new(mongoMocks.DeviceDB)
	defer deviceDB.AssertExpectations(t)

	t.Run("fail with error by id", func(t *testing.T) {
		deviceID := primitive.NewObjectID()
		deviceDB.On("ByID", ctx, deviceID).Return(nil, errors.CouldNotFindObject("device", deviceID.Hex())).Once()

		deviceController := NewDeviceService(deviceDB)
		resp, err := deviceController.GetDevice(ctx, deviceID)
		require.Nil(t, resp)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+deviceID.Hex()+" could not be found")
	})

	t.Run("ok - by id", func(t *testing.T) {
		device := model.Device{
			ID:    primitive.NewObjectID(),
			Name:  "Marte",
			Brand: "brand2",
		}

		deviceDB.On("ByID", ctx, device.ID).Return(&device, nil).Once()
		deviceController := NewDeviceService(deviceDB)
		resp, err := deviceController.GetDevice(ctx, device.ID)
		require.NoError(t, err)
		require.NotEmpty(t, resp)
		require.Equal(t, device, *resp)
	})
}

func Test_DeviceController(t *testing.T) {
	errMock := fmt.Errorf("errMock")

	ctx := context.Background()

	deviceDB := new(mongoMocks.DeviceDB)
	defer deviceDB.AssertExpectations(t)

	device := model.Device{
		ID:    primitive.NewObjectID(),
		Name:  "saturno",
		Brand: "brand1",
	}

	t.Run("fail on db save", func(t *testing.T) {
		deviceDB.On("Create", mock.Anything, &device).Return(errors.CreateError("device", "unexpected error")).Once()

		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Create(ctx, &device)
		require.EqualError(t, err, "result: false; code: 1500003; message: error creating device reason unexpected error")
	})

	t.Run("ok - device created", func(t *testing.T) {
		deviceDB.On("Create", mock.Anything, &device).Return(nil).Run(func(args mock.Arguments) {
			device := args[1].(*model.Device)
			device.CreatedAt = time.Now().UTC().Truncate(time.Second)
		}).Once()

		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Create(ctx, &device)
		require.NoError(t, err)
	})

	t.Run("ok - update", func(t *testing.T) {
		deviceDB.On("Update", mock.Anything, &device).Return(&device, nil).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Update(ctx, &device)
		require.NoError(t, err)
	})
	t.Run("update failed", func(t *testing.T) {
		deviceDB.On("Update", mock.Anything, &device).Return(nil, errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Update(ctx, &device)
		require.EqualError(t, err, errMock.Error())
	})

	t.Run("ok - update name", func(t *testing.T) {
		deviceDB.On("UpdateName", mock.Anything, device.ID, "plutao").Return(nil).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.UpdateName(ctx, device.ID, "plutao")
		require.NoError(t, err)
	})
	t.Run("update name failed", func(t *testing.T) {
		deviceDB.On("UpdateName", mock.Anything, device.ID, "plutao").Return(errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.UpdateName(ctx, device.ID, "plutao")
		require.EqualError(t, err, errMock.Error())
	})

	t.Run("ok - update brand", func(t *testing.T) {
		deviceDB.On("UpdateBrand", mock.Anything, device.ID, model.Brand("brand1")).Return(nil).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.UpdateBrand(ctx, device.ID, model.Brand("brand1"))
		require.NoError(t, err)
	})
	t.Run("update brand failed", func(t *testing.T) {
		deviceDB.On("UpdateBrand", mock.Anything, device.ID, model.Brand("brand1")).Return(errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.UpdateBrand(ctx, device.ID, model.Brand("brand1"))
		require.EqualError(t, err, errMock.Error())
	})

	t.Run("ok - list devices", func(t *testing.T) {
		deviceDB.On("List", mock.Anything).Return([]model.Device{{
			ID:    device.ID,
			Name:  "venus",
			Brand: model.Brand("brand2"),
		}}, nil).Once()
		deviceController := NewDeviceService(deviceDB)
		devices, err := deviceController.GetDevices(ctx)
		require.NoError(t, err)
		require.Len(t, devices, 1)
		require.Equal(t, devices[0].Name, "venus")
		require.Equal(t, devices[0].Brand, model.Brand("brand2"))
	})
	t.Run("failed listing devices", func(t *testing.T) {
		deviceDB.On("List", mock.Anything).Return([]model.Device(nil), errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		devices, err := deviceController.GetDevices(ctx)
		require.EqualError(t, err, errMock.Error())
		require.Equal(t, []dto.DeviceDTO(nil), devices)
	})

	t.Run("ok - list devices by brand", func(t *testing.T) {
		deviceDB.On("ListByBrand", ctx, device.Brand).Return([]model.Device{{
			ID:    device.ID,
			Name:  "saturno",
			Brand: model.Brand("brand1"),
		}}, nil).Once()
		deviceController := NewDeviceService(deviceDB)
		devices, err := deviceController.GetDevicesByBrand(ctx, device.Brand)
		require.NoError(t, err)
		require.Len(t, devices, 1)
		require.Equal(t, devices[0].Name, "saturno")
		require.Equal(t, devices[0].Brand, model.Brand("brand1"))
	})
	t.Run("failed listing devices by brand", func(t *testing.T) {
		deviceDB.On("ListByBrand", ctx, device.Brand).Return([]model.Device(nil), errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		devices, err := deviceController.GetDevicesByBrand(ctx, device.Brand)
		require.EqualError(t, err, errMock.Error())
		require.Equal(t, []dto.DeviceDTO(nil), devices)
	})

	t.Run("ok - get device by id", func(t *testing.T) {
		deviceDB.On("ByID", mock.Anything, device.ID).Return(&device, nil).Once()
		deviceController := NewDeviceService(deviceDB)
		device, err := deviceController.GetDevice(ctx, device.ID)
		require.NoError(t, err)
		require.NotNil(t, device)
		require.Equal(t, device.Name, "saturno")
		require.Equal(t, device.Brand, model.Brand("brand1"))
	})
	t.Run("failed getting device by id", func(t *testing.T) {
		deviceDB.On("ByID", mock.Anything, device.ID).Return((*model.Device)(nil), errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		device, err := deviceController.GetDevice(ctx, device.ID)
		require.EqualError(t, err, errMock.Error())
		require.Equal(t, (*model.Device)(nil), device)
	})

	t.Run("ok - delete", func(t *testing.T) {
		deviceDB.On("Delete", mock.Anything, device.ID).Return(nil).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Delete(ctx, device.ID)
		require.NoError(t, err)
	})
	t.Run("delete failed", func(t *testing.T) {
		deviceDB.On("Delete", mock.Anything, device.ID).Return(errMock).Once()
		deviceController := NewDeviceService(deviceDB)
		err := deviceController.Delete(ctx, device.ID)
		require.EqualError(t, err, errMock.Error())
	})
}
