# 🚀 GO-TODO-APP

> A production-grade backend engineering project built from first principles using **Go**.  
> This is **not just a TODO app** — it is a journey of transforming a simple CRUD API into a **real-world scalable backend system** while learning professional backend engineering practices.

---

## 📌 Project Philosophy

Most developers build CRUD apps.

This project is different.

The goal is **not** to build a TODO application.

The goal is to use a simple TODO domain as a vehicle to deeply understand:

- Backend Engineering Fundamentals
- API Design
- Clean Architecture
- Database Engineering
- Infrastructure Engineering
- Observability
- Performance Engineering
- Production Systems Design

This repository progressively evolves through **5 stages**, eventually becoming a **production-grade backend system**.

---

## 🌍 Live Demo

🔗 **Live Website:**  
[Add Link Here]

---

## 🎥 Demo Video

📺 **Project Walkthrough:**  
[Add YouTube / Loom Video Here]

---

## 📸 Screenshots

### API Documentation

<!-- Add Screenshot -->

![API Docs](./assets/api-docs.png)

---

### Architecture Diagram

<!-- Add Screenshot -->

![Architecture](./assets/architecture.png)

---

### Monitoring Dashboard (Grafana)

<!-- Add Screenshot -->

![Monitoring](./assets/monitoring.png)

---

## 🏗 Final Architecture

```text
Client
   ↓
Rate Limiter
   ↓
Request ID Middleware
   ↓
Timeout Middleware
   ↓
Logger Middleware
   ↓
Recover Middleware
   ↓
HTTP Handler Layer
   ↓
Service Layer
   ↓
Repository Layer
   ↓
Redis Cache
   ↓
PostgreSQL Database

Metrics → Prometheus

Monitoring → Grafana

Tracing → OpenTelemetry

CI/CD → GitHub Actions
```

---

# ⚙️ Tech Stack

## Backend

- Go
- Fiber

## Database

- PostgreSQL
- pgx (No ORM)

## Caching

- Redis

## Infrastructure

- Docker
- Docker Compose

## Logging & Monitoring

- Zerolog
- Prometheus
- Grafana
- OpenTelemetry

## Testing

- Unit Testing
- Integration Testing
- Load Testing (k6)

## CI/CD

- GitHub Actions

---

# 📚 Development Roadmap

This project is intentionally built in stages.

---

## Stage 1 — Production Grade API Foundations

Implemented:

- REST API Design
- CRUD Operations
- Fiber HTTP Server
- DTOs
- Request Validation
- Structured Logging
- Middleware
- Panic Recovery Middleware
- Health Check Endpoint
- Centralized Error Handling

Concepts Learned:

- HTTP Lifecycle
- Request Parsing
- JSON Marshaling / Unmarshaling
- Middleware Pipeline
- Validation Architecture

---

## Stage 2 — Backend Architecture Engineering

Implemented:

- Handler Layer
- Service Layer
- Store Layer
- Dependency Injection
- Interfaces
- Mutex for Concurrency Safety

Concepts Learned:

- Separation of Concerns
- Race Conditions
- Single Responsibility Principle
- Encapsulation
- Concurrency Fundamentals

---

## Stage 3 — Database Engineering

Implemented:

- PostgreSQL Integration
- pgx Connection Pooling
- Raw SQL Queries
- Repository Pattern
- Migrations
- Environment Configuration

Concepts Learned:

- SQL Fundamentals
- Schema Design
- Transactions
- Query Optimization
- Database Architecture

---

## Stage 4 — Infrastructure Engineering

Implemented:

- Dockerization
- Docker Compose
- Graceful Shutdown
- Redis Caching
- Request Timeouts
- Rate Limiting
- Request IDs
- CORS
- Background Workers

Concepts Learned:

- Containerization
- Infrastructure Reliability
- Distributed Systems Basics
- Cache Design Patterns

---

## Stage 5 — Production Systems Engineering

Implemented:

- Prometheus Metrics
- Grafana Dashboards
- OpenTelemetry Tracing
- CI/CD Pipeline
- Integration Testing
- Load Testing
- Profiling
- Performance Benchmarking

Concepts Learned:

- Observability
- Production Monitoring
- Performance Engineering
- Reliability Engineering

---

# 📂 Project Structure

```text
GO-TODO-APP/

cmd/
  api/
    main.go

internal/

  handlers/
  service/
  repository/
  middleware/
  dto/
  models/
  utils/
  logger/
  db/
  config/

migrations/

Dockerfile

docker-compose.yml

README.md
```

---

# 🔌 API Endpoints

## Create Todo

```http
POST /api/todos
```

Request:

```json
{
  "body": "Learn Go"
}
```

---

## Get All Todos

```http
GET /api/todos
```

---

## Get Todo By ID

```http
GET /api/todos/:id
```

---

## Update Todo

```http
PATCH /api/todos/:id
```

---

## Delete Todo

```http
DELETE /api/todos/:id
```

---

## Health Check

```http
GET /health
```

---

# 🧠 Engineering Concepts Practiced

This project intentionally focuses on backend engineering fundamentals:

- HTTP Request Lifecycle
- Clean Architecture
- Structured Logging
- Validation Design
- Concurrency Safety
- Database Transactions
- Connection Pooling
- Docker Infrastructure
- Graceful Shutdown
- Redis Caching
- Rate Limiting
- API Security
- Monitoring & Observability
- Distributed Tracing
- CI/CD Pipelines
- Performance Profiling

---

# 🧪 Testing

Run unit tests:

```bash
go test ./...
```

Run benchmark tests:

```bash
go test -bench=.
```

Run load testing:

```bash
k6 run load-test.js
```

---

# 🐳 Running Locally

Clone repository:

```bash
git clone https://github.com/your-username/GO-TODO-APP.git
```

Run with Docker:

```bash
docker compose up --build
```

Or run manually:

```bash
go run cmd/api/main.go
```

---

# 📈 Future Improvements

Planned improvements:

- Authentication (JWT)
- OAuth Login
- Distributed Job Queue
- Event Driven Architecture
- WebSockets
- Kubernetes Deployment
- Blue-Green Deployment
- Multi-Service Architecture

---

# 📖 Learning Goal Behind This Project

This project exists because I wanted to deeply understand **backend engineering from first principles**.

Rather than building shallow tutorial projects, I chose to evolve one simple project into a production-grade backend system while understanding every engineering decision behind it.

---

# 🤝 Connect With Me

LinkedIn: [Add Link]

Portfolio: [Add Link]

Email: [Add Email]

---

# ⭐ Final Note

This repository represents a journey.

It starts as a simple TODO API.

It ends as a production-grade backend system.

The purpose is not CRUD.

The purpose is learning how real backend systems are built.
