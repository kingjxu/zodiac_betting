package dao

type TUserInfo struct {
	FId       int            `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FOpenId   string         `json:"f_open_id" xorm:"not null comment('用户名') VARCHAR(128) 'f_open_id'"`
	FNickName string         `json:"f_nick_name" xorm:"not null comment('昵称') VARCHAR(128) 'f_nick_name'"`
	FLevel   string `json:"f_level" xorm:"comment('用户等级') VARCHAR(128) 'f_level'"`
	FPhone    string `json:"f_phone" xorm:"comment('电话') VARCHAR(128) 'f_phone'"`
	FAvatar   string         `json:"f_avatar" xorm:"not null comment('头像') TEXT 'f_avatar'"`
}
func (t *TUserInfo)TableName() string {
	return "t_user_info"
}

func GetUserInfo(openid string) (*TUserInfo,bool,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,false,err
	}

	userInfo := new(TUserInfo)
	ok,err := x.Where("f_open_id = ?",openid).Get(userInfo)
	if err != nil {
		return nil,false,err
	}
	if !ok {
		return nil,false,nil
	}
	return userInfo,true,nil
}

func (t *TUserInfo)InsertOne() error {
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
func (t *TUserInfo)UpdateOne() error {
	x, err := GetXormEngine("zodiac_write")
	if err != nil {
		return err
	}

	_, err = x.Where("f_open_id = ?",t.FOpenId).Update(t)
	if err != nil {
		return err
	}

	return nil
}
