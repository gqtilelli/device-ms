package device

import (
	"context"
	"testing"

	"github.com/device-ms/client/device"
	"github.com/device-ms/itests"
	"github.com/device-ms/model"
	"github.com/device-ms/models"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_UpdateDevice(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail invalid id", func(t *testing.T) {
		id := "12345"
		dvUpd := models.UpdateDeviceRequest{
			Name:  "earth",
			Brand: "brand3",
		}
		params := device.NewUpdateDeviceParams().WithID(id).WithDeviceUpdateRequestBody(&dvUpd)
		_, err := iti.ServiceClient.Device.UpdateDevice(params)
		require.EqualError(t, err, "[PUT /{id}][400] updateDeviceBadRequest {\"code\":1500002,\"message\":\"parameter 'id' is invalid 'invalid object id [12345]'\"}")
	})

	t.Run("fail no brand", func(t *testing.T) {
		id := primitive.NewObjectID()
		dvUpd := models.UpdateDeviceRequest{
			Name: "earth",
		}
		params := device.NewUpdateDeviceParams().WithID(id.Hex()).WithDeviceUpdateRequestBody(&dvUpd)
		_, err := iti.ServiceClient.Device.UpdateDevice(params)
		require.EqualError(t, err, "[PUT /{id}][400] updateDeviceBadRequest {\"code\":1500001,\"message\":\"parameter 'brand' in body is required\"}")
	})

	t.Run("fail invalid brand", func(t *testing.T) {
		id := primitive.NewObjectID()
		dvUpd := models.UpdateDeviceRequest{
			Name:  "earth",
			Brand: "brand two",
		}
		params := device.NewUpdateDeviceParams().WithID(id.Hex()).WithDeviceUpdateRequestBody(&dvUpd)
		_, err := iti.ServiceClient.Device.UpdateDevice(params)
		require.EqualError(t, err, "[PUT /{id}][400] updateDeviceBadRequest {\"code\":1500002,\"message\":\"parameter 'brand' is invalid 'invalid value [brand two]'\"}")
	})

	t.Run("device not found", func(t *testing.T) {
		id := primitive.NewObjectID()
		dvUpd := models.UpdateDeviceRequest{
			Name:  "earth",
			Brand: "brand2",
		}
		params := device.NewUpdateDeviceParams().WithID(id.Hex()).WithDeviceUpdateRequestBody(&dvUpd)
		_, err := iti.ServiceClient.Device.UpdateDevice(params)
		require.EqualError(t, err, "[PUT /{id}][500] updateDeviceInternalServerError {\"code\":1500005,\"message\":\"the device with id "+id.Hex()+" could not be found: mongo: no documents in result\"}")
	})

	t.Run("ok", func(t *testing.T) {
		dv := &model.Device{
			Brand: "brand3",
			Name:  "terra",
		}
		err := iti.DeviceRepository.Create(ctx, dv)
		require.NoError(t, err)
		dv2 := &model.Device{
			Brand: "brand1",
			Name:  "venus",
		}
		err = iti.DeviceRepository.Create(ctx, dv2)
		require.NoError(t, err)

		params := device.NewUpdateDeviceParams().WithID(dv2.ID.Hex()).WithDeviceUpdateRequestBody(&models.UpdateDeviceRequest{
			Name:  "marte",
			Brand: "brand3",
		})
		_, err = iti.ServiceClient.Device.UpdateDevice(params)
		require.NoError(t, err)

		dvBD1, err := iti.DeviceRepository.ByID(ctx, dv.ID)
		require.NoError(t, err)
		dvBD2, err := iti.DeviceRepository.ByID(ctx, dv2.ID)
		require.NoError(t, err)

		require.Equal(t, "terra", dvBD1.Name)
		require.Equal(t, model.Brand("brand3"), dvBD1.Brand)
		require.Equal(t, "marte", dvBD2.Name)
		require.Equal(t, model.Brand("brand3"), dvBD2.Brand)
	})
}
