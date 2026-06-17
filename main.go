package main

import (
    "log"

    "github.com/gofiber/fiber/v3"
)

func main() {
    // Initialize a new Fiber app
    app := fiber.New()

    // Define a route for the GET method on the root path '/'
    app.Get("/", func(c fiber.Ctx) error {
        return c.Status(200).JSON(fiber.Map{"msg":"hello world"})
    })

    // Start the server on port 3000
    log.Fatal(app.Listen(":3000"))
}