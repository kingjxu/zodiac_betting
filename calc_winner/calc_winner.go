package calc_winner

import (
	"fmt"
	"github.com/gogoods/x/random"
	"time"
	"zodiac_betting/conf"
	"zodiac_betting/dao"
	log "zodiac_betting/rolllog"
	"zodiac_betting/util"
	"zodiac_betting/wx"
)

const Zodiac_Cnt = 12

func GetLatestTermBettingInfo() (int, []*dao.ZodiacGroupAmount, error) {
	winRec, _, err := dao.GetLatestWinningRecord()
	if err != nil {
		log.Error("dao.GetLatestRecord failed,err:%v", err)
		return 0, nil, err
	}
	thisTerm := winRec.FTermId + 1
	//统计下用户再该期次下购买的总额分布
	groupAmount, err := dao.GetZodiacGroupAmount(thisTerm)
	if err != nil {
		log.Error("dao.GetZodiacGroupAmount failed,err:%v", err)
		return 0, nil, err
	}
	for _, ga := range groupAmount {
		log.Info("termId:%v,dao.GetZodiacGroupAmount,animalId:%v,name:%v,amount:%v", thisTerm, ga.AnimalId, ga.AnimalName, ga.SumPrice)
	}
	return thisTerm, groupAmount, nil
}

//计算一下选定这个生肖中奖，这轮收益多少
func calcAmount(groupAmount []*dao.ZodiacGroupAmount, animalId int) int {
	totalAmount := 0
	for _, ga := range groupAmount {
		if ga.AnimalId != animalId {
			totalAmount += ga.SumPrice
			continue
		}

		totalAmount = totalAmount - ga.SumPrice*conf.ConfItems.Common.WinMulti
	}
	return totalAmount
}

func UpdateWinningInfo(winRecord *dao.TWinningRecord) (err error) {
	err = winRecord.InsertOne()
	if err != nil {
		log.Error("winRecord.InsertOne failed,err:%v", err)
		return
	}

	sliceBuyRecord, err := dao.GetBuyRecordByTermId(winRecord.FTermId)
	if err != nil {
		log.Error("dao.GetBuyRecordByTermId failed,err:%v", err)
		return
	}

	for _, rec := range sliceBuyRecord {

		rec.FWinningStatus = dao.Winning_Status_Not_Winner
		if rec.FBuyAnimalId == winRecord.FWinningAnimalId {
			rec.FWinningStatus = dao.Winning_Status_Winner
		}
		err = rec.UpdateBuyRecord()
		if err != nil {
			log.Error("rec.UpdateBuyRecord failed,rec:%+v,err:%v", rec, err)
			continue
		}
	}
	return nil
}

/*该版本是能够获得最大收益的版本
*买的最少的或者没买的生肖会是中奖的生肖
*就是说再该版本下，我们绝对是不会亏钱的
 */
func CalcWinnerV1() (int, *dao.TWinningRecord, error) {
	thisTerm, groupAmount, err := GetLatestTermBettingInfo()
	if err != nil {
		log.Error("GetLatestTermBettingInfo failed,err:%v", err)
		return 0, nil, err
	}
	winRecord := new(dao.TWinningRecord)
	winRecord.FTermId = thisTerm
	winRecord.FCreateTime = time.Now().Unix()

	if len(groupAmount) == Zodiac_Cnt { //每个生肖都有选，就选最后一个投注少的
		winRecord.FWinningAnimalId = groupAmount[Zodiac_Cnt-1].AnimalId
	} else {
		for zid := range dao.ZodiacIdName {
			zidNotIn := true
			for _, ga := range groupAmount {
				if zid == ga.AnimalId {
					zidNotIn = false
					break
				}
			}
			if zidNotIn {
				winRecord.FWinningAnimalId = zid
				break
			}
		}
	}
	winRecord.FWinningAnimalName = dao.ZodiacIdName[winRecord.FWinningAnimalId]
	winAmount := calcAmount(groupAmount, winRecord.FWinningAnimalId)
	log.Info("[CalcWinnerV1]termId:%v,animal:%v,winAmount:%v", thisTerm, winRecord.FWinningAnimalName, winAmount)
	return winAmount, winRecord, nil
}

/*该版本保证每次都有人中奖
*且中奖的是买的最少的，
*但是有可能会导致亏损
 */
func CalcWinnerV2() (int, *dao.TWinningRecord, error) {
	thisTerm, groupAmount, err := GetLatestTermBettingInfo()
	if err != nil {
		log.Error("GetLatestTermBettingInfo failed,err:%v", err)
		return 0, nil, err
	}
	winRecord := new(dao.TWinningRecord)
	winRecord.FTermId = thisTerm
	winRecord.FCreateTime = time.Now().Unix()

	if len(groupAmount) == 0 { //没有任何投注,随便选一个
		for z := range dao.ZodiacIdName {
			winRecord.FWinningAnimalId = z
			break
		}
	} else {
		winRecord.FWinningAnimalId = groupAmount[len(groupAmount)-1].AnimalId
	}
	winRecord.FWinningAnimalName = dao.ZodiacIdName[winRecord.FWinningAnimalId]
	winAmount := calcAmount(groupAmount, winRecord.FWinningAnimalId)
	log.Info("[CalcWinnerV2]termId:%v,animal:%v,winAmount:%v", thisTerm, winRecord.FWinningAnimalName, winAmount)

	return winAmount, winRecord, nil
}

/*
*该版本为纯随机的算法
 */
func CalcWinnerV3() (int, *dao.TWinningRecord, error) {
	thisTerm, groupAmount, err := GetLatestTermBettingInfo()
	if err != nil {
		log.Error("GetLatestTermBettingInfo failed,err:%v", err)
		return 0, nil, err
	}
	winRecord := new(dao.TWinningRecord)
	winRecord.FTermId = thisTerm
	winRecord.FCreateTime = time.Now().Unix()

	winRecord.FWinningAnimalId = random.Range(0, 12) + 1000
	winRecord.FWinningAnimalName = dao.ZodiacIdName[winRecord.FWinningAnimalId]
	winAmount := calcAmount(groupAmount, winRecord.FWinningAnimalId)
	log.Info("[CalcWinnerV3]termId:%v,animal:%v,winAmount:%v", thisTerm, winRecord.FWinningAnimalName, winAmount)

	return winAmount, winRecord, nil
}

func CalcWinner() {
	amount, winRec := 0, &dao.TWinningRecord{}
	var err error
	algVer := "v3"
	if random.Range(0, 100) < conf.ConfItems.Common.ProbAlgV2 {
		amount, winRec, err = CalcWinnerV3()
		if err != nil { //用V3这个算法有问题
			log.Error("termId:%v,CalcWinnerV3 failed,err:%v will use CalcWinnerV1", winRec.FTermId, err)
			amount, winRec, err = CalcWinnerV1()
			algVer = "v1"
		}
	} else {
		amount, winRec, err = CalcWinnerV1()
		algVer = "v1"
	}

	if err != nil {
		log.Error("[%v]termId:%v,CalcWinner failed,err:%v", algVer, winRec.FTermId, err)
		return
	}
	err = UpdateWinningInfo(winRec)
	if err != nil {
		log.Error("termId:%v,CalcWinnerV2 UpdateWinningInfo failed,err:%v", winRec.FTermId, err)
	}
	log.Info("[Finally][%v]termId:%v,animal:%v,winAmount:%v", algVer, winRec.FTermId, winRec.FWinningAnimalName, amount)
	go sendRedPack(winRec)
	return
}

const Client_IP = "81.71.85.201"
const MAX_WX_REDPACK_AMOUNT = 20000

func sendRedPack(winRec *dao.TWinningRecord) {
	winRecs, err := dao.GetWinnerRecordByTermId(int64(winRec.FTermId))
	if err != nil {
		log.Error("dao.GetWinnerRecordByTermId failed,termId:%v,err:%v", winRec.FTermId, err)
		return
	}
	log.Info("termId:%v,winnerCnt:%v", winRec.FTermId, len(winRecs))
	req := &wx.SendWxRedpackReq{
		MchId:    wx.WXPAY_MCH_ID,
		WxAppId:  wx.WXPAY_APPID,
		SendName: "多乐十二生肖",
		TotalNum: 1,
		Wishing:  fmt.Sprintf("恭喜您在第 %v 期,开奖生肖 %v 中奖啦", winRec.FTermId, winRec.FWinningAnimalName),
		ClientIp: Client_IP,
		ActName:  "多乐十二生肖",
		Remark:   "下一期就快开奖喽,快去投注吧",
	}
	for _, rec := range winRecs {
		req.ReOpenid = rec.FOpenId
		totalAmount := rec.FPayPrice * conf.ConfItems.Common.WinMulti
		leftAmount := totalAmount
		for { //一次性最多发200块
			req.TotalAmount = leftAmount
			req.MchBillno = genBillno()
			req.NonceStr = util.CreateRandString(32)
			if leftAmount > MAX_WX_REDPACK_AMOUNT {
				req.TotalAmount = MAX_WX_REDPACK_AMOUNT
			}
			rsp, err := wx.SendWxRedpack(req)
			if err != nil {
				log.Error("wx.SendWxRedpack failed,termId:%v,req:%+v,err:%v", winRec.FTermId, req, err)
				continue
			}
			redpackRec := genRedpackRecord(rec, totalAmount, req.TotalAmount, req.MchBillno, rsp)
			err = redpackRec.InsertOne()
			if err != nil {
				log.Error("redpackRec.InsertOne failed,rec:%+v,err:%v", rec, err)
			}
			leftAmount -= MAX_WX_REDPACK_AMOUNT
			if leftAmount <= 0 {
				break
			}
		}

	}
}

func genBillno() string {
	return fmt.Sprintf("%s%s%s", wx.WXPAY_MCH_ID, time.Now().Format("20060102"), random.CustomChars(false, false, true, false, false, 10))
}
func genRedpackRecord(buyRec *dao.TUserBuyRecord, totalAmount, amount int, billno string, rsp *wx.SendWxRedpackRsp) *dao.TSendRedpackRecord {
	rec := new(dao.TSendRedpackRecord)

	rec.FTermId = buyRec.FTermId
	rec.FPayPrice = buyRec.FPayPrice
	rec.FOpenId = buyRec.FOpenId
	rec.FCreateTime = time.Now().Unix()
	rec.FOrderId = buyRec.FOrderId
	rec.FTotalAmount = totalAmount
	rec.FAmount = amount
	rec.FAnimalName = buyRec.FBuyAnimalName
	rec.FWxRetCode = rsp.ReturnCode
	rec.FWxRetMsg = rsp.ReturnMsg
	rec.FStatus = dao.WXPAY_SEND_REDPACK_SUCCESS
	rec.FBillno = billno
	rec.FListId = rsp.SendListid
	if rsp.ReturnCode == "SUCCESS" {
		rec.FStatus = dao.WXPAY_SEND_REDPACK_FAILED
		rec.FWxRetCode = fmt.Sprintf("%v_%v", rec.FWxRetCode, rsp.ResultCode)
		rec.FWxRetMsg = fmt.Sprintf("%v_%v_%v", rec.FWxRetMsg, rsp.ErrCode, rsp.ErrCodeDes)
	}

	return rec
}
