package db_service

import (
	"game_server/db"
	"game_server/game/model"
)

type RoleInfo struct{}

const role_info_table = "t_role_info"

//获取金额
func (gw *RoleInfo) GetRoleInfo(roleId int) (*model.RoleInfo, error) {
	data := &model.RoleInfo{}

	_, err := db.Mysql.Table(WalletTable).Where("user_id= ?", roleId).Get(data)
	return data, err
}

func (this *RoleInfo) GetAllData() ([]model.RoleInfo, error) {
	var data []model.RoleInfo
	err := db.Mysql.Table(role_info_table).Where("1=1").Find(&data)
	return data, err
}
