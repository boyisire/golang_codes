package netease

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"
)

var (
	_im    *im
	imOnce sync.Once
)

const (
	ImUrl = "https://api.netease.im"
	//用户token
	TokenAction      = ImUrl + "/nimserver/user/create.action"       //创建token
	RefreTokenAction = ImUrl + "/nimserver/user/refreshToken.action" //刷新token
	//群组
	MuteTeamAction       = ImUrl + "/nimserver/team/muteTlistAll.action" //禁言群组
	MuteTeamMemberAction = ImUrl + "/nimserver/team/muteTlist.action"    //禁言群成员
	CreateTeamAction     = ImUrl + "/nimserver/team/create.action"       //创建群组
	MemberTeamsAction    = ImUrl + "/nimserver/team/joinTeams.action"    //查询用户是否加入群组
	JoinTeamAction       = ImUrl + "/nimserver/team/add.action"          //用户加入群组
	//聊天室
	CreateChatroomAction = ImUrl + "/nimserver/chatroom/create.action"      //创建聊天室
	MuteRoomAction       = ImUrl + "/nimserver/chatroom/muteRoom.action"    //禁言聊天室
	ChatroomUrlAction    = ImUrl + "/nimserver/chatroom/requestAddr.action" //请求聊天室地址
)

type im struct{}

//网易im sdk单例
func Im() *im {
	imOnce.Do(func() {
		_im = &im{}
	})
	return _im
}

type TokenResponse struct {
	Code int               `json:"code"`
	Info map[string]string `json:"info"`
	Desc string            `json:"desc"`
}

type CodeResponse struct {
	Code int `json:"code"`
}

type CreateTeamResponse struct {
	Code int    `json:"code"`
	Tid  string `json:"tid"`
}

type MemberTeamsResponse struct {
	Code  int     `json:"code"`
	Count int     `json:"count"`
	Infos []Infos `json:"infos"`
}

type Infos struct {
	Owner    string `json:"owner"`
	Tname    string `json:"tname"`
	Maxusers int    `json:"maxusers"`
	Tid      int    `json:"tid"`
	Size     int    `json:"size"`
	Custom   string `json:"custom"`
}

type CreateChatroomResponse struct {
	Code     int      `json:"code"`
	Chatroom Chatroom `json:"chatroom"`
}

type Chatroom struct {
	RoomId       int    `json:"roomid"`
	Valid        string `json:"valid"`
	Announcement string `json:"announcement"`
	Name         string `json:"name"`
	Broadcasturl string `json:"broadcasturl"`
	Ext          string `json:"ext"`
	Creator      string `json:"creator"`
}

type ChatroomUrlResponse struct {
	Code int      `json:"code"`
	Addr []string `json:"addr"`
}

//获取token
func (me *im) GetToken(accid, realname, portrait string) (token string, err error) {
	data := "accid=" + accid + "&name=" + realname + "&icon=" + portrait
	res, _ := HttpPost(TokenAction, string(data))
	var response TokenResponse
	json.Unmarshal(res, &response)
	if response.Code == 414 && response.Desc == "already register" {
		res, _ = HttpPost(RefreTokenAction, "accid="+accid)
		json.Unmarshal(res, &response)
	}
	if response.Code != 200 {
		return "", errors.New(response.Desc + "，code：" + strconv.Itoa(response.Code))
	}
	return response.Info["token"], nil
}

//禁言／解禁群组
func (me *im) MuteTeam(tid, owner, mute string) (code int, err error) {
	data := "tid=" + tid + "&owner=" + owner + "&mute=" + mute
	muteres, err := HttpPost(MuteTeamAction, string(data))
	var res CodeResponse
	json.Unmarshal(muteres, &res)
	return res.Code, nil
}

//禁言／解禁群成员
func (me *im) MuteTeamMember(tid, owner, accid, mute string) (code int, err error) {
	data := "tid=" + tid + "&owner=" + owner + "&accid=" + accid + "&mute=" + mute
	muteres, _ := HttpPost(MuteTeamMemberAction, string(data))
	var res CodeResponse
	json.Unmarshal(muteres, &res)
	return res.Code, nil
}

//创建群组
func (me *im) CreateTeam(tname, owner, members string) (tid string, err error) {
	data := "tname=" + tname + "&owner=" + owner + "&members=" + members + "&msg=邀请您加入班级群聊&magree=0&joinmode=0"
	team, _ := HttpPost(CreateTeamAction, string(data))
	var res CreateTeamResponse
	json.Unmarshal(team, &res)
	return res.Tid, nil
}

//判断是否群成员
func (me *im) IsTeamMember(accid string, tid int) (res bool) {
	team, _ := HttpPost(MemberTeamsAction, "accid="+accid)
	var jointeamres MemberTeamsResponse
	json.Unmarshal(team, &jointeamres)
	if jointeamres.Code == 200 {
		//用户在该群
		for _, v := range jointeamres.Infos {
			if v.Tid == tid {
				return true
			}
		}
	}
	return false
}

//用户加群
func (me *im) JoinTeam(tid, owner, accid string) (res bool) {
	member, _ := json.Marshal([]string{accid})
	data := "tid=" + tid + "&owner=" + owner + "&members=" + string(member) + "&msg=邀请您加入班级群聊&magree=0"
	jointeam, _ := HttpPost(JoinTeamAction, string(data))
	var join CodeResponse
	json.Unmarshal(jointeam, &join)
	if join.Code == 200 {
		return true
	}
	return false
}

//创建聊天室
func (me *im) CreateChatroom(creator, name string) (roomid int, err error) {
	data := "creator=" + creator + "&name=" + name
	chatroom, _ := HttpPost(CreateChatroomAction, string(data))
	var res CreateChatroomResponse
	json.Unmarshal(chatroom, &res)
	return res.Chatroom.RoomId, nil
}

//禁言或解禁
func (me *im) MuteRoom(roomID, accidID, isMute string) (code int, err error) {
	data := "roomid=" + roomID + "&operator=" + accidID + "&mute=" + isMute
	muteres, err := HttpPost(MuteRoomAction, string(data))
	var res CodeResponse
	json.Unmarshal(muteres, &res)
	return res.Code, nil
}

//获取聊天室地址
func (me *im) RequestChatroomAddr(roomID, accidID string) (res ChatroomUrlResponse) {
	data := "roomid=" + roomID + "&accid=" + accidID
	chaturl, _ := HttpPost(ChatroomUrlAction, string(data))
	json.Unmarshal(chaturl, &res)
	return
}
