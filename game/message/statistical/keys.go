package statistical

// 统计数据用的的 redis key
const (
	EVERY_HOUR_ACTIVE_PEOPLE_NUM = "st:hour_active_people_num:"     //每天每小时活跃人数
	PLATFORM_ONLINE_TIME         = "st:platform_online_time:"       // 单一平台在线时长
	ONLINE_TIME                  = "st:online_time:"                // 在线时长
	PLATFORM_LOGIN_NUM           = "st:platform_login_num:"         //当日单一平台登录数
	PLATFORM_ACTIVE_PEOPLE_NUM   = "st:platform_active_people_num:" //每天单一平台活跃人数
	LOGIN_NUM                    = "st:login_num:"                  //当日登录数
	ACTIVE_PEOPLE_NUM            = "st:active_people_num:"          //每天活跃人数

	BOX_DAY_OPEN_NUM        = "st:treasure_box_day_num:"             //每天点击宝箱的人数(去重)
	BOX_DAY_OPEN_PEOPLE_NUM = "st:treasure_box_day_people_num:"      //每天点击宝箱的人数(去重)
	BOX_DAY_FINISH_NUM      = "st:treasure_box_day_finish_num:"      // 每天完成宝箱的次数
	BOX_DAY_FINISH_HOUR_NUM = "st:treasure_box_day_finish_hour_num:" // 每天没消失完成宝箱的次数

	DOUBLEYEAR_FUWA              = "st:double_year_fuwa:"             // 每天兑换福娃数
	DOUBLEYEAR_FUWA_PV           = "st:double_year_fuwa_pv:"          // 每天福娃pv
	DOUBLEYEAR_FUWA_UV           = "st:double_year_fuwa_uv:"          // 每天福娃uv
	DOUBLEYEAR_DAY_PV            = "st:double_year_day_pv:"           // 新春日排行榜pv
	DOUBLEYEAR_DAY_UV            = "st:double_year_day_uv:"           // 新春日排行榜uv
	DOUBLEYEAR_TOTAL_PV          = "st:double_year_total_pv:"         // 新春总排行榜活动pv
	DOUBLEYEAR_TOTAL_UV          = "st:double_year_total_uv:"         // 新春总排行榜活动uv
	DOUBLEYEAR_DAILY_RANKING_CDT = "st:double_year_dailyranking_cdt:" // 新春日排行榜每日CDT
)
