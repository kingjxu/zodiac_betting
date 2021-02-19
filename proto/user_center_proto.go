package proto

import "zodiac_betting/dao"

type GetBuyRecordReq struct {
	Start int `json:"start"`
	PageSize int `json:"page_size"`
}

type GetBuyRecordRsp struct {
	Ret int `json:"ret"`
	PageCnt int `json:"page_cnt"`
	Records []*dao.UserBuyRecordAndWinnerAnimal `json:"records"`
}

type WinningRecordsReq struct {
	Start int `json:"start"`
	PageSize int `json:"page_size"`
}
type WinningRecordsRsp struct {
	Ret int `json:"ret"`
	PageCnt int `json:"page_cnt"`
	Records []*dao.TUserBuyRecord `json:"records"`
}

type GetUserInfoRsp struct {
	Ret int `json:"ret"`
	UserInfo *dao.TUserInfo `json:"user_info"`
}

type SetUserInfoReq struct {
	NickName string `json:"nick_name"`
	Phone string `json:"phone"`
	AvatarId string `json:"avatar_id"`
}

type SetUserInfoRsp struct {
	Ret int `json:"ret"`
}

type AvatarInfo struct {
	AvatarId string `json:"avatar_id"`
	AvatarUrl string `json:"avatar_url"`
}
type GetAvatarListRsp struct {
	Ret int `json:"ret"`
	Avatar []*AvatarInfo `json:"avatar"`
}

type AboutUsRsp struct {
	Ret int `json:"ret"`
	AboutUs string `json:"about_us"`
}

type AboutActRsp struct {
	Ret int `json:"ret"`
	AboutAct string `json:"about_act"`
}

