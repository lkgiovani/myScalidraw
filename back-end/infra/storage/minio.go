package storage

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
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
		Endpoint:  endpoint,
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		UseSSL:    useSSL,
	}
}

type MinIO struct {
	Client *minio.Client
	Bucket string
}

func NewMinIO(config MinIOConfig) (*MinIO, error) {

	client, err := minio.New(config.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(config.AccessKey, config.SecretKey, ""),
		Secure: config.UseSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	exists, err := client.BucketExists(context.Background(), config.Bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to check if bucket exists: %w", err)
	}

	if !exists {
		err = client.MakeBucket(context.Background(), config.Bucket, minio.MakeBucketOptions{})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket: %w", err)
		}

		policy := `{
			"Version": "2012-10-17",
			"Statement": [
				{
					"Effect": "Allow",
					"Principal": "*",
					"Action": ["s3:GetObject"],
					"Resource": ["arn:aws:s3:::` + config.Bucket + `/*"]
				}
			]
		}`
		err = client.SetBucketPolicy(context.Background(), config.Bucket, policy)
		if err != nil {
			return nil, fmt.Errorf("failed to set bucket policy: %w", err)
		}
	}

	log.Println("Connected to MinIO storage")
	return &MinIO{
		Client: client,
		Bucket: config.Bucket,
	}, nil
}

func (m *MinIO) UploadFile(fileID string, content []byte) (string, error) {
	contentType := "application/json"
	objectName := fileID + ".json"
	reader := bytes.NewReader(content)

	_, err := m.Client.PutObject(context.Background(), m.Bucket, objectName, reader, int64(len(content)),
		minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("http://%s/%s/%s", m.Client.EndpointURL().Host, m.Bucket, objectName), nil
}

func (m *MinIO) GetFile(fileID string) ([]byte, error) {
	objectName := fileID + ".json"

	object, err := m.Client.GetObject(context.Background(), m.Bucket, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}
	defer object.Close()

	content, err := io.ReadAll(object)
	if err != nil {
		return nil, fmt.Errorf("failed to read object: %w", err)
	}

	return content, nil
}

func (m *MinIO) DeleteFile(fileID string) error {
	objectName := fileID + ".json"

	err := m.Client.RemoveObject(context.Background(), m.Bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	return nil
}

func (m *MinIO) ListFiles() ([]string, error) {
	var files []string

	objectCh := m.Client.ListObjects(context.Background(), m.Bucket, minio.ListObjectsOptions{
		Recursive: true,
	})

	for object := range objectCh {
		if object.Err != nil {
			return nil, fmt.Errorf("error listing objects: %w", object.Err)
		}
		files = append(files, object.Key)
	}

	return files, nil
}

func (m *MinIO) GetFileURL(fileID string) string {
	objectName := fileID + ".json"

	presignedURL, err := m.Client.PresignedGetObject(context.Background(), m.Bucket, objectName, time.Hour*24*7, nil)
	if err != nil {
		return ""
	}

	return presignedURL.String()
}

var Module = fx.Options(
	fx.Provide(NewMinIO),
)
