package dao

import (
	"bytes"
	"context"
	"fmt"
	"github.com/minio/minio-go/v7"
	"mime/multipart"
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
	} else if fileType == "avatar" {
		object, err := db.MinioClient.GetObject(context.Background(), "avatar", fileName, minio.GetObjectOptions{})
		if err != nil {
			fmt.Println(err.Error())
			return nil, err
		}
		return object, nil
	}
	return nil, nil
}

func (FileManagerImpl) UploadFile(fileName string, fileType string, file []byte) error {
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
		buf := bytes.NewReader(file)
		fmt.Println(fileName)

		_, err = db.MinioClient.PutObject(background, "question-voice", "/"+fileName[0:2]+"/"+fileName,
			buf, int64(len(file)), minio.PutObjectOptions{ContentType: "audio/mpeg"})
		if err != nil {
			fmt.Println(err.Error())
			return err
		}
		return nil
	} else if fileType == "avatar" {

	}

	return nil
}

func (FileManagerImpl) UploadFileByMultipartFile(fileName string, fileType string, file multipart.File) error {
	if fileType == "avatar" {
		background := context.Background()

		exists, err := db.MinioClient.BucketExists(background, "avatar")
		if err != nil {
			fmt.Println(err)
			return err
		}
		if !exists {
			err := db.MinioClient.MakeBucket(background, "avatar", minio.MakeBucketOptions{})
			if err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
		fmt.Println(fileName)
		_, err = db.MinioClient.PutObject(background, "avatar", "/"+fileName[0:2]+"/"+fileName, file, -1, minio.PutObjectOptions{})

	}
	return nil
}

func (FileManagerImpl) DeleteFile(fileName string, fileType string) error {
	if fileType == "question-voice" {
		db.MinioClient.RemoveObject(context.Background(), "question-voice", fileName[0:2]+"/"+fileName, minio.RemoveObjectOptions{})
		return nil
	} else if fileType == "avatar" {
		db.MinioClient.RemoveObject(context.Background(), "avatar", fileName[0:2]+"/"+fileName, minio.RemoveObjectOptions{})
		return nil
	}
	return nil
}
