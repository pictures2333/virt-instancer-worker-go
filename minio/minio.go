package minio

import (
	"Instancer-worker-go/config"
	"context"
	"io"
	"log"
	"os"
	"sync"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

var (
	client *minio.Client
	once   sync.Once
)

func Init() {
	once.Do(func() {
		var (
			err    error
			useSSL bool = false
		)

		// create client
		client, err = minio.New(config.MinIOEndpoint, &minio.Options{
			Creds:  credentials.NewStaticV4(config.MinIOAccesskey, config.MinIOSecretkey, ""),
			Secure: useSSL,
		})
		if err != nil {
			log.Fatalf("Failed to connect MinIO : %v", err)
		}

		// success
		log.Printf("MinIO connected")
	})
}

func Download(
	bucketName string, objectName string, // src
	dstPath string,
) (err error) {
	ctx := context.Background()

	// get object
	var obj *minio.Object
	obj, err = client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return err
	}
	defer obj.Close()

	// download
	var dst *os.File
	dst, err = os.Create(dstPath)
	if err != nil {
		return err
	}
	defer dst.Close()

	if _, err = io.Copy(dst, obj); err != nil {
		return err
	}

	return nil
}
