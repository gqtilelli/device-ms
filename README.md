device-ms

The purpose of this microservice is to keep the record of devices.
The device entity is composed by:
o Device name
o Device brand
o Creation time
The supported operations are:
1. Add device;
2. Get device by identifier;
3. List all devices;
4. Update device (full and partial);
5. Delete a device;
6. Search device by brand;
The file swagger.yml contains the Restful API definition.

Database
The tests need that a MongoDB database is running in the local machine.
To run a MongoDB server using Docker, use the command:
> make runmongo
A docker manager, for example Docker Desktop, has to be installed and running in local machine.

go-swagger
go-swagger is used to generate code for integration tests (itests)
go-swagger executable "generate-swagger-client" is installed with the following command:
> go install github.com/go-swagger/go-swagger/cmd/swagger@latest
And its directory has to be added the PATH. For example, in a Mac, add the following line to the end of login script (~/.bashrc or ~/.zshrc, for example):
export PATH=$PATH:~/go/bin

mockery
mockery is used to generate mocks for interfaces for use in tests.
To install mockery, the following command is used:
go get github.com/vektra/mockery/v2/../
In a Mac, It can also be installed using brew:
brew install mockery

Automated Tests
To execute the tests, use the command: make test
A Mongodb has to be available in default port 27017 of local machine. See Database above.
There are unit tests for database operations (mongo), model and controller.
The handler is covered by the integration tests (itests).
Example:
make test
.
. (test code automatic generation from swagger (go-swagger) and mockery (mocks for interfaces) 
.
go test -cover ./controller ./mongo ./model ./itests/device
ok  	github.com/device-ms/controller	0.805s	coverage: 93.0% of statements
ok  	github.com/device-ms/mongo	1.596s	coverage: 80.7% of statements
ok  	github.com/device-ms/model	2.069s	coverage: 100.0% of statements
ok  	github.com/device-ms/itests/device	2.845s	coverage: [no statements]

Run device-ms
To run device-ms in local machine (with mongo runnning in local machine too - see Database above), run:
> make rundevice
If successful, the prompt won't return and the tests can be executed for example with curl from another terminal:
~ curl --location --request GET 'http://localhost:8080/device' 
[]
~ curl --location --request POST 'http://localhost:8080/device' \
--data-raw '{
    "name": "jupiter",
    "brand": "brand1" 
}'
{"id":"676b240a7bbab556f4a6b57b","name":"jupiter"}
˜ curl --location --request GET 'http://localhost:8080/device'   
[{"id":"676b240a7bbab556f4a6b57b","name":"jupiter","brand":"brand1","createdAt":"2024-12-24T21:13:46Z"}]
The accepted brands are defined by source code (device-ms/model/enums.go) and the accepted values are:
brand1
brand2
brand3
Type Ctrl-c to stop device-ms server and return to the prompt

Docker Run
It is possible also to run device-ms server using docker.
For this, create device-ms docker image with the following command:
> make builddockerdevice
Then, create a docker user network with the command:
> make devnet
Run a mongodb server in devnet with the following command:
> make rundevnetmongo
And run device-ms in devnet with the command:
> make rundevnetdevice

