# Summary
This is a tweeting app API I built with a couple of personal goals in mind: 1) familiarize myself with event-driven architecture and 2) learn to write Go. I also learned to use gRPC and RabbitMQ during the process. This API's functionality is straightfoward - you can create a user, log in, create a tweet, and follow other users to view their tweets - all stuff that could be built with a simpler monolithic REST API. However, I wanted to practice designing a different style of backend system while learning to write idiomatic Go. Technologies used include Go, gRPC, RabbitMQ, and MongoDB.

# Design Features:
  - [Event Driven Architecture (EDA)](https://en.wikipedia.org/wiki/Event-driven_architecture)
    - State-changing requests (i.e., creating a user, following a user, or creating a tweet) are processed via an event producer -> message queue -> event consumer.
    - The event producer service can "fire and forget" each request by publishing a message, which the consumer service then picks up from the queue to fulfill.
  - [Command Query Responsibility Segregation (CQRS)](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs):
    - This API separates read and write requests to optimize reads and prevent blocking of writes (see diagram below)
    - Reads are done via a Read View service, which stores a copy of all data in memory
    - Writes are done via the message queue
  - [Remote Procedure Call (RPC)](https://en.wikipedia.org/wiki/Remote_procedure_call):
    - gRPC was used for direct communication with the UI and between services
      - This allows for the flexibility to handle UI requests either synchronously or asynchronously
        - Reads are completed and provided synchronously in the gRPC response
        - Writes are completed asynchronously via the message queue
  - [Domain Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design):
    - Each microservice uses folder structure to separate logic into the following layers:
        - The Application layer defines the route handlers, i.e., the gRPC methods that a client is able to call. These handlers utilize domain objects and interfaces, but are not exposed to their implementation details.
        - The Domain layer defines the domain objects (i.e., `User`, `Follow`, and `Tweet`) and their respective repository interfaces. The Repository Pattern used by DDD is a means of using repository objects to abstract the getting/saving of domain objects.
        - The Infrastructure layer implements the repositories and/or other interface dependencies in the domain. Note the repository implementations differ by service, but do one of the following:
          - make a request to the database and/or Read View service to fetch domain objects
          - make a request to the database and/or Read View to save domain objects
          - produce an event that will save domain objects
  - [Dependency Injection (DI)](https://en.wikipedia.org/wiki/Dependency_injection):
    - The servers in this mono-repo depend on other objects (e.g., repositories, event producers) to handle routes. Instead of the servers constructing those objects themselves, they only possess interfaces of those objects and receive the implementations when instantiated - this allows for greater separation of concerns.
    - For example, the EventProducerServer (`cmd/eventproducer/internal/application/server.go`) takes in an `event.Producer` interface. The server object only knows that this injection will have a `Produce(Event)` method - the actual implementation is injected in `main.go` when the server starts.

# Project Structure:
  - `/cmd`: contains subdirectories, each containing the following code for one microservice:
    - `/internal`: code only used by the microservice (i.e., within the same `/cmd/<MICROSERVICE>` directory). Includes the following:
      - `/application`: Route handlers. Includes a `server.go` file that defines the server's gRPC methods. The Event Consumer is an exception, as it only listens to the message queue and is not a gRPC server.
      - `/domain`: Business logic. Includes type definitions for domain objects and repository interfaces
      - `/infrastructure`: Data persistence logic. Includes repository implementations
      - `main.go`: Root file used to run the service. Here, I load env variables, instantiate dependencies, and start the server.
    - `/proto`: contains protocol buffer definitions used by the gRPC server and client(s).
      - .proto files are the source files
      - .pb.go files are generated during the build
  - `/scripts/db`
    - JS scripts to create and upgrade a Mongo database, all run in order by the `upgrade.sh` script
  - `/test`
    - TO DO

# API Architecture:
![API Architecture](tweet-api-architecture.svg)

# Pre-Requisites:
  - Go (https://golang.org/doc/install)
  - Protocol Buffer (https://grpc.io/docs/languages/go/quickstart/)
  - RabbitMQ (https://www.rabbitmq.com/download.html)
  - MongoDB(https://docs.mongodb.com/manual/installation/)

# Getting Started:
  - Run MongoDB and RabbitMQ daemons:
    - In MacOS, something like `brew services start mongodb` and `brew services start rabbitmq`
  - Initialize MongoDB database:
    - In a terminal window, navigate to the /scripts/db directory and run `upgrade.sh`
  - Build services:
    - In a terminal window, run `make build-all`
  - Run services:
    - To run services locally, open a terminal window for each service and run the make run script (e.g., `make run BIN=eventproducer`)
    - Start the services in the following order to avoid connection timeout errors: 1) databaseaccess, 2) readview, 3) eventconsumer, 4) eventproducer, 5) apigateway
  - Ping the API gateway (via an RPC client tool such as BloomRPC) to create a user, log in, write a tweet, etc.

# Resources:
  - https://golang.org/doc/effective_go.html
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/ddd-oriented-microservice
  - https://docs.microsoft.com/en-us/dotnet/architecture/microservices/microservice-ddd-cqrs-patterns/infrastructure-persistence-layer-design
  - https://github.com/golang-standards/project-layout
  - https://github.com/vardius/go-api-boilerplate
