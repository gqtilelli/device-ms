swagger-test: swagger.yml
	CGO_ENABLED=0 swagger generate client -A swagger-test -f ./swagger.yml

mock-test: controller/mocks/*.go mongo/mocks/*.go

controller/mocks/*.go: controller/*.go
	mockery --all --dir ./controller/ --output ./controller/mocks --case underscore --disable-version-string --exported 

mongo/mocks/*.go: mongo/*.go
	mockery --all --dir ./mongo/ --output ./mongo/mocks --case underscore --disable-version-string --exported

test: swagger-test mock-test
	go test -cover ./controller ./mongo ./model ./itests/device

testclean:
	go clean -testcache

runmongo:
	docker run --name mongodb -d -p 27017:27017 mongo

rundevice:
	MONGO_URI="mongodb://localhost:27017/" go run main.go

.PHONY: mongo testclean
