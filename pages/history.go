package pages

import (
	"net/http"
	"strconv"
	"time"
	"zodiac_betting/common"
	"zodiac_betting/dao"
	"zodiac_betting/proto"
	log "zodiac_betting/rolllog"
)

func LotteryHistory(w http.ResponseWriter, req *http.Request)  {
	startTime := time.Now()
	defer func() {
		log.Info("LotteryHistory spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.LotteryHistoryRsp)

	err := req.ParseForm()
	if err != nil {
		log.Error("req.ParseForm failed, err:%v,\n",err)
	}
	start := req.FormValue("start")
	pageSize := req.FormValue("page_size")

	iStart,_ := strconv.Atoi(start)
	iPageSize,_ := strconv.Atoi(pageSize)

	cnt,err := dao.GetLotteryCnt()
	if err != nil {
		log.Error("json.Unmarshal failed, err:%v,start:%v,pageSize\n",err,start,pageSize)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"LotteryHistory")
		return
	}
	rsp.PageCnt = int(cnt) / iPageSize
	if int(cnt) % iPageSize != 0 {
		rsp.PageCnt += 1
	}

	sliceLottery,err := dao.GetLotteryHistory(iStart * iPageSize ,iPageSize)
	if err != nil {
		log.Error("json.Unmarshal failed, err:%v,start:%v,pageSize\n",err,start,pageSize)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"LotteryHistory")
		return
	}
	rsp.Ret = 0
	rsp.History = sliceLottery
	MakeRspData(w,rsp,"LotteryHistory")
	return
}




func WinnerProbability(w http.ResponseWriter, req *http.Request)  {
	startTime := time.Now()
	defer func() {
		log.Info("WinnerProbability spends:%v",time.Since(startTime).Seconds())
	}()

	w.Header().Set("content-type", "application/json")
	rsp := new(proto.WinnerProbabilityRsp)
	zodiacFreq,err := dao.GetZodiacWinnerFreq()
	if err != nil {
		log.Error("dao.GetZodiacWinnerFreq failed, err:%v",err)
		rsp.Ret = common.DB_READ_ERROR
		MakeRspData(w,rsp,"WinnerProbability")
		return
	}
	rsp.Ret = 0
	totalCnt := 0
	for _,zf := range zodiacFreq {
		totalCnt += zf.Count
	}

	probCur := 0
	for index,zf := range zodiacFreq {

		rsp.Probs = append(rsp.Probs,&proto.ZodiacWinnerProbability{
			AnimalId:zf.AnimalId,
			AnimalName:dao.ZodiacIdName[zf.AnimalId],
			Prob:zf.Count * 10000 / totalCnt,
		})

		if index == len(zodiacFreq) - 1 {
			rsp.Probs[index].Prob = 10000 - probCur
		}
		probCur += zf.Count * 10000 / totalCnt
	}
	MakeRspData(w,rsp,"WinnerProbability")
}
