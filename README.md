# Tweet App API
- Event-Driven Microservices API for a simple tweeting app
- Built with Go, gRPC, RabbitMQ, MongoDB
- UI client TBD

# Pre-Requisites:
  TBD

# Getting Started:
  - Run MongoDB and RabbitMQ daemons:
    - In MacOS, something like `brew services start mongodb` and `brew services start rabbitmq`
  - Create MongoDB database:
    - In a terminal window, navigate to the /scripts/db directory and run `upgrade.sh`
  - Build and run services:
    - To run services locally, open a terminal window and run the `build-and-run` Makefile script for each service (e.g., `make build-and-run BIN=eventproducer`)
    - Services must be started in the following order to avoid fatal connection errors: 1) databaseaccess, 2) readview, 3) eventconsumer, 4) eventproducer, 5) apigateway
  - Ping the API gateway (via an RPC client tool such as BloomRPC) to create a user, log in, post a tweet, etc.

# Features:
  - [Event Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture) (EDA)
    - Message-Queue used to produce/consume state-changing events (in this case, database writes)
  - [Command Query Responsibility Segregation](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs) (CQRS):
    - This API handles read and write requests separately
    - Reads are done via a Read View service:
      - This is an-memory data view used to optimize read speeds and to prevent blocking of writes
      - Only queries the DB on cold starts
      - Gets updated by the consumer after each successful write
    - Writes are done via the queue
      - The Events Producer publishes a write event to the message queue
      - The Events Consumer subscribes to and processes the write event
  - [Remote Procedure Call](https://en.wikipedia.org/wiki/Remote_procedure_call) (RPC):
    - gRPC used for direct communication with the UI and between services
      - The only exception being the Events Producer and Events Consumer, which communicate via the Message Queue
      - This allows for flexibility to handle UI requests either synchronously or asynchronously
        - Reads are completed and provided synchronously in the gRPC response
        - Writes are completed asynchronously via the message queue
  <!-- - [Domain Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design) (DDD): ### TBD ###
  - This API follows DDD principles in two levels:
    - The overall architecture utilizes microservices to separate contexts and functionality (see API architecture diagram below)
    - Each microservice uses folder structure to separate abstraction layers (see `/cmd/internal/` Project Structure notes below) -->

# Notes:
  - Project Structure
    - `/cmd`: contains subdirectories, each containing the following code for one microservice:
      - `/internal`: code only used by the microservice (i.e., within the same `/cmd/<MICROSERVICE>` directory). Includes at least one of the following subdirectories
        - `/application`: gRPC server implementation defining RPC handlers
        - `/domain`: business logic used by the RPC handlers
        - `/infrastructure`: logic to handle data persistence
        - `main.go`: root file used to run this service
      - `/proto`: contains protocol buffer definitions
        - .proto files are the source files
        - .pb.go files are generated during the build
    - `.env`: environment variables
    - `Makefile`: scripts used to build and run services
  - Services Overview
    TBD
  - RabbitMQ has more robust functionality such as routing messages to queues via exchanges. However, just one queue and no exchange was used here due to this API's simplicity.

# API Architecture:
![API Architecture](https://gitbuckets.s3-us-west-1.amazonaws.com/tweet-app-api/Screen+Shot+2020-11-25+at+1.17.23+PM.png)

# Resources:
  - https://golang.org/doc/effective_go.html
  - https://github.com/golang-standards/project-layout
  - https://github.com/vardius/go-api-boilerplate
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice