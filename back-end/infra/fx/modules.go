package fx

import (
	"myScalidraw/infra/config/environment"
	"myScalidraw/infra/database"
	"myScalidraw/infra/storage"
	"myScalidraw/internal/delivery/handlers"
	"myScalidraw/internal/delivery/handlers/fileHandlers"
	"myScalidraw/internal/delivery/handlers/userHandlers"
	"myScalidraw/internal/delivery/httpserver"
	"myScalidraw/internal/delivery/middleware"
	"myScalidraw/internal/domain/repository/file_repository"
	"myScalidraw/internal/domain/repository/user_repository"
	"myScalidraw/internal/domain/useCase/file"
	"myScalidraw/internal/domain/useCase/user"
	"myScalidraw/pkg/jwt"

	"go.uber.org/fx"
)

var ConfigModule = fx.Options(
	fx.Provide(environment.NewConfig),
)

var DatabaseModule = fx.Options(
	fx.Provide(
		func(config *environment.Config) database.PostgresConfig {
			return database.NewPostgresConfig(
				config.DB.URL_DB,
			)
		},
	),
	fx.Provide(database.NewDB),
	fx.Invoke(RegisterDatabaseHooks),
)

var StorageModule = fx.Options(
	fx.Provide(
		func(config *environment.Config) storage.MinIOConfig {
			return storage.NewMinIOConfig(
				config.MINIO.Endpoint,
				config.MINIO.AccessKey,
				config.MINIO.SecretKey,
				config.MINIO.Bucket,
				config.MINIO.UseSSL,
			)
		},
	),
	fx.Provide(storage.NewMinIO),
)

var ServerModule = fx.Options(
	fx.Provide(
		func(config *environment.Config) *httpserver.Server {
			return httpserver.NewServer(config.HTTP.Port)
		},
	),
	fx.Invoke(RegisterFiberServerHooks),
)

var RepositoryModule = fx.Options(
	fx.Provide(
		func(db *database.DB) file_repository.FileMetadataRepository {
			return file_repository.NewFileMetadataRepository(db)
		},
	),

	fx.Provide(
		func(minioClient *storage.MinIO, db *database.DB, metadataRepo file_repository.FileMetadataRepository) file_repository.FileRepository {
			return file_repository.NewFileRepositoryMinio(minioClient, db, metadataRepo)
		},
	),

	fx.Provide(
		func(db *database.DB) user_repository.UserRepository {
			return user_repository.NewUserRepository(db)
		},
	),
)

var UseCaseModule = fx.Options(
	fx.Provide(
		func(fileRepo file_repository.FileRepository, metadataRepo file_repository.FileMetadataRepository) *file.FileUseCase {
			return file.NewFileUseCase(fileRepo, metadataRepo)
		},
	),

	fx.Provide(
		func(userRepo user_repository.UserRepository) *user.UserUseCase {
			return user.NewUserUseCase(userRepo)
		},
	),
)

var AuthModule = fx.Options(
	fx.Provide(
		func(config *environment.Config) *jwt.JWTManager {
			return jwt.NewJWTManager(config.JWT_SECRET)
		},
	),

	fx.Provide(
		func(jwtManager *jwt.JWTManager, userUseCase *user.UserUseCase) *middleware.AuthMiddleware {
			return middleware.NewAuthMiddleware(jwtManager, userUseCase)
		},
	),
)

var HandlersModule = fx.Options(
	fx.Provide(fileHandlers.NewFileHandler),
	fx.Provide(userHandlers.NewUserHandler),
	fx.Provide(handlers.NewPingHandler),
	fx.Invoke(
		func(server *httpserver.Server, fileHandler *fileHandlers.FileHandler, userHandler *userHandlers.UserHandler, authMiddleware *middleware.AuthMiddleware, pingHandler *handlers.PingHandler) {
			fileHandler.RegisterRoutes(server.App)
			userHandler.RegisterRoutes(server.App, authMiddleware)
			pingHandler.RegisterRoutes(server.App)
		},
	),
)

var AllModules = fx.Options(
	ConfigModule,
	DatabaseModule,
	StorageModule,
	ServerModule,
	RepositoryModule,
	UseCaseModule,
	AuthModule,
	HandlersModule,
)
