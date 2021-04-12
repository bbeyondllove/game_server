package statistical_data

// 统计数据用的的 redis key
const (
	CITY_ICON_PV       = "set:city:icon:pv:"             // + 当前日期 日期格式 20060102
	CITY_ICON_UV       = "hset:city:icon:uv:"            // + 当前日期 日期格式 20060102
	BAOX_ADVERTIS      = "hset:effective_advertis:"      // 统计欢乐宝箱有效广告 当前日期 日期格式 20060102
	BAOX_ADVERTIS_USER = "hset:effective_advertis_user:" // 统计用户的 每天的 日期格式 20060102
)
