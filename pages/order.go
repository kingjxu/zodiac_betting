package pages

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
	"zodiac_betting/common"
	"zodiac_betting/dao"
	"zodiac_betting/proto"
	log "zodiac_betting/rolllog"
	"zodiac_betting/wx"
)

const Client_IP  = "81.71.85.201"
func CreateOrder(w http.ResponseWriter, req *http.Request) {
	startTime := time.Now()
	defer func() {
		log.Info("CreateOrder spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.CreateOrderRsp)

	req.ParseForm()
	log.Info("CreateOrder req:%#v",req.Form)
	orderId := fmt.Sprintf("%v",time.Now().UnixNano())
	openId,prodId,buyCnt := req.FormValue("openid"),req.FormValue("prod_id"),req.FormValue("buy_cnt")

	iProdId,_ := strconv.Atoi(prodId)
	iBuyCnt,_ := strconv.Atoi(buyCnt)

	if iBuyCnt <= 0 || openId == "" {
		rsp.Ret = common.REQ_PARAM_ERROR
		MakeRspData(w,rsp,"CreateOrder")
		return
	}

	prodPrice ,err :=dao.GetProductPrice(iProdId)
	if err != nil {
		log.Error("dao.GetProductPrice failed,prodId:%v,err:%v",iProdId,err)
		rsp.Ret = common.REQ_PARAM_ERROR
		MakeRspData(w,rsp,"CreateOrder")
		return
	}
	payPrice := prodPrice * iBuyCnt
	prepayId,err := wx.CreateWxOrder(orderId,Client_IP,openId,fmt.Sprintf("十二生肖-%v",dao.ZodiacIdName[iProdId]),payPrice)
	if err != nil {
		log.Error("wx.CreateWxOrder failed, err:%v\n",err)
		rsp.Ret = common.CREATE_WX_ORDER_FAILED
		MakeRspData(w,rsp,"CreateOrder")
		return
	}

	rsp.Ret = 0
	rsp.OrderId = orderId
	rsp.PrepayId = prepayId

	lRecord,ok,err := dao.GetLatestWinningRecord()
	if err != nil {
		log.Error("dao.GetLatestRecord failed,err:%v,ok:%v",err,ok)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"CreateOrder")
		return
	}

	buyRecord := new(dao.TUserBuyRecord)
	buyRecord.FTermId = lRecord.FTermId + 1
	buyRecord.FCreateTime = time.Now().Unix()
	buyRecord.FOrderId = orderId
	buyRecord.FOrderStatus = dao.Order_Status_To_Be_Pay
	buyRecord.FBuyAnimalId = iProdId
	buyRecord.FBuyAnimalName = dao.ZodiacIdName[iProdId]
	buyRecord.FOpenId = openId
	buyRecord.FBuyCnt = iBuyCnt
	buyRecord.FPayPrice = payPrice
	buyRecord.FWinningStatus = dao.Winning_Status_To_Be_Winning
	err = buyRecord.InsertOne()
	if err !=nil {
		log.Error("buyRecord.InsertOne failed,err:%v",err)
		rsp.Ret = common.DB_WRITE_ERROR
		MakeRspData(w,rsp,"CreateOrder")
		return
	}

	rsp.TermId = lRecord.FTermId+1
	MakeRspData(w,rsp,"CreateOrder")
}

//微信订单的发货回调
func Delivery(w http.ResponseWriter, req *http.Request) {
	reqBody,err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll failed, err:%v\n",err)
		w.Write([]byte("ok"))
		return
	}
	log.Info("Delivery body:%v",string(reqBody))
	rspNotify := new(wx.WxDeliveryNotifyRsp)
	reqNotify := new(wx.WxDeliveryNotifyReq)
	err = xml.Unmarshal([]byte(reqBody),reqNotify)
	if err != nil {
		log.Error("xml.Unmarshal failed,err:%v",err)
		rspNotify.ReturnCode = "FAIL"
		rspData := MakeWxRspData(rspNotify)
		w.Write(rspData)
		return
	}

	buyRecord,err := dao.GetUserBuyRecordByOrderId(reqNotify.Openid,reqNotify.OutTradeNo)
	if err != nil {
		log.Error("dao.GetUserBuyRecord failed,err:%v",err)
		rspNotify.ReturnCode = "FAIL"
		rspData := MakeWxRspData(rspNotify)
		log.Info("[Delivery]rspBody:%v",string(rspData))
		w.Write(rspData)
		return
	}

	if len(buyRecord.FWxTransactionId) > 0 {
		log.Warn("openid:%v,orderId:%v,wx_trans_id:%v payed.ignore",reqNotify.Openid,buyRecord.FOrderId,buyRecord.FWxTransactionId)
		rspNotify.ReturnCode = "SUCCESS"
		rspData := MakeWxRspData(rspNotify)
		log.Info("[Delivery]rspBody:%v",string(rspData))
		w.Write(rspData)
		return
	}

	buyRecord.FOrderStatus = dao.Order_Status_Payed
	buyRecord.FWxTransactionId = reqNotify.TransactionId
	err = buyRecord.UpdateBuyRecord()
	if err != nil {
		log.Error("buyRecord.UpdateBuyRecord failed,err:%v",err)
		rspNotify.ReturnCode = "FAIL"
		rspData := MakeWxRspData(rspNotify)
		w.Write(rspData)
		log.Info("[Delivery]rspBody:%v",string(rspData))
		return
	}
	rspNotify.ReturnCode = "SUCCESS"
	rspData := MakeWxRspData(rspNotify)
	w.Write(rspData)
	log.Info("[Delivery]rspBody:%v",string(rspData))
	return
}