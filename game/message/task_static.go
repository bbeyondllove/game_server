package message

import (
	"game_server/core/utils"
	"game_server/game/model"
	"game_server/game/proto"
)

func checkFinished(taskInfo *model.Tasks, taskvalue int, userTask *proto.UserTask) bool {
	bFlag, _, data := utils.SetKeyValue(taskInfo.TaskKey, userTask.UserId, int64(taskvalue), true, utils.ITEM_DAY)
	if !bFlag {
		return false
	}

	if data >= int64(taskInfo.TaskValue) {
		return true
	}
	return false
}

//任务检测
func taskStatic(id string, userTask *proto.UserTask) bool {
	/*
		task, awardList := getTaskById(userTask.TaskId)
		if task == nil || awardList == nil || len(awardList) == 0 {
			return false
		}

		mssageIdAry := strings.Split(task.MessageId, "|")
		msgId := strconv.Itoa(proto.MSG_INVITE_USER)
		if utils.IsExistInArrs(msgId, mssageIdAry) {
			err, invite := GetInviteUsers(userTask.UserId, 1, 1, 500, 0)
			if err != nil || len(invite.InviteUsers) == 0 {
				return false
			}

			for _, v := range invite.InviteUsers {
				if utils.GetUnixDay(v.InviteTime) == utils.GetCurDay() {
					user_info := db.RedisMgr.HGetAll(v.UserId)
					if user_info != nil {
						isFinished := checkFinished(task, 1, userTask)
						if isFinished {
							userTask.Status = 1
							updateTask(userTask, true)
						}
						return true
					}
				}
			}
		}
	*/

	return false
}
