package minioutils

import (
	"context"
	"fmt"
	"io"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"

	"github.com/raulaguila/go-api/pkg/helper"
)

// Minio provides a wrapper for interacting with the MinIO object storage service using a client instance.
type Minio struct {
	client *minio.Client
}

// NewMinioClient initializes and returns a new Minio client configured using environment variables.
// It connects to the Minio server specified by MINIO_HOST and MINIO_API_PORT, using credentials from MINIO_USER and MINIO_PASS.
func NewMinioClient() *Minio {
	minioClient, err := minio.New(
		fmt.Sprintf("%v:%v", os.Getenv("MINIO_HOST"), os.Getenv("MINIO_API_PORT")),
		&minio.Options{
			Creds: credentials.NewStaticV4(
				os.Getenv("MINIO_USER"),
				os.Getenv("MINIO_PASS"),
				"",
			),
			Secure: false,
		},
	)
	helper.PanicIfErr(err)

	return &Minio{
		client: minioClient,
	}
}

// UpdatePolicies updates the bucket policy to allow 's3:GetObject' access for specified resources in the given bucket.
// It constructs an access policy allowing public read permissions and sets it for the specified resources.
func (s *Minio) UpdatePolicies(ctx context.Context, bucketName string, resources ...string) error {
	var resourcesString string

	for i, resource := range resources {
		if !strings.HasPrefix(resource, "/") {
			resource = "/" + resource
		}
		if i > 0 {
			resourcesString += ","
		}
		resourcesString += `"arn:aws:s3:::` + bucketName + resource + `"`
	}

	policy := `{"Version": "2012-10-17","Statement": [{"Action": ["s3:GetObject"],"Effect": "Allow","Principal": {"AWS": ["*"]},"Resource": [` + resourcesString + `],"Sid": ""}]}`
	return s.client.SetBucketPolicy(ctx, bucketName, policy)
}

// InitBucket creates a bucket if it doesn't exist, enables versioning, and updates its access policies.
func (s *Minio) InitBucket(ctx context.Context, bucketName string, resources ...string) error {
	exist, err := s.client.BucketExists(ctx, bucketName)
	if err != nil {
		return err
	}

	if !exist {
		if err := s.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}

	if err := s.client.EnableVersioning(ctx, bucketName); err != nil {
		return err
	}

	return s.UpdatePolicies(ctx, bucketName, resources...)
}

// PutObject uploads an object to the specified bucket in MinIO using the provided reader and object name. It returns an error if the operation fails.
func (s *Minio) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader) error {
	_, err := s.client.PutObject(ctx, bucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})

	return err
}

// GenerateObjectURL generates a presigned URL for accessing a specified object in a Minio bucket.
// It takes the context, bucket name, object name, and original name as parameters.
// The presigned URL includes a content disposition response header to suggest a file download with the given original name.
// The URL is valid for 60 minutes after generation.
// Returns the presigned URL as a string, along with an error if any occurs during URL generation.
func (s *Minio) GenerateObjectURL(ctx context.Context, bucketName, objectName, originalName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+originalName+"\"")

	fileURL, err := s.client.PresignedGetObject(ctx, bucketName, objectName, 60*time.Minute, reqParams)
	if err != nil {
		return "", err
	}

	return fileURL.String(), nil
}

// GetObject retrieves an object from a specified bucket in Minio storage. Returns a pointer to the object and any error encountered.
func (s *Minio) GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

// DeleteObject removes an object from a specified bucket in MinIO.
// It returns an error if the object cannot be removed.
// Parameters include the context, bucket name, and object name to be deleted.
func (s *Minio) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	return s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
