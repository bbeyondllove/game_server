package db_service

import (
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"reflect"
)

//更新记录
func UpdateFields(table_name string, key string, seqId string, fields map[string]interface{}) (int64, error) {
	update := " set "
	i := 0
	update_field := make(map[string]interface{})

	for k, v := range fields {
		if v == nil || v == "" {
			continue
		}
		update_field[k] = v
	}

	for k, v := range update_field {
		update += k
		update += "="
		ntype := reflect.TypeOf(v).String()
		if ntype == "string" || ntype == "time.Time" {
			update += "'"
		}
		update += utils.Strval(v)
		if ntype == "string" || ntype == "time.Time" {
			update += "'"
		}

		if i < len(update_field)-1 {
			update += ","
		}
		i++
	}
	sql := "update " + table_name + update + " where " + key + "='" + seqId + "'"
	//fmt.Printf(sql)
	r, err := db.Mysql.Exec(sql)
	if err != nil {
		logger.Errorf("[Error]: %v", err)
		return 0, err
	}
	count, _ := r.RowsAffected()
	if count == 0 || err != nil {
		logger.Errorf("[Error]: %v", err)
		return 0, err
	}

	return count, err
}

