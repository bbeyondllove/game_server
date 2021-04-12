package message

import (
	"encoding/json"
	"game_server/core/logger"
	"game_server/core/network"
	"game_server/core/utils"
	"game_server/db"
	"game_server/game/db_service"
	"game_server/game/errcode"
	"game_server/game/model"
	"game_server/game/proto"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	G_Awards         sync.Map
	G_Task           sync.Map
	G_TaskAward      sync.Map
	G_CleanSignInTag sync.Map
	G_Task_State     proto.TaskStart
)

func init() {
	GetItemList()
	getTaskAndAwarad()
}

//获取任务和奖励数据
func getTaskAndAwarad() {
	AwardInfo, aerr := db_service.AwardsIns.GetAllData()
	if aerr != nil {
		return
	}

	TaskInfo, terr := db_service.TasksIns.GetAllData()
	if terr != nil {
		return
	}

	TaskAward, twerr := db_service.TaskAwardIns.GetAllData()
	if twerr != nil {
		return
	}

	awardMap := make(map[int]map[int]*proto.AwardItem)
	for _, value := range AwardInfo {
		if _, ok := awardMap[value.AwardId]; !ok {
			awardMap[value.AwardId] = make(map[int]*proto.AwardItem, 0)
		}

		if subvalue, subok := G_ItemList.Load(value.ItemId); subok {
			node := subvalue.(*proto.ProductItem)
			idx, _ := strconv.Atoi(node.Attr1)
			sex := (idx + 1) % 2
			isGift := false
			if node.IsGift == 1 {
				isGift = true
			}
			award := &proto.AwardItem{
				ItemId:   value.ItemId,
				ItemNum:  value.ItemNum,
				IsGift:   isGift,
				ImgUrl:   node.ImgUrl,
				Desc:     node.Desc,
				Attr1:    node.Attr1,
				ItemName: node.ItemName,
				Sex:      sex,
			}
			awardMap[value.AwardId][value.ItemId] = award
		}
	}

	for k, v := range awardMap {
		G_Awards.Store(k, v)
	}

	getAward := func(taskId int) (awardId, awardNum string) {
		for _, v := range TaskAward {
			if v.TaskId == taskId {
				return v.AwardId, v.AwardNum
			}
		}
		return "", ""
	}

	status := 0
	for _, value := range TaskInfo {
		status = value.Status
		G_Task.Store(value.Id, value)
		awardId, awardNum := getAward(value.Id)
		if awardId == "" {
			continue
		}

		awardInfo := GetAwardInfo(awardId, awardNum)
		G_TaskAward.Store(value.Id, awardInfo)
	}

	SetTaskStatus(status)
}

func SetTaskStatus(status int) {
	G_Task_State.Mt.Lock()
	defer G_Task_State.Mt.Unlock()
	G_Task_State.Status = status
}

func GetTaskStatus() int {
	G_Task_State.Mt.RLock()
	defer G_Task_State.Mt.RUnlock()
	return G_Task_State.Status
}

//开启或关闭活动
func updateTasks(taskInfo *proto.ChangeTaskStatus) bool {
	SetTaskStatus(taskInfo.Status)

	return true
}

//设置用户任务数据到redis
func setUserTask() {
	userTask, err := db_service.UserTaskIns.GetData()
	if err != nil {
		logger.Errorf("setUserTask   failed(), err=", err.Error())
		return
	}

	//同步redis
	for _, v := range userTask {
		UserTaskInfo := &proto.UserTask{
			Id:         v.Id,
			UserId:     v.UserId,
			TaskId:     v.TaskId,
			TaskType:   v.TaskType,
			SigninType: v.SigninType,
			AwardInfo:  v.AwardInfo,
			Status:     v.Status,
			CreateTime: utils.Time2Str(v.CreateTime),
			UpdateTime: utils.Time2Str(v.UpdateTime),
		}
		updateTask(UserTaskInfo, false)
	}
}
func getTaskById(taskId int) (*model.Tasks, *proto.AwardInfo) {
	var task *model.Tasks

	item, ok := G_Task.Load(taskId)
	if !ok {
		return task, nil
	}

	task = item.(*model.Tasks)
	list, subok := G_TaskAward.Load(task.Id)
	if !subok {
		return task, nil
	}

	return task, list.(*proto.AwardInfo)
}

func broadTaskInfo(msg *proto.ChangeTaskStatus) {
	//通知玩家更新任务信息
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_UPDATE_TASK_INFO)
	responseMessage := &proto.S2CChangeTaskStatus{}
	responseMessage.Status = msg.Status
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)

	logger.Debugf(string(rsp.Bytes()))
	for _, v := range G_BaseCfg.LocationId {
		user_info := db.RedisMgr.HGetAll(strconv.Itoa(v))
		if user_info == nil || len(user_info) == 0 || user_info["role_id"] == "0" {
			continue
		}

		Sched.BroadCastMsg(int32(v), "0", rsp)
	}

}

//清除前一天任务和上一周签到数据
func updateTaskData() {
	curdate := utils.GetCurDay()
	if _, ok := G_CleanSignInTag.Load(curdate); ok {
		//logger.Debugf("has cleaned")
		return
	}

	G_CleanSignInTag.Range(func(k, v interface{}) bool {
		G_CleanSignInTag.Delete(k)
		return true
	})
	G_CleanSignInTag.Store(curdate, true)
	db_service.UserTaskIns.CleanData(proto.TASK_DAILY)
	logger.Debugf("has cleaned TASK_DAILY")
	weekday := utils.GetWeekDay()

	if weekday == 1 {
		db_service.UserTaskIns.CleanData(proto.TASK_SIGNIN)
		logger.Debugf("has cleaned TASK_SIGNIN")
		getTaskAndAwarad()
	}

	logger.Debugf("HandleFinishEvent end")
}

func CheckTaskData() {
	G_UndoTaskMap.Range(func(key interface{}, value interface{}) bool {
		bFinished := taskStatic(key.(string), value.(*proto.UserTask))
		if bFinished {
			G_UndoTaskMap.Delete(key)
			logger.Debugf("userstask has finished:%+v,%+v", key, value)
		}
		return true
	})

}

func getGiftItem(awardId int) []*proto.AwardItem {
	ret := make([]*proto.AwardItem, 0)
	if value, ok := G_Awards.Load(awardId); ok {
		item := value.(map[int]*proto.AwardItem)
		for _, v := range item {
			ret = append(ret, v)
		}
	}
	return ret
}

func GetAwardInfo(awardId, awardNum string) *proto.AwardInfo {
	awards := new(proto.AwardInfo)
	awardIdAry := strings.Split(awardId, "|")
	awardNumAry := strings.Split(awardNum, "|")
	if len(awardIdAry) == 0 || len(awardNumAry) == 0 || len(awardIdAry) != len(awardNumAry) {
		logger.Debugf("GetAwardAry  awardIdAry  or awardIdAry  setting error :%+v,%+v", awardId, awardNum)
		return awards
	}

	for key, value := range awardIdAry {
		if value == "" || awardNumAry[key] == "" {
			continue
		}
		itemid, _ := strconv.Atoi(value)
		itemnum, _ := strconv.ParseFloat(awardNumAry[key], 32)

		if subvalue, ok := G_ItemList.Load(itemid); ok {
			node := subvalue.(*proto.ProductItem)
			awards.ItemId = itemid
			awards.ItemNum = float32(itemnum)
			awards.IsGift = false
			if node.IsGift == 1 {
				awards.IsGift = true
			}

			awards.Desc = node.Desc
			awards.Attr1 = node.Attr1
			awards.ItemName = node.ItemName
			awards.ImgUrl = node.ImgUrl
			idx, _ := strconv.Atoi(node.Attr1)
			awards.Sex = (idx + 1) % 2

			if node.IsGift > 0 {
				awards.AwardList = getGiftItem(itemid)
			}
		}
	}

	return awards
}

func getUserTaskById(taskId int, taskType int, taskInfo map[string]*proto.UserTask) *proto.UserTask {
	var uTask *proto.UserTask
	for _, v := range taskInfo {
		if v.TaskId == taskId {
			if taskType != proto.TASK_SIGNIN {
				return v
			}
			uTask = v
		}
	}
	return uTask
}

//返回任务列表
func getTaskList(userId string, bSignin bool) (map[int]map[string]*proto.UserTask, map[int]*proto.TaskItem, int, bool, bool, bool) {
	taskList := make(map[int]*proto.TaskItem)
	taskInfo := make(map[int]map[string]*proto.UserTask)
	curDayStatus := false
	curWeekStatus := false
	FinishNum := 0
	canAddSignin := false
	weekday := utils.GetWeekDay()

	status := GetTaskStatus()
	if status == proto.TASK_STOP {
		//logger.Debufg("getTaskList task is not start")
		return taskInfo, taskList, weekday, curDayStatus, curWeekStatus, canAddSignin
	}
	for i := proto.TASK_NORMAL; i <= proto.TASK_SIGNIN; i++ {
		_, taskInfo[i] = getUserTask(userId, i)
	}

	G_Task.Range(func(key interface{}, value interface{}) bool {
		node := value.(*model.Tasks)

		if (bSignin && node.TaskType == proto.TASK_SIGNIN) ||
			(!bSignin && node.TaskType != proto.TASK_SIGNIN) {

			taskItem := new(proto.TaskItem)
			taskItem.Id = key.(int)
			taskItem.TaskType = node.TaskType
			taskItem.Title = node.Title
			taskItem.Desc = node.Desc
			taskItem.TaskKey = node.TaskKey
			taskItem.TaskValue = node.TaskValue
			taskItem.MessageId = node.MessageId
			taskItem.EventId = node.EventId
			taskItem.FrontTaskId = node.FrontTaskId
			if list, ok := G_TaskAward.Load(taskItem.Id); ok {
				taskItem.Awards = list.(*proto.AwardInfo)
			}

			uTask := getUserTaskById(taskItem.Id, node.TaskType, taskInfo[node.TaskType])
			if uTask != nil {
				taskItem.Status = uTask.Status

				if taskItem.Status > 0 {
					updateTime := utils.Str2Time(uTask.UpdateTime)
					if uTask.SigninType != 2 && utils.GetTimeDay(updateTime) == utils.GetCurDay() {
						curDayStatus = true
					}
					FinishNum++
					if uTask.SigninType == 2 {
						curWeekStatus = true
					}
				}
			}

			taskList[taskItem.Id] = taskItem
		}

		return true
	})

	undoDay := weekday - FinishNum
	if !curWeekStatus && undoDay >= 1 {
		canAddSignin = true
	}
	return taskInfo, taskList, weekday, curDayStatus, curWeekStatus, canAddSignin
}

//返回用户已做任务信息,要做任务信息,对应的奖励信息
func geUndoTask(userId string, bSignin bool) (map[int]map[string]*proto.UserTask, map[int]*proto.TaskItem, int, bool, bool, bool) {
	taskInfo, taskList, weekday, curDayStatus, curWeekStatus, canAddSignin := getTaskList(userId, bSignin)
	undoList := make(map[int]*proto.TaskItem)
	if len(taskList) == 0 {
		return taskInfo, undoList, weekday, curDayStatus, curWeekStatus, canAddSignin
	}

	for k, v := range taskList {
		if v.Status > 0 {

			continue
		}

		undoList[k] = v
	}

	return taskInfo, undoList, weekday, curDayStatus, curWeekStatus, canAddSignin
}

func getUserTask(userId string, taskType int) (bool, map[string]*proto.UserTask) {
	taskKey := proto.STATIC_KEY[taskType] + userId
	taskInfo := make(map[string]*proto.UserTask)
	//token 检验
	ret := db.RedisMgr.HGetAll(taskKey)
	if ret == nil {
		logger.Errorf("getUserTask error:%+v", taskKey)
		return false, taskInfo
	}

	for k, v := range ret {
		item := make(map[string]interface{})
		err := json.Unmarshal([]byte(v), &item)
		if err != nil || len(item) == 0 {
			continue
		}

		taskId := int(item["task_id"].(float64))
		taskType := int(item["task_type"].(float64))
		signinType := int(item["signin_type"].(float64))
		status := int(item["status"].(float64))
		createTime := item["create_time"].(string)
		updateTime := item["update_time"].(string)

		node := &proto.UserTask{
			Id:         item["id"].(string),
			UserId:     item["user_id"].(string),
			TaskId:     taskId,
			TaskType:   taskType,
			SigninType: signinType,
			AwardInfo:  item["award_info"].(string),
			Status:     status,
			CreateTime: createTime,
			UpdateTime: updateTime,
		}
		taskInfo[k] = node

	}
	return true, taskInfo
}

func addTask(conn *network.Conn, userId string, taskId int, taskType int, signinType int, bNeedFinished bool) (bool, int32) {
	status := 0
	awardList, ok := G_TaskAward.Load(taskId)
	if !ok {
		return false, errcode.ERROR_MYSQL
	}

	awardInfo := awardList.(*proto.AwardInfo)
	awards := make([]*proto.AwardItem, 0)
	if awardInfo.IsGift {
		awards = awardInfo.AwardList
	} else {
		awards = append(awards, &awardInfo.AwardItem)
	}
	if taskType == proto.TASK_SIGNIN {
		status = 2
	} else if bNeedFinished {
		status = 1
	}

	buf, _ := json.Marshal(awards)
	userdata := &model.UserTask{
		UserId:     userId,
		TaskId:     taskId,
		TaskType:   taskType,
		SigninType: signinType,
		AwardInfo:  string(buf),
		Status:     status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	_, err, id := db_service.UserTaskIns.Add(userdata)
	if err != nil {
		logger.Errorf("db_service.UserTask Add error=%+v", err)
		return false, errcode.ERROR_REDIS
	}

	UserTaskInfo := &proto.UserTask{
		Id:         userdata.Id,
		UserId:     userId,
		TaskId:     taskId,
		TaskType:   taskType,
		SigninType: signinType,
		AwardInfo:  userdata.AwardInfo,
		Status:     status,
		CreateTime: utils.Time2Str(userdata.CreateTime),
		UpdateTime: utils.Time2Str(userdata.UpdateTime),
	}
	taskKey := proto.STATIC_KEY[taskType] + userId
	taskbuf, _ := json.Marshal(UserTaskInfo)

	itemType := utils.ITEM_NORMAL
	if taskType == proto.TASK_DAILY {
		itemType = utils.ITEM_DAY
	} else if taskType == proto.TASK_SIGNIN {
		itemType = utils.ITEM_WEEK
	}

	bRet, _, _ := utils.SetKeyValue(taskKey, UserTaskInfo.Id, taskbuf, false, itemType)
	if !bRet {
		return false, errcode.ERROR_REDIS
	}

	if status == 1 {
		pushTaskFinish(conn, UserTaskInfo)
	}
	if !bNeedFinished {
		G_UndoTaskMap.Store(id, userdata)
	}
	return true, errcode.MSG_SUCCESS
}

func updateTask(userTask *proto.UserTask, bUpdateTable bool) bool {
	taskKey := proto.STATIC_KEY[userTask.TaskType] + userTask.UserId
	taskbuf, _ := json.Marshal(userTask)

	itemType := utils.ITEM_NORMAL
	if userTask.TaskType == proto.TASK_DAILY {
		itemType = utils.ITEM_DAY
	} else if userTask.TaskType == proto.TASK_SIGNIN {
		itemType = utils.ITEM_WEEK
	}

	bRet, _, _ := utils.SetKeyValue(taskKey, userTask.Id, taskbuf, false, itemType)
	if !bRet {
		return false
	}

	if bUpdateTable {
		data_map := make(map[string]interface{})
		data_map["update_time"] = time.Now()
		data_map["status"] = userTask.Status
		_, err := db_service.UpdateFields(db_service.UserTaskTable, "id", userTask.Id, data_map)
		if err != nil {
			return false
		}
	}
	return true
}

func taskProcess(conn *network.Conn, userId string, messageId int, eventId int, changeValue int, bNeedFinished bool) {
	updateTaskData()

	_, taskList, _, _, _, _ := geUndoTask(userId, false)
	if len(taskList) == 0 {
		return
	}

	for k, v := range taskList {
		mssageIdAry := strings.Split(v.MessageId, "|")
		msgId := strconv.Itoa(messageId)
		if !utils.IsExistInArrs(msgId, mssageIdAry) {
			continue
		}
		if v.EventId != eventId {
			continue
		}

		taskKey := v.TaskKey + ":" + userId
		if eventId != 0 {
			taskKey = v.TaskKey + ":" + strconv.Itoa(eventId) + ":" + userId
		}
		var bFlag bool
		var data int64
		if v.TaskType == proto.TASK_DAILY {
			bFlag, _, data = utils.SetKeyValue(taskKey, "static", int64(changeValue), true, utils.ITEM_DAY)
		}
		if !bFlag {
			continue
		}
		if (bNeedFinished && data >= int64(v.TaskValue)) || !bNeedFinished {
			addTask(conn, userId, k, v.TaskType, 0, bNeedFinished)
		}
	}

}

//获取签到列表
func (s *CSession) HandleGetSignInList(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetSignInList in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_SIGNIN_LIST_RSP)
	responseMessage := &proto.S2CTaskList{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	if G_ActivityManage.GetActivityStatus(proto.ACTIVITY_SIGN_IN) == proto.ACTIVITY_NOT_START {
		responseMessage.Code = errcode.ERROR_NOT_START
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}
	if G_ActivityManage.GetActivityStatus(proto.ACTIVITY_SIGN_IN) == proto.ACTIVITY_END {
		responseMessage.Code = errcode.ERROR_END
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return

	}
	_, taskList, weekday, curDayStatus, curWeekStatus, _ := getTaskList(payLoad.UserId, true)
	taskAry := make([]*proto.TaskItem, 0)
	taskIdAry := make([]int, 0)

	for k, _ := range taskList {
		taskIdAry = append(taskIdAry, k)
	}
	sort.Ints(taskIdAry)
	for _, v := range taskIdAry {
		taskAry = append(taskAry, taskList[v])
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.TaskList = taskAry
	responseMessage.SigninDayNo = weekday
	responseMessage.CurDayStatus = curDayStatus
	responseMessage.CurWeekStatus = curWeekStatus
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf(string(rsp.Bytes()))

	logger.Debugf("HandleGetSignInList end")
	return
}

func sendAward(eventId int, userId string, userInfo map[string]string, awardInfo *proto.AwardInfo) (float32, map[int][]*proto.AwardItem, int) {
	times := 1 // 倍数
	cdt := float32(0)
	itemInfos := make(map[int][]*proto.AwardItem, 0)
	awardList := make([]*proto.AwardItem, 0)
	if awardInfo.IsGift {
		awardList = awardInfo.AwardList
	} else {
		awardList = append(awardList, &awardInfo.AwardItem)
	}
	for _, v := range awardList {
		node, ok := G_ItemList.Load(v.ItemId)
		if !ok {
			continue
		}
		itemInfo := node.(*proto.ProductItem)
		if itemInfo.ItemType != proto.ITEM_CDT {
			err := UpdateData(userId, userInfo, "item_info", itemInfo, int(v.ItemNum))
			if err != nil {
				continue
			}

			itemInfos[itemInfo.ItemType] = append(itemInfos[itemInfo.ItemType], v)

		} else {
			// todo 圣诞老人 检查当前角色是否是圣诞老人并且没有过期
			isTrue, _ := CheckCurrentRole(userId)
			var code int32
			if isTrue {
				// 使用圣诞老人 双倍cdt
				code, _ = db_service.NewCdt().UpdateUserCdt(userId, v.ItemNum, proto.MSG_CHRISMAS_ROLE)
				// 活动的倍数
				times = 2
			} else {
				code, _ = db_service.NewCdt().UpdateUserCdt(userId, v.ItemNum, eventId)
			}
			if code == 0 {
				cdt += v.ItemNum
			}
		}
	}
	return cdt, itemInfos, times
}

//签到
func (s *CSession) HandleSignIn(requestMsg *utils.Packet) {
	logger.Debugf("HandleSignIn in request:", requestMsg.GetBuffer())

	updateTaskData()
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_SIGN_IN_RSP)
	responseMessage := &proto.S2CTaskAward{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SSignIn{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, userInfo := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	userTask, toDoTask, _, curDayStatus, _, canAddSignin := geUndoTask(payLoad.UserId, true)
	if len(toDoTask) == 0 {
		logger.Errorf("HandleSignIn setting error toDoTask=%+v", toDoTask)
		responseMessage.Code = errcode.ERROR_NONEED_SIGNIN
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		return
	}

	if msg.SigninType == proto.SIGNIN_LOST {
		if !curDayStatus {
			logger.Errorf("当天还没签到，不能补签")
			responseMessage.Code = errcode.ERROR_MUST_SINGIN_CURRDAY
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
			return
		}

		if !canAddSignin {
			logger.Errorf("HandleSignIn lostDay=0 不需要补签", payLoad.UserId)
			responseMessage.Code = errcode.ERROR_NONEED_ADDSIGNIN
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
			return
		}
		bSuccess, code := UseItem(payLoad.UserId, userInfo, msg.ItemId)
		if !bSuccess {
			responseMessage.Code = code
			responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
			rsp.WriteData(responseMessage)
			logger.Debugf(string(rsp.Bytes()))
			s.sendPacket(rsp)
			return
		}
	} else if curDayStatus {
		logger.Errorf("currday Task Is Finished userTask=%+v", userTask)
		responseMessage.Code = errcode.ERROR_CURDAY_IS_FINISHED
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		return
	}

	taskIdAry := make([]int, 0)
	for k, _ := range toDoTask {
		taskIdAry = append(taskIdAry, k)
	}
	sort.Ints(taskIdAry)
	var task *proto.TaskItem
	for _, v := range taskIdAry {
		task = toDoTask[v]
		break
	}

	bFlag, code := addTask(&s.conn, payLoad.UserId, task.Id, task.TaskType, msg.SigninType, true)
	if !bFlag {
		responseMessage.Code = code
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	awardList, ok := G_TaskAward.Load(task.Id)
	if !ok {
		responseMessage.Code = code
		responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	responseMessage.Awards = awardList.(*proto.AwardInfo)
	_, _, times := sendAward(proto.MSG_SIGN_IN, payLoad.UserId, userInfo, responseMessage.Awards)
	// 圣诞老人乘 * 2
	if times == 2 {
		// todo 需要复制一份值来修改， 否则会修改全局变量
		temp := *responseMessage.Awards
		temp.ItemNum = temp.ItemNum * 2
		responseMessage.Awards = &temp
	}
	go taskProcess(&s.conn, payLoad.UserId, proto.MSG_SIGN_IN, 0, 1, true)
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	s.sendPacket(rsp)
	logger.Debugf("HandleSignIn end")
	return
}

//获取任务列表
func (s *CSession) HandleGetTaskList(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetTaskList in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_TASK_LIST_RSP)
	responseMessage := &proto.S2CTaskList{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return

	}
	_, taskList, _, _, _, _ := getTaskList(payLoad.UserId, false)
	taskAry := make([]*proto.TaskItem, 0)
	taskIdAry := make([]int, 0)

	for k, _ := range taskList {
		taskIdAry = append(taskIdAry, k)
	}
	sort.Ints(taskIdAry)
	for _, v := range taskIdAry {
		taskAry = append(taskAry, taskList[v])
	}

	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]
	responseMessage.TaskList = taskAry
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf(string(rsp.Bytes()))

	logger.Debugf("HandleGetTaskList end")
	return
}

//领取任务奖励
func (s *CSession) HandleGetTaskAward(requestMsg *utils.Packet) {
	logger.Debugf("HandleGetTaskAward in request:", requestMsg.GetBuffer())
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_GET_TASK_AWARD_RSP)
	responseMessage := &proto.S2CTaskAward{}
	responseMessage.Code = errcode.ERROR_REQUEST_NOT_ALLOW
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	msg := &proto.C2STaskAward{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return
	}

	flag, payLoad, userInfo := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		rsp.WriteData(responseMessage)
		s.sendPacket(rsp)
		return

	}

	if _, ok := G_Task.Load(msg.TaskId); !ok {
		logger.Errorf("HandleGetTaskAward taskid error:msg=%+v", msg)
		rsp.WriteData(responseMessage)
		logger.Debugf(string(rsp.Bytes()))
		s.sendPacket(rsp)
		return
	}

	bFlag, taskInfo := getUserTask(payLoad.UserId, proto.TASK_DAILY)
	responseMessage.Code = errcode.MSG_SUCCESS
	responseMessage.Message = errcode.ERROR_MSG[responseMessage.Code]

	var myTask *proto.UserTask
	if bFlag {
		for _, v := range taskInfo {
			if v.TaskId == msg.TaskId {
				myTask = v
				break
			}
		}
	}
	if myTask != nil && myTask.Status == 1 {
		_, responseMessage.Awards = getTaskById(msg.TaskId)
		myTask.Status = 2
		updateTask(myTask, true)
		_, _, times := sendAward(proto.MSG_GET_TASK_AWARD, payLoad.UserId, userInfo, responseMessage.Awards)
		// 圣诞老人乘2
		if times == 2 {
			// todo 需要复制一份值来修改， 否则回修改全局变量
			temp := *responseMessage.Awards
			temp.ItemNum = temp.ItemNum * 2
			responseMessage.Awards = &temp
		}
	}
	rsp.WriteData(responseMessage)
	s.sendPacket(rsp)
	logger.Debugf(string(rsp.Bytes()))
	logger.Debugf("HandleGetTaskAward end")
	return
}

//推送任务完成
func pushTaskFinish(conn *network.Conn, userTask *proto.UserTask) {
	logger.Debugf("PushTaskFinish start:", userTask)
	rsp := &utils.Packet{}
	rsp.Initialize(proto.MSG_BROAD_USER_TASK)
	responseMessage := &proto.PushTaskStatus{}
	responseMessage.Message = ""
	responseMessage.Code = errcode.MSG_SUCCESS

	responseMessage.TaskId = userTask.TaskId

	_, responseMessage.Awards = getTaskById(userTask.TaskId)
	rsp.WriteData(responseMessage)
	logger.Debugf(string(rsp.Bytes()))
	SendPacket(*conn, rsp)
	logger.Debugf("PushTaskFinish end")
	return
}

//进入商店
func (s *CSession) HandleEnterShop(requestMsg *utils.Packet) {
	logger.Debugf("HandleEnterShop in request:", requestMsg.GetBuffer())
	msg := &proto.C2SBase{}
	err := json.Unmarshal(requestMsg.Bytes(), msg)
	if err != nil {
		logger.Errorf("json.Unmarshal error, err=", err.Error())
		return
	}

	flag, payLoad, _ := utils.GetUserByToken(msg.Token)
	if !flag {
		logger.Errorf("token error", msg.Token)
		return

	}

	//任务处理
	taskProcess(&s.conn, payLoad.UserId, proto.MSG_ENTER_SHOP, 0, 1, true)

	logger.Debugf("HandleEnterShop end")
	return
}
