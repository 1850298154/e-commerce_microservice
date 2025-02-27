package img

import (
	"bytes"
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

func ConvertAndCompressImage(ctx context.Context, imgData []byte, fileName string, dstPath string) error {
	// 确保目标目录存在
	dir := "data"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		// 创建目录
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return fmt.Errorf("创建目录失败: %w", err)
		}
	}
	tmpFilePath := fmt.Sprintf("data/%s", fileName)
	dst, err := os.Create(tmpFilePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %w", err)
	}
	defer func() {
		err := dst.Close()
		if err != nil {
			klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
		}
		if err := os.Remove(tmpFilePath); err != nil {
			klog.CtxErrorf(ctx, "删除临时文件失败: %s", err)
		}
	}()

	// 将字节切片写入文件
	_, err = io.Copy(dst, bytes.NewReader(imgData))
	if err != nil {
		return fmt.Errorf("写入文件失败: %w", err)
	}
	srcFile, err := os.Open(tmpFilePath)
	if err != nil {
		return err
	}
	defer func() {
		err := srcFile.Close()
		if err != nil {
			klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
		}
	}()
	// 解码图像
	srcImg, _, err := image.Decode(srcFile)
	if err != nil {
		return fmt.Errorf("failed to decode image: %w", err)
	}

	f, err := os.Create(dstPath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			klog.CtxErrorf(ctx, "关闭文件失败: %s", err)
		}
	}(f)

	// 压缩并保存图像为 JPEG
	if err := jpeg.Encode(f, srcImg, &jpeg.Options{Quality: 100}); err != nil {
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
