package wechat

import (
	"os"

	_ "github.com/joho/godotenv/autoload"
)

var API_URL = os.Getenv("API_URL")

//微信相关api
const (
	retryNum          = 1
	AccessTokenURL    = "https://api.weixin.qq.com/cgi-bin/token"                        //获取access_token的接口                                                     //接口重试次数
	API_Template      = "https://api.weixin.qq.com/cgi-bin/message/template/send"        //模板消息
	API_MINI_TEMPLATE = "https://api.weixin.qq.com/cgi-bin/message/wxopen/template/send" //小程序模板消息
	API_CS            = "https://api.weixin.qq.com/cgi-bin/message/custom/send"          //客服消息
	API_MEDIA         = "https://api.weixin.qq.com/cgi-bin/media/get"                    //图片|普通语音
	API_SPEEX         = "https://api.weixin.qq.com/cgi-bin/media/get/jssdk"              //高清语音
	API_Group         = "https://api.weixin.qq.com/cgi-bin/message/mass/send"            //群发接口
	API_Subscribe     = "https://api.weixin.qq.com/cgi-bin/message/template/subscribe"   //一次性订阅
)

//模板消息id
var WechatTemplate = map[int]map[string]string{
	4: {
		"startlive":         os.Getenv("TEMPLATE_Startlive_4"),         //课程开课通知
		"previewlive":       os.Getenv("TEMPLATE_Previewlive_4"),       //讲座预告
		"applynotice":       os.Getenv("TEMPLATE_Applynotice_4"),       //用户申请提醒
		"subscribe":         os.Getenv("TEMPLATE_Subscribe_4"),         //预约提醒
		"joinclass":         os.Getenv("TEMPLATE_Joinclass_4"),         //加入班级
		"payment":           os.Getenv("TEMPLATE_Payment_4"),           //成为会员
		"paymentcourse":     os.Getenv("TEMPLATE_PaymentCourse_4"),     //购买课程
		"livepreviewremind": os.Getenv("TEMPLATE_LivePreviewRemind_4"), //课程预习提醒
		"livenotice":        os.Getenv("TEMPLATE_Livenotice_4"),        //课程通知
		"liveremind":        os.Getenv("TEMPLATE_Liveremind_4"),        //课程提醒

	},
	5: {
		"startlive":         os.Getenv("TEMPLATE_Startlive_5"),         //课程开课通知
		"previewlive":       os.Getenv("TEMPLATE_Previewlive_5"),       //讲座预告
		"applynotice":       os.Getenv("TEMPLATE_Applynotice_5"),       //用户申请提醒
		"subscribe":         os.Getenv("TEMPLATE_Subscribe_5"),         //预约提醒
		"joinclass":         os.Getenv("TEMPLATE_Joinclass_5"),         //加入班级
		"payment":           os.Getenv("TEMPLATE_Payment_5"),           //成为会员
		"paymentcourse":     os.Getenv("TEMPLATE_PaymentCourse_5"),     //购买课程
		"livepreviewremind": os.Getenv("TEMPLATE_LivePreviewRemind_5"), //课程预习提醒
		"livenotice":        os.Getenv("TEMPLATE_Livenotice_5"),        //课程通知
		"liveremind":        os.Getenv("TEMPLATE_Liveremind_5"),        //课程提醒
	},
	6: {
		"startlive":         os.Getenv("TEMPLATE_Startlive_6"),         //课程开课通知
		"previewlive":       os.Getenv("TEMPLATE_Previewlive_6"),       //讲座预告
		"applynotice":       os.Getenv("TEMPLATE_Applynotice_6"),       //用户申请提醒
		"subscribe":         os.Getenv("TEMPLATE_Subscribe_6"),         //预约提醒
		"joinclass":         os.Getenv("TEMPLATE_Joinclass_6"),         //加入班级
		"payment":           os.Getenv("TEMPLATE_Payment_6"),           //成为会员
		"paymentcourse":     os.Getenv("TEMPLATE_PaymentCourse_6"),     //购买课程
		"livepreviewremind": os.Getenv("TEMPLATE_LivePreviewRemind_6"), //课程预习提醒
		"livenotice":        os.Getenv("TEMPLATE_Livenotice_6"),        //课程通知
		"liveremind":        os.Getenv("TEMPLATE_Liveremind_6"),        //课程提醒
	},
	8: {
		"startlive":         os.Getenv("TEMPLATE_Startlive_8"),         //课程开课通知
		"previewlive":       os.Getenv("TEMPLATE_Previewlive_8"),       //讲座预告
		"applynotice":       os.Getenv("TEMPLATE_Applynotice_8"),       //用户申请提醒
		"subscribe":         os.Getenv("TEMPLATE_Subscribe_8"),         //预约提醒
		"joinclass":         os.Getenv("TEMPLATE_Joinclass_8"),         //加入班级
		"payment":           os.Getenv("TEMPLATE_Payment_8"),           //成为会员
		"paymentcourse":     os.Getenv("TEMPLATE_PaymentCourse_8"),     //购买课程
		"livepreviewremind": os.Getenv("TEMPLATE_LivePreviewRemind_8"), //课程预习提醒
		"livenotice":        os.Getenv("TEMPLATE_Livenotice_8"),        //课程通知
		"liveremind":        os.Getenv("TEMPLATE_Liveremind_8"),        //课程提醒
	},
	14: {
		"liveremind": os.Getenv("TEMPLATE_Liveremind_14"), //课程提醒
	},
	15: {
		"startlive":         os.Getenv("TEMPLATE_Startlive_15"),         //课程开课通知
		"previewlive":       os.Getenv("TEMPLATE_Previewlive_15"),       //讲座预告
		"applynotice":       os.Getenv("TEMPLATE_Applynotice_15"),       //用户申请提醒
		"subscribe":         os.Getenv("TEMPLATE_Subscribe_15"),         //预约提醒
		"joinclass":         os.Getenv("TEMPLATE_Joinclass_15"),         //加入班级
		"payment":           os.Getenv("TEMPLATE_Payment_15"),           //成为会员
		"paymentcourse":     os.Getenv("TEMPLATE_PaymentCourse_15"),     //购买课程
		"livepreviewremind": os.Getenv("TEMPLATE_LivePreviewRemind_15"), //课程预习提醒
		"livenotice":        os.Getenv("TEMPLATE_Livenotice_15"),        //课程通知
		"liveremind":        os.Getenv("TEMPLATE_Liveremind_15"),        //课程提醒
	},
}
