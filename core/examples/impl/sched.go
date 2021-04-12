package impl

import (
	"sync"
	"sync/atomic"
	"time"

	"game_server/core/utils"
)

const SCHED_SLEEP_CONST = int64(50 * time.Millisecond)

var Sched = NewSched()

type Scheduler struct {
	sessList      map[int64]*CSession
	addReconQueue *utils.SyncQueue
	stopEvent     int32
	wg            sync.WaitGroup
}

func NewSched() *Scheduler {
	return &Scheduler{
		addReconQueue: utils.NewSyncQueue(),
		sessList:      make(map[int64]*CSession),
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

func (s *Scheduler) isStopped() bool {
	return atomic.LoadInt32(&s.stopEvent) == 1
}

func (s *Scheduler) run() {
	defer s.wg.Done()

	var prevSleepTime int64
	PreTime := time.Now().UnixNano()

	for !s.isStopped() {
		Time := time.Now().UnixNano()
		diff := Time - PreTime
		PreTime = Time

		s.Update()

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
	WorldList.SyncProcess()
}

//调度所有的session消息
func (s *Scheduler) updateSessions() {
	for k, sess := range s.sessList {
		if !sess.Update() {
			delete(s.sessList, k)
		}
	}
}

func (s *Scheduler) addSession(cs *CSession) bool {
	var ret bool
	if _, ok := s.sessList[cs.UserId]; ok {
		//发送挤用户下线消息

		delete(s.sessList, cs.UserId)
	}
	s.sessList[cs.UserId] = cs

	return ret

}
