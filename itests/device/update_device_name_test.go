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

func Test_UpdateDeviceName(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail decode error", func(t *testing.T) {
		id := "12345"
		params := device.NewUpdateDeviceNameParams().WithID(id)
		_, err := iti.ServiceClient.Device.UpdateDeviceName(params)
		require.EqualError(t, err, "[PUT /{id}/name][400] updateDeviceNameBadRequest {\"code\":1500008,\"message\":\"decode error: EOF\"}")
	})

	t.Run("fail invalid id", func(t *testing.T) {
		id := "12345"
		params := device.NewUpdateDeviceNameParams().WithID(id).WithDeviceNameUpdate(&models.DeviceNameUpdateRequest{
			Name: "terra",
		})
		_, err := iti.ServiceClient.Device.UpdateDeviceName(params)
		require.EqualError(t, err, "[PUT /{id}/name][400] updateDeviceNameBadRequest {\"code\":1500002,\"message\":\"parameter 'id' is invalid 'invalid object id [12345]'\"}")
	})

	t.Run("device not found", func(t *testing.T) {
		id := primitive.NewObjectID()
		params := device.NewUpdateDeviceNameParams().WithID(id.Hex()).WithDeviceNameUpdate(&models.DeviceNameUpdateRequest{
			Name: "jupiter",
		})
		_, err := iti.ServiceClient.Device.UpdateDeviceName(params)
		require.EqualError(t, err, "[PUT /{id}/name][500] updateDeviceNameInternalServerError {\"code\":1500005,\"message\":\"the device with id "+id.Hex()+" could not be found: mongo: no documents in result\"}")
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

		params := device.NewUpdateDeviceNameParams().WithID(dv2.ID.Hex()).WithDeviceNameUpdate(&models.DeviceNameUpdateRequest{
			Name: "jupiter",
		})
		_, err = iti.ServiceClient.Device.UpdateDeviceName(params)
		require.NoError(t, err)

		dvBD1, err := iti.DeviceRepository.ByID(ctx, dv.ID)
		require.NoError(t, err)
		dvBD2, err := iti.DeviceRepository.ByID(ctx, dv2.ID)
		require.NoError(t, err)

		require.Equal(t, "terra", dvBD1.Name)
		require.Equal(t, model.Brand("brand3"), dvBD1.Brand)
		require.Equal(t, "jupiter", dvBD2.Name)
		require.Equal(t, model.Brand("brand1"), dvBD2.Brand)
	})
}
