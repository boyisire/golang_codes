package sms

import (
	"encoding/json"
	"sync"

	"17jzh.com/core/aliyun/base"
	log "github.com/Sirupsen/logrus"

	"github.com/aliyun/alibaba-cloud-sdk-go/sdk"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
)

var (
	_aliSms    *aliSms
	aliSmsOnce sync.Once
)

type aliSms struct{}

//阿里云Sms单例
func AliSms() *aliSms {
	aliSmsOnce.Do(func() {
		_aliSms = &aliSms{}
	})
	return _aliSms
}

//创建频道
func (me *aliSms) SendSms(signName, templateCode, phone string, templateParam map[string]string) (success bool) {
	success = false
	client, err := sdk.NewClientWithAccessKey("cn-hangzhou", base.AccessKeyId, base.AccessSecret)
	if err != nil {
		log.WithFields(log.Fields{
			"PhoneNumbers":  phone,
			"SignName":      signName,
			"TemplateCode":  templateCode,
			"TemplateParam": templateParam,
			"err":           err,
		}).Error("Send SMS NewClientWithAccessKey error")
		return
	}

	request := requests.NewCommonRequest()
	request.Method = "POST"
	request.Scheme = "https" // https | http
	request.Domain = "dysmsapi.aliyuncs.com"
	request.Version = "2017-05-25"
	request.ApiName = "SendSms"

	request.QueryParams["PhoneNumbers"] = phone
	request.QueryParams["SignName"] = signName
	request.QueryParams["TemplateCode"] = templateCode
	if templateParam != nil {
		str, _ := json.Marshal(templateParam)
		request.QueryParams["TemplateParam"] = string(str)

	}

	response, err := client.ProcessCommonRequest(request)
	if err != nil {
		log.WithFields(log.Fields{
			"PhoneNumbers":  phone,
			"SignName":      signName,
			"TemplateCode":  templateCode,
			"TemplateParam": templateParam,
			"err":           err,
		}).Error("Send SMS error")
		return
	}
	if response.IsSuccess() {
		success = true
		return
	}
	log.WithFields(log.Fields{
		"PhoneNumbers":  phone,
		"SignName":      signName,
		"TemplateCode":  templateCode,
		"TemplateParam": templateParam,
		"response":      response.GetHttpContentString(),
	}).Error("Send SMS response false")
	return
}
