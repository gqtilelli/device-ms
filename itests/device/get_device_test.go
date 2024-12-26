package device

import (
	"context"
	"testing"
	"time"

	"github.com/device-ms/client/device"
	"github.com/device-ms/itests"
	"github.com/device-ms/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_GetDevice(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail invalid id", func(t *testing.T) {
		dvIDStr := "12345"
		params := device.NewGetDeviceParams().WithID(dvIDStr)
		_, err := iti.ServiceClient.Device.GetDevice(params)
		require.EqualError(t, err, "[GET /{id}][400] getDeviceBadRequest {\"code\":1500002,\"message\":\"parameter 'id' is invalid 'invalid object id [12345]'\"}")
	})

	t.Run("fail device not found", func(t *testing.T) {
		dvID := primitive.NewObjectID()
		params := device.NewGetDeviceParams().WithID(dvID.Hex())
		_, err := iti.ServiceClient.Device.GetDevice(params)
		require.EqualError(t, err, "[GET /{id}][500] getDeviceInternalServerError {\"code\":1500005,\"message\":\"the device with id "+dvID.Hex()+" could not be found\"}")
	})

	t.Run("ok", func(t *testing.T) {
		dv := &model.Device{
			Brand: "brand3",
			Name:  "earth",
		}
		now := time.Now().UTC().Truncate(time.Second)
		err := iti.DeviceRepository.Create(ctx, dv)
		require.NoError(t, err)

		params := device.NewGetDeviceParams().WithID(dv.ID.Hex())
		res, err := iti.ServiceClient.Device.GetDevice(params)
		require.NoError(t, err)

		require.Equal(t, "earth", res.Payload.Name)
		require.Equal(t, "brand3", res.Payload.Brand)
		nowStr := now.Format(time.RFC3339)[:19]
		createdAtStr := res.Payload.CreatedAt.String()[:19]
		require.Equal(t, nowStr, createdAtStr)
	})
}
