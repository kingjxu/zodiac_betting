package pages

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gogoods/x/random"
	"net/http"
	"strings"
	"time"
	"zodiac_betting/common"
	"zodiac_betting/conf"
	"zodiac_betting/dao"
	"zodiac_betting/proto"
	log "zodiac_betting/rolllog"
	"zodiac_betting/wx"
)

const ImagePrefix = "http://duole.site/pic/"

func UserAccount(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("UserAccount spends:%v", time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	code := req.URL.Query().Get("code")
	log.Info("UserAccount code:%v", code)
	accessToken, err := wx.GetAccessToken(code)
	rsp := new(proto.UserAccountRsp)
	if err != nil {
		log.Error("wx.GetAccessToken failed,code:%v,err:%v", code, err)
		rsp.Ret = common.GET_ACCESSTOKEN_FAILED
		MakeRspData(w, rsp, "UserAccount")
		return
	}
	rsp.Ret = 0
	rsp.Openid = accessToken.Openid

	MakeRspData(w, rsp, "UserAccount")
}

func GetTicket(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("GetTicket spends:%v", time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")

	rsp := new(proto.WxTicketRsp)
	rsp.Ret = 0
	rsp.Ticket = wx.GetTicket()
	rsp.HB = fmt.Sprintf("%v%v", ImagePrefix, "hb.png")

	MakeRspData(w, rsp, "GetTicket")
}

func GetMarqueeInfo(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("GetMarqueeInfo spends:%v", time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.MarqueeInfoRsp)
	rec, ok, err := dao.GetLatestWinningRecord()
	if !ok || err != nil {
		log.Error("dao.GetLatestRecord failed,ok:%v,err:%v", ok, err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w, rsp, "GetMarqueeInfo")
		return
	}
	rsp.Ret = 0
	rsp.LastTerm = rec.FTermId
	rsp.LastLotteryTime = rec.FCreateTime
	rsp.LastAnimal = rec.FWinningAnimalName
	rsp.NextLotteryLeft = rec.FCreateTime + int64(conf.ConfItems.Common.LotteryInterval) - time.Now().Unix()
	rsp.WinningInfos = genWinnerUserInfo()
	MakeRspData(w, rsp, "GetMarqueeInfo")
	return
}
func genWinnerUserInfo() []string {
	var sliceUser []string
	len := len(dao.UserName)
	for i := 0; i < 5; i++ {
		index := random.Range(0, len)
		betting := random.Range(10, 21)
		sliceUser = append(sliceUser, fmt.Sprintf("%v 中奖 %v 注", dao.UserName[index], betting))
	}
	return sliceUser
}
func GetProductInfo(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("GetProductInfo spends:%v", time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")

	rsp := new(proto.GetProductInfoRsp)
	sliceZodiac, err := dao.GetZodiacInfos()
	if err != nil {
		log.Error("dao.GetZodiacInfos failed,err:%v", err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w, rsp, "GetProductInfo")
		return
	}
	rsp.Ret = 0
	for _, zodiac := range sliceZodiac {
		slicePic := strings.Split(zodiac.FPic, ".")
		rsp.ProdInfos = append(rsp.ProdInfos, &proto.ProductInfo{
			Id:            zodiac.FZodiacId,
			Name:          zodiac.FName,
			Price:         zodiac.FPrice,
			Image:         fmt.Sprintf("%v%v", ImagePrefix, zodiac.FPic),
			ImageSelected: fmt.Sprintf("%v%v_selected.%v", ImagePrefix, slicePic[0], slicePic[1]),
		})
	}

	MakeRspData(w, rsp, "GetProductInfo")
}

func GetClientInfo(req *http.Request) (string, string) {
	openId := req.Header.Get("openid")
	clientAddr := req.RemoteAddr
	return openId, clientAddr
}

func MakeRspData(w http.ResponseWriter, inter interface{}, cmd string) []byte {
	rsp, err := json.Marshal(inter)
	if err != nil {
		log.Error("json.Marshal() failed,data:%+v,err:%v", inter, err)
		return nil
	}
	log.Debug("[%v] rspBody:%v", cmd, string(rsp))
	w.Write(rsp)
	return rsp
}

func MakeWxRspData(inter interface{}) []byte {
	rspData, err := xml.MarshalIndent(inter, "", "	")
	if err != nil {
		log.Error("xml.MarshalIndent failed,err:%v\n", err)
		return nil
	}
	return rspData
}
