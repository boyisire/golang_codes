package live

import (
	"fmt"
	"os"
	"sync"

	"17jzh.com/core/aliyun/base"
	"17jzh.com/core/util"

	log "github.com/Sirupsen/logrus"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/live"
)

var (
	_aliVideo       *aliVideo
	aliVideoOnce    sync.Once
	pushDomain      = os.Getenv("ALI_LIVE_PUSH_DOMAIN")   //直播推流域名 如 live-aliyun.nxdev.cn
	pullDomain      = os.Getenv("ALI_LIVE_PULL_DOMAIN")   //直播拉流域名 如 live-aliyun-player.nxdev.cn
	livePushMainKey = os.Getenv("ALI_LIVE_PUSH_MAIN_KEY") //鉴权推流主key
	livePullMainKey = os.Getenv("ALI_LIVE_PULL_MAIN_KEY") //鉴权拉流主key
	liveAppName     = os.Getenv("ALI_LIVE_APP_NAME")      //直播流所属应用名称
	//sharpness       = os.Getenv("ALI_LIVE_SHARPNESS")     //直播流畅度_lld 流畅 _lsd 标清 _lhd高清

)

type aliVideo struct{}

//创建频道响应
type CreateChannelResponse struct {
	base.CommonResponse
	Data channelData `json:"data"`
}

type channelData struct {
	Cid         string `json:"cid"`
	Name        string `json:"name"`
	PushUrl     string `json:"pushUrl"`
	HttpPullUrl string `json:"httpPullUrl"`
	HlsPullUrl  string `json:"hlsPullUrl"`
	RtmpPullUrl string `json:"rtmpPullUrl"`
}

//阿里云直播单例
func AliVideo() *aliVideo {
	aliVideoOnce.Do(func() {
		_aliVideo = &aliVideo{}
	})
	return _aliVideo
}

//生成鉴权key
func authKey(uri string, stopTime string, mainKey string) (authkey string) {
	nonce := util.RandStringBytesMaskImprSrc(10)
	str := uri + "-" + stopTime + "-" + nonce + "-0-" + mainKey
	sign := util.Md5Sum(str)
	authkey = stopTime + "-" + nonce + "-0-" + sign
	return
}

//创建频道
func (me *aliVideo) CreateChannel(streamName string, stopTime string, sharpness string) (res CreateChannelResponse) {
	uri := "/" + liveAppName + "/" + streamName
	res.Code = 200
	res.Msg = ""
	res.Data.Name = uri
	res.Data.PushUrl = "rtmp://" + pushDomain + uri + "?auth_key=" + authKey(uri, stopTime, livePushMainKey)
	uri = uri + sharpness
	res.Data.HttpPullUrl = "https://" + pullDomain + uri + ".flv?auth_key=" + authKey(uri+".flv", stopTime, livePullMainKey)
	res.Data.HlsPullUrl = "https://" + pullDomain + uri + ".m3u8?auth_key=" + authKey(uri+".m3u8", stopTime, livePullMainKey)
	res.Data.RtmpPullUrl = "rtmp://" + pullDomain + uri + "?auth_key=" + authKey(uri, stopTime, livePullMainKey)
	return res
}

//添加录制配置
func (me *aliVideo) AddLiveAppRecordConfig() error {
	//仅需要在上线之前手动调用一次 创建录制模板
	return nil
	client, err := live.NewClientWithAccessKey("cn-hangzhou", base.AccessKeyId, base.AccessSecret)

	request := live.CreateAddLiveAppRecordConfigRequest()

	request.OssBucket = base.OssBucket
	request.OssEndpoint = base.LIVE_OSS_ENDPOINT
	request.DomainName = pullDomain
	request.OnDemand = "7" //7为推流不自动录制，程序调用
	request.AppName = liveAppName
	request.RecordFormat = &[]live.AddLiveAppRecordConfigRecordFormat{
		{
			Format:               "mp4",
			OssObjectPrefix:      "nxdev/record/{AppName}/{StreamName}/{Sequence}{EscapedStartTime}{EscapedEndTime}",
			SliceOssObjectPrefix: "nxdev/record/{AppName}/{StreamName}/{UnixTimestamp}_{Sequence}",
			CycleDuration:        "10800",
		},
		{
			Format:               "m3u8",
			OssObjectPrefix:      "nxdev/record/{AppName}/{StreamName}/{Sequence}{EscapedStartTime}{EscapedEndTime}",
			SliceOssObjectPrefix: "nxdev/record/{AppName}/{StreamName}/{UnixTimestamp}_{Sequence}",
			CycleDuration:        "10800",
		},
	}

	response, err := client.AddLiveAppRecordConfig(request)
	if err != nil {
		return err
	}
	log.Info("添加录制视频配置 " + fmt.Sprintf("response is %#v\n", response))
	return nil
}

//控制视频录制 command=start 开始 stop 停止
func (me *aliVideo) RealTimeRecordCommand(command, stream string) error {
	client, err := live.NewClientWithAccessKey("cn-hangzhou", base.AccessKeyId, base.AccessSecret)

	request := live.CreateRealTimeRecordCommandRequest()

	request.Command = command
	request.StreamName = stream
	request.AppName = liveAppName
	request.DomainName = pullDomain

	response, err := client.RealTimeRecordCommand(request)
	if err != nil {
		return err
	}
	log.Info("控制视频录制 " + fmt.Sprintf("response is %#v\n", response))
	return nil
}

/*
func main() {
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", "{your_access_key_id}", "{your_access_key_id}")
	if err != nil {
		panic(err)
	}

	request := requests.NewCommonRequest()
	request.Domain = "ecs.aliyuncs.com"
	request.Version = "2014-05-26"
	// 因为是RPC接口，因此需指定ApiName(Action)
	request.ApiName = "DescribeInstanceStatus"
	request.QueryParams["PageNumber"] = "1"
	request.QueryParams["PageSize"] = "30"
	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		panic(err)
	}
	fmt.Print(response.GetHttpContentString())
}
*/
