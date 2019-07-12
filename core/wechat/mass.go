package wechat

import (
	"encoding/json"

	"core/util"
)

type Mass struct {
	*Wechat
}

func NewMass(wechat *Wechat) *Mass {
	return &Mass{wechat}
}

/**
 * 群发消息
 * @param {[type]} msg interface{} [description]
 */
func (me *Mass) Send(appid string, msg interface{}) {
	access_token, err := me.GetAccessToken(appid)
	if err != nil {
		return
	}
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_Group+"?access_token="+access_token, msg)
		if err != nil {
			if i < retryNum {
				continue
			}
			return
		}

		var res util.CommonError
		json.Unmarshal([]byte(resJson), &res)
		if res.ErrCode != 0 {
			//require subscribe
			if res.ErrCode == 43004 {
				return
			}
			if i < retryNum {
				continue
			}
		}
		return
	}
}
