package httpserver

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"go.uber.org/fx"
)

type Server struct {
	App  *fiber.App
	Port int
}

func NewServer(port int) *Server {
	app := fiber.New(fiber.Config{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,PATCH,OPTIONS",
		AllowHeaders:     "Origin,Content-Type,Accept,Authorization,X-Requested-With",
		AllowCredentials: false,
		ExposeHeaders:    "Content-Length,Access-Control-Allow-Origin,Access-Control-Allow-Headers,Content-Type",
	}))

	return &Server{
		App:  app,
		Port: port,
	}
}

func (s *Server) Start() error {
	addr := fmt.Sprintf(":%d", s.Port)
	log.Printf("Starting Fiber server on %s", addr)
	return s.App.Listen(addr)
}

func RegisterHooks(lc fx.Lifecycle, server *Server) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := server.Start(); err != nil {
					log.Fatalf("Failed to start Fiber server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down Fiber server")
			return server.App.Shutdown()
		},
	})
}

var Module = fx.Options(
	fx.Provide(NewServer),
	fx.Invoke(RegisterHooks),
)
