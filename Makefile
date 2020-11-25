# The following scripts are used to build and run the services in this mono-repo
# The BIN param for each script specifies the service (e.g., `apigateway`, `databaseaccess`)
# Build scripts will run for all services

build-proto: # Example: `make build-proto BIN=apigateway`
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cmd/$(BIN)/proto/server.proto
build: # Example: `make build BIN=apigateway`
	go get github.com/martinmhan/tweet-app-api/cmd/$(BIN)/internal && make build-proto BIN=$(BIN)
run: # Example: `make run BIN=apigateway`
	go run cmd/$(BIN)/internal/main.go
