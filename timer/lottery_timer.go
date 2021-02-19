package timer

import (
	"time"
	"zodiac_betting/calc_winner"
	"zodiac_betting/conf"
	rlog "zodiac_betting/rolllog"
)

func LotteryTimer()  {
	cnt := 0
	for   {
		time.Sleep(time.Second * time.Duration(conf.ConfItems.Common.LotteryInterval))
		cnt ++
		rlog.Info("now:%v,this is %vth lottery",time.Now().Format("2006-01-02 15:04:05"),cnt)
		calc_winner.CalcWinner()

	}
}
