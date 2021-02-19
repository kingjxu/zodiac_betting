package wx

import (
	log "zodiac_betting/rolllog"
	"zodiac_betting/third_part/wxpay"
)

const (
	WXPAY_SEND_REDPACK_URL = "https://api.mch.weixin.qq.com/mmpaymkttransfers/sendredpack"
	WXPAY_CERT_FILE        = "../conf/cert/apiclient_cert.pem"
	WXPAY_KEY_FILE         = "../conf/cert/apiclient_key.pem"
	WXPAY_ROOTCA_FILE      = "../conf/cert/rootca.pem"
)

type SendWxRedpackReq struct {
	NonceStr    string
	Sign        string
	MchBillno   string
	MchId       string
	WxAppId     string
	SendName    string
	ReOpenid    string
	TotalAmount int
	TotalNum    int
	Wishing     string
	ClientIp    string
	ActName     string
	Remark      string
}
type SendWxRedpackRsp struct {
	ReturnCode  string
	ReturnMsg   string
	ResultCode  string
	ErrCode     string
	ErrCodeDes  string
	MchBillno   string
	MchId       string
	Wxappid     string
	ReOpenid    string
	TotalAmount int64
	SendListid  string
}

func SendWxRedpack(req *SendWxRedpackReq) (*SendWxRedpackRsp, error) {
	client := wxpay.NewClient(WXPAY_APPID, WXPAY_MCH_ID, WXPAY_CREATE_ORDER_APPKEY)
	err := client.WithCert(WXPAY_CERT_FILE, WXPAY_KEY_FILE, WXPAY_ROOTCA_FILE)
	if err != nil {
		log.Error("client.WithCert failed,err:%v", err)
		return nil, err
	}

	params := make(wxpay.Params)
	params.SetString("nonce_str", req.NonceStr)
	params.SetString("mch_billno", req.MchBillno)
	params.SetString("mch_id", req.MchId)
	params.SetString("wxappid", req.WxAppId)
	params.SetString("send_name", req.SendName)
	params.SetString("re_openid", req.ReOpenid)
	params.SetInt64("total_amount", int64(req.TotalAmount))
	params.SetInt64("total_num", int64(req.TotalNum))
	params.SetString("wishing", req.Wishing)
	params.SetString("client_ip", req.ClientIp)
	params.SetString("act_name", req.ActName)
	params.SetString("remark", req.Remark)
	params.SetString("sign", client.Sign(params))

	ret, err := client.Post(WXPAY_SEND_REDPACK_URL, params, true)
	if err != nil {
		log.Error("client.Post failed,param:%+v,err:%v", params, err)
	}
	log.Info("openid:%v,amount:%v,post SendWxRedpack ret:%+v", ret)
	rsp := new(SendWxRedpackRsp)
	rsp.ReturnCode = ret.GetString("return_code")
	rsp.ReturnMsg = ret.GetString("return_msg")

	if rsp.ReturnCode == "SUCCESS" {
		rsp.ResultCode = ret.GetString("result_code")
		rsp.ErrCode = ret.GetString("err_code")
		rsp.ErrCodeDes = ret.GetString("err_code_des")

		if rsp.ResultCode == "SUCCESS" {
			rsp.MchBillno = ret.GetString("mch_billno")
			rsp.MchId = ret.GetString("mch_id")
			rsp.Wxappid = ret.GetString("wxappid")
			rsp.ReOpenid = ret.GetString("re_openid")
			rsp.TotalAmount = ret.GetInt64("total_amount")
			rsp.SendListid = ret.GetString("send_listid")
		}

	}
	return rsp, nil
}
