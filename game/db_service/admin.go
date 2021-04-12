package db_service

import (
	"game_server/db"
	"game_server/game/model"
	"time"
)

type AdminUser struct{}

const (
	AdminUserTable = "t_admin_user"
)

//添加记录
func (this *AdminUser) Add(data_map *model.AdminUser) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *AdminUser) Update(account string, username, password, department *string, role_id *int64) (int, error) {
	data := map[string]interface{}{}
	if username != nil {
		data["username"] = *username
	}
	if password != nil {
		data["password"] = *password
	}
	if department != nil {
		data["department"] = *department
	}
	if role_id != nil {
		data["role_id"] = *role_id
	}
	number, err := db.Mysql.Table(AdminUserTable).Where("account= ? ", account).Update(data)
	return int(number), err
}

func (this *AdminUser) UpdateStatus(account string, status int64) (int, error) {
	data := map[string]interface{}{
		"status": status,
	}
	number, err := db.Mysql.Table(AdminUserTable).Where("account= ? ", account).Update(data)
	return int(number), err
}

func (this *AdminUser) Delete(account string) (int, error) {
	var data model.AdminUser
	number, err := db.Mysql.Table(AdminUserTable).Where("account= ? ", account).Limit(1, 0).Unscoped().Delete(&data)
	return int(number), err
}

func (this *AdminUser) GetPageData(page, size int) (int, int, []model.AdminUser, error) {
	if page <= 0 {
		page = 1
	}
	count, err := db.Mysql.Table(AdminUserTable).AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}
	var data []model.AdminUser
	err = db.Mysql.Table(AdminUserTable).OrderBy("id").
		Limit(size, (page-1)*size).
		Find(&data)
	return int(count), page, data, err
}

func (this *AdminUser) GetDataById(userId int64) (model.AdminUser, error) {
	data := model.AdminUser{}
	_, err := db.Mysql.Table(AdminUserTable).Where("id= ? ", userId).Get(&data)
	return data, err
}

func (this *AdminUser) GetDataByUid(userId string) (model.AdminUser, error) {
	data := model.AdminUser{}
	_, err := db.Mysql.Table(AdminUserTable).Where("id= ? ", userId).Get(&data)
	return data, err
}

func (this *AdminUser) GetDataByUname(userName string) (model.AdminUser, error) {
	data := model.AdminUser{}
	_, err := db.Mysql.Table(AdminUserTable).Where("account = ? ", userName).Get(&data)
	return data, err
}

func (this *AdminUser) Count(role_id int64) (int64, error) {
	num, err := db.Mysql.Table(AdminUserTable).Where("role_id = ? ", role_id).AllCols().Count()
	return num, err
}

type AdminRole struct{}

const (
	AdminRoleTable = "t_admin_role"
)

func (this *AdminRole) Add(data_map *model.AdminRole) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *AdminRole) Update(uesrname string, data model.AdminRole) (int, error) {
	data_map := map[string]interface{}{
		"remark":      data.Remark,
		"limits":      data.Limits,
		"update_time": time.Now(),
	}
	number, err := db.Mysql.Table(AdminRoleTable).Where("name= ? ", uesrname).Update(data_map)
	return int(number), err
}

func (this *AdminRole) GetDataById(id int64) (model.AdminRole, error) {
	data := model.AdminRole{}
	_, err := db.Mysql.Table(AdminRoleTable).Where("id = ? ", id).Get(&data)
	return data, err
}

func (this *AdminRole) GetDataByName(name string) (model.AdminRole, error) {
	data := model.AdminRole{}
	_, err := db.Mysql.Table(AdminRoleTable).Where("name = ? ", name).Get(&data)
	return data, err
}

func (this *AdminRole) GetAllData() ([]model.AdminRole, error) {
	var data []model.AdminRole
	err := db.Mysql.Table(AdminRoleTable).
		Find(&data)
	return data, err
}

func (this *AdminRole) Delete(role_name string) (int64, error) {
	var data model.AdminRole
	number, err := db.Mysql.Table(AdminRoleTable).Where("name = ? ", role_name).Limit(1, 0).Unscoped().Delete(&data)
	return number, err
}

type AdminLimit struct{}

const (
	AdminLimitTable = "t_admin_limit"
)

func (this *AdminLimit) Add(data_map *model.AdminLimit) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *AdminLimit) GetAllData() ([]model.AdminLimit, error) {
	var data []model.AdminLimit
	err := db.Mysql.Table(AdminLimitTable).OrderBy("id").
		Find(&data)
	return data, err
}

func (this *AdminLimit) GetData(internal int) ([]model.AdminLimit, error) {
	var data []model.AdminLimit
	err := db.Mysql.Table(AdminLimitTable).Where("internal = ?", internal).OrderBy("id").
		Find(&data)
	return data, err
}

type AdminOperationLog struct{}

const (
	AdminOperationLogTable = "t_admin_limit"
)

func (this *AdminOperationLog) Add(data_map *model.AdminOperationLog) (bool, error) {
	session := db.Mysql.NewSession()
	defer session.Close()
	err := session.Begin()

	data_map.CreateTime = time.Now()
	data_map.UpdateTime = time.Now()

	_, err = db.Mysql.Insert(data_map)
	if err != nil {
		_ = session.Rollback()
		return false, err
	}
	err = session.Commit()
	if err != nil {
		return false, err
	}
	return true, nil
}

func (this *AdminOperationLog) GetPageData(optionId string, beginTime, endTime string, page, size int) (int, int, []model.AdminOperationLog, error) {
	var data []model.AdminOperationLog
	conn := db.Mysql.Table(AdminOperationLogTable)
	if len(optionId) > 0 {
		conn = conn.Where("operation_id = ?", optionId).And("create_time >= ? and create_time <= ?", beginTime, endTime)
	} else {
		conn = conn.Where("create_time >= ? and create_time <= ?", beginTime, endTime)
	}

	if page <= 0 {
		page = 1
	}

	count, err := conn.AllCols().Count()
	if err != nil {
		return 0, 0, nil, err
	}

	conn = db.Mysql.Table(AdminOperationLogTable)
	if len(optionId) > 0 {
		conn = conn.Where("operation_id = ?", optionId).And("create_time >= ? and create_time <= ?", beginTime, endTime)
	} else {
		conn = conn.Where("create_time >= ? and create_time <= ?", beginTime, endTime)
	}
	err = conn.Limit(size, (page-1)*size).Find(&data)
	return int(count), page, data, err
}
