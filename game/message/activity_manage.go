package message

import (
	"fmt"
	"game_server/core/base"
	coretimer "game_server/core/timer"
	"game_server/core/utils"
	dao "game_server/game/db_service"
	"game_server/game/model"
	"game_server/game/proto"
	"sync"
	"time"
)

type ActivityManage struct {
	ActivityStatus map[int]int
	Activitys      map[int]model.ActivityConfig
	Timer          *coretimer.TimeWheel
	lock           sync.Mutex
	satus_lock     sync.Mutex
	time_lock      sync.Mutex
}

var (
	G_ActivityManage = ActivityManage{
		ActivityStatus: make(map[int]int, 0),
		Activitys:      make(map[int]model.ActivityConfig, 0),
		Timer:          coretimer.NewTimeWheel(),
	}
)

func (this *ActivityManage) Init() error {
	configs, err := dao.ActivityConfigIns.GetAllData()
	if err != nil {
		return err
	}
	for _, config := range configs {
		this.ChangeTime(config, false)
	}
	go this.Timer.Start()
	return nil
}

func (this *ActivityManage) GetActivityConfig(activity_type int) *model.ActivityConfig {
	this.lock.Lock()
	defer this.lock.Unlock()

	config, ok := this.Activitys[activity_type]
	if ok {
		return &config
	}
	return nil
}

func (this *ActivityManage) SetActivityConfig(config model.ActivityConfig) {
	this.lock.Lock()
	defer this.lock.Unlock()

	this.Activitys[config.ActivityType] = config
}

func (this *ActivityManage) GetActivityStatus(activity_type int) int {
	// 春节活动
	if activity_type == 8000 {
		return G_DoubleYearEvent.getActivityState()
	}

	this.satus_lock.Lock()
	defer this.satus_lock.Unlock()

	status, ok := this.ActivityStatus[activity_type]
	if ok {
		return status
	}
	return 0
}

func (this *ActivityManage) SetActivityStatus(activity_type int, status int) {
	this.satus_lock.Lock()
	defer this.satus_lock.Unlock()

	this.ActivityStatus[activity_type] = status
}

// 设置活动状态并更新活动
func (this *ActivityManage) SetActivityStatusAndUpdate(activity_type int, status int, config model.ActivityConfig, update bool) {
	this.SetActivityStatus(activity_type, status)
	if update {
		this.UpdateActivity(activity_type, status, config)
	}
}

// 更新活动
func (this *ActivityManage) UpdateActivity(activity_type int, status int, config model.ActivityConfig) {
	// 更新状态
	{
		// 开关宝箱活动
		if activity_type == proto.ACTIVITY_TREASURE_BOX {
			if status == proto.ACTIVITY_END {
				// 停止并清理
				G_TreasureBoxEvent.StopAndClean()
			} else if status == proto.ACTIVITY_START {
				if G_TreasureBoxEventRunning {
				} else {
					// 未执行则启动
					G_TreasureBoxEvent.Start()
				}
			} else {
				if G_TreasureBoxEventRunning {
					// 停止并清理
					G_TreasureBoxEvent.StopAndClean()
				}
			}
		} else if activity_type == proto.ACTIVITY_SPRING_FESTIVAL {
			activity_state := G_DoubleYearEvent.getActivityState()
			// 春节活动
			base.Setting.Springfestival.ActivityStartDatetime = utils.Time2Str(config.StartTime)
			base.Setting.Springfestival.ActivityEndDatetime = utils.Time2Str(config.FinishTime)

			// 只有停止后才需要重置
			if activity_state == proto.ACTIVITY_END {
				if status != proto.ACTIVITY_END {
					WorldMapInit(true)
				}
			}
		} else if activity_type == proto.ACTIVITY_SIGN_IN {
		}
	}
}

func (this *ActivityManage) ChangeTime(config model.ActivityConfig, update bool) error {
	activity_type := config.ActivityType
	this.SetActivityConfig(config)

	now := time.Now()
	begin_interval := config.StartTime.Unix() - now.Unix()
	end_interval := config.FinishTime.Unix() - now.Unix()

	// 清理旧定时器
	this.StopWaitActivityBegin(config)
	this.StopWaitActivityEnd(config)

	if end_interval <= 0 {
		this.SetActivityStatusAndUpdate(activity_type, proto.ACTIVITY_END, config, update)
	} else if begin_interval > 0 {
		this.SetActivityStatusAndUpdate(activity_type, proto.ACTIVITY_NOT_START, config, update)
		// 开启定时器
		return this.StartWaitActivityBegin(begin_interval, config)
	} else {
		this.SetActivityStatusAndUpdate(activity_type, proto.ACTIVITY_START, config, update)
		return this.StartWaitActivityEnd(end_interval, config)
	}
	return nil
}

func (this *ActivityManage) StartWaitActivityBegin(time_interval int64, config model.ActivityConfig) error {
	this.time_lock.Lock()
	defer this.time_lock.Unlock()

	key := fmt.Sprintf("CheckStart:%d", config.ActivityType)
	this.Timer.SetTimer(key, uint32(time_interval), this.CheckActivityStart, config)
	return nil
}

func (this *ActivityManage) StopWaitActivityBegin(config model.ActivityConfig) error {
	this.time_lock.Lock()
	defer this.time_lock.Unlock()

	key := fmt.Sprintf("CheckStart:%d", config.ActivityType)
	coretimer.Delete(this.Timer.TimerMap[key])
	return nil
}

func (this *ActivityManage) CheckActivityStart(args interface{}) {
	config := args.(model.ActivityConfig)
	now := time.Now()
	begin_interval := config.StartTime.Unix() - now.Unix()
	end_interval := config.FinishTime.Unix() - now.Unix()

	if begin_interval > 0 {
		this.StartWaitActivityBegin(begin_interval, config)
	} else {
		if end_interval > 0 {
			this.SetActivityStatusAndUpdate(config.ActivityType, proto.ACTIVITY_START, config, true)
			this.StartWaitActivityEnd(end_interval, config)
		} else {
			this.SetActivityStatusAndUpdate(config.ActivityType, proto.ACTIVITY_END, config, true)
		}
	}
}

func (this *ActivityManage) CheckActivityEnd(args interface{}) {
	config := args.(model.ActivityConfig)
	now := time.Now()
	end_interval := config.FinishTime.Unix() - now.Unix()

	if end_interval > 0 {
		this.StartWaitActivityEnd(end_interval, config)
	} else {
		this.SetActivityStatusAndUpdate(config.ActivityType, proto.ACTIVITY_END, config, true)
	}
}

func (this *ActivityManage) StartWaitActivityEnd(time_interval int64, config model.ActivityConfig) error {
	this.time_lock.Lock()
	defer this.time_lock.Unlock()

	key := fmt.Sprintf("CheckStop:%d", config.ActivityType)
	this.Timer.SetTimer(key, uint32(time_interval), this.CheckActivityEnd, config)
	return nil
}

func (this *ActivityManage) StopWaitActivityEnd(config model.ActivityConfig) error {
	this.time_lock.Lock()
	defer this.time_lock.Unlock()

	key := fmt.Sprintf("CheckStop:%d", config.ActivityType)
	coretimer.Delete(this.Timer.TimerMap[key])
	return nil
}
