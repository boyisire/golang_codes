package netease

import (
	"encoding/json"
	"strings"
)

//公共响应参数
type commonResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

//创建频道响应
type CreateChannelResponse struct {
	commonResponse
	Ret channelRet `json:"ret"`
}

type channelRet struct {
	Cid         string `json:"cid"`
	Name        string `json:"name"`
	PushUrl     string `json:"pushUrl"`
	HttpPullUrl string `json:"httpPullUrl"`
	HlsPullUrl  string `json:"hlsPullUrl"`
	RtmpPullUrl string `json:"rtmpPullUrl"`
}

//设置录制状态响应
type SetRecordResponse struct {
	commonResponse
}

//设置录制回调地址响应
type SetCallbackResponse struct {
	Code int             `json:"code"`
	Ret  map[string]bool `json:"ret"`
}

//获取录制视频列表响应
type VideoListResponse struct {
	commonResponse
	Ret videoRet `json:"ret"`
}

type videoRet struct {
	Pnum         int     `json:"pnum"`         //当前页
	TotalRecords int     `json:"totalRecords"` //总记录数
	TotalPnum    int     `json:"totalPnum"`    //总页数
	Records      int     `json:"records"`      //单页记录数
	VideoList    []Video `json:"videoList"`
}

type Video struct {
	VideoName    string `json:"video_name"`
	OrigVideoKey string `json:"orig_video_key"`
	Vid          int64  `json:"vid"`
}

func (l *Video) UnmarshalJSON(j []byte) error {
	var vl map[string]interface{}

	err := json.Unmarshal(j, &vl)
	if err != nil {
		return err
	}
	for k, v := range vl {
		if strings.ToLower(k) == "video_name" {
			l.VideoName = v.(string)
		} else if strings.ToLower(k) == "vid" {
			l.Vid = int64(v.(float64))
		} else if strings.ToLower(k) == "orig_video_key" {
			l.OrigVideoKey = VideoHost + v.(string)
		}
	}
	return nil
}

//设置录制合并响应
type SetMergeResponse struct {
	Code int             `json:"code"`
	Ret  map[string]bool `json:"ret"`
}

type GetChannelStatusResponse struct {
	commonResponse
	Ret recodeStatus `json:"ret"`
}

type recodeStatus struct {
	Cid        string `json:"cid"`        //频道ID，32位字符串
	Status     int    `json:"status"`     //频道状态（0：空闲； 1：直播； 2：禁用； 3：直播录制）
	NeedRecord int    `json:"needRecord"` //1-开启录制； 0-关闭录制
	Format     int    `json:"format"`     //1-flv； 0-mp4
	Duration   int    `json:"duration"`   //录制切片时长(分钟)，默认120分钟
}

//设置视频截图
type VideoImageResponse struct {
	commonResponse
	Ret videoImage `json:"ret"`
}

type videoImage struct {
	Url    string `json:"url"`
	Width  string `json:"width"`
	Height int64  `json:"height"`
}
