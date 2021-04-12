package model

import "time"

type ActivityDate struct {
	Id             int       `xorm:"int autoincr pk" json:"id"`
	Date           time.Time `xorm:"date" json:"date"`
	ActivityType   int       `xorm:"int(2)" json:"activity_type"`
	InvolvNum      int       `xorm:"int" json:"involv_num"`
	DataDayNum     string    `xorm:"varchar(256)" json:"data_day_num"`
	MakeUpCardNum  int       `xorm:"int" json:"make_up_card_num"`
	ActivityUv     int       `xorm:"int" json:"activity_uv"`
	ActivityPv     int       `xorm:"int" json:"activity_pv"`
	CreateTime     time.Time `xorm:"timestamp notnull" json:"create_time"`
	TreasureBoxNum int       `xorm:"int" json:"treasure_box_num"`
}
