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

type Minio struct {
	client *minio.Client
}

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

func (s *Minio) PutObject(ctx context.Context, bucketName, objectName string, reader io.Reader) error {
	_, err := s.client.PutObject(ctx, bucketName, objectName, reader, -1, minio.PutObjectOptions{
		ContentType: "image/jpeg",
	})

	return err
}

func (s *Minio) GenerateObjectURL(ctx context.Context, bucketName, objectName, originalName string) (string, error) {
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\""+originalName+"\"")

	fileURL, err := s.client.PresignedGetObject(ctx, bucketName, objectName, 60*time.Minute, reqParams)
	if err != nil {
		return "", err
	}

	return fileURL.String(), nil
}

func (s *Minio) GetObject(ctx context.Context, bucketName, objectName string) (*minio.Object, error) {
	return s.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
}

func (s *Minio) DeleteObject(ctx context.Context, bucketName, objectName string) error {
	return s.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
}
