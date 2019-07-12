package wechat

//log "github.com/Sirupsen/logrus"

type Message struct {
}

func NewMessage() *Message {
	return &Message{}
}

//客服文本消息
type CsTextMsg struct {
	Touser  string            `json:"touser"`
	Msgtype string            `json:"msgtype"`
	Text    map[string]string `json:"text"`
}

//客服图片消息
type CsImgMsg struct {
	Touser  string            `json:"touser"`
	Msgtype string            `json:"msgtype"`
	Image   map[string]string `json:"image"`
}

//客服图文消息
type CsNewsMsg struct {
	Touser  string                 `json:"touser"`
	Msgtype string                 `json:"msgtype"`
	News    map[string]interface{} `json:"news"`
}

//模板消息
type TemplateMsg struct {
	Touser      string      `json:"touser"`
	TemplateId  string      `json:"template_id"`
	Url         string      `json:"url"`
	Miniprogram interface{} `json:"miniprogram"`
	Data        interface{} `json:"data"`
}

//模板消息小程序参数
type Miniprogram struct {
	Appid    string `json:"appid"`
	Pagepath string `json:"pagepath"`
}

//小程序模板消息
type MiniTemplateMsg struct {
	Touser          string      `json:"touser"`
	TemplateId      string      `json:"template_id"`
	Page            string      `json:"page"`
	FormId          string      `json:"form_id"`
	Data            interface{} `json:"data"`
	EmphasisKeyword string      `json:"emphasis_keyword"`
}

//客服图文消息
type GroupPicMsg struct {
	Touser  []string          `json:"touser"`
	Image   map[string]string `json:"image"`
	Msgtype string            `json:"msgtype"`
}

type OnceSubscribeMsg struct {
	Touser      string      `json:"touser"`
	TemplateId  string      `json:"template_id"`
	Url         string      `json:"url"`
	Title       string      `json:"title"`
	Scene       int         `json:"scene"`
	Miniprogram interface{} `json:"miniprogram"`
	Data        interface{} `json:"data"`
}

//客服文本
func (me *Message) GetCsTextMsg(touser, content string) *CsTextMsg {
	m := new(CsTextMsg)

	m.Touser = touser
	m.Msgtype = "text"
	m.Text = map[string]string{
		"content": content,
	}
	return m
}

//客服图文
func (me *Message) GetCsNewsMsg(openid string, content []interface{}) *CsNewsMsg {
	m := new(CsNewsMsg)

	m.Touser = openid
	m.Msgtype = "news"
	m.News = map[string]interface{}{
		"articles": content,
	}
	return m
}

//群推图片
func (me *Message) GetGrouppicMsg(touser []string, mediaID string) *GroupPicMsg {
	m := new(GroupPicMsg)

	m.Touser = touser
	m.Msgtype = "image"
	m.Image = map[string]string{
		"media_id": mediaID,
	}
	return m
}

//推图片
func (me *Message) GetpicMsg(touser string, mediaID string) *CsImgMsg {
	m := new(CsImgMsg)

	m.Touser = touser
	m.Msgtype = "image"
	m.Image = map[string]string{
		"media_id": mediaID,
	}
	return m
}

//上课提醒 || 讲座预览 || 订单完成 || 加入班级完成 || 新入班申请applynotice
func (me *Message) GetTemplateMsg(acccunt_id int, openid, url, template_key string, content interface{}, miniprogram interface{}) *TemplateMsg {
	m := new(TemplateMsg)
	m.Touser = openid
	if _, ok := WechatTemplate[acccunt_id][template_key]; ok {
		m.TemplateId = WechatTemplate[acccunt_id][template_key]
	} else {
		m.TemplateId = template_key
	}
	m.Url = url
	if miniprogram != nil {
		m.Miniprogram = miniprogram
	}
	m.Data = content
	return m
}

//小程序模板消息
func (me *Message) GetMiniTemplateMsg(tid int, openid, page, formId, templateKey string, content interface{}) *MiniTemplateMsg {
	m := new(MiniTemplateMsg)
	m.Touser = openid
	m.TemplateId = WechatTemplate[tid][templateKey]
	m.Page = page
	m.FormId = formId
	m.Data = content
	m.EmphasisKeyword = "keyword2.DATA"
	return m
}

func (me *Message) GetOnceSubscribeMsg(openid, template_id, url, title string, scene int, content interface{}, miniprogram interface{}) *OnceSubscribeMsg {
	m := new(OnceSubscribeMsg)
	m.Touser = openid
	m.TemplateId = template_id
	m.Url = url
	if miniprogram != nil {
		m.Miniprogram = miniprogram
	}
	m.Title = title
	m.Scene = scene
	m.Data = content
	return m
}

type VType struct {
	Value string `json:"value"`
	Color string `json:"color"`
}

type SixKwData struct {
	Keyword1 VType `json:"keyword1"`
	Keyword2 VType `json:"keyword2"`
	Keyword3 VType `json:"keyword3"`
	Keyword4 VType `json:"keyword4"`
	Keyword5 VType `json:"keyword5"`
	Keyword6 VType `json:"keyword6"`
	First    VType `json:"first"`
	Remark   VType `json:"remark"`
}

type FiveKwData struct {
	Keyword1 VType `json:"keyword1"`
	Keyword2 VType `json:"keyword2"`
	Keyword3 VType `json:"keyword3"`
	Keyword4 VType `json:"keyword4"`
	Keyword5 VType `json:"keyword5"`
	First    VType `json:"first"`
	Remark   VType `json:"remark"`
}

type FourKwData struct {
	Keyword1 VType `json:"keyword1"`
	Keyword2 VType `json:"keyword2"`
	Keyword3 VType `json:"keyword3"`
	Keyword4 VType `json:"keyword4"`
	First    VType `json:"first"`
	Remark   VType `json:"remark"`
}

type ThreeKwData struct {
	Keyword1 VType `json:"keyword1"`
	Keyword2 VType `json:"keyword2"`
	Keyword3 VType `json:"keyword3"`
	First    VType `json:"first"`
	Remark   VType `json:"remark"`
}

type TwoKwData struct {
	Keyword1 VType `json:"keyword1"`
	Keyword2 VType `json:"keyword2"`
	First    VType `json:"first"`
	Remark   VType `json:"remark"`
}

type OnceSubscribeData struct {
	Content VType `json:"content"`
}
