package fx

import (
	"context"
	"fmt"
	"log"
	"myScalidraw/infra/config/environment"
	"myScalidraw/infra/database"
	"myScalidraw/internal/delivery/httpserver"
	"myScalidraw/internal/domain/models/file_model"
	"myScalidraw/internal/domain/models/user_model"
	"myScalidraw/pkg/projectError"
	"net/http"

	"go.uber.org/fx"
)

func RegisterFiberServerHooks(
	lc fx.Lifecycle,
	server *httpserver.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Fiber server running on port %d\n", server.Port)
				if err := server.Start(); err != nil && err != http.ErrServerClosed {
					projectError.FatalError(&projectError.Error{
						Code:      projectError.ECONFLICT,
						Message:   "Error starting Fiber server",
						Path:      "infra/fx/lifecycle.go",
						PrevError: err,
					})
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Stopping Fiber server...")
			return server.App.Shutdown()
		},
	})
}

func RegisterDatabaseHooks(
	lc fx.Lifecycle,
	db *database.DB,
	config *environment.Config,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			log.Printf("Connected to PostgreSQL database at %s:%d/%s\n",
				config.DB.URL_DB)

			log.Println("Running automatic migrations...")
			if err := db.AutoMigrate(
				&file_model.FileMetadata{},
				&user_model.User{},
				&user_model.UserFilePermission{},
			); err != nil {
				return fmt.Errorf("failed to execute migrations: %w", err)
			}
			log.Println("Migrations completed successfully")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Closing database connection...")
			return db.Close()
		},
	})
}
