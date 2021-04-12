package model

import "time"

//小区房子表
type House struct {
	Id         string `xorm:"int(20) notnull pk" json:"id"`
	Position_x string `xorm:"int(6) notnull" json:"position_x"`
	Position_y string `xorm:"int(6) notnull" json:"position_y"`
	House_seq  string `xorm:"varchar(11)" json:"house_seq"`
	IsSale     int    `xorm:"int(1) notnull" json:"is_sale"`

	CreateTime time.Time `xorm:"timestamp notnull" json:"createTime"`
	UpdateTime time.Time `xorm:"timestamp notnull" json:"updateTime"`
}
