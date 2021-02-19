package proto

type UserAccountRsp struct {
	Ret    int    `json:"ret"`
	Openid string `json:"openid"`
}

type WxTicketRsp struct {
	Ret    int    `json:"ret"`
	Ticket string `json:"ticket"`
	HB     string `json:"hb"`
}

type ActivityIntroductionRsp struct {
	Ret          int    `json:"ret"`
	LotteryAlert string `json:"lottery_alert"`
	Introduction string `json:"introduction"`
}
type MarqueeInfoRsp struct {
	Ret             int      `json:"ret"`
	LastTerm        int      `json:"last_term"`
	LastLotteryTime int64    `json:"last_lottery_time"`
	LastAnimal      string   `json:"last_animal"`
	NextLotteryLeft int64    `json:"next_lottery_left"`
	WinningInfos    []string `json:"winning_infos"`
}

type ProductInfo struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Price         int    `json:"price"`
	Image         string `json:"image"`
	ImageSelected string `json:"image_selected"`
}
type GetProductInfoRsp struct {
	Ret       int            `json:"ret"`
	ProdInfos []*ProductInfo `json:"prod_infos"`
}
