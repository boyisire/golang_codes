package wechat

import (
	"encoding/json"
	"errors"
	"strconv"

	"17jzh.com/core/util"
)

type Staff struct {
	*Wechat
}

func NewStaff(wechat *Wechat) *Staff {
	return &Staff{wechat}
}

/**
 * 发送客服消息
 * @param {[type]} msg interface{} [description]
 */
func (me *Staff) SendStaffSingle(appid string, msg interface{}) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return
	}
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_CS+"?access_token="+accessToken, msg)
		if err != nil {
			if i < retryNum {
				continue
			}
			util.FailOnError(err, "staff api reply timeout")
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
			util.FailOnError(errors.New(strconv.FormatInt(res.ErrCode, 10)), res.ErrMsg)
		}
		return
	}
}

/**
 * 群发【文本 | 图文】
 * @param {[type]} appid  		string  [description]
 * @param {[type]} pushID       int           [description]
 * @param {[type]} openid       string        [description]
 * @param {[type]} msgs         interface{} [description]
 */
func (me *Staff) SendStaffGroup(appid string, pushID int, userID int, msgs interface{}) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return
	}
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_CS+"?access_token="+accessToken, msgs)
		if err != nil {
			if i < retryNum {
				continue
			}
			me.addPushErrRedisLog("wechat:push:"+strconv.Itoa(pushID)+":failed", err, 9999, userID)
			return
		}

		var res util.CommonError
		json.Unmarshal([]byte(resJson), &res)
		if res.ErrCode != 0 {
			if res.ErrCode == 43004 {
				return
			}
			if i < retryNum {
				continue
			}
		}
		//群推图文成功写入集合
		me.addPushNewsSuccessLog("wechat:push:"+strconv.Itoa(pushID)+":success", userID)
		return
	}
}
