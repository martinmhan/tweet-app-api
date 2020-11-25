# The following scripts are used to build and run the services in this mono-repo
# The BIN parameter is needed for each script and specifies the service (spelled exactly like its cmd/ subdirectory, e.g., `apigateway`, `databaseaccess`)
# Scripts will error with an invalid or no BIN parameter
# The build-proto script first checks if the provided BIN parameter has a /proto subdirectory, since the eventconsumer service does not have one

build-proto: # Example: `make build-proto BIN=apigateway`
	if [ -d "cmd/$(BIN)/proto" ]; then \
		protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cmd/$(BIN)/proto/server.proto; \
	fi
build: # Example: `make build BIN=apigateway`
	go get github.com/martinmhan/tweet-app-api/cmd/$(BIN)/internal && make build-proto BIN=$(BIN)
run: # Example: `make run BIN=apigateway`
	go run cmd/$(BIN)/internal/main.go
