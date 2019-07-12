package netease

import (
	"encoding/json"
	"sync"
)

var (
	_vcloud    *vcloud
	vcloudOnce sync.Once
)

const (
	VcloudHost             = "https://vcloud.163.com"
	CreateChannelAction    = VcloudHost + "/app/channel/create"          //创建频道
	SetAlwaysRecordAction  = VcloudHost + "/app/channel/setAlwaysRecord" //设置录制状态
	SetCallbackAction      = VcloudHost + "/app/record/setcallback"      //设置录制视频回调地址
	GetVideoListAction     = VcloudHost + "/app/videolist"               //视频列表
	SetVideoMergeAction    = VcloudHost + "/app/video/merge"             //视频合并
	CallbackUrlQueryAction = VcloudHost + "/app/record/callbackQuery"    //查询回调地址
	GetChannelStatusAction = VcloudHost + "/app/channelstats"            //获取频道状态
	CreateVideoImageAction = VcloudHost + "/app/vod/snapshot/create"     //获取视频截图
)

type vcloud struct{}

//网易云sdk单例
func Vcloud() *vcloud {
	vcloudOnce.Do(func() {
		_vcloud = &vcloud{}
	})
	return _vcloud
}

//创建频道
func (me *vcloud) CreateChannel(title string) (res CreateChannelResponse, err error) {
	jsonData, _ := json.Marshal(map[string]string{"name": title})
	channel, err := HttpJson(CreateChannelAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(channel, &res)
	return
}

//获取频道状态
func (me *vcloud) GetChannelStatus(cid string) (res GetChannelStatusResponse, err error) {
	data := map[string]interface{}{"cid": cid}
	jsonData, _ := json.Marshal(data)
	record, err := HttpJson(GetChannelStatusAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(record, &res)
	if err != nil {
		return
	}
	return
}

//设置频道为录制状态
func (me *vcloud) SetRecord(cid string, needRecord, format, duration int) (res SetRecordResponse, err error) {
	data := map[string]interface{}{"cid": cid, "needRecord": needRecord, "format": format, "duration": duration}
	jsonData, _ := json.Marshal(data)
	record, err := HttpJson(SetAlwaysRecordAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(record, &res)
	if err != nil {
		return
	}
	return
}

//设置录制回调
func (me *vcloud) SetRecordCallback(url string) (res SetCallbackResponse, err error) {
	data := map[string]interface{}{"recordClk": url}
	jsonData, _ := json.Marshal(data)
	recordClk, err := HttpJson(SetCallbackAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(recordClk, &res)
	return
}

//获取录制视频文件列表
func (me *vcloud) GetVideoList(cid string, pnum, pagesize int) (res VideoListResponse, err error) {
	data := map[string]interface{}{"cid": cid, "records": pagesize, "pnum": pnum}
	jsonData, _ := json.Marshal(data)
	record, err := HttpJson(GetVideoListAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(record, &res)
	return
}

//录制文件合并
func (me *vcloud) SetVideoMerge(outputName string, vidList []int64) (res SetMergeResponse, err error) {
	data := map[string]interface{}{"outputName": outputName, "vidList": vidList}
	jsonData, _ := json.Marshal(data)
	recordClk, err := HttpJson(SetVideoMergeAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(recordClk, &res)
	return
}

//设置视频截图
func (me *vcloud) SetVideoImage(vid int64, size, offset int) (res VideoImageResponse, err error) {
	data := map[string]interface{}{"vid": vid, "size": size, "offset": offset}
	jsonData, _ := json.Marshal(data)
	recordClk, err := HttpJson(CreateVideoImageAction, jsonData)
	if err != nil {
		return
	}
	err = json.Unmarshal(recordClk, &res)
	return
}
