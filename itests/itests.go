package itests

import (
	"context"
	"log"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/device-ms/client"
	"github.com/device-ms/client/device"
	"github.com/device-ms/controller"
	"github.com/device-ms/handler"
	"github.com/device-ms/mongo"
	"github.com/go-openapi/runtime"
	httptransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
	mongodriver "go.mongodb.org/mongo-driver/mongo"
)

type (
	// IntTestInfra is the infrastructure for integration tests
	IntTestInfra struct {
		DB               *mongodriver.Database
		DeviceRepository *mongo.DeviceRepository
		ServerAddress    string
		Router           handler.Router
		CloseServices    func()
		AuthDelegate     runtime.ClientAuthInfoWriter
		ServiceClient    *client.Swagger
		ValidVenueID     primitive.ObjectID
		ValidDeviceID    primitive.ObjectID
		ValidProfileID   primitive.ObjectID
		Controller       controller.Service
	}
)

var (
	// TestMutex is a mutex for tests
	TestMutex = sync.Mutex{}
)

// NewStr transforms a string into *string
func NewStr(s string) *string {
	return &s
}

// NewITests creates a new itests struct
func NewITests(ctx context.Context, t *testing.T) IntTestInfra {
	iti := IntTestInfra{}
	conn := `{"addresses":["localhost:27017"],"ssl":false,"directConnection":true}`
	err := os.Setenv("RS_DB_MONGO_CONN", conn)
	require.NoError(t, err)

	var drop func()
	iti.DeviceRepository, drop = mongo.CreateDeviceTestRepo(ctx, t)
	drop()

	iti.ValidVenueID = primitive.NewObjectID()
	iti.ValidDeviceID = primitive.NewObjectID()

	iti.Controller = controller.New(
		ctx,
		iti.DeviceRepository,
	)

	iti.Router = handler.NewDeviceRouter(iti.Controller)
	iti.CloseServices = func() {
	}

	return iti
}

// StartTestServer starts a test server
func (iti *IntTestInfra) StartTestServer(ctx context.Context, t *testing.T) (generatedClient *device.ClientService, closeServers func()) {
	TestMutex.Lock()
	log.Printf("starting device test service")
	server := httptest.NewServer(iti.Router)
	TestMutex.Unlock()

	iti.ServerAddress = strings.TrimPrefix(server.URL, "http://")

	unauthenticatedTransport := httptransport.New(iti.ServerAddress, handler.URLPath, []string{"http"})
	iti.ServiceClient = client.New(unauthenticatedTransport, strfmt.Default)

	closeServers = func() {
		server.Close()
		iti.CloseServices()
	}

	return
}
