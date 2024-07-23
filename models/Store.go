package models

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type StoreInterface interface {
	List(path string)
	Upload(bucketName string, path string, file string, contentType string)
	Delete(bucketName string, path string)
	Download(bucketName string, path string) string
	Stream(bucketName string, path string) string
	Exist(bucketName string, path string)
}

type Store struct {
	EndPoint   string
	AccessId   string
	AccessPass string
	UseSSL     bool
}

func (m *Store) List(path string) {
	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	objectCh := minioClient.ListObjects(ctx, path, minio.ListObjectsOptions{
		Prefix:    "",
		Recursive: true,
	})
	for object := range objectCh {
		if object.Err != nil {
			fmt.Println(object.Err)
		}
		fmt.Println(object)
	}
}

func (m *Store) Upload(bucketName string, path string, file string, contentType string) (ErrorFound bool) {
	errFound := false
	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		errFound = true
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	info, err := minioClient.FPutObject(ctx, bucketName, path, file, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		errFound = true
		fmt.Println(err)
	}
	_ = info
	// fmt.Println(info)
	return errFound
}

func (m *Store) Delete(bucketName string, path string) {
	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err = minioClient.RemoveObject(ctx, bucketName, path, opts)
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (m *Store) Download(bucketName string, path string) string {
	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		fmt.Println(err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	expiry := time.Second * 24 * 60 * 60 // 1 day.
	reqParams := make(url.Values)
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, path, expiry, reqParams)
	if err != nil {
		fmt.Println(err)
	}
	return presignedURL.String()
}

func (m *Store) Stream(bucketName string, path string) (*minio.Object, bool) {
	errFound := false
	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		errFound = true
		fmt.Println(err)
	}

	object, err := minioClient.GetObject(context.Background(), bucketName, path, minio.GetObjectOptions{})
	if err != nil {
		errFound = true
		fmt.Println(err)
	}
	return object, errFound
}

func (m *Store) Exist(bucketName string, path string) (minio.ObjectInfo, bool) {
	errFound := false

	minioClient, err := minio.New(m.EndPoint, &minio.Options{Creds: credentials.NewStaticV4(m.AccessId, m.AccessPass, ""), Secure: m.UseSSL})
	if err != nil {
		errFound = true
		fmt.Println(err)
	}

	object, oErr := minioClient.StatObject(context.Background(), bucketName, path, minio.StatObjectOptions{})
	if oErr != nil {
		if oErr.(minio.ErrorResponse).Code == "NoSuchKey" {
			return minio.ObjectInfo{}, false
		}
		errFound = true
		fmt.Println(oErr)
	}

	return object, errFound
}
