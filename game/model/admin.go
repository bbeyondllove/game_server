package model

import "time"

// 管理员
type AdminUser struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	RoleName    string    `xorm:"-" json:"role_name" desc:"角色名"`
	RoleId      int64     `xorm:"int(20) notnull comment('角色ID')" json:"role_id" desc:"角色ID"`
	Department  string    `xorm:"varchar(255) notnull comment('部门')" json:"department" desc:"部门"`
	Username    string    `xorm:"varchar(255) notnull comment('用户名')" json:"username" desc:"用户名"`
	Account     string    `xorm:"varchar(255) notnull unique comment('账户')" json:"account" desc:"账户"`
	Password    string    `xorm:"varchar(255) comment('账号密码')" json:"password" desc:"账号密码"`
	CreaterName string    `xorm:"-" json:"creater_name" desc:"创建人"`
	Creater     int64     `xorm:"int(20) notnull comment('创建人')" json:"creater" desc:"创建人"`
	Status      int       `xorm:"int(20) notnull comment('状态,0:禁用 1:启用')" json:"status" desc:"状态,0:可用 1:禁用"`
	CreateTime  time.Time `xorm:"timestamp notnull comment('创建时间')" json:"create_time"`
	UpdateTime  time.Time `xorm:"timestamp notnull comment('更新时间')" json:"update_time"`
}

// 管理员角色
type AdminRole struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Name       string    `xorm:"varchar(255) notnull unique comment('角色名称')" json:"name" desc:"角色名称"`
	Remark     string    `xorm:"varchar(255) notnull comment('角色说明')" json:"remark" desc:"角色说明"`
	Limits     string    `xorm:"varchar(255) notnull comment('角色权限')" json:"limits" desc:"角色权限"`
	CreateTime time.Time `xorm:"timestamp notnull comment('创建时间')" json:"create_time"`
	UpdateTime time.Time `xorm:"timestamp notnull comment('更新时间')" json:"update_time"`
}

// 管理员权限
type AdminLimit struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Name       string    `xorm:"varchar(255) notnull comment('权限名')" json:"label" desc:"权限名"`
	Group      string    `xorm:"varchar(255) notnull comment('权限组')" json:"group" desc:"权限组"`
	Code       string    `xorm:"varchar(255) notnull comment('权限代码')" json:"code" desc:"权限代码"`
	Internal   int       `xorm:"int(2) notnull comment('内部/外部权限,内部权限不可被修改')" json:"-" desc:"权限方式"`
	CreateTime time.Time `xorm:"timestamp notnull comment('创建时间')" json:"-"`
	UpdateTime time.Time `xorm:"timestamp notnull comment('更新时间')" json:"-"`
}

// 操作日志
type AdminOperationLog struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	OperationId string    `xorm:"int(20) notnull comment('操作人')" json:"operation_id" desc:"操作人"`
	Object      string    `xorm:"varchar(255) notnull comment('操作对象')" json:"object" desc:"操作对象"`
	Action      string    `xorm:"varchar(255) notnull comment('操作动作')" json:"action" desc:"操作动作"`
	Befor       string    `xorm:"varchar(255) notnull comment('修改前')" json:"befor" desc:"修改前"`
	After       string    `xorm:"varchar(255) notnull comment('修改后')" json:"after" desc:"修改后"`
	Ip          string    `xorm:"varchar(64) notnull comment('ip')" json:"ip" desc:"ip地址"`
	CreateTime  time.Time `xorm:"timestamp notnull comment('创建时间')" json:"create_time"`
	UpdateTime  time.Time `xorm:"timestamp notnull comment('更新时间')" json:"update_time"`
}
