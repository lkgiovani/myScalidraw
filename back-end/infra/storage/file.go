package storage

import (
	"os"
)

func ReadFileFromDisk(path string) (string, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
