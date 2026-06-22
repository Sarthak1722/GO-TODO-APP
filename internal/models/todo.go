package models

type Todo struct {
	ID        int    `json:"id"`
	Body      string `json:"body"`
	Completed bool   `json:"completed"`
}




// Now that all 5 migrations are written and applied, here is how a senior engineer uses Tern to control the database state.

// Scenario A: Reverting a single mistake.
// Let's say dropping the completed column in Version 5 broke the frontend. I need to undo only that specific migration.
// Command: tern migrate -d 4
// Result: Tern runs the down migration of 005, putting the completed column back.

// Scenario B: Complete database reset.
// I am writing integration tests and need to wipe the database entirely to start fresh.
// Command: tern migrate -d 0
// Result: Tern walks backward from Version 5 down to Version 1, running every single down script in reverse order, ending with a completely empty database.

// Scenario C: Back to Production State.
// After wiping the database, I need all my tables back immediately.
// Command: tern migrate
// Result: Tern applies all 5 up scripts in sequential order.















// 1. The Database Layer (The Source of Truth)
// We write 002_add_priority.sql and run tern migrate. Our Postgres database now has a priority integer column. But right now, the Go application is completely blind to it.

// 2. The Domain Model Layer
// We must update your core data blueprint so Go knows what a Todo looks like now.

// Action: Go into models/todo.go and add Priority int to the Todo struct.

// 3. The DTO Layer (Data Transfer Objects)
// If a user is going to send you a priority when they create a task, your API contract must allow it.

// Action: Go into dto/todo_request.go and add Priority int to your CreateTodoRequest struct, likely adding validation tags (e.g., ensuring priority is between 1 and 5).

// 4. The Repository Layer (The Bridge)
// Your SQL queries are hardcoded strings. If you don't update them, Postgres will never receive the priority data from Go.

// Action: Go into repository/todo_repository.go. You must update your INSERT query to include the new column (e.g., INSERT INTO todos (title, completed, priority) VALUES ($1, $2, $3)). You also have to update your SELECT queries and ensure your Scan() functions map the new database column to your updated Go struct.

// 5. The Service & Handler Layers
// Finally, you wire the user's input down to the database.

// Action: Your Handler extracts the priority from the JSON payload, passes it to the Service, and the Service hands it to the Repository.