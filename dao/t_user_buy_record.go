package dao

import (
	"fmt"
)

type TUserBuyRecord struct {
	FId            int    `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FTermId        int    `json:"f_term_id" xorm:"not null comment('期次id') INT(11) 'f_term_id'"`
	FOpenId        string `json:"f_open_id" xorm:"not null comment('用户id') VARCHAR(128) 'f_open_id'"`
	FUserName      string `json:"f_user_name" xorm:"not null default '' comment('用户名') TEXT 'f_user_name'"`
	FBuyAnimalId   int    `json:"f_buy_animal_id" xorm:"not null default 0 comment('生肖id') INT(11) 'f_buy_animal_id'"`
	FBuyCnt        int    `json:"f_buy_cnt" xorm:"not null default 0 comment('购买注数') INT(11) 'f_buy_cnt'"`
	FWinningStatus int    `json:"f_winning_status" xorm:"not null default 0 comment('中奖状态') INT(11) 'f_winning_status'"`
	FPayPrice      int    `json:"f_pay_price" xorm:"not null default 0 comment('总共支付价格') INT(11) 'f_pay_price'"`
	FOrderId       string `json:"f_order_id" xorm:"not null comment('订单id') VARCHAR(128) 'f_order_id'"`
	FOrderStatus   int    `json:"f_order_status" xorm:"not null comment('订单状态') INT(11) 'f_order_status'"`
	FCreateTime    int64  `json:"f_create_time" xorm:"not null comment('创建时间') BIGINT(20) 'f_create_time'"`
	//	FUpdateTime      time.Time `json:"f_update_time" xorm:"not null default 'current_timestamp()' comment('更新时间') TIMESTAMP 'f_update_time'"`
	FWxTransactionId string `json:"f_wx_transaction_id" xorm:"not null comment('微信订单号') VARCHAR(128) 'f_wx_transaction_id'"`
	FBuyAnimalName   string `json:"f_buy_animal_name" xorm:"not null default 0 comment('生肖名') VARCHAR(128) 'f_buy_animal_name'"`
}

func (t *TUserBuyRecord) TableName() string {
	return "t_user_buy_record"
}

func (t *TUserBuyRecord) InsertOne() error {
	x, err := GetXormEngine("zodiac_write")
	if err != nil {
		return err
	}

	_, err = x.InsertOne(t)
	if err != nil {
		return err
	}

	return nil
}

func GetUserBuyRecordByOrderId(openid, orderId string) (*TUserBuyRecord, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	record := new(TUserBuyRecord)
	ok, err := x.Where("f_open_id = ? and f_order_id = ?", openid, orderId).Get(record)
	if err != nil || !ok {
		return nil, fmt.Errorf("ok:%v,err:%v", ok, err)
	}
	return record, nil
}

func GetWinnerRecordByTermId(termId int64) ([]*TUserBuyRecord, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	var recs []*TUserBuyRecord
	err = x.Where("f_term_id = ? and f_winning_status = ?", termId, Winning_Status_Winner).Find(&recs)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func GetUserBuyCnt(openid string) (int64, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return 0, err
	}
	rec := new(TUserBuyRecord)
	cnt, err := x.Where("f_open_id = ? and f_order_status = ?", openid, Order_Status_Payed).Count(rec)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
func GetUserBuyRecord(openid string, limit, offset int) ([]*TUserBuyRecord, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	var recs []*TUserBuyRecord
	err = x.Where("f_open_id = ? and f_order_status = ?", openid, Order_Status_Payed).Desc("f_id").Limit(limit, offset).Find(&recs)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func GetUserWinningCnt(openid string) (int64, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return 0, err
	}

	record := new(TUserBuyRecord)
	cnt, err := x.Where("f_open_id = ? and f_order_status = ? and f_winning_status = ?", openid, Order_Status_Payed, Winning_Status_Winner).Count(record)
	if err != nil {
		return 0, err
	}

	return cnt, nil
}
func GetUserWinningRecords(openid string, limit, offset int) ([]*TUserBuyRecord, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	var recs []*TUserBuyRecord
	err = x.Where("f_open_id = ? and f_order_status = ? and f_winning_status = ?", openid, Order_Status_Payed, Winning_Status_Winner).Desc("f_id").Limit(limit, offset).Find(&recs)
	if err != nil {
		return nil, err
	}
	return recs, nil
}

func (t *TUserBuyRecord) UpdateBuyRecord() error {
	x, err := GetXormEngine("zodiac_write")
	if err != nil {
		return err
	}

	_, err = x.Where("f_open_id = ? and f_order_id = ?", t.FOpenId, t.FOrderId).Update(t)
	if err != nil {
		return err
	}

	return nil
}

func GetBuyRecordByTermId(termId int) ([]*TUserBuyRecord, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	var sliceRecord []*TUserBuyRecord
	err = x.Where("f_term_id = ? and f_order_status = ?", termId, Order_Status_Payed).Find(&sliceRecord)
	if err != nil {
		return nil, err
	}
	return sliceRecord, nil
}

type ZodiacGroupAmount struct {
	AnimalId   int    `xorm:"f_buy_animal_id"`
	AnimalName string `xorm:"f_buy_animal_name"`
	SumPrice   int    `xorm:"sum_price"`
}

const ZodiacGroupAmountSql = `select f_buy_animal_id,f_buy_animal_name,sum(f_pay_price) as sum_price from t_user_buy_record where f_term_id = ? group by f_buy_animal_id order by sum_price desc`

func GetZodiacGroupAmount(termId int) ([]*ZodiacGroupAmount, error) {
	x, err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil, err
	}

	var zodiacAmount []*ZodiacGroupAmount
	err = x.SQL(ZodiacGroupAmountSql, termId).Find(&zodiacAmount)
	if err != nil {
		return nil, err
	}
	return zodiacAmount, nil
}
