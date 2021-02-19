package pages

import (
	"io/ioutil"
	"net/http"
	log "zodiac_betting/rolllog"
)

func ReportLog(w http.ResponseWriter, req *http.Request)  {
	openId,_ := GetClientInfo(req)
	reqBody,err := ioutil.ReadAll(req.Body)
	if err != nil {
		log.Error("ioutil.ReadAll failed, err:%v\n",err)
	}
	log.Info("openid:%v,log:%+v",openId,string(reqBody))
	w.Write([]byte("ok"))
}
