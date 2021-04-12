package message

import (
	kk_core "game_server/core"
	"game_server/core/utils"
	"game_server/game/proto"
	"sync"
	"sync/atomic"
	"time"

	"game_server/core/logger"

	plug "github.com/syyongx/php2go"
)

const SCHED_SLEEP_CONST = int64(50 * time.Millisecond)

var Sched = NewSched()

type Scheduler struct {
	ConnMap map[int32]map[string]*CSession //key1城市编号,key2用户ID value:*CSession
	// sessList      map[uint32]*CSession
	addReconQueue *utils.SyncQueue
	stopEvent     int32
	wg            sync.WaitGroup
	bInit         map[int32]bool
}

func NewSched() *Scheduler {
	return &Scheduler{
		ConnMap:       make(map[int32]map[string]*CSession, 0),
		addReconQueue: utils.NewSyncQueue(),
		bInit:         make(map[int32]bool),

		// sessList:      make(map[uint32]*CSession),
	}
}

/*
启动主调度协程
主调度协程中不能有阻塞任务，阻塞任务需要放到对应协程池执行
*/
func Start() {
	Sched.wg.Add(1)
	go Sched.run()
	Sched.wg.Wait()
}

func (s *Scheduler) InitConnMap(locateId int32) {
	if v, ok := s.bInit[locateId]; ok && v {
		return
	}
	s.ConnMap[locateId] = make(map[string]*CSession, 0)
	s.bInit[locateId] = true
}

func (s *Scheduler) ExistLoacateId(locateId int32) bool {
	_, ok := s.ConnMap[locateId]
	return ok
}

func (s *Scheduler) isStopped() bool {
	return atomic.LoadInt32(&s.stopEvent) == 1
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	var prevSleepTime int64
	PreTime := time.Now().UnixNano()
	// tick := time.NewTicker(time.Duration(G_BaseCfg.EventInterval) * time.Second)

	for !s.isStopped() {
		Time := time.Now().UnixNano()
		diff := Time - PreTime
		PreTime = Time

		select {
		case <-G_NormalEventTimer.C:
			//普通事件处理
			G_NormalEvent.EventProcess()
		case <-G_TreasureBoxTimer.C:
			//宝箱事件
			G_TreasureBoxEvent.EventProcess()
		default:
			s.Update()
		}

		// diff (D0) include time of previous sleep (d0) + tick time (t0)
		// we want that next d1 + t1 == WORLD_SLEEP_CONST
		// we can't know next t1 and then can use (t0 + d1) == WORLD_SLEEP_CONST requirement
		// d1 = WORLD_SLEEP_CONST - t0 = WORLD_SLEEP_CONST - (D0 - d0) = WORLD_SLEEP_CONST + d0 - D0
		if diff <= SCHED_SLEEP_CONST+prevSleepTime {
			prevSleepTime = SCHED_SLEEP_CONST + prevSleepTime - diff
			time.Sleep(time.Duration(prevSleepTime))
		} else {
			prevSleepTime = 0
		}

	}
}

func Destroy() {
	atomic.StoreInt32(&Sched.stopEvent, 1)
	Sched.wg.Wait()
}

//整体调度
func (s *Scheduler) Update() {
	s.updateSessions()
	kk_core.WorldList.SyncProcess()
	G_OffLine.HandleMsg()
}

//调度所有的session消息
func (s *Scheduler) updateSessions() {
	msg := &utils.Packet{}
	msg.Initialize(proto.MSG_BROAD_USER_OFFLINE)

	offMsg := &proto.CityUser{}
	// msg.MsgData = make(map[string]interface{}, 0)
	userList := make(map[int32][]string, 0)
	flag := false

	for k, sess := range s.ConnMap {
		userList[k] = make([]string, 0)
		// logger.Debugf("updateSessions() len(sess)=", len(sess))
		for key, session := range sess {
			// logger.Debugf("updateSessions() uid=", session.UserId)
			if !session.Update() {
				logger.Debugf("updateSessions delete session uid=", session.UserId)
				delete(s.ConnMap[k], key)
				flag = true
				userList[k] = append(userList[k], key)
			}
		}
	}
	if len(userList) > 0 && flag {
		offMsg.UserList = userList
		msg.WriteData(offMsg)
		G_OffLine.QueuePacket(msg)
	}

}

func (s *Scheduler) GetUser(locationId int32, userId string) *CSession {
	if city, ok := s.ConnMap[locationId]; ok {
		if user, isok := city[userId]; isok {
			return user
		}
	} else {
		return nil
	}
	return nil
}

//获取同一城市的所有玩家
func (s *Scheduler) GetSameCityUser(locationId int32) []string {
	ret := make([]string, 0)
	if _, ok := s.ConnMap[locationId]; !ok {
		return ret
	}

	for k, _ := range s.ConnMap[locationId] {
		ret = append(ret, k)
	}
	return ret
}

//广播消息给同一场景的所有玩家
//locationId	//城市id
//leftTop   	//左上角坐标
//rightBottom   //右下角坐标
func (s *Scheduler) SendScreenUser(locationId int32, userId string, leftTopX, leftTopY, rightBottomX, rightBottomY int, msg *utils.Packet) {
	if _, ok := s.ConnMap[locationId]; !ok {
		return
	}

	for _, v := range s.ConnMap[locationId] {
		if userId == v.UserId {
			continue
		}
		if v.X >= leftTopX && v.X <= rightBottomX && v.Y >= leftTopY && v.Y <= rightBottomY {
			v.sendPacket(msg)
		}
	}
}

func (s *Scheduler) addSession(cs *CSession, bReconnect bool) bool {
	var ret bool
	if cityMap, isok := s.ConnMap[cs.LocateId]; isok {
		if session, ok := cityMap[cs.UserId]; ok {
			if !bReconnect && !session.BReconnect {
				//挤用户下线
				msg := &utils.Packet{}
				msg.Initialize(proto.MSG_LOGINANOTHER)
				session.sendPacket(msg)

				go session.CloseSession()
				delete(s.ConnMap[cs.LocateId], cs.UserId)
			}
		}
		cs.BReconnect = bReconnect
		s.ConnMap[cs.LocateId][cs.UserId] = cs
		ret = true
	}
	return ret

}

func (s *Scheduler) delSession(cs *CSession) bool {
	if cityMap, isok := s.ConnMap[cs.LocateId]; isok {
		if _, ok := cityMap[cs.UserId]; ok {
			delete(s.ConnMap[cs.LocateId], cs.UserId)
			return true
		}
	}
	return false
}

//消息广播，向指定城市locateId除了指定玩家userId外的所有玩家广播消息
//如果userId=""，则向所有玩家广播消息
func (s *Scheduler) BroadCastMsg(locateId int32, userId string, msg *utils.Packet) {
	logger.Debugf("broadCastMsg in:", msg.GetBuffer())

	if cityMap, isok := s.ConnMap[locateId]; isok {
		// logger.Debugf("len(cityMap)=", len(cityMap))
		for _, session := range cityMap {
			if userId == session.UserId {
				continue
			}
			session.sendPacket(msg)
			// logger.Debugf("sendPacket to uid=", session.UserId)
		}
	} else {
		logger.Errorf("locateId is not exist locateId=", locateId)
	}
	logger.Debugf("broadCastMsg end")
}

func (s *Scheduler) BroadCastMsgToPlatform(locateId int32, userId string, msg *utils.Packet, allow_platform int) {
	logger.Debugf("broadCastMsg in:", msg.GetBuffer())

	if cityMap, isok := s.ConnMap[locateId]; isok {
		// logger.Debugf("len(cityMap)=", len(cityMap))
		for _, session := range cityMap {
			if userId == session.UserId {
				continue
			}
			if (session.Platform != 0) && (session.Platform != allow_platform) {
				continue
			}
			session.sendPacket(msg)
			// logger.Debugf("sendPacket to uid=", session.UserId)
		}
	} else {
		logger.Errorf("locateId is not exist locateId=", locateId)
	}
	logger.Debugf("broadCastMsg end")
}

func (s *Scheduler) BroadCastMsgByVersion(locateId int32, userId string, msg *utils.Packet, version string) {
	logger.Debugf("broadCastMsg in:", msg.GetBuffer())

	if cityMap, isok := s.ConnMap[locateId]; isok {
		// logger.Debugf("len(cityMap)=", len(cityMap))
		for _, session := range cityMap {
			if userId == session.UserId {
				continue
			}
			if len(version) > 0 {
				if (len(session.Version) != 0) && (plug.VersionCompare(version, session.Version, "<=")) {
					continue
				}
			}

			session.sendPacket(msg)
			// logger.Debugf("sendPacket to uid=", session.UserId)
		}
	} else {
		logger.Errorf("locateId is not exist locateId=", locateId)
	}
	logger.Debugf("broadCastMsg end")
}

//推送消息给指定用户
func (s *Scheduler) SendToUser(userId string, msg *utils.Packet) {
	logger.Debugf("SendToUser in:", msg.GetBuffer())
	for _, v := range s.ConnMap {
		if session, ok := v[userId]; ok {
			session.sendPacket(msg)
			break
		}
	}

	logger.Debugf("SendToUser end")
}

//推送消息给指定用户
func (s *Scheduler) SendToAllUser(msg *utils.Packet) {
	logger.Debugf("SendToAllUser in:", msg.GetBuffer())
	for _, v := range s.ConnMap {
		for _, session := range v {
			session.sendPacket(msg)
		}
	}

	logger.Debugf("SendToAllUser end")
}
