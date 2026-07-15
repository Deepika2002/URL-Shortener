# AI Agent Development Rules & Quality Standards

## 1. Architectural Compliance
- Never combine the Read Service and Write Service; they must remain isolated microservices.
- Read Service MUST use a distributed mutex lock on cache misses to prevent cache stampedes.
- Write Service MUST implement the Redis Bloom Filter check before querying primary storage.

## 2. Code Quality & Performance
- All asynchronous tracking events must utilize a RoundRobinPartitioner into Redpanda.
- All database operations must be optimized for horizontal scaling (e.g., proper partition keys).
- Code must pass automated lints and include comprehensive unit/integration tests.

## 3. Development Workflow (AI Scope & Regression Rules)
- **Scope Containment:** NEVER execute steps beyond what the user explicitly requested from `implementation_plan.md`. If the user asks for Step 1, do Step 1 and STOP. Do not proactively start Step 2.
- **Regression Prevention:** Before finalizing a step and ending your turn, you MUST run all existing tests in the current microservice or package (e.g., `go test ./...`) to ensure your new code hasn't broken functionality from previous steps.
- DO NOT write code without verifying the existing database schemas in `.ai/schemas.md`.
- Write structural code iteratively. Test one microservice end-to-end using Docker Compose before moving to the next.
- Database schemas MUST be applied automatically upon infrastructure startup (e.g., using an initialization script or container). Manual schema creation is prohibited.

## 4. Project Scaffolding & Monorepo Structure
We will use a monorepo approach for the microservices. The workspace MUST adhere to the following directory structure:
- `/cmd/`: Contains the main applications (e.g., `/cmd/read-service`, `/cmd/write-service`, `/cmd/kgs`).
- `/internal/`: Private application and library code specific to individual services.
- `/pkg/`: Public library code that is shared across multiple services.
- `/deploy/`: Infrastructure configuration (e.g., DB initialization, Dockerfiles).
- `/scripts/`: Helper scripts for local testing and DB migrations.

## 5. Avoiding Common AI Failures
- **Dependency Management:** Whenever you import a new third-party Go package, you MUST explicitly run `go get <package>` before attempting to run `go test` or `go build`.
- **Unit Test Isolation:** Unit tests MUST use interfaces and mocks (e.g., `testify/mock` or standard library interfaces) for external I/O (ScyllaDB, Redis, ZooKeeper, Kafka). Standard `go test ./...` must succeed without requiring live docker containers unless explicitly separated as integration tests.
- **Error Looping:** If you encounter a failing test or compilation error that you cannot fix after 2 consecutive attempts, STOP and explain the issue to the user. Do not silently loop or guess.
