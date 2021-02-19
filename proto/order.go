package proto

type CreateOrderReq struct {
	ProdId int `json:"prod_id"`
	BuyCnt int `json:"buy_cnt"`
}

type CreateOrderRsp struct {
	Ret int `json:"ret"`
	OrderId string `json:"order_id"`
	PrepayId string `json:"prepay_id"`
	TermId int `json:"term_id"`
}
