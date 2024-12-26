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

func Test_CreateDevice(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail with missing required field brand", func(t *testing.T) {
		params := device.NewCreateDeviceParams().WithDeviceCreationRequestBody(&models.CreateDeviceRequest{})
		_, err := iti.ServiceClient.Device.CreateDevice(params, iti.AuthDelegate)
		require.EqualError(t, err, "[POST /][400] createDeviceBadRequest {\"code\":1500001,\"message\":\"parameter 'brand' in body is required\"}")
	})

	t.Run("fail with invalid brand", func(t *testing.T) {
		params := device.NewCreateDeviceParams().WithDeviceCreationRequestBody(&models.CreateDeviceRequest{
			Brand: "brand one",
		})
		_, err := iti.ServiceClient.Device.CreateDevice(params, iti.AuthDelegate)
		require.EqualError(t, err, "[POST /][400] createDeviceBadRequest {\"code\":1500002,\"message\":\"parameter 'brand' is invalid 'invalid value [brand one]'\"}")
	})

	t.Run("ok", func(t *testing.T) {
		params := device.NewCreateDeviceParams().WithDeviceCreationRequestBody(&models.CreateDeviceRequest{
			Brand: "brand3",
			Name:  "earth",
		})
		deviceCreated, err := iti.ServiceClient.Device.CreateDevice(params, iti.AuthDelegate)
		require.NoError(t, err)
		require.Equal(t, "earth", deviceCreated.Payload.Name)
		dvID, err := primitive.ObjectIDFromHex(deviceCreated.Payload.ID)
		require.NoError(t, err)

		dv, err := iti.DeviceRepository.ByID(ctx, dvID)
		require.NoError(t, err)
		require.Equal(t, "earth", dv.Name)
		require.Equal(t, model.Brand("brand3"), dv.Brand)
	})
}
