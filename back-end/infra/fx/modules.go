package fx

import (
	"myScalidraw/infra/config/environment"
	"myScalidraw/infra/database"
	"myScalidraw/infra/storage"
	"myScalidraw/internal/delivery/handlers/fileHandlers"
	"myScalidraw/internal/delivery/httpserver"
	"myScalidraw/internal/domain/repository"
	"myScalidraw/internal/domain/repository/impl"
	"myScalidraw/internal/domain/useCase/file"

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
		func(db *database.DB) repository.FileMetadataRepository {
			return impl.NewFileMetadataRepository(db)
		},
	),

	fx.Provide(
		func(minioClient *storage.MinIO, db *database.DB, metadataRepo repository.FileMetadataRepository) repository.FileRepository {
			return impl.NewFileRepositoryMinio(minioClient, db, metadataRepo)
		},
	),
)

var UseCaseModule = fx.Options(
	fx.Provide(
		func(fileRepo repository.FileRepository, metadataRepo repository.FileMetadataRepository) *file.FileUseCase {
			return file.NewFileUseCase(fileRepo, metadataRepo)
		},
	),
)

var HandlersModule = fx.Options(
	fx.Provide(fileHandlers.NewFileHandler),
	fx.Invoke(
		func(server *httpserver.Server, fileHandler *fileHandlers.FileHandler) {
			fileHandler.RegisterRoutes(server.App)
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
	HandlersModule,
)
