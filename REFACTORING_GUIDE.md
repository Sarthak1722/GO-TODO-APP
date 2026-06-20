# Production-Style Todo App - Refactoring Guide

## 🎯 What Was Changed

Your Todo app has been refactored to follow **production-level architecture patterns**. The code logic remains the same, but the structure now follows **clean architecture principles**.

### New Project Structure

```
todo-app/
├── cmd/
│   └── api/
│       └── main.go                 # Entry point, dependency injection
├── internal/
│   ├── handlers/                   # HTTP request handlers
│   │   └── todo_handler.go         # Updated to use service layer
│   ├── service/                    # Business logic layer (NEW)
│   │   └── todo_service.go         # Validation + business rules
│   ├── store/                      # Data access layer (NEW)
│   │   └── todo_store.go           # In-memory storage
│   ├── models/                     # Data models
│   │   └── todo.go
│   ├── dto/                        # Data transfer objects
│   │   └── todo_request.go
│   ├── middleware/                 # HTTP middleware
│   │   ├── logger.go               # Request logging
│   │   ├── requestid.go            # Request ID tracking
│   │   └── recover.go              # Panic recovery (NEW)
│   ├── utils/                      # Shared utilities (NEW)
│   │   └── response.go             # Consistent API responses
│   ├── validator/                  # Validation setup
│   │   └── validator.go
│   ├── errors/                     # Error handling
│   │   └── validationErrors.go
│   └── logger/                     # Logging setup
│       └── logger.go
├── go.mod
├── go.sum
└── air.toml                        # Hot reload config

```

## 📚 Architecture Layers (Clean Architecture)

### Layer 1: HTTP Handler Layer (internal/handlers/)
```
Handler receives HTTP request
↓
Validates basic request binding
↓
Calls service layer
↓
Returns HTTP response using utils/response.go
```

**File**: [internal/handlers/todo_handler.go](internal/handlers/todo_handler.go)

**What it does**:
- Accepts HTTP requests
- Binds JSON to structs
- Calls service layer functions
- Returns consistent API responses
- Logs with request ID

**Example**:
```go
func (h *TodoHandler) CreateTodo(c fiber.Ctx) error {
    requestID := c.Locals("request_id").(string)
    todoReq := dto.CreateTodoRequest{}
    
    if err := c.Bind().Body(&todoReq); err != nil {
        logger.Log.Error().Str("requestID", requestID).Err(err).Msg("bad request body")
        return utils.RespondError(c, fiber.StatusBadRequest, "invalid request body", nil)
    }
    
    createdTodo, validationErrs := h.service.CreateTodo(todoReq)
    if validationErrs != nil {
        return utils.RespondValidationError(c, validationErrs)
    }
    
    return utils.RespondSuccess(c, fiber.StatusCreated, createdTodo, "todo created")
}
```

---

### Layer 2: Service Layer (internal/service/)
```
Service receives validated data
↓
Performs business logic & validation
↓
Calls store layer for data operations
↓
Returns processed data
```

**File**: [internal/service/todo_service.go](internal/service/todo_service.go)

**What it does**:
- Contains all business logic
- Validates input data using validator
- Calls store for data persistence
- Can be reused by different handlers (CLI, API, etc.)

**Key Concept**: Your business rules are now **independent of HTTP**!

**Example**:
```go
func (s *TodoService) CreateTodo(req dto.CreateTodoRequest) (*models.Todo, map[string]string) {
    // Validate request
    if err := validator.Validate.Struct(req); err != nil {
        return nil, errors.FormatValidationErrors(err)
    }
    
    // Convert DTO to model
    todo := models.Todo{
        Body:      req.Body,
        Completed: req.Completed,
    }
    
    // Create in store
    createdTodo := s.store.CreateTodo(todo)
    return &createdTodo, nil
}
```

---

### Layer 3: Store Layer (internal/store/)
```
Store receives data
↓
Performs data persistence operations
↓
Returns data from storage
```

**File**: [internal/store/todo_store.go](internal/store/todo_store.go)

**What it does**:
- Handles data persistence (in-memory, database, etc.)
- No business logic
- No HTTP knowledge
- Easy to swap implementations

**Benefit**: Want to switch from in-memory to PostgreSQL? Just create a `PostgresStore` struct that implements the same `Store` interface!

**Interface Pattern** (notice this is what enables flexibility):
```go
type Store interface {
    GetAllTodos() []models.Todo
    GetTodoByID(id int) *models.Todo
    CreateTodo(todo models.Todo) models.Todo
    DeleteTodoByID(id int) bool
    PatchTodoByID(id int, todo models.Todo) *models.Todo
}
```

---

### Layer 4: Utilities (internal/utils/)
```
Handler returns response
↓
Uses utils/response.go for consistent format
```

**File**: [internal/utils/response.go](internal/utils/response.go)

**What it does**:
- Provides consistent response format across all endpoints
- Standardizes error responses
- Standardizes success responses

**Consistent API Response Format**:

✅ **Success Response**:
```json
{
  "success": true,
  "data": { "id": 1, "body": "Buy milk", "completed": false },
  "message": "todo created"
}
```

❌ **Error Response**:
```json
{
  "success": false,
  "error": "todo not found",
  "details": null
}
```

✔️ **Validation Error Response**:
```json
{
  "success": false,
  "error": "validation failed",
  "details": {
    "body": "Too short"
  }
}
```

---

## 🔄 Request Flow (From Start to End)

### When a client creates a todo with POST /api/todos:

```
1. main.go
   └─> Initializes Store, Service, Handler
   └─> Sets up routes
   └─> Starts server

2. Request arrives at handler
   POST /api/todos
   └─> middleware/requestid.go generates request ID
   └─> middleware/logger.go logs request
   └─> handlers/todo_handler.go:CreateTodo()

3. Handler CreateTodo
   └─> Binds JSON to DTO
   └─> Calls service.CreateTodo(dto)

4. Service CreateTodo
   └─> Validates using validator
   └─> Creates model
   └─> Calls store.CreateTodo(model)

5. Store CreateTodo
   └─> Assigns ID
   └─> Appends to in-memory slice
   └─> Returns created todo

6. Service returns to Handler
   └─> Handler gets todo back
   └─> Calls utils.RespondSuccess()

7. Response sent to client
   └─> middleware/logger.go logs response
   └─> 201 Created with consistent format
```

---

## 📝 New Response Format (Production Level)

All endpoints now return a **consistent response structure**:

### Before (Inconsistent):
```json
// Different formats for different endpoints
// GET /api/todos
[{ "id": 1, "body": "...", "completed": false }]

// POST /api/todos
{ "msg": "todo created", "data": {...} }

// DELETE /api/todos/:id
{ "msg": "Todo Deleted" }
```

### After (Consistent):
```json
// GET /api/todos
{
  "success": true,
  "data": [{ "id": 1, "body": "...", "completed": false }],
  "message": ""
}

// POST /api/todos
{
  "success": true,
  "data": { "id": 1, "body": "...", "completed": false },
  "message": "todo created"
}

// DELETE /api/todos/:id
{
  "success": true,
  "data": null,
  "message": "todo deleted"
}

// Error example
{
  "success": false,
  "error": "todo not found",
  "details": null
}
```

---

## 🛠️ How to Extend This (Learning Path)

### 1. Add a New Handler Method (Baby Step)

Want to add a **GET /api/todos/search** endpoint? Follow these steps:

**Step 1**: Add method to Service
```go
// internal/service/todo_service.go
func (s *TodoService) SearchTodos(query string) []models.Todo {
    allTodos := s.GetAllTodos()
    var results []models.Todo
    for _, todo := range allTodos {
        if strings.Contains(strings.ToLower(todo.Body), strings.ToLower(query)) {
            results = append(results, todo)
        }
    }
    return results
}
```

**Step 2**: Add method to Handler
```go
// internal/handlers/todo_handler.go
func (h *TodoHandler) SearchTodos(c fiber.Ctx) error {
    requestID := c.Locals("request_id").(string)
    query := c.Query("q")
    
    if query == "" {
        return utils.RespondError(c, fiber.StatusBadRequest, "query parameter 'q' is required", nil)
    }
    
    results := h.service.SearchTodos(query)
    logger.Log.Info().Str("requestID", requestID).Str("query", query).Msg("search performed")
    
    return utils.RespondSuccess(c, fiber.StatusOK, results, "")
}
```

**Step 3**: Register route in main.go
```go
// cmd/api/main.go
app.Get("/api/todos/search", todoHandler.SearchTodos)
```

---

### 2. Swap Store Implementation (Database Switch)

Want to use **PostgreSQL** instead of in-memory?

**Step 1**: Create new store implementation
```go
// internal/store/postgres_store.go
package store

type PostgresStore struct {
    db *sql.DB
}

func (p *PostgresStore) GetAllTodos() []models.Todo {
    rows, _ := p.db.Query("SELECT id, body, completed FROM todos")
    // Process rows...
    return todos
}

// Implement all Store interface methods...
```

**Step 2**: Use it in main.go
```go
// Before
todoStore := store.NewInMemoryStore()

// After
todoStore := store.NewPostgresStore(dbConnection)

// Rest of code stays exactly the same!
```

**Key Benefit**: Your handlers and service layers **don't change at all**! Only swapped the store.

---

### 3. Add Pagination (Learning Exercise)

Add optional `?page=1&limit=10` query parameters:

**In Service**:
```go
type PaginationParams struct {
    Page  int
    Limit int
}

func (s *TodoService) GetAllTodosPaginated(params PaginationParams) ([]models.Todo, int) {
    all := s.GetAllTodos()
    total := len(all)
    
    start := (params.Page - 1) * params.Limit
    end := start + params.Limit
    
    if start >= len(all) {
        return []models.Todo{}, total
    }
    if end > len(all) {
        end = len(all)
    }
    
    return all[start:end], total
}
```

**In Handler**:
```go
func (h *TodoHandler) GetAllTodos(c fiber.Ctx) error {
    page, _ := strconv.Atoi(c.Query("page", "1"))
    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    
    todos, total := h.service.GetAllTodosPaginated(service.PaginationParams{
        Page:  page,
        Limit: limit,
    })
    
    return utils.RespondSuccess(c, fiber.StatusOK, fiber.Map{
        "items": todos,
        "total": total,
        "page":  page,
        "limit": limit,
    }, "")
}
```

---

## ✅ Production Checklist

You now have:
- ✅ **Layered Architecture** (Clean Architecture pattern)
- ✅ **Dependency Injection** (handlers depend on services, not databases)
- ✅ **Consistent API Responses** (all endpoints return same format)
- ✅ **Request ID Tracking** (every request has unique ID for logging)
- ✅ **Proper Middleware** (request logging, panic recovery, ID generation)
- ✅ **Separation of Concerns** (Handler ≠ Service ≠ Store)
- ✅ **Interface-Based Design** (easy to swap implementations)
- ✅ **Proper Error Handling** (validation errors separate from system errors)

---

## 🎓 Key Learning Points

1. **Handler** = HTTP stuff only (binding, responses)
2. **Service** = Business logic (validation, calculations, decisions)
3. **Store** = Data persistence (no HTTP, no validation)
4. **Utils** = Reusable code (response formatting)

This structure means:
- You can test Service without HTTP
- You can swap Store without changing handlers
- You can reuse Service in CLI tools
- Everything is easier to debug and maintain

---

## 📂 File Reference

| File | Purpose | Learning Focus |
|------|---------|-----------------|
| [cmd/api/main.go](cmd/api/main.go) | Application entry point & dependency injection | How layers connect |
| [internal/handlers/todo_handler.go](internal/handlers/todo_handler.go) | HTTP request handling | HTTP specifics only |
| [internal/service/todo_service.go](internal/service/todo_service.go) | Business logic | Independent of HTTP |
| [internal/store/todo_store.go](internal/store/todo_store.go) | Data persistence | Interface pattern |
| [internal/utils/response.go](internal/utils/response.go) | Consistent responses | Standard formats |
| [internal/middleware/recover.go](internal/middleware/recover.go) | Panic recovery | Error safety |
| [internal/dto/todo_request.go](internal/dto/todo_request.go) | Data transfer objects | Input validation |
| [internal/models/todo.go](internal/models/todo.go) | Domain models | Core data structure |

---

## 🚀 Next Steps

1. **Test the refactored code**: Make API calls and verify responses are consistent
2. **Add a new feature** using the pattern above (search, sort, filter)
3. **Try swapping the store** by creating a different implementation
4. **Add pagination** using the example provided
5. **Add database layer** (PostgreSQL, MongoDB, etc.)

Happy learning! 🎉
