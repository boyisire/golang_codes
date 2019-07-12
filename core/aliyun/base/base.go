package base

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var (
	AccessKeyId  = os.Getenv("ALI_ACCESS_KEY_ID") //阿里accessKeyId
	AccessSecret = os.Getenv("ALI_ACCESS_SECRET") //阿里accessSecret
	OssBucket    = os.Getenv("ALI_OSS_BUCKET")    //阿里OssBucket 默认 17jzh-live
)

const (
	LIVE_OSS_ENDPOINT = "oss-cn-shanghai.aliyuncs.com"       //视频存储地址
	LIVE_OSS_DNS      = "http://storage-aliyun-sh.17jzh.com" //视频存储地址

)

//公共响应参数
type CommonResponse struct {
	Code int    `json:"code"` //200为正常
	Msg  string `json:"msg"`
}
