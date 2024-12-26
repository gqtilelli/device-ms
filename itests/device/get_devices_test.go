package device

import (
	"context"
	"testing"

	"github.com/device-ms/client/device"
	"github.com/device-ms/itests"
	"github.com/device-ms/model"
	"github.com/stretchr/testify/require"
)

func Test_GetDevices(t *testing.T) {
	ctx := context.Background()
	iti := itests.NewITests(ctx, t)
	_, closeServer := iti.StartTestServer(ctx, t)
	defer closeServer()

	t.Run("fail invalid brand", func(t *testing.T) {
		brand := "brand one"
		params := device.NewGetDevicesParams().WithBrand(&brand)
		_, err := iti.ServiceClient.Device.GetDevices(params)
		require.EqualError(t, err, "[GET /][400] getDevicesBadRequest {\"code\":1500002,\"message\":\"parameter 'brand' is invalid 'invalid value [brand one]'\"}")
	})

	t.Run("no device found with brand", func(t *testing.T) {
		brand := "brand2"
		params := device.NewGetDevicesParams().WithBrand(&brand)
		dvs, err := iti.ServiceClient.Device.GetDevices(params)
		require.NoError(t, err)
		require.Len(t, dvs.Payload, 0)
	})

	t.Run("no device found at all", func(t *testing.T) {
		params := device.NewGetDevicesParams()
		dvs, err := iti.ServiceClient.Device.GetDevices(params)
		require.NoError(t, err)
		require.Len(t, dvs.Payload, 0)
	})

	t.Run("ok", func(t *testing.T) {
		dv := &model.Device{
			Brand: "brand3",
			Name:  "earth",
		}
		err := iti.DeviceRepository.Create(ctx, dv)
		require.NoError(t, err)

		params := device.NewGetDevicesParams()
		dvs, err := iti.ServiceClient.Device.GetDevices(params)
		require.NoError(t, err)
		require.Len(t, dvs.Payload, 1)
		require.Equal(t, "earth", dvs.Payload[0].Name)
		require.Equal(t, "brand3", dvs.Payload[0].Brand)

		dv2 := &model.Device{
			Brand: "brand1",
			Name:  "venus",
		}
		err = iti.DeviceRepository.Create(ctx, dv2)
		require.NoError(t, err)

		params2 := device.NewGetDevicesParams()
		dvs2, err := iti.ServiceClient.Device.GetDevices(params2)
		require.NoError(t, err)
		require.Len(t, dvs2.Payload, 2)
		require.Equal(t, "earth", dvs2.Payload[0].Name)
		require.Equal(t, "brand3", dvs2.Payload[0].Brand)
		require.Equal(t, "venus", dvs2.Payload[1].Name)
		require.Equal(t, "brand1", dvs2.Payload[1].Brand)

		brand3 := "brand3"
		params3 := device.NewGetDevicesParams().WithBrand(&brand3)
		dvs3, err := iti.ServiceClient.Device.GetDevices(params3)
		require.NoError(t, err)
		require.Len(t, dvs3.Payload, 1)
		require.Equal(t, "earth", dvs3.Payload[0].Name)
		require.Equal(t, "brand3", dvs3.Payload[0].Brand)

		brand2 := "brand2"
		params4 := device.NewGetDevicesParams().WithBrand(&brand2)
		dvs4, err := iti.ServiceClient.Device.GetDevices(params4)
		require.NoError(t, err)
		require.Len(t, dvs4.Payload, 0)
	})
}
