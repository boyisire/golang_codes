package wechat

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"

	"17jzh.com/core/util"
)

type Subscribe struct {
	*Wechat
}

func NewSubscribe(wechat *Wechat) *Subscribe {
	return &Subscribe{wechat}
}

func (me *Subscribe) SendSubscribe(appid string, pushID, userID int, msg interface{}) (err error) {
	access_token, err := me.GetAccessToken(appid)
	if err != nil {
		util.FailOnError(err, "获取appid失败"+appid)
		return err
	}
	failedKey := "wechat:push:" + strconv.Itoa(pushID) + ":failed"
	for i := 0; i <= retryNum; i++ {
		resJson, err := util.PostJSON(API_Subscribe+"?access_token="+access_token, msg)
		if err != nil {
			if i < retryNum {
				continue
			}
			me.addPushErrRedisLog(failedKey, err, 9999, userID)
			return err
		}
		var res util.CommonError
		json.Unmarshal([]byte(resJson), &res)
		if res.ErrCode != 0 {
			err = errors.New(fmt.Sprintf("%d %s", res.ErrCode, res.ErrMsg))
			if res.ErrCode == 43004 {
				return err
			}
			if i < retryNum {
				continue
			}
			err := errors.New(strconv.FormatInt(res.ErrCode, 10))
			me.addPushErrRedisLog(failedKey, err, res.ErrCode, userID)
		}
		return err
	}
	return nil
}
