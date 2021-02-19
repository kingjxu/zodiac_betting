package dao

import "time"

var ZodiacIdName = map[int]string{
	1000: "鼠",
	1001: "牛",
	1002: "虎",
	1003: "兔",
	1004: "龙",
	1005: "蛇",
	1006: "马",
	1007: "羊",
	1008: "猴",
	1009: "鸡",
	1010: "狗",
	1011: "猪",
}

var AvatarIdName = map[string]string{
	"dushen":  "dushen.jpg",
	"duxia":   "duxia.jpg",
	"dusheng": "dusheng.jpg",
	"cny":     "cny.jpg",
	"zh":      "dzh.jpg",
	"tm":      "tm.jpg",
	"jn":      "jn.jpg",
	"zn":      "zn.jpg",
}

var UserName = []string{
	"D*",
	"K**",
	"菜****",
	"孙**",
	"J*****",
	"往***",
	"15***",
	"李**",
	"K***",
	"如**",
	"皮卡***",
	"小***",
	"158****",
	"拼***",
	"天***",
	"婷**",
	"文***",
	"张**",
	"皮卡***",
	"仰***",
	"159****",
	"lu***",
}

const (
	Order_Status_To_Be_Pay = 1 //待支付
	Order_Status_Payed     = 2 //已支付
)

const (
	Winning_Status_To_Be_Winning = 1 //待开奖
	Winning_Status_Not_Winner    = 2 //未中奖
	Winning_Status_Winner        = 3 //已中奖

)

const (
	WXPAY_SEND_REDPACK_SUCCESS = 1
	WXPAY_SEND_REDPACK_FAILED  = 2
)

type UserBuyRecordAndWinnerAnimal struct {
	FId              int       `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FTermId          int       `json:"f_term_id" xorm:"not null comment('期次id') INT(11) 'f_term_id'"`
	FOpenId          string    `json:"f_open_id" xorm:"not null comment('用户id') VARCHAR(128) 'f_open_id'"`
	FUserName        string    `json:"f_user_name" xorm:"not null default '' comment('用户名') TEXT 'f_user_name'"`
	FBuyAnimalId     int       `json:"f_buy_animal_id" xorm:"not null default 0 comment('生肖id') INT(11) 'f_buy_animal_id'"`
	FBuyCnt          int       `json:"f_buy_cnt" xorm:"not null default 0 comment('购买注数') INT(11) 'f_buy_cnt'"`
	FWinningStatus   int       `json:"f_winning_status" xorm:"not null default 0 comment('中奖状态') INT(11) 'f_winning_status'"`
	FPayPrice        int       `json:"f_pay_price" xorm:"not null default 0 comment('总共支付价格') INT(11) 'f_pay_price'"`
	FOrderId         string    `json:"f_order_id" xorm:"not null comment('订单id') VARCHAR(128) 'f_order_id'"`
	FOrderStatus     int       `json:"f_order_status" xorm:"not null comment('订单状态') INT(11) 'f_order_status'"`
	FCreateTime      int64     `json:"f_create_time" xorm:"not null comment('创建时间') BIGINT(20) 'f_create_time'"`
	FUpdateTime      time.Time `json:"f_update_time" xorm:"not null default 'current_timestamp()' comment('更新时间') TIMESTAMP 'f_update_time'"`
	FWxTransactionId string    `json:"f_wx_transaction_id" xorm:"not null comment('微信订单号') VARCHAR(128) 'f_wx_transaction_id'"`
	FBuyAnimalName   string    `json:"f_buy_animal_name" xorm:"not null default 0 comment('生肖名') INT(11) 'f_buy_animal_name'"`
	FWinAnimalName   string    `json:"f_win_animal_name" xorm:"not null default 0 comment('生肖名') INT(11) 'f_win_animal_name'"`
}
