package double_year

import (
	"encoding/json"
	"game_server/core/base"
	"game_server/core/logger"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/message/statistical"
	"game_server/game/model"
	"game_server/game/proto"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

const (
	// RankListDayKey 每日排行榜key.
	RankListDayKey = "rank_list_day"
	// RankListAllKey 总排行榜key.
	RankListAllKey = "rank_list_all"
	// RankListUser 个人每日收集道具key, 格式为: rank_list_prop_用户id.
	RankListUserPropKey = "rank_list_prop_"
	// RankListUserAllPropKey 个人总收集道具key, 格式为: rank_list_all_prop_用户id.
	RankListUserAllPropKey = "rank_list_all_prop_"

	// RankListBaseScores 最低上榜分数值.
	RankListBaseScores = float64(600)
	// RankListTotalBaseScores 总排行榜发放奖励.
	RankListTotalBaseScores = float64(6000)

	// PropHose 袜子. 红包
	PropHose = 6101
	// PropHose 帽子. 福
	PropHat = 6102
	// PropTruncheon 权杖. 铜币
	PropTruncheon = 6103
	// PropLantern 灯笼.  鞭炮
	PropLantern = 6104
	// PropGiftBox 礼盒. 饺子
	PropGiftBox = 6105
	// PropInviteFriend 邀请好友.
	PropInviteFriend = 6008

	// PropOnlineTime 累计在线时长, 这是作为累计在线时长总称，用于在心跳处理.
	PropOnlineTime = 6009
	// PropOnlineTimeTen 0~10分钟.
	PropOnlineTimeTen = 6010
	// PropOnlineTimeTwo 0~20分钟.
	PropOnlineTimeTwo = 6011
	// PropOnlineTimeThree 20~30分钟.
	PropOnlineTimeThree = 6012
	// PropOnlineTimeFour 30~40分钟.
	PropOnlineTimeFour = 6013
	// PropOnlineTimeFive 40~50分钟.
	PropOnlineTimeFive = 6014
	// PropOnlineTimeHundred 50~100分钟.
	PropOnlineTimeHundred = 6015
	// PropOnlineTimeMore 200+分钟.
	PropOnlineTimeMore = 6016
)

var (
	// 任务列表
	TaskList = []int{PropHose, PropHat, PropTruncheon, PropLantern, PropGiftBox, PropInviteFriend, PropOnlineTime, PropOnlineTimeTen, PropOnlineTimeTwo, PropOnlineTimeThree, PropOnlineTimeFour, PropOnlineTimeFive, PropOnlineTimeHundred, PropOnlineTimeMore}
)

// RankListDay 排行榜struct.
type RankList struct {
	B *Base
}

// UpdateProp 更新道具数据.
// 这里需要更新每天和总的道具信息.
func (r *RankList) UpdateProp(userId string, propNumber, count int) {
	if r.GetRandListStatus() != proto.ACTIVITY_START {
		// 活动未开启
		//logger.Errorf("UpdateProp is close\n")
		return
	}
	r.innerUpdateProp(RankListUserPropKey, userId, propNumber, count)
	r.innerUpdateProp(RankListUserAllPropKey, userId, propNumber, count)
}

// innerUpdateProp 公共更新新道具数据方法.
func (r *RankList) innerUpdateProp(key, userId string, propNumber, count int) {
	prop := r.getProp(key, userId, propNumber)

	lastCount := int(prop["count"].(float64)) + count
	var point int
	client := db.RedisMgr.GetRedisClient()
	// 计算每种道具增加的分数.
	switch propNumber {
	// 袜子,帽子,拐杖,灯笼,礼盒这几种道具是20个即可增加分数.
	case PropHose, PropHat, PropTruncheon, PropLantern, PropGiftBox:
		if base := r.requireNumber(propNumber); lastCount >= base {
			prop["count"] = lastCount % base
			point = r.propNumberToPoint(propNumber)
			prop["point"] = int(prop["point"].(float64)) + point
			prop["used"] = int(prop["used"].(float64)) + 1
		} else {
			prop["count"] = lastCount
		}
	// 邀请好友,增加一个好友就可以增加分数,所以count字段不需要累计.
	case PropInviteFriend:
		point = r.propNumberToPoint(propNumber)
		prop["point"] = int(prop["point"].(float64)) + point
		prop["used"] = int(prop["used"].(float64)) + 1
	// 累计在线增加分数，分数根据不同的在线时长计算.
	case PropOnlineTime:
		// 更新在线时长累加值.
		prop["count"] = lastCount

		// 具体时间段记录.
		specialPropNumber, changePoint := r.onlineTimeToPoint(lastCount)
		if changePoint > 0 {
			specialProp := r.getProp(key, userId, specialPropNumber)
			if specialProp["point"].(float64) == 0 {
				specialProp["point"] = changePoint
				specialProp["used"] = 1
				specialProp["number"] = specialPropNumber
				specialProp["name"] = r.propNumberToZhName()[specialPropNumber]
				specialProp["count"] = r.requireNumber(specialPropNumber)
				point = changePoint

				specialPropStr, _ := json.Marshal(specialProp)
				_, err := client.HSet(key+userId, strconv.Itoa(specialPropNumber), specialPropStr).Result()
				if err != nil {
					logger.Errorf("UpdateProp[%v]: update user[%v] specialOnline[%v] prop[%v] count[%v] fail:%v\n", key, userId, specialPropNumber, changePoint, err)
				}
			}
		}
	}

	propStr, _ := json.Marshal(prop)

	// 更新用户每日道具信息.
	_, err := client.HSet(key+userId, strconv.Itoa(propNumber), propStr).Result()
	if err != nil {
		logger.Errorf("UpdateProp[%v]: update user[%v] prop[%v] count[%v] point[%v] fail:%v\n", key, userId, propNumber, count, point, err)
	}
	logger.Infof("UpdateProp[%v]: update user[%v] prop[%v] count[%v] point[%v] success!\n", key, userId, propNumber, count, point)

	// 只需更新一次，在每日道具更新.
	if point > 0 && key == RankListUserPropKey {
		// 更新排行榜数据.
		model.NewRankList().UpdateDayRankList(userId, time.Now().Format("20060102"), point)
		//// 更新每日排行榜集合.
		//r.updateRankList(RankListDayKey, userId, point)
		//// 更新总排行榜集合.
		//r.updateRankList(RankListAllKey, userId, point)
	}
}

// getProp 获取个人单个道具收集信息.
func (r *RankList) getProp(key, userId string, propNumber int) map[string]interface{} {
	client := db.RedisMgr.GetRedisClient()
	propStr, err := client.HGet(key+userId, strconv.Itoa(propNumber)).Result()
	if err != nil {
		// key不存在，初始化key对应的结构.
		if err == redis.Nil {
			logger.Infof(key + userId + " key does not exists")
			initP := r.initProp(key, userId)
			for k, v := range initP {
				if k == propNumber {
					return v.(map[string]interface{})
				}
			}
		}
		logger.Errorf("getProp from redis fail:%v\n", err)
	}

	rsp := make(map[string]interface{})
	err = json.Unmarshal([]byte(propStr), &rsp)
	if err != nil && err != redis.Nil {
		logger.Errorf("json.Unmarshal of getProp fail:userId=%v,propNumber=%v,err=%v,propStr=%v\n", userId, propNumber, err, propStr)
		return map[string]interface{}{
			// 注：这里的转换为float64是因为从redis中取出的整数默认为float64,这里统一起来是为了兼容处理.
			"number": float64(0),
			"name":   r.propNumberToZhName()[propNumber],
			"count":  float64(0),
			"used":   float64(0),
			"point":  float64(0),
		}
	}

	return rsp
}

// initProp 初始化用户收集道具结构到redis缓存中.
// 累计时长和邀请好友也当成是一种道具.
// 结构如下：
// rank_list_userId:[
// 		"道具编号":[
// 			"number": 道具编号,
// 			"name": 在前端展示的名称.
// 			"count": 收集未得到奖励的数量，如果是累计在线时长，则是最新累计时长.
// 			"used": 已经对换得到积分次数.
// 			"point": 已经对换得到积分值.
// 		],
// 		...,
// ]
func (r *RankList) initProp(key, userId string) map[int]interface{} {
	p := r.propNumberToZhName()
	result := make(map[int]interface{})
	client := db.RedisMgr.GetRedisClient()
	for k, v := range p {
		kStr := strconv.Itoa(k)
		prop := map[string]interface{}{
			"number": float64(k),
			"name":   v,
			// 注：这里的转换为float64是因为从redis中取出的整数默认为float64,这里统一起来是为了兼容处理.
			"count": float64(0),
			"used":  float64(0),
			"point": float64(0),
		}
		propByte, err := json.Marshal(prop)
		if err != nil {
			logger.Errorf("json.Marshal initProp propNumber[%v] fail:%v\n", k, err)
		}
		_, err = client.HSet(key+userId, kStr, string(propByte)).Result()
		if err != nil {
			logger.Errorf("initProp cache to redis fail:%v\n", err)
		}

		result[k] = prop
	}
	// 设置每日和总的过期时间, 这里添加60秒是为了防止排行榜奖励发放时，道具数据不丢失.
	if key == RankListUserPropKey {
		client.Expire(key+userId, r.GetCacheExpireTime(1))
	} else {
		client.Expire(key+userId, r.GetCacheExpireTime(16))
	}

	return result
}

func (r *RankList) GetCacheExpireTime(day int) time.Duration {
	timeStr := time.Now().Format("2006-01-02")
	t, _ := time.ParseInLocation("2006-01-02", timeStr, time.Local)
	nextDay := t.Unix() + int64(86400*day)

	return time.Duration(nextDay-time.Now().Unix()+10) * time.Second
}

// GetDayProp 获取个人每日道具收集信息.
func (r *RankList) GetDayProp(userId string) map[string]interface{} {
	return r.getRankProp(RankListUserPropKey, userId)
}

// GetAllProp 获取个人活动时间总的道具收集信息.
func (r *RankList) GetAllProp(userId string) map[string]interface{} {
	return r.getRankProp(RankListUserAllPropKey, userId)
}

// getRankProp 获取每日和总双旦记录, 根据前端UI设计组装数据.
func (r *RankList) getRankProp(key, userId string) map[string]interface{} {
	var list []map[string]interface{}
	//for k := range r.propNumberToZhName() {

	// 这样写是为了按照道具编号排序.
	allProp := r.propNumberToZhName()
	//for k := PropHose; k <= PropOnlineTimeMore; k++ {
	for _, k := range TaskList {
		if k == PropOnlineTime {
			continue
		}

		// 由于是连续的整数，可能不存在，需要判断下.
		_, ok := allProp[k]
		if !ok {
			continue
		}

		data := r.getProp(key, userId, k)
		// 显示所有双旦记录，包括没有对换的.
		//if data["used"].(float64) == 0 {
		//	continue
		//}
		list = append(list,
			map[string]interface{}{
				"name":          data["name"],                                       // 道具说明.
				"point":         data["point"],                                      // 已加分数.
				"used":          data["used"],                                       // 已经完成次数.
				"award":         r.propNumberToPoint(int(data["number"].(float64))), // 道具对应的奖励.
				"requireNumber": r.requireNumber(int(data["number"].(float64))),     // 得到一次奖励需要完成的次数.
				//"completedNumber": math.Floor(data["count"].(float64) / 60),           // 得到一次奖励需已经完成的次数.
				"completedNumber": data["count"],                                       // 得到一次奖励需已经完成的次数.
				"isInviteFriend":  k == PropInviteFriend,                               // 是否为邀请好友类型.
				"isShowTag":       r.propNumberIsRepeat(int(data["number"].(float64))), // 是否展示“可重复”标签.
			},
		)
	}

	// 按道具编号进行排序.

	return map[string]interface{}{
		"list": list,
	}
}

// InviteRecord 邀请好友记录.
func (r *RankList) InviteRecord(userId string) map[string]interface{} {
	//data, _ := logic.UserInvitation{}.GetUserInvitList(userId)
	return map[string]interface{}{}
}

// getRankList 获取指定区间排行榜数据.
// 由于排行榜由两部分组成，一部分是分数排前10的用户，
// 另一部分是当前用户的分数、排名及距离前一名用户分数之差，
// 这里需要遍历整个集合, 只需遍历一次.
func (r *RankList) getRankList(rankKey, propKey, userId string, offset, limit int64, invitedFriend int) map[string]interface{} {
	var tmp, rankList []map[string]interface{}
	result := map[string]interface{}{
		"rankList":       rankList,
		"userRank":       0,
		"frontUserPoint": 0,
		"userPoint":      0,
	}

	// 排行榜数据,至少要邀请一名好友，而邀请一名好友需要200分数，这里只查找大于等于200分数的元素.
	client := db.RedisMgr.GetRedisClient()
	//zUsers, err := client.ZRevRangeByScoreWithScores(rankKey, redis.ZRangeBy{Min: "-inf", Max: "+inf"}).Result()
	var zUsers []map[string]interface{}
	var err error
	if rankKey == RankListDayKey {
		zUsers, err = model.NewRankList().GetDayRankList(false)
	} else {
		zUsers, err = model.NewRankList().GetTopRankList()
	}
	if err != nil {
		logger.Errorf("get rank list by key[%v] fail:%v\n", rankKey, err)
		return result
	}

	userRank := int64(0) // userId对应的排名.
	userPoint := 0       // userId对应的分数.
	formerPoint := 0     // 上一名的分数.
	frontUserPoint := 0  // 距离上一名分数之差.
	rank := int64(1)     // 排行榜名次.
	for _, user := range zUsers {
		if user["userId"].(string) == userId {
			userPoint = user["scores"].(int)
		}

		// 判断是否邀请了好友且达到指定数量.
		propStr, err := client.HGet(propKey+user["userId"].(string), strconv.Itoa(PropInviteFriend)).Result()
		if err != nil {
			logger.Errorf("get prop of user[%v] fail:%v\n", user["userId"].(string), err)
			continue
		}
		prop := make(map[string]interface{})
		_ = json.Unmarshal([]byte(propStr), &prop)
		if int(prop["used"].(float64)) < invitedFriend {
			logger.Errorf("user[%v] invitedFriend not enough:%v\n", user["userId"].(string), err)
			continue
		}

		// 获取nickname.
		nickname, _ := client.HGet(user["userId"].(string), "nick_name").Result()
		// 排行集合数据.
		tmp = append(tmp, map[string]interface{}{"rank": rank, "name": nickname, "scores": user["scores"].(int), "award": 0})

		if user["userId"].(string) == userId {
			userRank = rank
			frontUserPoint = formerPoint - userPoint
			if userRank == 1 {
				frontUserPoint = 0
			}
		}
		// 上一名的分数.
		formerPoint = int(user["scores"].(int))

		rank++
	}

	//调试信息，打印排行榜所有用户信息
	logger.Debugf("getRankList:[%v] %+v\n", userId, tmp)

	// 匹配名次对应的奖励, 只取前10名.
	for k, v := range tmp {
		if k == 10 {
			break
		}
		v["award"] = r.rankListAwardConfig(rankKey, int(v["rank"].(int64)))
		rankList = append(rankList, v)
	}

	return map[string]interface{}{
		"rankList":       rankList,
		"userRank":       userRank,
		"frontUserPoint": frontUserPoint,
		"userPoint":      userPoint,
	}
}

// GetDayRankList 获取每日排行榜.
func (r *RankList) GetDayRankList(userId string) map[string]interface{} {
	return r.getRankList(RankListDayKey, RankListUserPropKey, userId, 0, 9, 1)
}

// GetAllRankList 获取总排行榜.
func (r *RankList) GetAllRankList(userId string) map[string]interface{} {
	return r.getRankList(RankListAllKey, RankListUserAllPropKey, userId, 0, 9, 10)
}

// getUserRankIndex 返回用户在每日排行榜和总榜名次.
// 按照分数由高到低顺序取名次.
func (r *RankList) getUserRankIndex(key, userId string) int {
	client := db.RedisMgr.GetRedisClient()
	index, err := client.ZRevRank(key, userId).Result()
	if err != nil {
		// 当前用户没有排名.
		if err == redis.Nil {
			return -1
		}
		logger.Errorf("get user[%v] index of key[%v] fail:%v\n", userId, key, err)
		// 这里返回一个自定义值.
		return 15
	}
	return int(index)
}

// updateRankList 根据排行榜数据.
// 根据key不同，更新对应的数据，包括每日排行榜和总排行榜.
// 注：每日排行榜expire和总排行榜expire不一样.
func (r *RankList) updateRankList(key, userId string, point int) {
	client := db.RedisMgr.GetRedisClient()
	// 注：这里的redis.Z{}根据redis不同版本不一样，可能是&redis.Z{}
	_, err := client.ZIncr(key, redis.Z{
		Score:  float64(point),
		Member: userId,
	}).Result()
	if err != nil {
		logger.Errorf("updateRankList[%v] update user[%v] point[%v] fail:%v\n", key, userId, point, err)
	}

	// 设置过期时间
	n, _ := client.ZCard(key).Result()
	if n <= 1 {
		var expire time.Duration
		if key == RankListDayKey {
			expire = db_service.NewCdt().GetCacheExpireTime(1)
		} else {
			expire = db_service.NewCdt().GetCacheExpireTime(10)
		}
		client.Expire(key, expire)
	}

	logger.Infof("updateRankList[%v] update user[%v] point[%v] success\n", key, userId, point)
}

// GiveOutDayAward 发放每日排行榜奖励.
// 规则：每日排行榜前10名，且分数至少为600，至少邀请1名好友, 发放cdt最多不超过100个.
func (r *RankList) GiveOutDayAward() {
	r.award(RankListDayKey, RankListUserPropKey, "每日排行榜奖励.", "恭喜获得每日排行榜的第X名，以下是我们为你准备的排名奖励，请收下吧.", 1, int(r.limitDayTopCdt()), proto.MSG_RANK_LIST_DAY_CDT, int(r.limitDayTopCdt()), RankListBaseScores)
}

// GiveOutTotalAward 发放总排行榜奖励.
// 规则：累计排行榜前10名，且分数至少为600，至少邀请10名好友, 发放cdt最多不超过1000个.
func (r *RankList) GiveOutTotalAward() {
	r.award(RankListAllKey, RankListUserAllPropKey, "总排行榜奖励.", "恭喜获得总排行榜的第X名，以下是我们为你准备的排名奖励，请收下吧.", 10, int(r.limitAllTopCdt()), proto.MSG_RANK_LIST_ALL_CDT, int(r.limitAllTopCdt()), RankListTotalBaseScores)
}

// award 发放奖励.
// 奖励字段:
// "prizeList":[
// 		"prizeType": 0, // 奖励类型
// 		"prizeId": 0, // 奖励id,这个id为t_items表中的id
// 		"prizeName": 0, // 奖励名称
// 		"prizeNum": 0, // 奖励数量
// 		"prizeImg": 0, // 奖励图片
// ]
func (r *RankList) award(rankKey, propKey, title, content string, invitedFriend int, totalCdt int, eventType int, cdtLimit int, scores float64) {
	// 判断活动是否开启.
	currentTime := time.Now()
	if currentTime.Unix() < utils.Str2Time(base.Setting.Springfestival.RankingListStartDate+" "+"00:00:00").Unix() || currentTime.Unix() >= utils.Str2Time(base.Setting.Springfestival.RankingListEndDate+" 00:10:00").Unix() {
		return
	}

	award := map[string]interface{}{
		"userId":       "",
		"emailType":    2,
		"emailTitle":   title,
		"emailContent": content,
	}

	// 排行榜数据.
	client := db.RedisMgr.GetRedisClient()
	//zUsers, err := client.ZRevRangeWithScores(rankKey, 0, 9).Result()
	var zUsers []map[string]interface{}
	var err error
	day_ranklist_time := ""
	if rankKey == RankListDayKey {
		zUsers, err = model.NewRankList().GetDayRankList(true)
		day_ranklist_time = model.GetDate(true)
	} else {
		zUsers, err = model.NewRankList().GetTopRankList()
	}
	if err != nil {
		logger.Errorf("get rank list by key[%v] fail:%v\n", rankKey, err)
	}

	// 获取cdt图片.
	cdtImgUrl := r.getCdtImg()

	email := db_service.NewEmailLogic()
	var baseCdt int
	rank := 1 // 排行榜名次.
	for _, user := range zUsers {
		if rank > 10 {
			return
		}
		// 判断分数是否达到600.
		if float64(user["scores"].(int)) < scores {
			logger.Errorf("user[%v] scores is not enough:\n", user["userId"].(string))
			continue
		}

		// 判断是否邀请了好友且达到指定数量.
		propStr, err := client.HGet(propKey+user["userId"].(string), strconv.Itoa(PropInviteFriend)).Result()
		if err != nil {
			logger.Errorf("get prop of user[%v] fail:%v\n", user["userId"].(string), err)
			continue
		}
		prop := make(map[string]interface{})
		_ = json.Unmarshal([]byte(propStr), &prop)
		if int(prop["used"].(float64)) < invitedFriend {
			logger.Errorf("user[%v] invitedFriend[%v] not enough\n", user["userId"].(string), int(prop["used"].(float64)))
			continue
		}

		// 判断排行榜发放的cdt奖励是否达到上限值.
		awardNumber := r.rankListAwardConfig(rankKey, rank)
		if baseCdt >= cdtLimit {
			logger.Infof("rankListAward[%v] is limited!\n", rankKey)
			break
		}
		baseCdt += awardNumber

		award["userId"] = user["userId"].(string)
		msgType, _ := json.Marshal(map[string]interface{}{"eventType": eventType})
		prizeList := map[string]interface{}{
			"prizeId":   0,
			"prizeType": 2,
			"prizeName": "排名第" + strconv.Itoa(rank) + "奖励 CDT*1",
			"prizeNum":  awardNumber,
			"prizeImg":  cdtImgUrl,
			"extend":    string(msgType),
		}
		award["prizeList"] = []map[string]interface{}{prizeList}
		if rankKey == RankListDayKey {
			award["emailContent"] = "恭喜获得每日排行榜的第" + strconv.Itoa(rank) + "名，以下是我们为你准备的排名奖励，请收下吧"
		} else {
			award["emailContent"] = "恭喜获得总排行榜的第" + strconv.Itoa(rank) + "名，以下是我们为你准备的排名奖励，请收下吧"
		}
		rank++

		// 发送邮件.
		_, status := email.AddEmail(award)
		if !status {
			logger.Errorf("send award[%v] of user[%v] email fail\n", RankListDayKey, award["userId"])
			continue
		}

		statistical.StatisticsDotIns.DoubleYearDailyRankingCdt(day_ranklist_time, awardNumber)
	}
}

// RepairAwardCdt 补发排行榜奖励.
func (r *RankList) RepairAwardCdt(userId, emailTitle, emailContent string, eventType int, cdtValue int) bool {
	award := map[string]interface{}{
		"userId":       userId,
		"emailType":    2,
		"emailTitle":   emailTitle,
		"emailContent": emailContent,
	}
	msgType, _ := json.Marshal(map[string]interface{}{"eventType": eventType})
	prizeList := map[string]interface{}{
		"prizeId":   0,
		"prizeType": 2,
		"prizeName": "",
		"prizeNum":  cdtValue,
		"prizeImg":  r.getCdtImg(),
		"extend":    string(msgType),
	}
	award["prizeList"] = []map[string]interface{}{prizeList}
	_, status := db_service.NewEmailLogic().AddEmail(award)
	if !status {
		logger.Errorf("send repair user[%v] cdt[%v] email[%v] fail\n", userId, cdtValue, emailTitle)
		return false
	}
	return true
}

// getCdtImg 获取cdt图片，用于邮件中显示.
func (r *RankList) getCdtImg() string {
	items := &db_service.Items{}
	item, err := items.GetDataById(1001)
	if err != nil {
		return ""
	}
	return item.ImgUrl
}

// limitDayTopCdt 每日排行榜cdt奖励限制值.
func (r *RankList) limitDayTopCdt() float32 {
	return base.Setting.Doubleyear.TradeCdtDayTopLimit
}

// limitAllTopCdt 总排行榜cdt奖励限制值.
func (r *RankList) limitAllTopCdt() float32 {
	return float32(1000)
	//return base.Setting.Doubleyear.TradeCdtDayTopTotalLimit
}

// propNumberToZhName 道具编号映射中文活动名称，用于给前端展示.
func (r *RankList) propNumberToZhName() map[int]string {
	return map[int]string{
		PropHose:              "收集红包数量达到20个",
		PropHat:               "收集福数量达到20个",
		PropTruncheon:         "收集铜币数量达到20个",
		PropLantern:           "收集鞭炮数量达到20个",
		PropGiftBox:           "收集饺子数量达到20个",
		PropInviteFriend:      "成功邀请1名好友",
		PropOnlineTime:        "累计在线时长",
		PropOnlineTimeTen:     "累计在线时长-10分钟",
		PropOnlineTimeTwo:     "累计在线时长-20分钟",
		PropOnlineTimeThree:   "累计在线时长-30分钟",
		PropOnlineTimeFour:    "累计在线时长-40分钟",
		PropOnlineTimeFive:    "累计在线时长-50分钟",
		PropOnlineTimeHundred: "累计在线时长-100分钟",
		PropOnlineTimeMore:    "累计在线时长-200分钟",
	}
}

// propNumberToPoint 每种道具对应分数奖励.
func (r *RankList) propNumberToPoint(propNumber int) int {
	propPoint := map[int]int{
		PropHose:              10,
		PropHat:               5,
		PropTruncheon:         10,
		PropLantern:           20,
		PropGiftBox:           20,
		PropInviteFriend:      200,
		PropOnlineTime:        5,
		PropOnlineTimeTen:     5,
		PropOnlineTimeTwo:     5,
		PropOnlineTimeThree:   10,
		PropOnlineTimeFour:    15,
		PropOnlineTimeFive:    20,
		PropOnlineTimeHundred: 25,
		PropOnlineTimeMore:    50,
	}

	return propPoint[propNumber]
}

// propNumberIsRepeat 道具是否显示可重复对换.
func (r *RankList) propNumberIsRepeat(propNumber int) bool {
	propPoint := map[int]bool{
		PropHose:              true,
		PropHat:               true,
		PropTruncheon:         true,
		PropLantern:           true,
		PropGiftBox:           true,
		PropInviteFriend:      true,
		PropOnlineTime:        true,
		PropOnlineTimeTen:     false,
		PropOnlineTimeTwo:     false,
		PropOnlineTimeThree:   false,
		PropOnlineTimeFour:    false,
		PropOnlineTimeFive:    false,
		PropOnlineTimeHundred: false,
		PropOnlineTimeMore:    false,
	}

	return propPoint[propNumber]
}

// requireNumber 每种道具领取一次奖励需要达到次数.
func (r *RankList) requireNumber(propNumber int) int {
	propRequireNumber := map[int]int{
		PropHose:              20,
		PropHat:               20,
		PropTruncheon:         20,
		PropLantern:           20,
		PropGiftBox:           20,
		PropInviteFriend:      1,
		PropOnlineTime:        10,
		PropOnlineTimeTen:     10,
		PropOnlineTimeTwo:     20,
		PropOnlineTimeThree:   30,
		PropOnlineTimeFour:    40,
		PropOnlineTimeFive:    50,
		PropOnlineTimeHundred: 100,
		PropOnlineTimeMore:    200,
	}

	return propRequireNumber[propNumber]
}

// rankListAwardConfig 排行榜奖励配置.
func (r *RankList) rankListAwardConfig(rankType string, rank int) int {
	// rankDayAward 每日排行榜奖励配置.
	rankDayAward := map[int]int{
		1:  25,
		2:  20,
		3:  16,
		4:  12,
		5:  9,
		6:  7,
		7:  5,
		8:  3,
		9:  2,
		10: 1,
	}
	// rankAllAward 总排行榜奖励配置.
	rankAllAward := map[int]int{
		1:  500,
		2:  200,
		3:  100,
		4:  50,
		5:  40,
		6:  30,
		7:  20,
		8:  10,
		9:  8,
		10: 6,
	}

	var award int
	switch rankType {
	case RankListDayKey:
		award = rankDayAward[rank]
	case RankListAllKey:
		award = rankAllAward[rank]
	}
	return award
}

// onlineTimeToPoint 根据累计在线时长返回对应的时间段及分数.
func (r *RankList) onlineTimeToPoint(onlineTime int) (int, int) {
	switch {
	// 10~20分钟.
	case onlineTime >= 600 && onlineTime < 1200:
		return PropOnlineTimeTen, r.propNumberToPoint(PropOnlineTimeTen)
	// 20~30分钟.
	case onlineTime >= 1200 && onlineTime < 1800:
		return PropOnlineTimeTwo, r.propNumberToPoint(PropOnlineTimeTwo)
	// 30~40分钟.
	case onlineTime >= 1800 && onlineTime < 2400:
		return PropOnlineTimeThree, r.propNumberToPoint(PropOnlineTimeThree)
	// 40~50分钟.
	case onlineTime >= 2400 && onlineTime < 3000:
		return PropOnlineTimeFour, r.propNumberToPoint(PropOnlineTimeFour)
	// 50~100分钟.
	case onlineTime >= 3000 && onlineTime < 6000:
		return PropOnlineTimeFive, r.propNumberToPoint(PropOnlineTimeFive)
	// 100~200分钟.
	case onlineTime >= 6000 && onlineTime < 12000:
		return PropOnlineTimeHundred, r.propNumberToPoint(PropOnlineTimeHundred)
	// 200+分钟.
	case onlineTime > 12000:
		return PropOnlineTimeMore, r.propNumberToPoint(PropOnlineTimeMore)
	}
	return PropOnlineTimeTen, 0
}

//检查排行榜是否开启   true  开启  false 关闭
func (r *RankList) GetRandListStatus() int {
	startDate := utils.Str2Time(base.Setting.Springfestival.RankingListStartDate + "00:00:00")
	endDate := utils.Str2Time(base.Setting.Springfestival.RankingListEndDate + "00:00:00")
	curDate := time.Now()
	if curDate.Sub(startDate) < 0 {
		return proto.ACTIVITY_NOT_START
	}

	if curDate.Sub(endDate) >= 0 {
		return proto.ACTIVITY_END
	}

	return proto.ACTIVITY_START
}

// NewRankList 实例化RankList结构体.
func NewRankList() *RankList {
	return &RankList{
		B: NewBase(),
	}
}
