package db

import (
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

var MinioClient *minio.Client

func InitMinio() {
	endpoint := "192.168.31.75:9000"
	accessKeyID := "ROOTUSER"
	secretAccessKey := "CHANGEME123"
	useSSL := false
	var err error
	// Initialize minio client object.
	MinioClient, err = minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("%#v\n", MinioClient) // minioClient is now set up
	buckets, err := MinioClient.ListBuckets(context.Background())
	if err != nil {
		log.Fatalf("Error listing buckets: %v", err)
	}

	// Print the list of buckets
	fmt.Println("Buckets:")
	for _, bucket := range buckets {
		fmt.Println(bucket.Name)
	}
}
