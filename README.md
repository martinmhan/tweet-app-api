# tweet-app-api
- Event-Driven Microservices API for a simple tweeting app
- Built with Go, gRPC, RabbitMQ, MongoDB

# Pre-requisites:
  TBD

# Getting Started:
  TBD

# Features:
  Event Driven Architecture (EDA)
    - Message-Queue used to produce/consume state-changing events (in this case, database writes)
  Command Query Responsibility Segregation (CQRS):
    - This API handles read and write requests separately
    - Reads are done via a Read View service:
      - This is an-memory data view used to optimize read speeds and to prevent blocking of writes
      - Only queries the DB on cold starts
      - Gets updated by the consumer after each successful write
    - Writes are done via the queue
      - The Events Producer publishes a write event to the message queue
      - The Events Consumer subscribes to and processes the write event
  Domain Driven Design (DDD):
    - Each microservice uses folder structure and abstraction layers to follow DDD principles
      - (see `/cmd/internal/` Project Structure notes below)
    - The overall architecture utilizes microservices to separate contexts and functionality:
      - (see API architecture diagram below)
  Remote Procedure Call (RPC):
    - gRPC used for direct communication with the UI and between services
      - The only exception being the Events Producer and Events Consumer, which communicate via the Message Queue
      - This allows for flexibility to handle UI requests either synchronously or asynchronously
        - Reads are completed and provided synchronously in the gRPC response
        - Writes are completed asynchronously via the message queue
  Server-Sent Events (SSE)
    - Used to fan out new tweet notifications to a user's followers

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
  - Services
    - API Gateway: 
    - Events Producer:
  - RabbitMQ has more robust functionality such as routing messages to queues via exchanges. However, just one queue and no exchage was used here due to the API's simplicity.

<API architecture diagram>

# Resources:
  - https://golang.org/doc/effective_go.html
  - https://github.com/golang-standards/project-layout
  - https://github.com/vardius/go-api-boilerplate