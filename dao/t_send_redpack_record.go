package dao

import "time"

type TSendRedpackRecord struct {
	FId          int       `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FTermId      int       `json:"f_term_id" xorm:"not null comment('期次id') INT(11) 'f_term_id'"`
	FOpenId      string    `json:"f_open_id" xorm:"not null comment('用户id') VARCHAR(128) 'f_open_id'"`
	FAnimalName  string    `json:"f_animal_name" xorm:"not null default 0 comment('生肖名') VARCHAR(128) 'f_animal_name'"`
	FPayPrice    int       `json:"f_pay_price" xorm:"not null default 0 comment('总共支付价格') INT(11) 'f_pay_price'"`
	FOrderId     string    `json:"f_order_id" xorm:"not null comment('订单id') VARCHAR(128) 'f_order_id'"`
	FTotalAmount int       `json:"f_total_amount" xorm:"not null default 0 comment('总共奖励金额') INT(11) 'f_total_amount'"`
	FAmount      int       `json:"f_amount" xorm:"not null default 0 comment('此次奖励金额') INT(11) 'f_amount'"`
	FCreateTime  int64     `json:"f_create_time" xorm:"not null comment('创建时间') BIGINT(20) 'f_create_time'"`
	FUpdateTime  time.Time `json:"f_update_time" xorm:"not null default 'current_timestamp()' comment('更新时间') TIMESTAMP 'f_update_time'"`
	FStatus      int       `json:"f_status" xorm:"not null comment('发放状态') INT(11) 'f_status'"`
	FWxRetCode   string    `json:"f_wx_ret_code" xorm:"default null  comment('微信返回码') VARCHAR(128) 'f_wx_ret_code'"`
	FWxRetMsg    string    `json:"f_wx_ret_msg" xorm:"default null  comment('微信返信息') VARCHAR(128) 'f_wx_ret_msg'"`
	FBillno      string    `json:"f_billno" xorm:"not null  comment('商户单号') VARCHAR(128) 'f_billno'"`
	FListId      string    `json:"f_list_id" xorm:"default null  comment('微信单号') VARCHAR(128) 'f_list_id'"`
}

func (t *TSendRedpackRecord) TableName() string {
	return "t_send_redpack_record"
}

func (t *TSendRedpackRecord) InsertOne() error {
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
