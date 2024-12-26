package mongo

import (
	"context"
	"time"

	"github.com/device-ms/errors"
	"github.com/device-ms/model"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DeviceCollectionName is the base name for the collection
const (
	DeviceCollectionName = "device"
)

// DeviceDB Device database
type DeviceDB interface {
	Create(ctx context.Context, device *model.Device) error
	ByID(ctx context.Context, id primitive.ObjectID) (*model.Device, error)
	List(ctx context.Context) ([]model.Device, error)
	Update(ctx context.Context, device *model.Device) (*model.Device, error)
	UpdateName(ctx context.Context, id primitive.ObjectID, name string) error
	UpdateBrand(ctx context.Context, id primitive.ObjectID, brand model.Brand) error
	Delete(ctx context.Context, id primitive.ObjectID) error
	ListByBrand(ctx context.Context, brand model.Brand) ([]model.Device, error)
}

// DeviceRepository  repository
type DeviceRepository struct {
	Collection *mongo.Collection
}

// NewDeviceDB creates new  collection
func NewDeviceDB(ctx context.Context, db *mongo.Database) (*DeviceRepository, error) {
	Collection := db.Collection(DeviceCollectionName, nil)

	indexes := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "brand", Value: 1}},
			Options: options.Index(),
		},
	}

	_, err := Collection.Indexes().CreateMany(ctx, indexes)
	if err != nil {
		return nil, err
	}

	return &DeviceRepository{
		Collection: Collection,
	}, nil
}

// Create saves new device to db
func (dr DeviceRepository) Create(ctx context.Context, device *model.Device) error {
	device.CreatedAt = time.Now().UTC().Truncate(time.Second)
	res, err := dr.Collection.InsertOne(ctx, device)
	if err != nil {
		return errors.CreateError(DeviceCollectionName, err.Error())
	}
	device.ID = res.InsertedID.(primitive.ObjectID)
	return nil
}

// ByID gets device by its id
func (dr DeviceRepository) ByID(ctx context.Context, id primitive.ObjectID) (*model.Device, error) {
	device := new(model.Device)
	err := dr.Collection.FindOne(ctx, bson.M{"_id": id}).Decode(device)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, errors.CouldNotFindObject(DeviceCollectionName, id.Hex())
		}
		return nil, errors.CouldNotFindObjectError(DeviceCollectionName, id.Hex(), err)
	}
	return device, nil
}

// List lists all device in db
func (dr DeviceRepository) List(ctx context.Context) ([]model.Device, error) {
	cur, err := dr.Collection.Find(ctx, bson.D{})
	if err != nil {
		return nil, errors.ListError(DeviceCollectionName, err, "ALL")
	}

	devices := make([]model.Device, 0)
	err = cur.All(ctx, &devices)
	if err != nil {
		return nil, errors.ListError(DeviceCollectionName, err, "ALL")
	}

	return devices, nil
}

// Update updates an existing Device in the database.
// Timestamp updated on success.
func (dr DeviceRepository) Update(ctx context.Context, device *model.Device) (*model.Device, error) {
	now := time.Now().UTC().Truncate(time.Second)

	result, err := dr.Collection.UpdateByID(ctx, device.ID, bson.M{
		"$set": bson.M{
			"updatedAt": &now,
			"name":      device.Name,
			"brand":     device.Brand,
		}})
	if err != nil {
		return nil, errors.UpdateError(DeviceCollectionName, err.Error())
	}
	if result.MatchedCount == 0 {
		return nil, errors.CouldNotFindObjectError(DeviceCollectionName, device.ID.Hex(), mongo.ErrNoDocuments)
	}

	return device, nil
}

// UpdateName updates a device name by its id
func (dr DeviceRepository) UpdateName(ctx context.Context, id primitive.ObjectID, name string) error {
	now := time.Now().UTC().Truncate(time.Second)
	result, err := dr.Collection.UpdateByID(ctx, id,
		bson.M{
			"$set": bson.M{
				"updatedAt": &now,
				"name":      name,
			}})
	if err != nil {
		return errors.UpdateError(DeviceCollectionName, err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.CouldNotFindObjectError(DeviceCollectionName, id.Hex(), mongo.ErrNoDocuments)
	}
	return nil
}

// UpdateBrand updates a device brand by its id
func (dr DeviceRepository) UpdateBrand(ctx context.Context, id primitive.ObjectID, brand model.Brand) error {
	now := time.Now().UTC().Truncate(time.Second)
	if !brand.IsValid() {
		return errors.InvalidParameterError("brand", "invalid value")
	}
	result, err := dr.Collection.UpdateByID(ctx, id,
		bson.M{
			"$set": bson.M{
				"updatedAt": &now,
				"brand":     brand,
			}})
	if err != nil {
		return errors.UpdateError(DeviceCollectionName, err.Error())
	}
	if result.MatchedCount == 0 {
		return errors.CouldNotFindObjectError(DeviceCollectionName, id.Hex(), mongo.ErrNoDocuments)
	}
	return nil
}

// Delete deletes device from database
func (dr DeviceRepository) Delete(ctx context.Context, id primitive.ObjectID) error {
	result, err := dr.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return errors.DeleteError(DeviceCollectionName, err.Error())
	}
	if result.DeletedCount == 0 {
		return errors.CouldNotFindObjectError(DeviceCollectionName, id.Hex(), mongo.ErrNoDocuments)
	}
	return nil
}

// ListByBrand gets device by brand
func (dr DeviceRepository) ListByBrand(ctx context.Context, brand model.Brand) ([]model.Device, error) {
	if !brand.IsValid() {
		return nil, errors.InvalidParameterError("brand", "invalid value")
	}
	cur, err := dr.Collection.Find(ctx, bson.M{"brand": brand})
	if err != nil {
		return nil, errors.ListError(DeviceCollectionName, err, "brand", string(brand))
	}

	devices := make([]model.Device, 0)
	err = cur.All(ctx, &devices)
	if err != nil {
		return nil, errors.ListError(DeviceCollectionName, err, "brand", string(brand))
	}

	return devices, nil
}
