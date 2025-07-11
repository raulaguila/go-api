package minio

import (
	"context"
	"fmt"
	"os"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/raulaguila/go-api/pkg/packhub"
)

func initBucket(client *minio.Client) error {
	ctx := context.Background()

	exist, err := client.BucketExists(ctx, os.Getenv("MINIO_BUCKET_FILES"))
	if err != nil {
		return err
	}

	if !exist {
		if err := client.MakeBucket(ctx, os.Getenv("MINIO_BUCKET_FILES"), minio.MakeBucketOptions{}); err != nil {
			return err
		}
	}

	return client.EnableVersioning(ctx, os.Getenv("MINIO_BUCKET_FILES"))
}

func ConnectMinio() *minio.Client {
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

	packhub.PanicIfErr(err)
	packhub.PanicIfErr(initBucket(minioClient))

	return minioClient
}
