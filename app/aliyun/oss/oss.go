package oss

import (
	"encoding/json"
	alioss "github.com/aliyun/aliyun-oss-go-sdk/oss"
	"log"
	"sync"
)

//SDK
//https://github.com/aliyun/aliyun-oss-go-sdk/blob/master/README-CN.md

type AliyunOss struct {
	Endpoint        string
	AccessKeyId     string
	AccessKeySecret string
	BucketName      string //yourBucketName
	BucketDomain    string //Bucket 域名
	UploadPath      string //bucket 前缀
}

var (
	client *alioss.Client
	once   sync.Once
	bucket *alioss.Bucket
)

func (c AliyunOss) InitOssClient() {
	// 创建OSSClient实例
	cli, err := alioss.New(c.Endpoint, c.AccessKeyId, c.AccessKeySecret)
	if err != nil {
		log.Fatal(err.Error())
	}
	client = cli
	// 获取存储空间
	buck, err := client.Bucket(c.BucketName)
	if err != nil {
		log.Fatal(err.Error())
	}
	bucket = buck
}

//上传本地文件
//localFileName:本地文件
//objectName:oss文件名称
func (c AliyunOss) UploadFile(localFileName string, objectName string) (string, error) {
	once.Do(c.InitOssClient)
	b, err := bucket.IsObjectExist(objectName)
	if b {
		log.Printf("改文件已存在")
		return objectName, nil
	}

	// 上传文件
	err = bucket.PutObjectFromFile(objectName, localFileName)
	if err != nil {
		return "", err
	}

	return objectName, nil
}

func (c AliyunOss) RemoveFile(objectNames []string) {
	once.Do(c.InitOssClient)

	result, err := bucket.DeleteObjects(objectNames)
	if err != nil {
		println(err.Error())
	}
	log.Print(json.Marshal(result))
}

//获取文件 url link
func (c AliyunOss) GetObjectFileUrl(objectFile string) (string, error) {
	once.Do(c.InitOssClient)
	signedUrl, err := bucket.SignURL(objectFile, alioss.HTTPGet, 60)
	if err == nil {
		log.Print("signedUrl:\t" + signedUrl)
		return signedUrl, nil
	}
	return "", err
}
