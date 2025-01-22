package minio

import (
	"strings"

	"2501YTC/app/product/conf"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// MinioCreateTempDirServant 结构体定义
type MinioCreateTempDirServant struct {
	Client  *minio.Client
	Bucket  string
	Domain  string
	TempDir string
}

// MinioService 是全局 MinIO 服务实例
var MinioService *MinioCreateTempDirServant

// Init 创建并返回 MinIO 服务客户端实例
func Init() {
	// 从配置中获取 MinIO 配置信息
	address := conf.GetConf().Minio.Address
	accessKey := conf.GetConf().Minio.AccessKey
	secretKey := conf.GetConf().Minio.SecretKey
	secure := conf.GetConf().Minio.Secure
	bucket := conf.GetConf().Minio.Bucket
	domain := conf.GetConf().Minio.Domain
	tempDir := strings.Trim(conf.GetConf().Minio.TempDir, " /") + "/"

	// 初始化 MinIO 客户端对象
	client, err := minio.New(address, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: secure,
	})
	if err != nil {
		panic(err)
	}

	MinioService = &MinioCreateTempDirServant{
		Client:  client,
		Bucket:  bucket,
		Domain:  domain,
		TempDir: tempDir,
	}
}
