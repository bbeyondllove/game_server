package model

import "time"

// 实时日统计数据
type StatisticsRealDay struct {
	Id              int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date            time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Platform        int       `xorm:"int(20) notnull comment('平台')" json:"platform" desc:"平台"`
	DayRegisteCount int       `xorm:"int(20) notnull comment('当日注册')" json:"day_registe_count" desc:"当日注册"`
	DayNewlyAdded   int       `xorm:"int(20) notnull comment('当日新增')" json:"day_newly_added" desc:"当日新增"`
	DayLoginCount   int       `xorm:"int(20) notnull comment('登录次数')" json:"day_login_count" desc:"登录次数"`
	ActiveCount     int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	DayCdt          float32   `xorm:"decimal(12,4) notnull comment('当日产出CDT')" json:"day_Cdt" desc:"当日产出CDT"`
	CreateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 日统计数据
type StatisticsDay struct {
	Id                int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date              time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Platform          int       `xorm:"int(20) notnull comment('平台')" json:"platform" desc:"平台"`
	TotalRegCount     int       `xorm:"int(20) notnull comment('总注册数')" json:"total_reg_count" desc:"总注册数"`
	TotalNewlyAdded   int       `xorm:"int(20) notnull comment('总新增')" json:"total_newly_added" desc:"总新增"`
	DayRegisteCount   int       `xorm:"int(20) notnull comment('当日注册')" json:"day_registe_count" desc:"当日注册"`
	DayNewlyAdded     int       `xorm:"int(20) notnull comment('当日新增')" json:"day_newly_added" desc:"当日新增"`
	DayLoginCount     int       `xorm:"int(20) notnull comment('登录次数')" json:"day_login_count" desc:"登录次数"`
	ActiveCount       int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	OnlineTime        int       `xorm:"int(20) notnull comment('在线时长')" json:"online_count" desc:"在线时长"`
	AvgOnlineTime     int       `xorm:"int(20) notnull comment('人均在线时长')" json:"avg_online_count" desc:"人均在线时长"`
	PhoneRegisteCount int       `xorm:"int(20) notnull comment('手机累计注册数量')" json:"phone_registe_count" desc:"手机累计注册数量"`
	EmailRegisteCount int       `xorm:"int(20) notnull comment('邮箱累计注册数量')" json:"email_registe_count" desc:"邮箱累计注册数量"`
	RealNameCount     int       `xorm:"int(20) notnull comment('累计实名用户')" json:"real_name_count" desc:"累计实名用户"`
	DayCdt            float32   `xorm:"decimal(12,4) notnull comment('当日产出CDT')" json:"day_cdt" desc:"当日产出CDT"`
	TotalCdt          float32   `xorm:"decimal(12,4) notnull comment('累计产出CDT')" json:"total_cdt" desc:"累计产出CDT"`
	CreateTime        time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime        time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 活跃用户数
type StatisticsActiveCount struct {
	Id               int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date             time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Hour             int       `xorm:"int(20) notnull comment('小时')" json:"hour" desc:"小时"`
	ActiveCount      int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	TotalActiveCount int       `xorm:"int(20) notnull comment('总活跃用户数')" json:"total_active_count" desc:"总活跃用户数"`
	Ratio            int       `xorm:"int(20) notnull comment('占比')" json:"ratio" desc:"占比"`
	CreateTime       time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime       time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 每小时统计留存数
type StatisticsHourRetained struct {
	Id              int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date            time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Hour            int       `xorm:"int(20) notnull comment('小时')" json:"hour" desc:"小时"`
	ActiveCount     int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Date1           time.Time `xorm:"DATE notnull comment('前日')" json:"date1" desc:"前日"`
	Retained1       int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Added1          int       `xorm:"int(20) notnull comment('前日新增')" json:"added1" desc:"前日新增"`
	Retained1Ratio  int       `xorm:"int(20) notnull comment('次日留存率')" json:"retained1_ratio" desc:"次日留存率"`
	Date3           time.Time `xorm:"DATE notnull comment('前3日')" json:"date3" desc:"前3日"`
	Retained3       int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Added3          int       `xorm:"int(20) notnull comment('前3日新增')" json:"added3" desc:"前3日新增"`
	Retained3Ratio  int       `xorm:"int(20) notnull comment('3日留存率')" json:"retained3_ratio" desc:"3日留存率"`
	Date7           time.Time `xorm:"DATE notnull comment('前7日')" json:"date7" desc:"前7日"`
	Retained7       int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	Added7          int       `xorm:"int(20) notnull comment('前7日新增')" json:"added7" desc:"前7日新增"`
	Retained7Ratio  int       `xorm:"int(20) notnull comment('7日留存率')" json:"retained7_ratio" desc:"7日留存率"`
	Date15          time.Time `xorm:"DATE notnull comment('前15日')" json:"date15" desc:"前15日"`
	Retained15      int       `xorm:"int(20) notnull comment('15日留存')" json:"retained15" desc:"15日留存"`
	Added15         int       `xorm:"int(20) notnull comment('前15日新增')" json:"added15" desc:"前15日新增"`
	Retained15Ratio int       `xorm:"int(20) notnull comment('15日留存率')" json:"retained15_ratio" desc:"15日留存率"`
	Date30          time.Time `xorm:"DATE notnull comment('前30日')" json:"date30" desc:"前30日"`
	Retained30      int       `xorm:"int(20) notnull comment('30日留存')" json:"retained30" desc:"30日留存"`
	Added30         int       `xorm:"int(20) notnull comment('前30日新增')" json:"added30" desc:"前30日新增"`
	Retained30Ratio int       `xorm:"int(20) notnull comment('30日留存率')" json:"retained30_ratio" desc:"30日留存率"`
	CreateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 统计留存数
type StatisticsRetained struct {
	Id              int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date            time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	ActiveCount     int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Date1           time.Time `xorm:"DATE notnull comment('前日')" json:"date1" desc:"前日"`
	Retained1       int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Added1          int       `xorm:"int(20) notnull comment('前日新增')" json:"added1" desc:"前日新增"`
	Retained1Ratio  int       `xorm:"int(20) notnull comment('次日留存率')" json:"retained1_ratio" desc:"次日留存率"`
	Date3           time.Time `xorm:"DATE notnull comment('前3日')" json:"date3" desc:"前3日"`
	Retained3       int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Added3          int       `xorm:"int(20) notnull comment('前3日新增')" json:"added3" desc:"前3日新增"`
	Retained3Ratio  int       `xorm:"int(20) notnull comment('3日留存率')" json:"retained3_ratio" desc:"3日留存率"`
	Date7           time.Time `xorm:"DATE notnull comment('前7日')" json:"date7" desc:"前7日"`
	Retained7       int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	Added7          int       `xorm:"int(20) notnull comment('前7日新增')" json:"added7" desc:"前7日新增"`
	Retained7Ratio  int       `xorm:"int(20) notnull comment('7日留存率')" json:"retained7_ratio" desc:"7日留存率"`
	Date15          time.Time `xorm:"DATE notnull comment('前15日')" json:"date15" desc:"前15日"`
	Retained15      int       `xorm:"int(20) notnull comment('15日留存')" json:"retained15" desc:"15日留存"`
	Added15         int       `xorm:"int(20) notnull comment('前15日新增')" json:"added15" desc:"前15日新增"`
	Retained15Ratio int       `xorm:"int(20) notnull comment('15日留存率')" json:"retained15_ratio" desc:"15日留存率"`
	Date30          time.Time `xorm:"DATE notnull comment('前30日')" json:"date30" desc:"前30日"`
	Retained30      int       `xorm:"int(20) notnull comment('30日留存')" json:"retained30" desc:"30日留存"`
	Added30         int       `xorm:"int(20) notnull comment('前30日新增')" json:"added30" desc:"前30日新增"`
	Retained30Ratio int       `xorm:"int(20) notnull comment('30日留存率')" json:"retained30_ratio" desc:"30日留存率"`
	CreateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime      time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 宝箱统计
type StatisticsTreasureBox struct {
	Id               int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date             time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	OpenPv           int       `xorm:"int(20) notnull comment('点击宝箱PV')" json:"open_pv" desc:"点击宝箱PV"`
	OpenUv           int       `xorm:"int(20) notnull comment('点击宝箱UV')" json:"open_uv" desc:"点击宝箱UV"`
	DayOpenedCount   int       `xorm:"int(20) notnull comment('每日打开宝箱数量')" json:"day_opened_count" desc:"每日打开宝箱数量"`
	TotalOpenedCount int       `xorm:"int(20) notnull comment('累计打开宝箱数量')" json:"total_opened_count" desc:"累计打开宝箱数量"`
	DayCdtCount      int       `xorm:"int(20) notnull comment('每日获取CDT次数')" json:"day_cdt_count" desc:"每日获取CDT次数"`
	TotalCdtCount    int       `xorm:"int(20) notnull comment('累计获取CDT次数')" json:"total_cdt_count" desc:"累计获取CDT次数"`
	DayCdt           float32   `xorm:"decimal(12,4) notnull comment('每日产出CDT数量')" json:"day_cdt" desc:"每日产出CDT数量"`
	TotalDayCdt      float32   `xorm:"decimal(12,4) notnull comment('累计产出CDT数量')" json:"total_day_cdt" desc:"累计产出CDT数量"`
	CreateTime       time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime       time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 宝箱实时统计
type StatisticsRealTreasureBox struct {
	Id             int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date           time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Hour           int       `xorm:"int(20) notnull comment('小时')" json:"hour" desc:"小时"`
	OpenedCount    int       `xorm:"int(20) notnull comment('打开宝箱数量')" json:"opened_count" desc:"打开宝箱数量"`
	DayOpenedCount int       `xorm:"int(20) notnull comment('每日打开宝箱数量')" json:"day_opened_count" desc:"每日打开宝箱数量"`
	Ratio          int       `xorm:"int(20) notnull comment('占比')" json:"ratio" desc:"占比"`
	CreateTime     time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime     time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 签到统计
type StatisticsSignIn struct {
	Id                int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date              time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Pv                int       `xorm:"int(20) notnull comment('参与人数PV')" json:"pv" desc:"参与人数PV"`
	Uv                int       `xorm:"int(20) notnull comment('参与人数UV')" json:"uv" desc:"参与人数UV"`
	ActiveCount       int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	UsedMakeUpCardNum int       `xorm:"int(20) notnull comment('消耗补签卡')" json:"opened_count" desc:"消耗补签卡"`
	S1                int       `xorm:"int(20) notnull comment('第一天签到人数')" json:"s1" desc:"第一天签到人数"`
	S2                int       `xorm:"int(20) notnull comment('第二天签到人数')" json:"s2" desc:"第二天签到人数"`
	S3                int       `xorm:"int(20) notnull comment('第三天签到人数')" json:"s3" desc:"第三天签到人数"`
	S4                int       `xorm:"int(20) notnull comment('第四天签到人数')" json:"s4" desc:"第四天签到人数"`
	S5                int       `xorm:"int(20) notnull comment('第五天签到人数')" json:"s5" desc:"第五天签到人数"`
	S6                int       `xorm:"int(20) notnull comment('第六天签到人数')" json:"s6" desc:"第六天签到人数"`
	S7                int       `xorm:"int(20) notnull comment('第七天签到人数')" json:"s7" desc:"第七天签到人数"`
	CreateTime        time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime        time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春CDT兑换
type StatisticsDoubleYearCdt struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date        time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Pv          int       `xorm:"int(20) notnull comment('参与人数PV')" json:"pv" desc:"参与人数PV"`
	Uv          int       `xorm:"int(20) notnull comment('参与人数UV')" json:"uv" desc:"参与人数UV"`
	ActiveCount int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Retained1   int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Retained3   int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Retained7   int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	Cdt         float32   `xorm:"decimal(12,4) notnull comment('当日兑换CDT')" json:"cdt" desc:"当日兑换CDT"`
	CreateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春CDT兑换用户日记录
type StatisticsDoubleYearUserDayCdt struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date       time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	UserId     string    `xorm:"int(20) notnull comment('用户ID')" json:"user_id" `
	Cdt        float32   `xorm:"decimal(12,4) notnull comment('当日兑换CDT')" json:"cdt" desc:"当日兑换CDT"`
	UserName   string    `xorm:"varchar(255) notnull comment('用户名')" json:"user_name" `
	TotalCdt   float32   `xorm:"decimal(12,4) notnull comment('累计兑换CDT')" json:"total_cdt" desc:"累计兑换CDT"`
	CreateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动碎片兑换PVUV记录
type StatisticsDoubleYearFragment struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date        time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Pv          int       `xorm:"int(20) notnull comment('参与人数PV')" json:"pv" desc:"参与人数PV"`
	Uv          int       `xorm:"int(20) notnull comment('参与人数UV')" json:"uv" desc:"参与人数UV"`
	ActiveCount int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Retained1   int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Retained3   int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Retained7   int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	CreateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动用户碎片兑换记录
type StatisticsDoubleYearUserFragment struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date       time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	UserId     string    `xorm:"int(20) notnull comment('用户ID')" json:"user_id" `
	UserName   string    `xorm:"varchar(255) notnull comment('用户名')" json:"user_name" `
	Count      int       `xorm:"int(20) notnull comment('当日获得量')" json:"count" desc:"当日获得量"`
	TotalCount int       `xorm:"int(20) notnull comment('累计获得量')" json:"total_count" desc:"累计获得量"`
	CreateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动日排行PVUV
type StatisticsDoubleYearDailyRanking struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date        time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Pv          int       `xorm:"int(20) notnull comment('参与人数PV')" json:"pv" desc:"参与人数PV"`
	Uv          int       `xorm:"int(20) notnull comment('参与人数UV')" json:"uv" desc:"参与人数UV"`
	ActiveCount int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Retained1   int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Retained3   int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Retained7   int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	Cdt         float32   `xorm:"decimal(12,4) notnull comment('CDT奖励产出')" json:"cdt" desc:"CDT奖励产出"`
	CreateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动用户日排行
type StatisticsDoubleYearUserDailyRanking struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date       time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	UserId     string    `xorm:"int(20) notnull comment('用户ID')" json:"user_id" `
	UserName   string    `xorm:"varchar(255) notnull comment('用户名')" json:"user_name" `
	Scores     int       `xorm:"int(20) notnull comment('新春值')" json:"scores" desc:"新春值"`
	InvitedNum int       `xorm:"int(20) notnull comment('邀请人数')" json:"invited_num" desc:"邀请人数"`
	Ranking    int       `xorm:"int(20) notnull comment('排名')" json:"ranking" desc:"排名"`
	CreateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动总排行PVUV
type StatisticsDoubleYearTotalRanking struct {
	Id          int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date        time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	Pv          int       `xorm:"int(20) notnull comment('参与人数PV')" json:"pv" desc:"参与人数PV"`
	Uv          int       `xorm:"int(20) notnull comment('参与人数UV')" json:"uv" desc:"参与人数UV"`
	ActiveCount int       `xorm:"int(20) notnull comment('活跃用户数')" json:"active_count" desc:"活跃用户数"`
	Retained1   int       `xorm:"int(20) notnull comment('次日留存')" json:"retained1" desc:"次日留存"`
	Retained3   int       `xorm:"int(20) notnull comment('3日留存')" json:"retained3" desc:"3日留存"`
	Retained7   int       `xorm:"int(20) notnull comment('7日留存')" json:"retained7" desc:"7日留存"`
	CreateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime  time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}

// 新春活动用户总排行
type StatisticsDoubleYearUserTotalRanking struct {
	Id         int64     `xorm:"int(20) autoincr pk" json:"id"`
	Date       time.Time `xorm:"DATE notnull comment('日期')" json:"date" desc:"日期"`
	UserId     string    `xorm:"int(20) notnull comment('用户ID')" json:"user_id" `
	UserName   string    `xorm:"varchar(255) notnull comment('用户名')" json:"user_name" `
	Scores     int       `xorm:"int(20) notnull comment('新春值')" json:"scores" desc:"新春值"`
	InvitedNum int       `xorm:"int(20) notnull comment('邀请人数')" json:"invited_num" desc:"邀请人数"`
	Ranking    int       `xorm:"int(20) notnull comment('排名')" json:"ranking" desc:"排名"`
	CreateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('创建时间')" json:"-"`
	UpdateTime time.Time `xorm:"timestamp notnull default(CURRENT_TIMESTAMP) updated(CURRENT_TIMESTAMP) comment('更新时间')" json:"-"`
}
