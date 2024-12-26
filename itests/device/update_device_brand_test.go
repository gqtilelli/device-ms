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

func Test_UpdateDeviceBrand(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail decode error", func(t *testing.T) {
		id := "12345"
		params := device.NewUpdateDeviceBrandParams().WithID(id)
		_, err := iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.EqualError(t, err, "[PUT /{id}/brand][400] updateDeviceBrandBadRequest {\"code\":1500008,\"message\":\"decode error: EOF\"}")
	})

	t.Run("fail invalid id", func(t *testing.T) {
		id := "12345"
		params := device.NewUpdateDeviceBrandParams().WithID(id).WithDeviceBrandUpdate(&models.DeviceBrandUpdateRequest{
			Brand: "brand3",
		})
		_, err := iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.EqualError(t, err, "[PUT /{id}/brand][400] updateDeviceBrandBadRequest {\"code\":1500002,\"message\":\"parameter 'id' is invalid 'invalid object id [12345]'\"}")
	})

	t.Run("fail no brand", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := device.NewUpdateDeviceBrandParams().WithID(id.Hex()).WithDeviceBrandUpdate(&models.DeviceBrandUpdateRequest{})
		_, err := iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.EqualError(t, err, "[PUT /{id}/brand][500] updateDeviceBrandInternalServerError {\"code\":1500002,\"message\":\"parameter 'brand' is invalid 'invalid value'\"}")
	})

	t.Run("fail invalid brand", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := device.NewUpdateDeviceBrandParams().WithID(id.Hex()).WithDeviceBrandUpdate(&models.DeviceBrandUpdateRequest{
			Brand: "brand two",
		})
		_, err := iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.EqualError(t, err, "[PUT /{id}/brand][500] updateDeviceBrandInternalServerError {\"code\":1500002,\"message\":\"parameter 'brand' is invalid 'invalid value'\"}")
	})

	t.Run("device not found", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := device.NewUpdateDeviceBrandParams().WithID(id.Hex()).WithDeviceBrandUpdate(&models.DeviceBrandUpdateRequest{
			Brand: "brand2",
		})
		_, err := iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.EqualError(t, err, "[PUT /{id}/brand][500] updateDeviceBrandInternalServerError {\"code\":1500005,\"message\":\"the device with id "+id.Hex()+" could not be found: mongo: no documents in result\"}")
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

		params := device.NewUpdateDeviceBrandParams().WithID(dv2.ID.Hex()).WithDeviceBrandUpdate(&models.DeviceBrandUpdateRequest{
			Brand: "brand3",
		})
		_, err = iti.ServiceClient.Device.UpdateDeviceBrand(params)
		require.NoError(t, err)

		dvBD1, err := iti.DeviceRepository.ByID(ctx, dv.ID)
		require.NoError(t, err)
		dvBD2, err := iti.DeviceRepository.ByID(ctx, dv2.ID)
		require.NoError(t, err)

		require.Equal(t, "terra", dvBD1.Name)
		require.Equal(t, model.Brand("brand3"), dvBD1.Brand)
		require.Equal(t, "venus", dvBD2.Name)
		require.Equal(t, model.Brand("brand3"), dvBD2.Brand)
	})
}
