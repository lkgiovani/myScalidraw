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
		URL_DB string
	}
	MINIO struct {
		Endpoint  string
		AccessKey string
		SecretKey string
		Bucket    string
		UseSSL    bool
	}
	URL_SHORTENED_PREFIX string
	JWT_SECRET           string
	FRONTEND_URL         string
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

	dbURL, err := getString("URL_DB", "Error loading DB URL")
	if err != nil {
		return nil, err
	}

	urlShortenedPrefix, err := getString("URL_SHORTENED_PREFIX", "Error loading URL Shortened Prefix")
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
			URL_DB string
		}{
			URL_DB: dbURL,
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
		JWT_SECRET:           jwtSecret,
		FRONTEND_URL:         frontendUrl,
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
