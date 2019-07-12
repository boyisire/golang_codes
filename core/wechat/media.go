package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"

	"core/util"
)

type Media struct {
	*Wechat
}

func NewMedia(wechat *Wechat) *Media {
	return &Media{wechat}
}

/**
 * 微信多媒体下载
 * @param {[type]} media_id   [description]
 * @param {[type]} media_type string)       (data []byte, fileExt string, err error [description]
 */
func (me *Media) Download(appid, media_id, media_type string) (data []byte, fileExt string, err error) {
	var url string
	access_token, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return nil, "", err
	}
	// retry
	for i := 0; i <= retryNum; i++ {
		if media_type == "image" {
			url = API_MEDIA + "?access_token=" + access_token + "&media_id=" + media_id
		} else if media_type == "voice" {
			url = API_SPEEX + "?access_token=" + access_token + "&media_id=" + media_id
		} else {
			return nil, "", errors.New("media_type err : " + media_type)
		}
		resp, err := http.Get(url)
		if err != nil {
			if i < retryNum {
				continue
			}
			return nil, "", err
		}
		defer resp.Body.Close()
		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			if i < retryNum {
				continue
			}
			return nil, "", err
		}
		// json
		if resp.Header.Get("Content-Type") == "application/json; encoding=utf-8" {
			var rtn util.CommonError
			json.Unmarshal(data, &rtn)
			if i < retryNum {
				continue
			}
			err = errors.New(fmt.Sprintf("%d %s", rtn.ErrCode, rtn.ErrMsg))
			return nil, "", err
		}
		// media success
		fileExt = util.GetFileExt(resp.Header.Get("Content-disposition"))
		return data, fileExt, nil
	}
	return
}
