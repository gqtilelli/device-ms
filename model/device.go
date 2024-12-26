package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// Device is the device information model
type Device struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	Name      string             `bson:"name"`
	Brand     Brand              `bson:"brand"`
	CreatedAt time.Time          `bson:"createdAt"`
	UpdatedAt *time.Time         `bson:"updatedAt,omitempty"`
}
