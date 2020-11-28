# Summary
This is a tweeting app API I built with a couple of learning goals in mind: 1) familiarize myself with new architectural design principles (see Design Features below), and 2) learn new tech along the way (Go, gRPC, RabbitMQ). The API's functionality is simple - you can create a user, log in, create a tweet, and follow other users to view their tweet. However, I wanted to practice designing an efficient backend system all the while learning to write idiomatic Go. Technologies used include: Go, gRPC, RabbitMQ, and MongoDB.

# Design Features:
  - [Event Driven Architecture (EDA)](https://en.wikipedia.org/wiki/Event-driven_architecture)
    - Message-Queue used to produce/consume state-changing events (in this case, database writes)
  - [Command Query Responsibility Segregation (CQRS)](https://docs.microsoft.com/en-us/azure/architecture/patterns/cqrs):
    - This API handles read and write requests separately to optimize reads and prevent blocking of writes
    - Reads are done via a Read View service:
      - This is an in-memory data store
      - Only queries the DB to populate the data store on cold starts
      - Gets updated by the event consumer after each successful write
    - Writes are done via the queue
      - The Event Producer publishes a write event to the message queue
      - The Event Consumer executes the write event (i.e., writes to DB, then updates read view)
  - [Remote Procedure Call (RPC)](https://en.wikipedia.org/wiki/Remote_procedure_call):
    - gRPC used for direct communication with the UI and between services
      - The only exception being the Events Producer and Events Consumer, which communicate via the Message Queue
      - This allows for flexibility to handle UI requests either synchronously or asynchronously
        - Reads are completed and provided synchronously in the gRPC response
        - Writes are completed asynchronously via the message queue
  - [Domain Driven Design (DDD)](https://en.wikipedia.org/wiki/Domain-driven_design):
    - Each microservice uses folder structure to separate logic into the following layers:
        - In the Application layer, I define the route handlers, i.e., the gRPC methods that a client is able to call. These handlers utilize domain objects and interfaces, but are not exposed to their implementation details.
        - In the Domain layer, I define the domain objects (i.e., User, Follow, and Tweet) and their respective repository interfaces (if needed by the service). The Repository Pattern used by DDD is a means of abstracting the getting/saving of domain objects by utilizing repository objects.
        - In the Infrastructure layer, I implement the repositories, i.e., define the getting/saving functions. Please note the implementations will differ by service, but will do one of the following:
          - make a request to the database and/or Read View service to fetch domain objects
          - make a request to the database and/or Read View to save domain objects
          - produce an event that will save domain objects
  - [Dependency Injection (DI)](https://en.wikipedia.org/wiki/Dependency_injection):
    - The servers in this mono-repo depend on other objects (e.g., repositories, event producers) to handle their routes. Instead of the servers constructing those objects themselves, they only possess interfaces of those objects and receive the implementations when instantiated. This allows for greater separation of concerns.
    - For example, the EventProducerServer (`cmd/eventproducer/internal/application/server.go`) takes in an `event.Producer` interface. The server object only knows that this injection will have a `Produce(Event)` method - the actual implementation is injected in `main.go` when the server starts.

# Notes:
  - Project Structure
    - `/cmd`: contains subdirectories, each containing the following code for one microservice:
      - `/internal`: code only used by the microservice (i.e., within the same `/cmd/<MICROSERVICE>` directory). Includes the following: 
        - `/application`: Route handlers. Includes a `server.go` file that defines the server's gRPC methods. The Event Consumer is an exception, as it only listens to the message queue and is not a gRPC server.
        - `/domain`: Business logic. Includes type definitions for domain objects and repository interfaces
        - `/infrastructure`: Data persistence logic. Includes repository implementations
        - `main.go`: root file used to run this service
      - `/proto`: contains protocol buffer definitions used by the gRPC server
        - .proto files are the source files
        - .pb.go files are generated during the build
    - `.env`: environment variables
    - `Makefile`: scripts used to build and run services
  - Challenges
    - `Follow` Domain Object
      - A question I ran into was whether to make a "Follower" an entirely separate domain object, since an array of `User`s could carry the same info as a list of followers. However, this requires each array to be mapped to a User ID of the "Followee" and can otherwise cause confusion re: the direction of the Follower/Followee relationship.
      - I ended up making a domain object called "Follow" that represents a Follower/Followee relationship between two users. This allows for flat collections of `Follow`s and doesn't rely on external data structures to provide the relational info.

  - TBD
    - API tests
    - DB security (password hashing)
    - Swift iOS client
    - Fan out notifications to UI clients on tweet events

# API Architecture:
![API Architecture](https://gitbuckets.s3-us-west-1.amazonaws.com/tweet-app-api/Screen+Shot+2020-11-25+at+1.17.23+PM.png)


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
