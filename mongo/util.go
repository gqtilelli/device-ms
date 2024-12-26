package mongo

import (
	"context"
	"fmt"
	"os"
	"sync"
	"testing"
	"time"

	nrmongo "github.com/newrelic/go-agent/v3/integrations/nrmongo"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	databaseName = "device"
	envMongoURI  = "MONGO_URI"
	mongoTimeout = 20
)

var (
	db    *mongo.Database
	mutex sync.Mutex
)

func createTestDB(ctx context.Context, t *testing.T) *mongo.Database {
	var err error
	mutex.Lock()
	defer mutex.Unlock()
	if db == nil {
		db, err = initDB(ctx, "device-test", "mongodb://localhost:27017/")
		require.NoError(t, err)
	} else {
		db.Client().Database("device-test")
	}
	return db
}

func createDB(ctx context.Context) (*mongo.Database, error) {
	var err error
	if db == nil {
		mongoURI := os.Getenv(envMongoURI)
		if mongoURI == "" {
			return nil, fmt.Errorf("Mongo URI (%s) not defined", envMongoURI)
		}
		db, err = initDB(ctx, databaseName, mongoURI)
	}
	return db, err
}

// CreateDeviceRepo creates a device repository
func CreateDeviceRepo(ctx context.Context) (*DeviceRepository, error) {
	db, err := createDB(ctx)
	if err != nil {
		return nil, err
	}
	return NewDeviceDB(ctx, db)
}

// CreatDeviceTestRepo creates a device test repository
func CreateDeviceTestRepo(ctx context.Context, t *testing.T) (repo *DeviceRepository, drop func()) {
	db = createTestDB(ctx, t)

	repo, err := NewDeviceDB(ctx, db)
	require.NoError(t, err)

	drop = func() {
		_, err := repo.Collection.DeleteMany(ctx, bson.D{})
		require.NoError(t, err)
	}
	return
}

func initDB(ctx context.Context, name, mongoURI string) (*mongo.Database, error) {
	nrMon := nrmongo.NewCommandMonitor(nil)
	opts := options.Client().ApplyURI(mongoURI).SetAppName(name)
	opts = opts.SetMonitor(nrMon)
	opts = opts.SetConnectTimeout(time.Minute)
	err := opts.Validate()
	if err != nil {
		return nil, err
	}
	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, err
	}
	err = checkConnection(ctx, client)
	if err != nil {
		return nil, err
	}
	return client.Database(name), nil
}

func checkConnection(ctx context.Context, client *mongo.Client) error {
	ctx, cancel := context.WithTimeout(ctx, mongoTimeout*time.Second)
	defer cancel()

	return client.Ping(ctx, readpref.Primary())
}
