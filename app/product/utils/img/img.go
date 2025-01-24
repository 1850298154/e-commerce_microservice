package img

import (
	"context"
	"fmt"
	"image"
	"image/jpeg"
	"io"
	"os"
	"strings"
	"time"

	minioDal "2501YTC/app/product/biz/dal/minio"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/minio/minio-go/v7"
	uuid "github.com/satori/go.uuid"
)

func ConvertAndCompressImage(ctx context.Context, src io.Reader, dstPath string) error {
	srcImg, _, err := image.Decode(src)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	dstFile, err := os.Create(dstPath)
	if err != nil {
		return fmt.Errorf("failed to create destination file: %w", err)
	}
	defer func(dstFile *os.File) {
		err := dstFile.Close()
		if err != nil {
			klog.CtxErrorf(ctx, "failed to close destination file: %v", err)
		}
	}(dstFile)

	if err := jpeg.Encode(dstFile, srcImg, &jpeg.Options{Quality: 100}); err != nil {
		return fmt.Errorf("failed to encode JPEG: %w", err)
	}

	return nil
}

func GenerateObjectKey(uploadType string, fileExt string) string {
	return fmt.Sprintf("%s/%d/%s%s", uploadType, time.Now().Year(), uuid.NewV1().String(), fileExt)
}

// ms 是全局的 MinioService 实例
var ms = &minioDal.MinioService

// PutObject 用于上传对象
func PutObject(objectKey string, reader io.Reader, size int64, contentType string) (string, error) {
	opts := minio.PutObjectOptions{ContentType: contentType}
	_, err := (*ms).Client.PutObject(context.Background(), (*ms).Bucket, objectKey, reader, size, opts)
	if err != nil {
		return "", err
	}
	return (*ms).Domain + (*ms).Bucket + "/" + objectKey, nil
}

func GetObjectKeyFromUrl(fullUrl string) (objectKey string, ok bool) {
	objectKey = strings.TrimPrefix(fullUrl, (*ms).Domain+(*ms).Bucket+"/")
	if objectKey == fullUrl {
		return "", false
	}
	return objectKey, true
}

// DeleteObject 用于删除相应对象
func DeleteObject(objectKey string) error {
	err := (*ms).Client.RemoveObject(
		context.Background(),
		(*ms).Bucket,
		objectKey,
		minio.RemoveObjectOptions{ForceDelete: true},
	)
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}
	return nil
}

// DeleteObjectByUrlAsync 通过给定的 Url 异步删除对象
func DeleteObjectByUrlAsync(ctx context.Context, url string) {
	objectKey, ok := GetObjectKeyFromUrl(url)
	if ok {
		go func(objectKey string) {
			err := DeleteObject(objectKey)
			if err != nil {
				klog.CtxErrorf(ctx, "failed to delete object: %v", err)
			}
		}(objectKey)
	}
}
