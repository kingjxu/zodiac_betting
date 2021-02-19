package main

import (
	"net/http"
	"zodiac_betting/pages"
	log "zodiac_betting/rolllog"
	"zodiac_betting/timer"
)

func initCgi(mux *http.ServeMux) {

	mux.HandleFunc("/homepage/user_account", pages.UserAccount)
	mux.HandleFunc("/homepage/marquee_info", pages.GetMarqueeInfo)
	mux.HandleFunc("/homepage/product_info", pages.GetProductInfo)
	mux.HandleFunc("/homepage/ticket", pages.GetTicket)

	mux.HandleFunc("/history/lottery_history", pages.LotteryHistory)
	mux.HandleFunc("/history/winner_prob", pages.WinnerProbability)

	mux.HandleFunc("/order/create_order", pages.CreateOrder)
	mux.HandleFunc("/order/delivery", pages.Delivery)

	mux.HandleFunc("/user_center/avatar_list", pages.GetAvatarList)
	mux.HandleFunc("/user_center/get_user_info", pages.GetUserInfo)
	mux.HandleFunc("/user_center/set_user_info", pages.SetUserInfo)
	mux.HandleFunc("/user_center/buy_records", pages.GetBuyRecord)
	mux.HandleFunc("/user_center/winning_records", pages.WinningRecords)

	mux.HandleFunc("/user_center/about_us", pages.AboutUs)
	mux.HandleFunc("/user_center/about_act", pages.AboutAct)

	mux.HandleFunc("/report/report_log", pages.ReportLog)

}

func main() {

	log.SetLevel(log.LEVEL_DEBUG)

	mux := http.NewServeMux()
	mux.Handle("/pic/", http.StripPrefix("/pic/", http.FileServer(http.Dir("/usr/local/services/zodiac_betting/pic"))))

	initCgi(mux)

	go timer.LotteryTimer()
	err := http.ListenAndServe(":8080", mux)
	if err != nil {
		log.Error("Failed to start server, err:%v\n", err)
		return
	}

}
