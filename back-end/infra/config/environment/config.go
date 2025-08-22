package environment

import (
	"myScalidraw/pkg/env"
	"myScalidraw/pkg/projectError"
)

type Config struct {
	HTTP struct {
		Url  string
		Port int
	}
	DB struct {
		Host     string
		Port     int
		User     string
		Password string
		Name     string
		SSLMode  string
	}
	MINIO struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		UseSSL    bool
	}
	URL_SHORTENED_PREFIX string
	REDIS                struct {
		Address string
	}
	JWT_SECRET   string
	FRONTEND_URL string
}

func NewConfig() (*Config, error) {

	httpUrl, err := getString("URL", "Error loading HTTP URL")
	if err != nil {
		return nil, err
	}

	httpPort, err := getInt("PORT", "Error loading HTTP Port")
	if err != nil {
		return nil, err
	}

	dbHost, err := getString("DB_HOST", "Error loading DB Host")
	if err != nil {
		return nil, err
	}

	dbPort, err := getInt("DB_PORT", "Error loading DB Port")
	if err != nil {
		return nil, err
	}

	dbUser, err := getString("DB_USER", "Error loading DB User")
	if err != nil {
		return nil, err
	}

	dbPassword, err := getString("DB_PASSWORD", "Error loading DB Password")
	if err != nil {
		return nil, err
	}

	dbName, err := getString("DB_NAME", "Error loading DB Name")
	if err != nil {
		return nil, err
	}

	dbSSLMode, err := getString("DB_SSL_MODE", "Error loading DB SSL Mode")
	if err != nil {
		return nil, err
	}

	urlShortenedPrefix, err := getString("URL_SHORTENED_PREFIX", "Error loading URL Shortened Prefix")
	if err != nil {
		return nil, err
	}

	redisAddress, err := getString("REDIS_ADDRESS", "Error loading Redis Address")
	if err != nil {
		return nil, err
	}

	jwtSecret, err := getString("JWT_SECRET", "Error loading JWT Secret")
	if err != nil {
		return nil, err
	}

	frontendUrl, err := getString("FRONTEND_URL", "Error loading Frontend URL")
	if err != nil {
		return nil, err
	}

	minioEndpoint, err := getString("MINIO_ENDPOINT", "Error loading MinIO Endpoint")
	if err != nil {
		return nil, err
	}

	minioAccessKey, err := getString("MINIO_ACCESS_KEY", "Error loading MinIO Access Key")
	if err != nil {
		return nil, err
	}

	minioSecretKey, err := getString("MINIO_SECRET_KEY", "Error loading MinIO Secret Key")
	if err != nil {
		return nil, err
	}

	minioBucket, err := getString("MINIO_BUCKET", "Error loading MinIO Bucket")
	if err != nil {
		return nil, err
	}

	minioUseSSL, err := getBool("MINIO_USE_SSL", "Error loading MinIO UseSSL")
	if err != nil {
		return nil, err
	}

	return &Config{
		HTTP: struct {
			Url  string
			Port int
		}{
			Url:  httpUrl,
			Port: httpPort,
		},
		DB: struct {
			Host     string
			Port     int
			User     string
			Password string
			Name     string
			SSLMode  string
		}{
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			Name:     dbName,
			SSLMode:  dbSSLMode,
		},
		MINIO: struct {
			Endpoint  string
			AccessKey string
			SecretKey string
			Bucket    string
			UseSSL    bool
		}{
			Endpoint:  minioEndpoint,
			AccessKey: minioAccessKey,
			SecretKey: minioSecretKey,
			Bucket:    minioBucket,
			UseSSL:    minioUseSSL,
		},
		URL_SHORTENED_PREFIX: urlShortenedPrefix,
		REDIS: struct {
			Address string
		}{
			Address: redisAddress,
		},
		JWT_SECRET:   jwtSecret,
		FRONTEND_URL: frontendUrl,
	}, nil
}

func getInt(key, errorMessage string) (int, error) {
	value, err := env.GetEnvOrDieAsInt(key)

	if err != nil {
		return 0, &projectError.Error{
			Code:    projectError.EINVALID,
			Message: errorMessage,
		}
	}

	return value, nil

}

func getString(key, errorMessage string) (string, error) {
	value, err := env.GetEnvOrDie(key)
	if err != nil {
		return "", &projectError.Error{
			Code:    projectError.EINVALID,
			Message: errorMessage,
		}
	}
	return value, nil
}

func getBool(key, errorMessage string) (bool, error) {
	value, err := env.GetEnvOrDieAsBool(key)
	if err != nil {
		return false, &projectError.Error{
			Code:    projectError.EINVALID,
			Message: errorMessage,
		}
	}
	return value, nil
}
