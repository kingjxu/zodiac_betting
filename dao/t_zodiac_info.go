package dao

import "fmt"

type TZodiacInfo struct {
	FId       int    `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FZodiacId int    `json:"f_zodiac_id" xorm:"not null comment('生肖的id') INT(11) 'f_zodiac_id'"`
	FName     string `json:"f_name" xorm:"not null comment('生肖的名字') VARCHAR(128) 'f_name'"`
	FPrice    int    `json:"f_price" xorm:"not null comment('生肖的价格') INT(11) 'f_price'"`
	FPic      string `json:"f_pic" xorm:"not null comment('图片的地址') TEXT 'f_pic'"`
}

func (t *TZodiacInfo)TableName() string {
	return "t_zodiac_info"
}

func GetZodiacInfos() ([]*TZodiacInfo,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,err
	}

	var sliceZodiac []*TZodiacInfo
	err = x.Where("f_id >0").Find(&sliceZodiac)
	if err != nil {
		return nil,err
	}
	return sliceZodiac,nil
}

const DEFAULT_PRICE = 500
func GetProductPrice(prodId int) (int,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return DEFAULT_PRICE,err
	}

	zodiac := new(TZodiacInfo)
	ok,err := x.Where("f_zodiac_id = ?",prodId).Get(zodiac)
	if err != nil || !ok{
		return DEFAULT_PRICE,fmt.Errorf("ok:%v,err:%v",ok,err)
	}
	return zodiac.FPrice,nil
}
