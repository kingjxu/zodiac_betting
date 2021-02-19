package wx

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	rlog "zodiac_betting/rolllog"
	util "zodiac_betting/util"
)

const (
	WXPAY_CREATE_ORDER_URL    = "https://api.mch.weixin.qq.com/pay/unifiedorder"
	WXPAY_NOTIFY_URL          = "http://duole.site/order/delivery"
	WXPAY_CREATE_ORDER_APPKEY = "0123456789abcdefghijklmnopqrstuv"
	WXPAY_MCH_ID              = "1603344511"
)

type CreateWXOrderReq struct {
	XMLName        xml.Name `xml:"xml"`
	AppId          string   `xml:"appid"`
	Attach         string   `xml:"attach"`
	MchId          string   `xml:"mch_id"`
	NonceStr       string   `xml:"nonce_str"`
	Sign           string   `xml:"sign"`
	Body           string   `xml:"body"`
	OutTradeNo     string   `xml:"out_trade_no"`
	TotalFee       int      `xml:"total_fee"`
	SpbillCreateIp string   `xml:"spbill_create_ip"`
	NotifyUrl      string   `xml:"notify_url"`
	TradeType      string   `xml:"trade_type"`
	OpenId         string   `xml:"openid"`
}

type CreateWXOrderRsp struct {
	ReturnCode string `xml:"return_code"`
	ReturnMsg  string `xml:"return_msg"`
	Appid      string `xml:"appid"`
	MchId      string `xml:"mch_id"`
	NonceStr   string `xml:"nonce_str"`
	Openid     string `xml:"openid"`
	Sign       string `xml:"sign"`
	ResultCode string `xml:"result_code"`
	PrepayId   string `xml:"prepay_id"`
	TradeType  string `xml:"trade_type"`
}

type WxDeliveryNotifyReq struct {
	Appid         string `xml:"appid"`
	Attach        string `xml:"attach"`
	BankType      string `xml:"bank_type"`
	FeeType       string `xml:"fee_type"`
	MchId         string `xml:"mch_id"`
	NonceStr      string `xml:"nonce_str"`
	Openid        string `xml:"openid"`
	OutTradeNo    string `xml:"out_trade_no"`
	ResultCode    string `xml:"result_code"`
	ReturnCode    string `xml:"return_code"`
	TotalFee      string `xml:"total_fee"`
	TradeType     string `xml:"trade_type"`
	TransactionId string `xml:"transaction_id"`
}
type WxDeliveryNotifyRsp struct {
	XMLName    xml.Name `xml:"xml"`
	ReturnCode string   `xml:"return_code"`
	ReturnMsg  string   `xml:"return_msg"`
}

func CreateWxOrder(orderId, clientIp, openId, prodName string, price int) (prepayId string, err error) {
	rlog.Info("orderId:%v,ip:%v,openId:%v,price:%v", orderId, clientIp, openId, price)
	req := makeCreateOrderReq(orderId, clientIp, openId, prodName, price)
	reqBody, err := xml.MarshalIndent(req, "", "	")
	if err != nil {
		rlog.Error("xml.MarshalIndent failed,err:%v\n", err)
		return
	}
	rlog.Info("createOrder reqBody: %v\n", string(reqBody))
	rspBody, err := httpPostCreateOrder(string(reqBody))
	if err != nil {
		rlog.Error("httpPostCreateOrder failed,err:%v\n", err)
		return
	}
	rlog.Info("rspBody:%v\n", rspBody)
	rsp := new(CreateWXOrderRsp)
	err = xml.Unmarshal([]byte(rspBody), rsp)
	if err != nil {
		rlog.Error("xml.Unmarshal failed,err:%v", err)
		return
	}

	prepayId = rsp.PrepayId
	return
}

func httpPostCreateOrder(content string) (string, error) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", WXPAY_CREATE_ORDER_URL, strings.NewReader(content))
	if err != nil {
		return "", err
	}
	rsp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

func makeCreateOrderReq(orderId, clientIp, openId, prodName string, price int) *CreateWXOrderReq {
	req := new(CreateWXOrderReq)
	req.AppId = WXPAY_APPID
	req.Attach = "dajinye"
	req.MchId = WXPAY_MCH_ID
	req.NonceStr = util.CreateRandString(32)
	req.Body = prodName
	req.OutTradeNo = orderId
	req.TotalFee = price
	req.SpbillCreateIp = clientIp
	req.NotifyUrl = WXPAY_NOTIFY_URL
	req.TradeType = "JSAPI"
	req.OpenId = openId
	fillSign(req)
	return req
}
func fillSign(req *CreateWXOrderReq) {
	paramCombine := fmt.Sprintf("appid=%v&attach=%v&body=%v&mch_id=%v&nonce_str=%v&notify_url=%v&openid=%v&out_trade_no=%v&spbill_create_ip=%v&total_fee=%v&trade_type=%v", req.AppId, req.Attach, req.Body, req.MchId, req.NonceStr, req.NotifyUrl, req.OpenId, req.OutTradeNo, req.SpbillCreateIp, req.TotalFee, req.TradeType)
	signTmp := fmt.Sprintf("%v&key=%v", paramCombine, WXPAY_CREATE_ORDER_APPKEY)
	rlog.Info("signTmp:%v\n", signTmp)
	rlog.Info("md5V:%v\n", md5V(signTmp))
	req.Sign = strings.ToUpper(md5V(signTmp))
}

func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
