package device

import (
	"context"
	"testing"

	"github.com/device-ms/client/device"
	"github.com/device-ms/itests"
	"github.com/device-ms/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_DeleteDevice(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail invalid id", func(t *testing.T) {
		dvIDStr := "12345"
		params := device.NewDeleteDeviceParams().WithID(dvIDStr)
		_, err := iti.ServiceClient.Device.DeleteDevice(params, iti.AuthDelegate)
		require.EqualError(t, err, "[DELETE /{id}][400] deleteDeviceBadRequest {\"code\":1500002,\"message\":\"parameter 'id' is invalid 'invalid object id ["+dvIDStr+"]'\"}")
	})

	t.Run("fail device not found", func(t *testing.T) {
		dvID := primitive.NewObjectID()
		params := device.NewDeleteDeviceParams().WithID(dvID.Hex())
		_, err := iti.ServiceClient.Device.DeleteDevice(params, iti.AuthDelegate)
		require.EqualError(t, err, "[DELETE /{id}][500] deleteDeviceInternalServerError {\"code\":1500005,\"message\":\"the device with id "+dvID.Hex()+" could not be found: mongo: no documents in result\"}")
	})

	t.Run("ok", func(t *testing.T) {
		dv := &model.Device{
			Brand: "brand3",
			Name:  "earth",
		}
		err := iti.DeviceRepository.Create(ctx, dv)
		require.NoError(t, err)

		params := device.NewDeleteDeviceParams().WithID(dv.ID.Hex())
		_, err = iti.ServiceClient.Device.DeleteDevice(params, iti.AuthDelegate)
		require.NoError(t, err)

		_, err = iti.DeviceRepository.ByID(ctx, dv.ID)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+dv.ID.Hex()+" could not be found")
	})
}
