package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/fx"
)

type MinIOConfig struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	Bucket    string
	UseSSL    bool
}

func NewMinIOConfig(endpoint, accessKey, secretKey, bucket string, useSSL bool) MinIOConfig {
	return MinIOConfig{
		Endpoint:  strings.TrimSpace(endpoint),
		AccessKey: strings.TrimSpace(accessKey),
		SecretKey: strings.TrimSpace(secretKey),
		Bucket:    strings.TrimSpace(bucket),
		UseSSL:    useSSL,
	}
}

func maskAccessKey(accessKey string) string {
	if len(accessKey) <= 4 {
		return strings.Repeat("*", len(accessKey))
	}
	return accessKey[:2] + strings.Repeat("*", len(accessKey)-4) + accessKey[len(accessKey)-2:]
}

type MinIO struct {
	Client *minio.Client
	Bucket string
}

func NewMinIO(config MinIOConfig) (*MinIO, error) {
	log.Printf("Connecting to MinIO at %s with access key: %s", config.Endpoint, maskAccessKey(config.AccessKey))

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client with provided credentials: %w", err)
	}

	ctx := context.Background()

	exists, err := client.BucketExists(ctx, config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to verify bucket access with provided credentials: %w", err)
	}

	if !exists {
		log.Printf("Creating bucket: %s", config.Bucket)
		err = client.MakeBucket(ctx, config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket '%s' - check access key permissions: %w", config.Bucket, err)
		}

		policy := fmt.Sprintf(`{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::%s/*"]
				}
			]
		}`, config.Bucket)

		err = client.SetBucketPolicy(ctx, config.Bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy - check access key permissions: %w", err)
		}
		log.Printf("Bucket '%s' created successfully with public read policy", config.Bucket)
	} else {
		log.Printf("Using existing bucket: %s", config.Bucket)
	}

	log.Printf("Successfully connected to MinIO storage at %s", config.Endpoint)
	return &MinIO{
		Client: client,
		Bucket: config.Bucket,
	}, nil
}

func (m *MinIO) UploadFile(fileID string, content []byte) (string, error) {
	if fileID == "" {
		return "", fmt.Errorf("file ID cannot be empty")
	}
	if len(content) == 0 {
		return "", fmt.Errorf("file content cannot be empty")
	}

	contentType := "application/json"
	objectName := fileID + ".json"
	reader := bytes.NewReader(content)

	ctx := context.Background()
	_, err := m.Client.PutObject(ctx, m.Bucket, objectName, reader, int64(len(content)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload file '%s' to MinIO: %w", fileID, err)
	}

	log.Printf("Successfully uploaded file: %s (size: %d bytes)", fileID, len(content))
	return fmt.Sprintf("http://%s/%s/%s", m.Client.EndpointURL().Host, m.Bucket, objectName), nil
}

func (m *MinIO) GetFile(fileID string) ([]byte, error) {
	if fileID == "" {
		return nil, fmt.Errorf("file ID cannot be empty")
	}

	objectName := fileID + ".json"
	ctx := context.Background()

	object, err := m.Client.GetObject(ctx, m.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get file '%s' from MinIO: %w", fileID, err)
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read file '%s' content: %w", fileID, err)
	}

	log.Printf("Successfully retrieved file: %s (size: %d bytes)", fileID, len(content))
	return content, nil
}

func (m *MinIO) DeleteFile(fileID string) error {
	if fileID == "" {
		return fmt.Errorf("file ID cannot be empty")
	}

	objectName := fileID + ".json"
	ctx := context.Background()

	err := m.Client.RemoveObject(ctx, m.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file '%s' from MinIO: %w", fileID, err)
	}

	log.Printf("Successfully deleted file: %s", fileID)
	return nil
}

func (m *MinIO) ListFiles() ([]string, error) {
	var files []string
	ctx := context.Background()

	objectCh := m.Client.ListObjects(ctx, m.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects in MinIO bucket '%s': %w", m.Bucket, object.Err)
		}
		files = append(files, object.Key)
	}

	log.Printf("Successfully listed %d files from bucket: %s", len(files), m.Bucket)
	return files, nil
}

func (m *MinIO) GetFileURL(fileID string) string {
	if fileID == "" {
		log.Printf("Warning: GetFileURL called with empty file ID")
		return ""
	}

	objectName := fileID + ".json"
	ctx := context.Background()

	presignedURL, err := m.Client.PresignedGetObject(ctx, m.Bucket, objectName, time.Hour*24*7, nil)
	if err != nil {
		log.Printf("Failed to generate presigned URL for file '%s': %v", fileID, err)
		return ""
	}

	log.Printf("Generated presigned URL for file: %s (expires in 7 days)", fileID)
	return presignedURL.String()
}

func (m *MinIO) CreateFolder(folderPath string) error {
	if folderPath == "" {
		return fmt.Errorf("folder path cannot be empty")
	}

	if !strings.HasSuffix(folderPath, "/") {
		folderPath += "/"
	}

	if strings.HasPrefix(folderPath, "/") {
		folderPath = folderPath[1:]
	}

	ctx := context.Background()
	reader := bytes.NewReader([]byte{})

	_, err := m.Client.PutObject(ctx, m.Bucket, folderPath, reader, 0,
		minio.PutObjectOptions{ContentType: "application/x-directory"})
	if err != nil {
		return fmt.Errorf("failed to create folder '%s' in MinIO: %w", folderPath, err)
	}

	log.Printf("Successfully created folder: %s", folderPath)
	return nil
}

func (m *MinIO) MoveFile(oldFileID, newFileID string) error {
	if oldFileID == "" || newFileID == "" {
		return fmt.Errorf("file IDs cannot be empty")
	}

	oldObjectName := oldFileID + ".json"
	newObjectName := newFileID + ".json"

	ctx := context.Background()

	src := minio.CopySrcOptions{
		Bucket: m.Bucket,
		Object: oldObjectName,
	}

	dst := minio.CopyDestOptions{
		Bucket: m.Bucket,
		Object: newObjectName,
	}

	_, err := m.Client.CopyObject(ctx, dst, src)
	if err != nil {
		return fmt.Errorf("failed to copy file from '%s' to '%s': %w", oldFileID, newFileID, err)
	}

	err = m.Client.RemoveObject(ctx, m.Bucket, oldObjectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to remove old file '%s': %w", oldFileID, err)
	}

	log.Printf("Successfully moved file from %s to %s", oldFileID, newFileID)
	return nil
}

var Module = fx.Options(
	fx.Provide(NewMinIO),
)
