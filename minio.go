package goxi_v2

import (
	"context"
	"fmt"
	"github.com/minio/madmin-go/v3"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
	"mime/multipart"
	"time"
)

type Minio struct {
	*minio.Client
	*madmin.AdminClient
	ctx context.Context
}

func NewMinio(endpoint, accessKeyID, secretAccessKey string, useSSL bool) (*Minio, error) {
	log.Printf("初始化连接Minio")
	ctx := context.Background()

	// Initialize minio client object.
	mc, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKeyID, secretAccessKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		log.Printf("初始化连接Minio失败:%s", err.Error())
		return nil, err
	}
	log.Printf("初始化连接Minio成功")
	return &Minio{
		Client: mc,
		ctx:    ctx,
	}, nil
}

// Bucket Operations

// CreateBucket Create a new bucket
func (m *Minio) CreateBucket(bucketName string) error {
	err := m.MakeBucket(m.ctx, bucketName, minio.MakeBucketOptions{
		Region: "us-east-1",
	})
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully created %s\n", bucketName)
	return nil
}

// ClearBucket Remove a bucket
func (m *Minio) ClearBucket(bucketName string) error {
	err := m.RemoveBucket(m.ctx, bucketName)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully removed %s\n", bucketName)
	return nil
}

// AllBuckets List all buckets
func (m *Minio) AllBuckets() ([]minio.BucketInfo, error) {
	buckets, err := m.ListBuckets(m.ctx)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	return buckets, nil
}

// BucketHasExists Check if a bucket exists
func (m *Minio) BucketHasExists(bucketName string) (bool, error) {
	found, err := m.BucketExists(m.ctx, bucketName)
	if err != nil {
		log.Fatal(err.Error())
		return false, err
	}
	return found, nil
}

// Policy

// SetUpBucketPolicy Set a bucket policy
func (m *Minio) SetUpBucketPolicy(bucketName string, policy string) error {
	err := m.SetBucketPolicy(m.ctx, bucketName, policy)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully set policy on %s\n", bucketName)
	return nil
}

// InquireBucketPolicy Get a bucket policy
func (m *Minio) InquireBucketPolicy(bucketName string) (string, error) {
	policy, err := m.GetBucketPolicy(m.ctx, bucketName)
	if err != nil {
		log.Fatal(err.Error())
		return "", err
	}
	return policy, nil
}

// File Object Operations

// GetObjectsByBucket Get all objects in a bucket
func (m *Minio) GetObjectsByBucket(bucketName string, prefix string, recursive bool) <-chan minio.ObjectInfo {
	//return m.ListObjects(m.ctx, bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: recursive})
	opts := minio.ListObjectsOptions{
		Recursive: recursive,
		Prefix:    prefix,
	}

	// List all objects from a bucket-name with a matching prefix.
	//for object := range m.ListObjects(context.Background(), bucketName, opts) {
	//	if object.Err != nil {
	//		fmt.Println(object.Err)
	//		return nil
	//	}
	//	fmt.Println(object)
	//}
	//return nil
	return m.ListObjects(context.Background(), bucketName, opts)
}

// GetObjectStat 获取文件对象详细
func (m *Minio) GetObjectStat(bucketName string, objectName string) minio.ObjectInfo {
	stat, err := m.StatObject(m.ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		log.Fatalln(err)
	}
	log.Println(stat)
	return stat
}

// UploadObject file to minio
func (m *Minio) UploadObject(bucketName string, objectName string, file *multipart.FileHeader, contentType string) error {
	objectSize := file.Size
	src, err := file.Open()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(src)
	info, err := m.PutObject(m.ctx, bucketName, objectName, src, objectSize, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully PutObject %s of size %d\n", objectName, info.Size)
	return nil
}

// UploadByFPutObject file to minio
func (m *Minio) UploadByFPutObject(bucketName string, objectName string, file *multipart.FileHeader, contentType string) error {
	src, err := file.Open()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer func(src multipart.File) {
		err = src.Close()
		if err != nil {
			log.Fatal(err.Error())
		}
	}(src)
	info, err := m.FPutObject(m.ctx, bucketName, objectName, file.Filename, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully PutObject %s of size %d\n", objectName, info.Size)
	return nil
}

// DownloadObject file from minio
func (m *Minio) DownloadObject(bucketName string, objectName string, filePath string) error {
	err := m.FGetObject(m.ctx, bucketName, objectName, filePath, minio.GetObjectOptions{})
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

// ClearObject file from minio
func (m *Minio) ClearObject(bucketName string, objectName string) error {
	opts := minio.RemoveObjectOptions{
		GovernanceBypass: true,
	}
	err := m.RemoveObject(m.ctx, bucketName, objectName, opts)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	log.Printf("Successfully removed %s\n", objectName)
	return nil
}

// ClearObjects file from minio
func (m *Minio) ClearObjects(bucketName string, prefix string) error {
	log.Printf("Removing all objects in %s with prefix %s\n", bucketName, prefix)
	objectsCh := make(chan minio.ObjectInfo)
	// Send object names that are needed to be removed to objectsCh
	go func() {
		defer close(objectsCh)
		// List all objects from a bucket-name with a matching prefix.
		for object := range m.ListObjects(m.ctx, bucketName, minio.ListObjectsOptions{Prefix: prefix, Recursive: true}) {
			if object.Err != nil {
				log.Fatalln(object.Err)
			}
			objectsCh <- object
		}
	}()

	opts := minio.RemoveObjectsOptions{
		GovernanceBypass: true,
	}

	for rErr := range m.RemoveObjects(m.ctx, bucketName, objectsCh, opts) {
		fmt.Println("Error detected during deletion: ", rErr)
	}
	return nil
}

// PreviewURL 预览文件
func (m *Minio) PreviewURL(bucketName string, objectName string, expires int64) (string, error) {
	url, err := m.PresignedGetObject(m.ctx, bucketName, objectName, time.Duration(expires)*time.Hour, nil)
	if err != nil {
		log.Fatal(err)
		return "", err
	}
	return url.String(), nil
}
