package logic

import (
	"bytes"
	"context"
	"fmt"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"jijizhazha1024/go-mall/services/product/internal/config"
	"time"
)

func UploadImage(image []byte, zone *storage.Zone, config config.Config) (url string, err error) {
	accessKey := config.QiNiu.AccessKey
	secretKey := config.QiNiu.SecretKey
	bucket := config.QiNiu.Bucket
	domain := config.QiNiu.Domain // 七牛云存储空间绑定的域名
	// 2. 初始化七牛云认证信息
	mac := qbox.NewMac(accessKey, secretKey)
	putPolicy := storage.PutPolicy{
		Scope: bucket,
	}
	upToken := putPolicy.UploadToken(mac)

	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}

	formUploader := storage.NewFormUploader(&cfg)
	ret := storage.PutRet{}
	putExtra := storage.PutExtra{}
	// 生成一个唯一的文件名，这里简单使用时间戳
	filename := fmt.Sprintf("%d.jpg", time.Now().UnixNano())
	// 将 []byte 转换为 io.Reader
	reader := bytes.NewReader(image)
	err = formUploader.Put(context.Background(), &ret, upToken, filename, reader, int64(len(image)), &putExtra)
	if err != nil {
		return "", fmt.Errorf("上传到七牛云失败: %v", err)
	}
	// 3. 生成七牛云 URL
	return fmt.Sprintf("http://%s/%s", domain, ret.Key), nil
}
