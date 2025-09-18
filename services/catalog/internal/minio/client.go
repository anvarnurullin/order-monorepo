package minio

import (
	"context"
	"fmt"
	"io"
	"log"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type Client struct {
	client *minio.Client
	bucket string
}

func NewClient(endpoint, accessKey, secretKey, bucket string) (*Client, error) {
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	return &Client{
		client: client,
		bucket: bucket,
	}, nil
}

func (c *Client) InitBucket(ctx context.Context) error {
	exists, err := c.client.BucketExists(ctx, c.bucket)
	if err != nil {
		return fmt.Errorf("failed to check bucket existence: %w", err)
	}

	if !exists {
		err = c.client.MakeBucket(ctx, c.bucket, minio.MakeBucketOptions{})
		if err != nil {
			return fmt.Errorf("failed to create bucket: %w", err)
		}
		log.Printf("Bucket %s created successfully", c.bucket)
	}

	policy := fmt.Sprintf(`{
		"Version": "2012-10-17",
		"Statement": [
			{
				"Effect": "Allow",
				"Principal": {"AWS": "*"},
				"Action": ["s3:GetObject"],
				"Resource": ["arn:aws:s3:::%s/*"]
			}
		]
	}`, c.bucket)

	err = c.client.SetBucketPolicy(ctx, c.bucket, policy)
	if err != nil {
		return fmt.Errorf("failed to set bucket policy: %w", err)
	}

	log.Printf("Bucket policy set for %s", c.bucket)
	return nil
}

func (c *Client) UploadFile(ctx context.Context, objectName string, reader io.Reader, objectSize int64, contentType string) (string, error) {
	_, err := c.client.PutObject(ctx, c.bucket, objectName, reader, objectSize, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	return fmt.Sprintf("/product-images/%s", objectName), nil
}

func (c *Client) DeleteFile(ctx context.Context, objectName string) error {
	err := c.client.RemoveObject(ctx, c.bucket, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete file: %w", err)
	}
	return nil
}

func (c *Client) GetObject(ctx context.Context, objectName string) (*minio.Object, error) {
	return c.client.GetObject(ctx, c.bucket, objectName, minio.GetObjectOptions{})
}