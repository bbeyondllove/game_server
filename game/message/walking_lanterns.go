package message

import (
	"fmt"
	"game_server/core/logger"
	coretimer "game_server/core/timer"
	"game_server/core/utils"
	dao "game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"sync"
	"time"
)

var (
	G_WalkingLantenrns = WalkingLantenrns{
		WalkingLantenrnsTimers: make(map[int64]*coretimer.TimeWheel, 0),
	}
)

type WalkingLantenrns struct {
	lock                   sync.Mutex
	WalkingLantenrnsTimers map[int64]*coretimer.TimeWheel
}

func (this *WalkingLantenrns) Start() {
}

func (this *WalkingLantenrns) Stop() {
	this.lock.Lock()
	defer this.lock.Unlock()

	for _, wheel := range this.WalkingLantenrnsTimers {
		if wheel != nil {
			wheel.Stop()
		}
	}
	this.WalkingLantenrnsTimers = make(map[int64]*coretimer.TimeWheel, 0)
}

func (this *WalkingLantenrns) NoticeLantenrns(args interface{}) {
	data := args.(model.TimeNotice)

	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_PUSH_NOTICE_LANTERNS_RSP)
	pushMessage := &proto.S2CNoticeLaterns{}
	//推送给所有用户
	pushMessage.Code = errcode.MSG_SUCCESS
	pushMessage.Message = ""
	pushMessage.Notice = proto.LanternInfo{
		Level:   G_BaseCfg.Backstage.NoticeLanternLevel,
		Content: data.Content,
	}
	rsp.WriteData(pushMessage)
	// 广播给在线的用户
	Sched.BroadCastMsg(int32(3302), "0", rsp)
}

func (this *WalkingLantenrns) CleanTimer(notice_id int64) {
	this.lock.Lock()
	defer this.lock.Unlock()

	wheel, ok := this.WalkingLantenrnsTimers[notice_id]
	if ok {
		if wheel != nil {
			wheel.Stop()
			delete(this.WalkingLantenrnsTimers, notice_id)
		}
	}
}

func (this *WalkingLantenrns) CompleteNotify(notice_info model.Notice) error {
	notice_id := notice_info.Id
	this.Clean(notice_id)

	// 更新通知状态
	_, err := dao.NoticeIns.UpdateStatus(notice_id, 1)
	if err != nil {
		return err
	}
	return nil
}

func (this *WalkingLantenrns) StartCompleteNotify(args interface{}) {
	notice_info := args.(model.Notice)
	this.CompleteNotify(notice_info)
}

func (this *WalkingLantenrns) AddNotice(notice_time time.Time, data model.Notice) {
	space := notice_time.Unix() - time.Now().Unix()
	if space <= 0 {
		err := this.CompleteNotify(data)
		if err != nil {
			logger.Error(err)
		}
	} else {
		this.lock.Lock()
		defer this.lock.Unlock()
		notice_id := data.Id
		wheel, ok := this.WalkingLantenrnsTimers[notice_id]
		if ok || wheel == nil {
			wheel = coretimer.NewTimeWheel()
			wheel.Start()
			this.WalkingLantenrnsTimers[notice_id] = wheel
		}
		wheel.SetTimer(fmt.Sprintf("%d-", notice_id), uint32(space), this.StartCompleteNotify, data)
	}
}

// 添加定时器
func (this *WalkingLantenrns) AddTimeNotice(notice_id int64, data model.TimeNotice) int {
	space := data.StartTime.Unix() - time.Now().Unix()
	fmt.Println("AddTimeNotice", space, data)
	if space > 0 {
		this.lock.Lock()
		defer this.lock.Unlock()
		wheel, ok := this.WalkingLantenrnsTimers[notice_id]
		if ok || wheel == nil {
			wheel = coretimer.NewTimeWheel()
			wheel.Start()
			this.WalkingLantenrnsTimers[notice_id] = wheel
		}
		data.NoticeId = notice_id
		wheel.SetTimer(fmt.Sprintf("%d-%d", notice_id, data.Id), uint32(space), this.NoticeLantenrns, data)
		return 0
	} else if space == 0 {
		this.NoticeLantenrns(data)
		return 0
	} else {
		logger.Errorf("WalkingLantenrns is drop ", space, data)
		return -1
	}
}

func (this *WalkingLantenrns) Clean(notice_id int64) {
	this.CleanTimer(notice_id)
}
