package dao

import (
	"git.code.oa.com/SNG_EDU_COMMON_PKG/bingo/handler"
	_ "github.com/go-sql-driver/mysql"
	"github.com/xormplus/xorm"
	"log"
	"sync"
	"time"
	"zodiac_betting/conf"
)

var engines *sync.Map
var lock sync.Mutex

func init() {
	engines = new(sync.Map)
}

//按数据源获取 xorm实例
func GetXormEngine(source string) (engine *xorm.Engine, err error) {

	x, ok := engines.Load(source)
	if ok {
		engine = x.(*xorm.Engine)
		return
	}

	lock.Lock()
	defer lock.Unlock()

	engine, err = xorm.NewEngine("mysql",  conf.ConfItems.MySQL.DSN)
	engine.SetConnMaxLifetime(time.Duration(conf.ConfItems.MySQL.MaxLifeTime) * time.Second)
	engine.SetMaxIdleConns(2)
	engine.SetMaxOpenConns(30)

	engine.SetLogger(xorm.NewSimpleLogger2(handler.DefaultLogWriter, "", xorm.DEFAULT_LOG_FLAG))
	//探活
	err = engine.Ping()
	if err != nil {
		return
	}
	engines.Store(source, engine)
	go keepAlive(engine)

	return
}

//保活
func keepAlive(engine *xorm.Engine) {
	for {
		<-time.After(30 * time.Second)
		if engine.Ping() != nil {
			log.Println("xorm ping error")
		}
	}
}
