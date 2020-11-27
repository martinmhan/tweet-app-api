# Summary
  - This is an API for a simple tweeting app I built with the following goals in mind:
    - Familiarize myself with backend design principles, namely Event-Driven Architecture, CQRS, and Domain-Driven Design (see Design Feaatures below)
    - Learn to write effective and idiomatic Go
    - Learn to use new technologies, i.e., RabbitMQ and gRPC
  - Technologies Used:
    - Go, gRPC, RabbitMQ, MongoDB
  - (Swift iOS Client TBD)


# Design Features:
  - [Event Driven Architecture](https://en.wikipedia.org/wiki/Event-driven_architecture) (EDA)
    - Message-Queue used to produce/consume state-changing events (in this case, database writes)
  - [Command Query Responsibility Segregation](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs) (CQRS):
    - This API handles read and write requests separately to optimize reads and prevent blocking of writes
    - Reads are done via a Read View service:
      - This is an in-memory data store
      - Only queries the DB to populate the data store on cold starts
      - Gets updated by the event consumer after each successful write
    - Writes are done via the queue
      - The Event Producer publishes a write event to the message queue
      - The Event Consumer executes the write event (i.e., writes to DB, then updates read view)
  - [Remote Procedure Call](https://en.wikipedia.org/wiki/Remote_procedure_call) (RPC):
    - gRPC used for direct communication with the UI and between services
      - The only exception being the Events Producer and Events Consumer, which communicate via the Message Queue
      - This allows for flexibility to handle UI requests either synchronously or asynchronously
        - Reads are completed and provided synchronously in the gRPC response
        - Writes are completed asynchronously via the message queue
  - [Domain Driven Design](https://en.wikipedia.org/wiki/Domain-driven_design) (DDD):
    - Each microservice uses folder structure to separate logic into the following layers:
        - Application: In this layer, I define the route handlers which in gRPC are the methods that a gRPC client is able to call. These handlers utilize interfaces from the Domain and are not exposed to their implementations under the hood.
        - Domain: In this layer, I define my domain objects (i.e., User, Follower, and Tweet) and their respective Repository interfaces (if needed by the service). The Repository Pattern used by DDD basically utilizes repository objects to abstract the getting/saving of domain objects.
        - Infrastructure: In this layer, I implement the repositories, so this is where I define the getting/saving functions. The implementations will differ by service, but will do one of the following:
          - make a request to the database or read view to fetch domain objects
          - make a request to save a new domain object to the database and/or read view
          - produce an event to save a domain object

# Notes:
  - Project Structure
    - `/cmd`: contains subdirectories, each containing the following code for one microservice:
      - `/internal`: code only used by the microservice (i.e., within the same `/cmd/<MICROSERVICE>` directory). Includes the following: 
        - `/application`: Implementation of route handlers. This folder will contain one `server.go` file that defines all of the methods that gRPC clients can call.
        - `/domain`: Business logic and type/interface definitions used by the gRPC methods
        - `/infrastructure`: Data persistence logic and implementations repository pattern
        - `main.go`: root file used to run this service
      - `/proto`: contains protocol buffer definitions
        - .proto files are the source files
        - .pb.go files are generated during the build
    - `.env`: environment variables
    - `Makefile`: scripts used to build and run services

# API Architecture:
![API Architecture](https://gitbuckets.s3-us-west-1.amazonaws.com/tweet-app-api/Screen+Shot+2020-11-25+at+1.17.23+PM.png)


# Pre-Requisites:
  TBD

# Getting Started:
  - Run MongoDB and RabbitMQ daemons:
    - In MacOS, something like `brew services start mongodb` and `brew services start rabbitmq`
  - Create MongoDB database:
    - In a terminal window, navigate to the /scripts/db directory and run `upgrade.sh`
  - Build and run services:
    - To run services locally, open a terminal window and run the `build-and-run` Makefile script for each service (e.g., `make build-and-run BIN=eventproducer`)
    - Start the services in the following order to avoid fatal connection errors: 1) databaseaccess, 2) readview, 3) eventconsumer, 4) eventproducer, 5) apigateway
  - Ping the API gateway (via an RPC client tool such as BloomRPC) to create a user, log in, write a tweet, etc.

# Resources:
  - https://golang.org/doc/effective_go.html
  - https://github.com/golang-standards/project-layout
  - https://github.com/vardius/go-api-boilerplate
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design