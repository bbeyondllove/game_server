package utils

import (
	"errors"
	"fmt"
	"game_server/core/logger"
	"game_server/db"
	"game_server/game/proto"
	rand2 "math/rand"
	"reflect"
	"runtime/debug"
	"strconv"
	"strings"
	"sync"
	"time"
)

const (
	ITEM_NORMAL = iota
	ITEM_DAY
	ITEM_WEEK
)
const (
	timeLayout = "2006-01-02 15:04:05"
)

//判断是否数字
func IsNum(s string) bool {
	_, err := strconv.ParseFloat(s, 64)
	return err == nil
}

/**
 * 生成单号sn
 */
func CreateOrderSn(prefix string) string {
	if prefix == "" {
		prefix = "L"
	}
	code := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	year := time.Now().Year()
	index := int(year) - 2017
	if index >= len(code) {
		index = 0
	}

	month := int(time.Now().Month())
	mCode := strings.ToUpper(fmt.Sprintf("%01X", month))

	day := strconv.Itoa(time.Now().Day())
	microTime := strconv.FormatInt(time.Now().UnixNano()/1e3, 10)

	r := rand2.New(rand2.NewSource(time.Now().UnixNano()))
	spri := fmt.Sprintf("%02d", r.Intn(100))

	orderSn := prefix + code[index] + mCode + day + microTime[5:15] + spri
	return orderSn
}

// addslashes() 函数返回在预定义字符之前添加反斜杠的字符串。
// 预定义字符是：
// 单引号（'）
// 双引号（"）
// 反斜杠（\）
func Addslashes(str string) string {
	tmpRune := []rune{}
	strRune := []rune(str)
	for _, ch := range strRune {
		switch ch {
		case []rune{'\\'}[0], []rune{'"'}[0], []rune{'\''}[0]:
			tmpRune = append(tmpRune, []rune{'\\'}[0])
			tmpRune = append(tmpRune, ch)
		default:
			tmpRune = append(tmpRune, ch)
		}
	}
	return string(tmpRune)
}

// stripslashes() 函数删除由php addslashes() 函数添加的反斜杠。
func Stripslashes(str string) string {
	dstRune := []rune{}
	strRune := []rune(str)
	strLenth := len(strRune)
	for i := 0; i < strLenth; i++ {
		if strRune[i] == []rune{'\\'}[0] {
			i++
		}
		dstRune = append(dstRune, strRune[i])
	}
	return string(dstRune)
}

func IsExistInArrs(val string, arrs []string) bool {
	if len(arrs) == 0 {
		return false
	}

	for _, v := range arrs {
		if v == val {
			return true
		}
	}
	return false
}

func Contain(obj interface{}, target interface{}) (bool, error) {
	targetValue := reflect.ValueOf(target)
	switch reflect.TypeOf(target).Kind() {
	case reflect.Slice, reflect.Array:
		for i := 0; i < targetValue.Len(); i++ {
			if targetValue.Index(i).Interface() == obj {
				return true, nil
			}
		}
	case reflect.Map:
		if targetValue.MapIndex(reflect.ValueOf(obj)).IsValid() {
			return true, nil
		}
	}

	return false, errors.New("not in array")
}

func ToSlice(arr interface{}) []interface{} {
	v := reflect.ValueOf(arr)
	if v.Kind() != reflect.Slice {
		panic("toslice arr not slice")
	}
	l := v.Len()
	ret := make([]interface{}, l)
	for i := 0; i < l; i++ {
		ret[i] = v.Index(i).Interface()
	}
	return ret
}

//slice去重
func RemoveRepByLoop(slc []string) []string {
	result := []string{} // 存放结果
	for i := range slc {
		flag := true
		for j := range result {
			if slc[i] == result[j] {
				flag = false // 存在重复元素，标识为false
				break
			}
		}
		if flag { // 标识为false，不添加进结果
			result = append(result, slc[i])
		}
	}
	return result
}

func GetMapLen(data *sync.Map) int {
	len := 0
	if data == nil {
		return len
	}
	data.Range(func(key interface{}, value interface{}) bool {
		len++
		return true
	})
	return len
}

//获取具体长度
func GetMapItemIdLen(data *sync.Map) map[int]int {
	lenMap := make(map[int]int, 0)
	if data == nil {
		return lenMap
	}
	data.Range(func(key interface{}, value interface{}) bool {
		eventVlue := value.(*proto.EventData)
		_, ok := lenMap[eventVlue.EventId]
		if ok {
			lenMap[eventVlue.EventId]++
		} else {
			lenMap[eventVlue.EventId] = 1
		}
		return true
	})
	return lenMap
}

func GetCurDay() string {
	t := time.Now()
	return t.Format("2006-01-02")
}

func GetUnixDay(timeValue int64) string {
	return time.Unix(timeValue, 0).Format("2006-01-02")
}

func GetTimeDay(t time.Time) string {
	return t.Format("2006-01-02")
}

func GetWeekDay() int {
	t := time.Now()
	weekDay := (int(t.Weekday()))
	if weekDay == 0 {
		weekDay = 7
	}
	return weekDay
}

/**字符串->时间对象*/
func Str2Time(formatTimeStr string) time.Time {
	loc, _ := time.LoadLocation("Local")
	theTime, _ := time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型
	return theTime

}

func ParseTime(formatTimeStr string) (time.Time, error) {
	loc, _ := time.LoadLocation("Local")
	return time.ParseInLocation(timeLayout, formatTimeStr, loc) //使用模板在对应时区转化为time.time类型
}

//获得当前时间，返顺字符串格式
func Time2Str(time_data time.Time) string {
	str := time_data.Format(timeLayout)
	return str
}

/**
获取本周周一的日期
*/
func GetFirstDateOfWeek() (weekMonday string) {
	now := time.Now()

	offset := int(time.Monday - now.Weekday())
	if offset > 0 {
		offset = -6
	}

	weekStartDate := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).AddDate(0, 0, offset)
	weekMonday = weekStartDate.Format(timeLayout)
	return
}

/**
获取上周的周一日期
*/
func GetLastWeekFirstDate() (weekMonday string) {
	thisWeekMonday := GetFirstDateOfWeek()
	TimeMonday := Str2Time(thisWeekMonday)
	weekMonday = TimeMonday.AddDate(0, 0, -7).Format(timeLayout)
	return
}

/**
获取下周的周一日期
*/
func GetNextWeekFirstDate() (weekMonday string) {
	thisWeekMonday := GetFirstDateOfWeek()
	TimeMonday := Str2Time(thisWeekMonday)
	weekMonday = TimeMonday.AddDate(0, 0, 7).Format(timeLayout)

	return
}

func SetKeyValue(key string, field string, value interface{}, bAdd bool, itemType int) (bool, int64, int64) {
	//logger.Debugf("SetExpInfo SetKeyValue(%+v:%+v:%+v)  ", key, field, value)
	valueRet, err := db.RedisGame.HGet(key, field).Result()
	oldValue := int64(0)
	newValue := int64(0)
	if err != nil && err.Error() != "redis: nil" {
		logger.Errorf("TaskProcess Get(", key, ") failed:%+v", err.Error())
		return false, oldValue, newValue
	}

	if valueRet != "" {
		oldValue, _ = strconv.ParseInt(valueRet, 10, 64)
	}

	if bAdd {
		newValue, err = db.RedisGame.HIncrBy(key, field, value.(int64)).Result()
		if err != nil {
			logger.Errorf("SetExpInfo IncrBy(%+v %+v:,) failed, err=%+v", key, err.Error())
			return false, oldValue, newValue
		}
	} else {
		_, err = db.RedisGame.HSet(key, field, value).Result()
		if err != nil {
			logger.Errorf("SetExpInfo Set(%+v %+v:,) failed, err=%+v", key, err.Error())
			return false, oldValue, newValue
		}
	}

	if itemType > ITEM_NORMAL {
		SetKeyExpInfo(key, itemType)
	}
	return true, oldValue, newValue
}

func SetKeyExpInfo(key string, itemType int) bool {
	var err error
	if itemType == ITEM_DAY {
		//设置当天晚上24点过期
		now := time.Now()
		next := now.Add(time.Hour * 24)
		next = time.Date(next.Year(), next.Month(), next.Day(), 0, 0, 0, 0, next.Location())
		tm := now.Add(next.Sub(now))
		_, err = db.RedisGame.ExpireAt(key, tm).Result()
	} else {
		//设置下周一0点过期
		nextMonday := GetNextWeekFirstDate()
		TimeMonday := Str2Time(nextMonday)
		_, err = db.RedisGame.ExpireAt(key, TimeMonday).Result()
	}
	if err != nil {
		logger.Errorf("SetExpInfo ExpireAt(%+v %+v:,) failed, err=%+v", key, err.Error())
		return false
	}

	return true
}

func SysRecoverWrap(f func()) func() {
	return func() {
		defer func() {
			if r := recover(); r != nil {
				logger.Errorf("SYSTEM ACTION PANIC: %v, stack: %v", r, string(debug.Stack()))
			}
		}()
		f()
	}
}

func ShowStack() {
	logger.Errorf("%+V", string(debug.Stack()))
}

func SetKeyTiemExp(key string, expTime string) bool {
	nextTime := Str2Time(expTime)
	_, err := db.RedisGame.ExpireAt(key, nextTime).Result()
	if err != nil {
		logger.Errorf("SetKeyTiemExp ExpireAt(%+v %+v:,) failed, err=%+v", key, err.Error())
		return false
	}

	return true
}
