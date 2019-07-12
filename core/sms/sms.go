// 短信平台控制中心
package sms

import (
	"sync"

	alisms "core/aliyun/sms"
)

var (
	_Sms    *Sms
	SmsOnce sync.Once
)

type Sms struct {
	Type string //短信平台类型
}

//Sms单例
func NewSms() *Sms {
	SmsOnce.Do(func() {
		_Sms = &Sms{}
	})
	_Sms.Type = "aliyun"
	return _Sms
}

// 短信平台类型设置
func (me *Sms) SetType(smsType string) *Sms {
	me.Type = smsType
	return me
}

// 发送短信
func (me *Sms) Send(param map[string]interface{}) bool {
	if me.Type == "aliyun" {
		mobile := param["mobile"].(string)                           //手机号
		signName := param["sign_name"].(string)                      //签名
		templateCode := param["template_code"].(string)              //模板ID
		templateParam := param["template_param"].(map[string]string) //短信内容
		if mobile != "" && signName != "" && templateCode != "" {
			return alisms.AliSms().SendSms(signName, templateCode, mobile, templateParam)
		}
	}
	return false
}
