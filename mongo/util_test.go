package mongo

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_createDB(t *testing.T) {
	t.Run("mongo env var not defined", func(t *testing.T) {
		ctx := context.Background()

		mutex.Lock()
		dbSave := db
		db = nil
		saveMongoURI := os.Getenv(envMongoURI)

		os.Unsetenv(envMongoURI)
		_, err := createDB(ctx)
		require.EqualError(t, err, "Mongo URI (MONGO_URI) not defined")

		db = dbSave
		os.Setenv(envMongoURI, saveMongoURI)
		mutex.Unlock()
	})
}

func Test_CreateDeviceRepo(t *testing.T) {
	t.Run("invalid mongo env var", func(t *testing.T) {
		ctx := context.Background()

		mutex.Lock()
		dbSave := db
		db = nil
		saveMongoURI := os.Getenv(envMongoURI)
		os.Setenv(envMongoURI, "ftp://0.0.0.0")

		_, err := CreateDeviceRepo(ctx)
		require.EqualError(t, err, "error parsing uri: scheme must be \"mongodb\" or \"mongodb+srv\"")

		db = dbSave
		os.Setenv(envMongoURI, saveMongoURI)
		mutex.Unlock()
	})
}
