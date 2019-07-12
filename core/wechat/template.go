package wechat

import (
	"encoding/json"
	"errors"
	"strconv"

	"17jzh.com/core/util"
)

type Template struct {
	*Wechat
}

func NewTemplate(wechat *Wechat) *Template {
	return &Template{wechat}
}

/**
 * 群发模板消息
 * @param {[type]} appid  string [description]
 * @param {[type]} pushID       int           [description]
 * @param {[type]} openid       string        [description]
 * @param {[type]} msgs         interface{} [description]
 */
func (me *Template) SendGroupTemplate(appid string, pushID int, userID int, msgs interface{}) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return
	}
	failedKey := "wechat:push:" + strconv.Itoa(pushID) + ":failed"
	// retry
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_Template+"?access_token="+accessToken, msgs)
		if err != nil {
			if i < retryNum {
				continue
			}
			me.addPushErrRedisLog(failedKey, err, 9999, userID)
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
			err := errors.New(strconv.FormatInt(res.ErrCode, 10))
			me.addPushErrRedisLog(failedKey, err, res.ErrCode, userID)
		}
		return
	}
}

/**
 * 单发模板消息
 * @param {[type]} msg interface{} [description]
 */
func (me *Template) SendTemplateSingle(appid string, msg interface{}) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return
	}
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_Template+"?access_token="+accessToken, msg)
		if err != nil {
			if i < retryNum {
				continue
			}
			util.FailOnError(err, "post weixin api err SendTemplateSingle")
			return
		}

		var res util.CommonError
		json.Unmarshal([]byte(resJson), &res)
		if res.ErrCode != 0 {
			if i < retryNum {
				continue
			}
			util.FailOnError(errors.New("post weixin api err SendTemplateSingle!"), res.ErrMsg)
		}
		return
	}
}

/**
 * 新入班申请提醒
 * @param {[type]} appid string        [description]
 * @param {[type]} msgs         interface{} [description]
 */
func (me *Template) SendApplyNotice(appid string, msgs interface{}) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return
	}
	// retry
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_Template+"?access_token="+accessToken, msgs)
		if err != nil {
			if i < retryNum {
				continue
			}
			util.FailOnError(err, "post weixin api err sendApplyNotice")
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
			util.FailOnError(errors.New(strconv.FormatInt(res.ErrCode, 10)), res.ErrMsg)
		}
		return
	}
}

func (me *Template) SendMiniTemplate(appid string, msgs interface{}) (err error) {
	accessToken, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return err
	}
	// retry
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_MINI_TEMPLATE+"?access_token="+accessToken, msgs)
		if err != nil {
			if i < retryNum {
				continue
			}
			util.FailOnError(err, "post weixin api err SendMiniTemplate")
			return err
		}

		var res util.CommonError
		json.Unmarshal([]byte(resJson), &res)
		if res.ErrCode != 0 {
			if res.ErrCode == 43004 {
				return err
			}
			if i < retryNum {
				continue
			}
			util.FailOnError(errors.New(strconv.FormatInt(res.ErrCode, 10)), res.ErrMsg)
		}
		return err
	}
	return nil

}
