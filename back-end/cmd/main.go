package main

import (
	"sync"

	"github.com/gofiber/fiber/v2"
)

var store = make(map[string]string)
var mu sync.Mutex

func main() {
	app := fiber.New()

	// GET /api/sala/:id  -> retorna o desenho salvo
	app.Get("/api/sala/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		mu.Lock()
		data, exists := store[id]
		mu.Unlock()

		if !exists {
			return c.Status(404).JSON(fiber.Map{"error": "sala não encontrada"})
		}

		return c.Type("json").SendString(data)
	})

	// POST /api/sala/:id -> salva o desenho
	app.Post("/api/sala/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")

		body := c.Body() // pega o JSON bruto enviado
		if len(body) == 0 {
			return c.Status(400).JSON(fiber.Map{"error": "body vazio"})
		}

		mu.Lock()
		store[id] = string(body)
		mu.Unlock()

		return c.JSON(fiber.Map{"status": "salvo com sucesso"})
	})

	app.Listen(":8080")
}
