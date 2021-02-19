package dao


type TWinningRecord struct {
	FId                int       `json:"f_id" xorm:"not null pk autoincr INT(11) 'f_id'"`
	FTermId            int       `json:"f_term_id" xorm:"not null comment('期次id') INT(11) 'f_term_id'"`
	FCreateTime        int64     `json:"f_create_time" xorm:"not null comment('开奖时间') BIGINT(20) 'f_create_time'"`
	//	FUpdateTime        time.Time `json:"f_update_time" xorm:"not null default 'current_timestamp()' comment('更新时间') TIMESTAMP 'f_update_time'"`
	FWinningAnimalId   int       `json:"f_winning_animal_id" xorm:"not null default 0 comment('中奖的生肖的属性id') INT(11) 'f_winning_animal_id'"`
	FWinningAnimalName string    `json:"f_winning_animal_name" xorm:"not null comment('生肖的名字') VARCHAR(128) 'f_winning_animal_name'"`
	FWinningUserCnt    int       `json:"f_winning_user_cnt" xorm:"not null default 0 comment('中奖人数') INT(11) 'f_winning_user_cnt'"`
	FWinningBetCnt     int       `json:"f_winning_bet_cnt" xorm:"not null default 0 comment('中奖注数') INT(11) 'f_winning_bet_cnt'"`
}


func (t *TWinningRecord)TableName() string {
	return "t_winning_record"
}

func GetLatestWinningRecord() (*TWinningRecord,bool,error)  {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,false,err
	}

	rec := new(TWinningRecord)
	exist,_,err := x.Desc("f_id").GetFirst(rec).GetResult()
	if err != nil {
		return nil,false,err
	}

	if !exist {
		return rec,false,nil
	}

	return rec,true,nil
}

func GetWinningRecordByTerms(terms []int) ([]*TWinningRecord,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,err
	}

	var recs []*TWinningRecord
	err = x.In("f_term_id",terms).Find(&recs)
	if err != nil {
		return nil,err
	}

	return recs,nil
}
func (t *TWinningRecord)InsertOne() (error)  {
	x,err := GetXormEngine("zodiac_write")
	if err != nil {
		return err
	}

	_,err = x.InsertOne(t)
	if err != nil {
		return err
	}

	return nil
}

func GetLotteryHistory(offset,limit int) ([]*TWinningRecord,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,err
	}

	var rec []*TWinningRecord
	err = x.Desc("f_id").Limit(limit,offset).Find(&rec)
	if err != nil {
		return nil,err
	}

	return rec,nil
}
func GetLotteryCnt() (int64,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return 0,err
	}

	rec := new(TWinningRecord)
	cnt,err := x.Count(rec)
	if err != nil {
		return 0,err
	}

	return cnt,nil
}
type ZodiacWinnerFreq struct {
	AnimalId int `xorm:"f_winning_animal_id"`
	Count int `xorm:"cc"`
}
const zodiacWinnerFreqSql = "select f_winning_animal_id,count(f_winning_animal_id) as cc from t_winning_record group by f_winning_animal_id"
func GetZodiacWinnerFreq() ([]*ZodiacWinnerFreq,error) {
	x,err := GetXormEngine("zodiac_read")
	if err != nil {
		return nil,err
	}

	var rec []*ZodiacWinnerFreq
	err = x.SQL(zodiacWinnerFreqSql).Find(&rec)
	if err != nil {
		return nil,err
	}

	return rec,nil
}
