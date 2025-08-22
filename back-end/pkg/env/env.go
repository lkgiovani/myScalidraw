package env

import (
	"fmt"
	"myScalidraw/pkg/projectError"

	"os"
	"strconv"
)

func GetEnvOrDie(key string) (string, error) {
	value := os.Getenv(key)

	if value == "" {
		return "", &projectError.Error{
			Code:    projectError.ECONFLICT,
			Message: fmt.Sprintf("Missing environment variable %s", key),
		}
	}

	return value, nil
}

func GetEnvOrDieAsInt(key string) (int, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return -1, &projectError.Error{
			Code:    projectError.ECONFLICT,
			Message: fmt.Sprintf("Environment variable %s not set\n", key),
		}
	}

	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return -1, &projectError.Error{
			Code:    projectError.ECONFLICT,
			Message: fmt.Sprintf("Error converting %s to int: %v\n", key, err),
		}
	}

	return value, nil
}

func GetEnvOrDieAsBool(key string) (bool, error) {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return false, &projectError.Error{
			Code:    projectError.ECONFLICT,
			Message: fmt.Sprintf("Environment variable %s not set\n", key),
		}
	}

	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return false, &projectError.Error{
			Code:    projectError.ECONFLICT,
			Message: fmt.Sprintf("Error converting %s to bool: %v\n", key, err),
		}
	}

	return value, nil
}
