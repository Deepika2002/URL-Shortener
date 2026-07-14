# AI Agent Development Rules & Quality Standards

## 1. Architectural Compliance
- Never combine the Read Service and Write Service; they must remain isolated microservices.
- Read Service MUST use a distributed mutex lock on cache misses to prevent cache stampedes.
- Write Service MUST implement the Redis Bloom Filter check before querying primary storage.

## 2. Code Quality & Performance
- All asynchronous tracking events must utilize a RoundRobinPartitioner into Redpanda.
- All database operations must be optimized for horizontal scaling (e.g., proper partition keys).
- Code must pass automated lints and include comprehensive unit/integration tests.

## 3. Development Workflow
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
