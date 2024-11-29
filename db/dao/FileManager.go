package dao

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"wxcloudrun-golang/db"
)

type FileManagerImpl struct {
	FileManager
}

func (FileManagerImpl) GetFile(fileName string, fileType string) (*minio.Object, error) {
	if fileType == "question-voice" {
		object, err := db.MinioClient.GetObject(context.Background(), "question-voice",
			fileName, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return object, nil
	}
	return nil, nil
}

func (FileManagerImpl) UploadFile(fileName string, fileType string, voice []byte) error {
	background := context.Background()
	if fileType == "question-voice" {
		exists, err := db.MinioClient.BucketExists(background, "question-voice")
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !exists {
			err := db.MinioClient.MakeBucket(background, "question-voice", minio.MakeBucketOptions{})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
		buf := bytes.NewReader(voice)

		_, err = db.MinioClient.PutObject(background, "question-voice", "/"+fileName[0:2]+"/"+fileName,
			buf, int64(len(voice)), minio.PutObjectOptions{ContentType: "audio/mpeg"})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	}
	return nil
}

func (FileManagerImpl) DeleteFile(fileName string, fileType string) error {
	if fileType == "question-voice" {
		db.MinioClient.RemoveObject(context.Background(), "question-voice", fileName[0:2]+"/"+fileName, minio.RemoveObjectOptions{})
		return nil
	}
	return nil
}
