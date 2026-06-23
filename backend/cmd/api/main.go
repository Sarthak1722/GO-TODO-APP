package main

import (
	"log"
	"os"
	"path/filepath"

	_ "github.com/Sarthak1722/todo_app/docs" // Import the generated docs!
	"github.com/Sarthak1722/todo_app/internal/config"
	"github.com/Sarthak1722/todo_app/internal/database"
	"github.com/Sarthak1722/todo_app/internal/handlers"
	"github.com/Sarthak1722/todo_app/internal/logger"
	"github.com/Sarthak1722/todo_app/internal/middleware"
	"github.com/Sarthak1722/todo_app/internal/repository"
	"github.com/Sarthak1722/todo_app/internal/service"
	"github.com/Sarthak1722/todo_app/internal/validator"
	"github.com/clerk/clerk-sdk-go/v2"
	swaggo "github.com/gofiber/contrib/v3/swaggo"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/static"
	"github.com/joho/godotenv"
)

// THESE COMMENTS ARE WRITTEN FOR SWAGGO, IT USES THIS FOR CREATING THE OPENAPI.YAML FILE HEADERS.
// @title Simple Todo API
// @version 1.0
// @description A small authenticated todo API built with Go and Fiber.
// @host localhost:8080
// @BasePath /api
func main() {
	// load .env
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found; using environment variables and defaults")
	}
	// Initialize logger
	logger.Init("dev")

	// Initialize validator
	validator.Init()

	// Load configuration
	cfg := config.Load()

	var todoStore repository.Store

	// Initialize store based on config
	if cfg.DBType == "postgres" {
		// Initialize the PostgreSQL database connection pool
		dbPool, err := database.NewPostgresDB(cfg.PostgresDSN)
		if err != nil {
			log.Fatalf("Fatal: Could not connect to the database: %v", err)
		}

		// Ensure the pool is closed gracefully when the app shuts down
		defer dbPool.Close()

		// Initialize PostgreSQL repository
		todoStore = repository.NewPostgresTodoRepository(dbPool)
		log.Println("Using PostgreSQL database")
	} else {
		// Initialize in-memory repository
		todoStore = repository.NewInMemoryTodoRepository()
		log.Println("Using in-memory database")
	}

	// Initialize service (business logic layer)
	todoService := service.NewTodoService(todoStore)

	// Initialize handlers with injected service
	todoHandler := handlers.NewTodoHandler(todoService)

	// Backend initialization complete
	log.Println("Backend initialization complete. Server is ready!")

	// (Your HTTP server startup code like `http.ListenAndServe` will go here later)

	clerkSecret := os.Getenv("CLERK_SECRET_KEY")
	if clerkSecret == "" {
		logger.Log.Fatal().Msg("CLERK_SECRET_KEY is not set in environment")
	}
	// Set the global key for the Clerk SDK
	clerk.SetKey(clerkSecret)

	// Initialize Fiber app
	app := fiber.New()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	// Initialize cors config
	app.Use(cors.New(cors.Config{
		AllowOrigins: []string{"http://localhost:5173"},
		AllowMethods: []string{"GET","POST","PUT","PATCH","DELETE"},
		AllowHeaders: []string{"Origin", "Content-Type", "Accept", "Authorization"},
	}))

	// Middleware
	app.Use(middleware.GetRequestID())
	app.Use(middleware.RequestLogger())
	app.Use(middleware.RecoverPanic())

	// Attach the Swagger UI to your router
	app.Get("/swagger/*", swaggo.HandlerDefault)

	// Health check endpoint
	app.Get("/health", func(c fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "healthy",
		})
	})

	api := app.Group("/api", middleware.ClerkAuth())

	// Todo routes
	api.Get("/todos", todoHandler.GetAllTodos)
	api.Get("/todos/:id", todoHandler.GetTodoByID)
	api.Post("/todos", todoHandler.CreateTodo)
	api.Delete("/todos/:id", todoHandler.DeleteTodoByID)
	api.Patch("/todos/:id", todoHandler.PatchTodoByID)

	frontendDist := filepath.Clean("../frontend/dist")
	if _, err := os.Stat(frontendDist); err == nil {
		app.Get("/*", static.New(frontendDist))
		app.Get("*", static.New(filepath.Join(frontendDist, "index.html")))
	} else {
		log.Printf("Frontend build not found at %s; serving API only", frontendDist)
	}

	// Start server
	log.Fatal(app.Listen(":" + port))
}
