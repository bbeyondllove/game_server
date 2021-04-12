package db

import (
	"game_server/core/base"
	"game_server/core/logger"

	"github.com/go-xorm/xorm"

	"xorm.io/core"
)

var (
	Mysql *xorm.Engine
)

func init() {
	defer logger.Flush()
	//连接数据库
	url := base.Setting.Mysql.Username + ":" + base.Setting.Mysql.Password + "@tcp(" + base.Setting.Mysql.Host + ":" + base.Setting.Mysql.Port + ")/" + base.Setting.Mysql.DbName + "?" + base.Setting.Mysql.Parameter
	var err error
	Mysql, err = xorm.NewEngine("mysql", url)
	if err != nil {
		logger.Errorf("An error occurred while connecting to Mysql ", err)
	}
	// 设置表名映射
	tableMapper := core.NewPrefixMapper(core.SnakeMapper{}, base.Setting.Mysql.TablePrefix)
	Mysql.SetTableMapper(tableMapper)
	Mysql.SetMaxOpenConns(base.Setting.Mysql.Connects)
	Mysql.SetMaxOpenConns(base.Setting.Mysql.Connects)
	//engine.TZLocation, _ = time.LoadLocation("Asia/Shanghai")
}
