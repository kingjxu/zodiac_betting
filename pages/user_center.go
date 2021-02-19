package pages

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"
	"zodiac_betting/common"
	"zodiac_betting/conf"
	"zodiac_betting/dao"
	"zodiac_betting/proto"
	log "zodiac_betting/rolllog"
)

func GetBuyRecord(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("GetBuyRecord spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	req.ParseForm()
	rsp := new(proto.GetBuyRecordRsp)
	log.Info("GetBuyRecord req:%#v",req.Form)
	err := req.ParseForm()
	if err != nil {
		log.Error("req.ParseForm failed, err:%v,\n",err)
	}

	openId := req.FormValue("openid")
	start := req.FormValue("start")
	pageSize := req.FormValue("page_size")
	iStart,_ := strconv.Atoi(start)
	iPageSize,_ := strconv.Atoi(pageSize)

	if iStart < 0 || openId == "" {
		rsp.Ret = common.REQ_PARAM_ERROR
		MakeRspData(w,rsp,"GetBuyRecord")
		return
	}

	cnt,err := dao.GetUserBuyCnt(openId)
	if err != nil {
		log.Error("dao.GetUserBuyCnt failed, err:%v\n",err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"GetBuyRecord")
		return
	}
	log.Debug("dao.GetUserBuyCnt cnt:%v",cnt)
	rsp.PageCnt = int(cnt) / iPageSize
	if int(cnt) % iPageSize  !=0 {
		rsp.PageCnt += 1
	}
	offset := iStart * iPageSize
	recs,err := dao.GetUserBuyRecord(openId,iPageSize,offset)
	if err != nil {
		log.Error("dao.GetUserBuyRecord failed, err:%v\n",err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"GetBuyRecord")
		return
	}
	log.Info("openid:%v,recCnt:%v",openId, len(recs))
	recsWin,err := convertRecords(recs)
	if err != nil {
		log.Error("dao.GetUserBuyRecord failed, err:%v\n",err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"GetBuyRecord")
		return
	}
	rsp.Ret = 0
	rsp.Records = recsWin
	MakeRspData(w,rsp,"GetBuyRecord")
	return
}
func convertRecords(recs []*dao.TUserBuyRecord) ([]*dao.UserBuyRecordAndWinnerAnimal,error)  {
	if len(recs) == 0 {

	}
	var sliceRec []*dao.UserBuyRecordAndWinnerAnimal
	var sliceTermId []int
	for _,rec := range recs {
		if rec.FWinningStatus == dao.Winning_Status_To_Be_Winning {
			continue
		}
		sliceTermId = append(sliceTermId,rec.FTermId)
	}

	winRecs,err := dao.GetWinningRecordByTerms(sliceTermId)
	if err != nil {
		log.Error("dao.GetWinningRecordByTerms failed,err:%v",err)
		return nil,err
	}

	for _,rec := range recs {
		byteData,err := json.Marshal(rec)
		if err != nil {
			log.Error("json.Marshal failed,err:%v")
			return nil,err
		}
		recWin := new(dao.UserBuyRecordAndWinnerAnimal)
		err = json.Unmarshal(byteData,recWin)
		if err != nil {
			log.Error("json.Unmarshal failed,err:%v")
			return nil,err
		}
		for _,wr := range winRecs {
			if wr.FTermId == rec.FTermId {
				recWin.FWinAnimalName = wr.FWinningAnimalName
				break
			}
		}
		sliceRec = append(sliceRec,recWin)
	}
	return sliceRec,nil
}
func WinningRecords(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("WinningRecords spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.WinningRecordsRsp)

	err := req.ParseForm()
	if err != nil {
		log.Error("req.ParseForm failed, err:%v,\n",err)
	}
	log.Info("GetBuyRecord req:%#v",req.Form)

	openId := req.FormValue("openid")
	start := req.FormValue("start")
	pageSize := req.FormValue("page_size")
	iStart,_ := strconv.Atoi(start)
	iPageSize,_ := strconv.Atoi(pageSize)
	offset := iStart * iPageSize

	if iStart < 0 || openId == "" {
		rsp.Ret = common.REQ_PARAM_ERROR
		MakeRspData(w,rsp,"WinningRecords")
		return
	}

	cnt,err := dao.GetUserWinningCnt(openId)
	if err != nil {
		log.Error("dao.GetUserBuyCnt failed, err:%v\n",err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"WinningRecords")
		return
	}
	rsp.PageCnt = int(cnt) / iPageSize
	if int(cnt) % iPageSize  !=0 {
		rsp.PageCnt += 1
	}

	recs,err := dao.GetUserWinningRecords(openId,iPageSize,offset)
	if err != nil {
		log.Error("dao.GetUserBuyRecord failed, err:%v\n",err)
		rsp.Ret = common.IO_READ_ERROR
		MakeRspData(w,rsp,"WinningRecords")
		return
	}
	rsp.Ret = 0
	rsp.Records = recs
	MakeRspData(w,rsp,"WinningRecords")
}

func BuildDefaultUserInfoRsp(rsp *proto.GetUserInfoRsp)  {

	rsp.UserInfo.FLevel = "一级赌童"
	rsp.UserInfo.FNickName = "请设置您的昵称"
	rsp.UserInfo.FAvatar = fmt.Sprintf("%v%v",ImagePrefix,"default.jpg")
	rsp.UserInfo.FPhone = ""
}
func GetUserInfo(w http.ResponseWriter, req *http.Request)  {
	startTime := time.Now()
	defer func() {
		log.Info("GetUserInfo spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	openId := req.URL.Query().Get("openid")
	rsp := new(proto.GetUserInfoRsp)
	rsp.UserInfo = new(dao.TUserInfo)
	userInfo,ok,err := dao.GetUserInfo(openId)
	if err != nil || !ok{
		log.Error("dao.GetUserInfo failed, ok:%v,err:%v\n",ok,err)
		rsp.Ret = 0
		BuildDefaultUserInfoRsp(rsp)
		MakeRspData(w,rsp,"GetUserInfo")
		return
	}
	userInfo.FAvatar = fmt.Sprintf("%v%v",ImagePrefix,userInfo.FAvatar)
	rsp.Ret = 0
	rsp.UserInfo = userInfo
	MakeRspData(w,rsp,"GetUserInfo")
}

func SetUserInfo(w http.ResponseWriter, req *http.Request)  {
	startTime := time.Now()
	defer func() {
		log.Info("SetUserInfo spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	req.ParseForm()
	log.Info("SetUserInfo req:%#v",req.Form)
	openId := req.FormValue("openid")
	nickName := req.FormValue("nick_name")
	phone := req.FormValue("phone")
	avatarId := req.FormValue("avatar_id")

	rsp := new(proto.SetUserInfoRsp)
	_,exist,err := dao.GetUserInfo(openId)
	if err != nil {
		log.Error("dao.GetUserInfo failed, err:%v,openid:%v\n",err,openId)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"SetUserInfo")
		return
	}
	userInfo :=&dao.TUserInfo{
		FPhone: phone,
		FAvatar:dao.AvatarIdName[avatarId],
		FNickName: nickName,
		FOpenId: openId,
	}
	if !exist {
		userInfo.FLevel = "一级赌童"
		rsp.Ret = 0
		err = userInfo.InsertOne()
		if err != nil {
			log.Error("userInfo.InsertOne failed, err:%v\n",err)
			rsp.Ret = common.DB_WRITE_ERROR
		}
		MakeRspData(w,rsp,"SetUserInfo")
		return
	}

	rsp.Ret = 0
	err = userInfo.UpdateOne()
	if err != nil {
		log.Error("userInfo.UpdateOne failed, err:%v,reqBody:%+v\n",err,req.Form)
		rsp.Ret = common.DB_WRITE_ERROR
	}
	MakeRspData(w,rsp,"SetUserInfo")
	return
}

func GetAvatarList(w http.ResponseWriter, req *http.Request)  {
	startTime := time.Now()
	defer func() {
		log.Info("GetAvatarList spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.GetAvatarListRsp)
	rsp.Ret = 0
	for k,v := range dao.AvatarIdName {
		rsp.Avatar = append(rsp.Avatar,&proto.AvatarInfo{
			AvatarId:k,
			AvatarUrl:fmt.Sprintf("%v%v",ImagePrefix,v),
		})
	}
	MakeRspData(w,rsp,"GetAvatarList")
}

func AboutUs(w http.ResponseWriter, req *http.Request)  {
	w.Header().Set("content-type", "application/json")
	rsp := new(proto.AboutUsRsp)
	rsp.Ret = 0
	rsp.AboutUs = conf.ConfItems.Common.AboutUs
	MakeRspData(w,rsp,"AboutUs")
}


func AboutAct(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("content-type", "application/json")
	rsp := new(proto.AboutActRsp)
	rsp.Ret = 0
	rsp.AboutAct = conf.ConfItems.Common.AboutAct
	MakeRspData(w,rsp,"AboutAct")
}