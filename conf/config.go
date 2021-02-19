package conf

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"path/filepath"
)

type stSqlConf struct {
	DSN     string `toml:"dsn"`
	MaxLifeTime int `toml:"max_life_time"`
}

type stRedisConf struct {
	Addr     string `toml:"addr"`
	Password     string `toml:"password"`
}

type stCommonConf struct {
	ProbAlgV2     int `toml:"prob_algv2"`
	WinMulti     int `toml:"win_multi"`
	LotteryInterval int `toml:"lottery_interval"`
	AboutUs string `toml:"about_us"`
	AboutAct string `toml:"about_act"`
}

type stLogConf struct {
	Path        string `toml:"path"`
	Level       uint32 `toml:"level"`
	FileMaxSize uint32 `toml:"file_max_size"`
	FileMaxCnt  uint32 `toml:"file_max_count"`
}

type stConfItems struct {
	MySQL   stSqlConf      `toml:"mysql"`
	Redis 	stRedisConf	   `toml:"redis"`
	Log     stLogConf     `toml:"log"`
	Common  stCommonConf   `toml:"common"`
}


var ConfItems stConfItems

func init() {
	fmt.Printf("module[config] init...\n")
	loadConfig()
}

func fileExist(_path string) bool {
	_, err := os.Stat(_path)
	if err != nil && os.IsNotExist(err) {
		return false
	}
	return true
}

func getCurrPath() string {
	exeAbsPathName, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Printf("filepath.Abs failed! err:%v\n",err)
		return ""
	}

	currPath := filepath.Dir(exeAbsPathName)
	return currPath
}

func loadConfig() {

	currPath := getCurrPath()
	if currPath == "" {
		return
	}

	fileName := currPath + "/../conf/zodiac_betting.conf"
	fmt.Printf("load config file path:%v\n", fileName)

	if bExist := fileExist(fileName); !bExist {
		fmt.Printf("load config faild, file_name:%v\n", fileName)
		return
	}

	buf, err := ioutil.ReadFile(fileName)
	if err != nil {
		fmt.Printf("ioutil.ReadFile faild, file_name:%v\n", fileName)
		return
	}

	err = toml.Unmarshal(buf, &ConfItems)
	if err != nil {
		fmt.Printf("toml.Unmarshal error:%v ", err)
	}

	if ConfItems.Log.FileMaxSize < 50*1024*1024 {
		ConfItems.Log.FileMaxSize = 50 * 1024 * 1024
	}

	if ConfItems.Log.FileMaxCnt < 1 {
		ConfItems.Log.FileMaxSize = 1
	}
	fmt.Printf("log level:%v,mysql dsn:%v,conf%+v\n",ConfItems.Log.Level,ConfItems.MySQL.DSN,ConfItems)
}