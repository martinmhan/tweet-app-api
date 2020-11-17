build-proto-api-gateway:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cmd/api-gateway/proto/server.proto
build-api-gateway:
	make build-proto-api-gateway && go get github.com/martinmhan/tweet-app-api/cmd/api-gateway/internal
run-api-gateway:
	go run cmd/api-gateway/internal/main.go

build-proto-database-access:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cmd/database-access/proto/server.proto
build-database-access:
	make build-proto-database-access && go get github.com/martinmhan/tweet-app-api/cmd/database-access/internal
run-database-access:
	go run cmd/database-access/internal/main.go

build-proto-events-producer:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative cmd/events-producer/proto/server.proto
build-events-producer:
	make build-proto-events-producer && go get github.com/martinmhan/tweet-app-api/cmd/events-producer/internal
run-database-access:
	go run cmd/events-producer/internal/main.go