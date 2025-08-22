package fx

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.uber.org/fx"

	"myScalidraw/infra/config/environment"
	"myScalidraw/infra/database"
	"myScalidraw/internal/delivery/httpserver"
	"myScalidraw/internal/domain/models"
	"myScalidraw/pkg/projectError"
)

func RegisterFiberServerHooks(
	lc fx.Lifecycle,
	server *httpserver.Server,
) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Printf("Servidor Fiber rodando na porta %d\n", server.Port)
				if err := server.Start(); err != nil && err != http.ErrServerClosed {
					projectError.FatalError(&projectError.Error{
						Code:      projectError.ECONFLICT,
						Message:   "Erro ao iniciar servidor Fiber",
						Path:      "infra/fx/lifecycle.go",
						PrevError: err,
					})
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Parando servidor Fiber...")
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
			log.Printf("Conectado ao banco de dados PostgreSQL em %s:%d/%s\n",
				config.DB.Host, config.DB.Port, config.DB.Name)

			log.Println("Executando migrações automáticas...")
			if err := db.AutoMigrate(&models.FileMetadata{}); err != nil {
				return fmt.Errorf("falha ao executar migrações: %w", err)
			}
			log.Println("Migrações concluídas com sucesso")

			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Fechando conexão com o banco de dados...")
			return db.Close()
		},
	})
}
