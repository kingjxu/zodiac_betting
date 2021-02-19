package proto

import "zodiac_betting/dao"

type LotteryHistoryReq struct {
	Start int `json:"start"`
	PageSize int `json:"page_size"`
}

type LotteryHistoryRsp struct {
	Ret int `json:"ret"`
	PageCnt int `json:"page_cnt"`
	History []*dao.TWinningRecord `json:"history"`
}

type ZodiacWinnerProbability struct {
	AnimalId int `json:"animal_id"`
	AnimalName string `json:"animal_name"`
	Prob int `json:"prob"` //万分之多少，相当于是保留小数点后两位
}
type WinnerProbabilityRsp struct {
	Ret int `json:"ret"`
	Probs []*ZodiacWinnerProbability `json:"probs"`
}
