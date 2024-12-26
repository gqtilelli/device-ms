package mongo

import (
	"context"
	"testing"
	"time"

	"github.com/device-ms/model"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func Test_Device_Create_ByID(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	device := model.Device{
		Name:  "netuno",
		Brand: "brand2",
	}

	t.Run("create and by id ok", func(t *testing.T) {
		err := repo.Create(ctx, &device)
		require.NoError(t, err)
		require.False(t, device.ID.IsZero())
		require.False(t, device.CreatedAt.IsZero())
		device, err := repo.ByID(ctx, device.ID)
		require.NoError(t, err)
		require.NotEmpty(t, device)
		require.Equal(t, "netuno", device.Name)
	})

	t.Run("by id - not found", func(t *testing.T) {
		dID := primitive.NewObjectID()
		_, err := repo.ByID(ctx, dID)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+dID.Hex()+" could not be found")
	})
}

func Test_DeviceList(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	devices := []model.Device{
		{
			Name:  "mercurio",
			Brand: "brand1",
		},
		{
			Name:  "marte",
			Brand: "brand2",
		},
		{
			Name:  "saturno",
			Brand: "brand3",
		},
		{
			Name:  "venus",
			Brand: "brand2",
		},
		{
			Name:  "plutao",
			Brand: "brand1",
		},
	}

	for i := range devices {
		_, err := repo.Collection.InsertOne(ctx, devices[i])
		require.NoError(t, err)
	}

	t.Run("list all", func(t *testing.T) {
		devices, err := repo.List(ctx)
		require.NoError(t, err)
		require.Len(t, devices, 5)
	})
}

func Test_DeviceUpdate(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	device := model.Device{
		Name:  "jupiter",
		Brand: "brand2",
	}
	require.NoError(t, repo.Create(ctx, &device))

	t.Run("success update", func(t *testing.T) {
		createdAt := time.Now().AddDate(0, 0, -1).Truncate(time.Second).UTC()
		device := model.Device{
			ID:        device.ID,
			Name:      "marte",
			Brand:     "brand1",
			CreatedAt: createdAt,
			UpdatedAt: &createdAt,
		}
		updatedDevice, err := repo.Update(ctx, &device)
		require.NoError(t, err)

		dv, err := repo.ByID(ctx, updatedDevice.ID)
		require.NoError(t, err)
		require.NotEqual(t, createdAt, device.UpdatedAt) // updatedAt is redefined by Update
		require.Equal(t, dv.ID, device.ID)
		require.NotEqual(t, dv.CreatedAt, device.CreatedAt) // createdAt is not updated
		require.Equal(t, dv.Name, device.Name)
		require.Equal(t, dv.Brand, device.Brand)
		require.NotEqual(t, dv.UpdatedAt, updatedDevice.UpdatedAt) // updatedAt is redefined by Update
	})
	t.Run("update failed", func(t *testing.T) {
		device.ID = primitive.NewObjectID()
		_, err := repo.Update(ctx, &device)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+device.ID.Hex()+" could not be found: mongo: no documents in result")
	})
}

func Test_DeviceUpdateName(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	device := model.Device{
		Name:  "jupiter",
		Brand: "brand2",
	}
	require.NoError(t, repo.Create(ctx, &device))

	t.Run("success update", func(t *testing.T) {
		err := repo.UpdateName(ctx, device.ID, "mercurio")
		require.NoError(t, err)

		dv, err := repo.ByID(ctx, device.ID)
		require.NoError(t, err)
		require.NotEqual(t, dv.UpdatedAt, device.UpdatedAt) // updatedAt is redefined by UpdateName
		require.Equal(t, "mercurio", dv.Name)
	})
	t.Run("update failed", func(t *testing.T) {
		device.ID = primitive.NewObjectID()
		err := repo.UpdateName(ctx, device.ID, "saturno")
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+device.ID.Hex()+" could not be found: mongo: no documents in result")
	})
}

func Test_DeviceUpdateBrand(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	device := model.Device{
		Name:  "jupiter",
		Brand: "brand2",
	}
	require.NoError(t, repo.Create(ctx, &device))

	t.Run("success update", func(t *testing.T) {
		err := repo.UpdateBrand(ctx, device.ID, "brand3")
		require.NoError(t, err)

		dv, err := repo.ByID(ctx, device.ID)
		require.NoError(t, err)
		require.NotEqual(t, dv.UpdatedAt, device.UpdatedAt) // updatedAt is redefined by UpdateName
		require.Equal(t, model.Brand("brand3"), dv.Brand)
	})
	t.Run("invalid brand", func(t *testing.T) {
		err := repo.UpdateBrand(ctx, device.ID, "new brand")
		require.EqualError(t, err, "result: false; code: 1500002; message: parameter 'brand' is invalid 'invalid value'")
	})
	t.Run("update failed", func(t *testing.T) {
		device.ID = primitive.NewObjectID()
		err := repo.UpdateBrand(ctx, device.ID, "brand1")
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+device.ID.Hex()+" could not be found: mongo: no documents in result")
	})
}

func Test_DeviceDelete(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	device := model.Device{
		Name:  "jupiter",
		Brand: "brand2",
	}
	require.NoError(t, repo.Create(ctx, &device))

	t.Run("failed and success delete", func(t *testing.T) {
		id := primitive.NewObjectID()
		err := repo.Delete(ctx, id)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+id.Hex()+" could not be found: mongo: no documents in result")

		id = device.ID
		err = repo.Delete(ctx, id)
		require.NoError(t, err)

		_, err = repo.ByID(ctx, id)
		require.EqualError(t, err, "result: false; code: 1500005; message: the device with id "+id.Hex()+" could not be found")
	})
}

func Test_DeviceSearchByBrand(t *testing.T) {
	ctx := context.Background()
	repo, drop := NewTestDeviceRepo(t)
	defer drop()

	devices := []model.Device{
		{
			Name:  "mercurio",
			Brand: "brand1",
		},
		{
			Name:  "marte",
			Brand: "brand2",
		},
		{
			Name:  "saturno",
			Brand: "brand2",
		},
		{
			Name:  "venus",
			Brand: "brand2",
		},
		{
			Name:  "plutao",
			Brand: "brand1",
		},
	}

	for i := range devices {
		_, err := repo.Collection.InsertOne(ctx, devices[i])
		require.NoError(t, err)
	}

	t.Run("invalid brand", func(t *testing.T) {
		_, err := repo.ListByBrand(ctx, "new brand")
		require.EqualError(t, err, "result: false; code: 1500002; message: parameter 'brand' is invalid 'invalid value'")
	})
	t.Run("empty", func(t *testing.T) {
		devices, err := repo.ListByBrand(ctx, "brand3")
		require.NoError(t, err)
		require.Len(t, devices, 0)
	})
	t.Run("success", func(t *testing.T) {
		devices, err := repo.ListByBrand(ctx, "brand2")
		require.NoError(t, err)
		require.Len(t, devices, 3)
	})
}

func NewTestDeviceRepo(t *testing.T) (repo *DeviceRepository, drop func()) {
	ctx := context.Background()
	db := createTestDB(ctx, t)

	repo, err := NewDeviceDB(ctx, db)
	require.NoError(t, err)

	return repo, func() {
		_, err := repo.Collection.DeleteMany(ctx, bson.D{})
		require.NoError(t, err)
	}
}
