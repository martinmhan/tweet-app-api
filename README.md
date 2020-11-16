# tweet-app-api
Event-Driven, Microservice API for Tweeting app built with Go, gRPC, RabbitMQ

#Features:
  Event Driven Architecture (EDA)
    - Message-Queue with Pub/Sub used to process state-changing events
  Command Query Responsibility Segregation (CQRS):
    - The API handles read and write requests separately
    - Read Optimized View service:
      - Reads are done via an in-memory data view to optimize speed and to prevent blocking of writes
      - Only queries the DB on cold starts
      - Gets updated by the consumer after each successful write
  Domain Driven Design (DDD):
    - Each microservice uses folder structure and abstraction layers to follow DDD principles
      - Each `/cmd/<SERVICE>/internal` directory contains three subdirectories
        - `/application`: 
        - `/domain`: 
        - `/infrastructure`: 
    - Additionally, the overall architecture utilizes microservices to  the following:
      - API Gateway to directly handle all RPC requests from UI
      - Database Access service that handles all database reads and writes
  gRPC:
    - Request/Response-style RPC communication for direct interactions with the user via an API Gateway
    - Enables flexibility for both synchronous (in this case, reads) and asynchronous (writes) processing of requests
  Server-Sent Events (SSE)
    - Used to fan out new tweet notifications to a user's followers
  API Gateway:
    - JWT-based authentication
    - Request routing

# Notes:
  - Project Structure
    - `/cmd`: contains subdirectories, each containing the following code for one microservice:
      - `/internal`: code only used by the microservice (i.e., within the same `/cmd/<MICROSERVICE>` directory)
        - `/application`: gRPC server implementation defining RPC handlers
        - `/domain`: objects containing business logic used by the RPC handlers
        - `/infrastructure`: if necessary, objects that handle data persistence
      - `/proto`: contains protocol buffer definitions
        - .proto files are the source files
        - .pb.go files are generated during the build
      - `main.go`: entry point file that starts the microservice
    - `/util`
    - `.env`: 
    - `Makefile`