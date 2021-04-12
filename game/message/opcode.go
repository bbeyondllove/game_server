package message

import (
	"game_server/core/utils"
	"game_server/game/proto"
)

var OpCodeTable = make(map[uint16]OpCodeHandler, 0)
var AgentCodeTable = make(map[uint16]AgentHandler, 0)

const (
	STATUS_NOT_AUTHED = 0
	STATUS_AUTHED     = 1
	STATUS_UNHANDLED  = 2
)

//消息路由器
type OpCodeHandler struct {
	Name    string
	Status  uint8
	Handler func(*CSession, *utils.Packet)
}

type AgentHandler struct {
	Name    string
	Status  uint8
	Handler func(*agent, *utils.Packet)
}

func init() {

	OpCodeTable[proto.MSG_POSITION_CHANGE] = OpCodeHandler{"MSG_POSITION_CHANGE", STATUS_AUTHED, (*CSession).HandlePositionChange}
	OpCodeTable[proto.MSG_GET_POSITION] = OpCodeHandler{"MSG_GET_POSITION", STATUS_AUTHED, (*CSession).HandleGetPosition}
	OpCodeTable[proto.MSG_GET_AMOUNT] = OpCodeHandler{"MSG_GET_AMOUNT", STATUS_AUTHED, (*CSession).HandleGetAmount}
	OpCodeTable[proto.MSG_GET_KNAPSACK] = OpCodeHandler{"MSG_GET_KNAPSACK", STATUS_AUTHED, (*CSession).HandleUserKnapsack}
	OpCodeTable[proto.MSG_UPDATE_ITEM_INFO] = OpCodeHandler{"MSG_UPDATE_ITEM_INFO", STATUS_AUTHED, (*CSession).HandleUpdateItemInfo}
	OpCodeTable[proto.MSG_FINISH_EVENT] = OpCodeHandler{"MSG_FINISH_EVENT", STATUS_AUTHED, (*CSession).HandleFinishEvent}
	OpCodeTable[proto.MSG_GET_INVITE_USERS] = OpCodeHandler{"MSG_GET_INVITE_USERS", STATUS_AUTHED, (*CSession).HandleGetInviteUsers}
	OpCodeTable[proto.MSG_GET_GRAB_COMRADES] = OpCodeHandler{"MSG_GET_GRAB_COMRADES", STATUS_AUTHED, (*CSession).HandleGetGrabComrades}
	OpCodeTable[proto.MSG_GET_USERINFO] = OpCodeHandler{"MSG_GET_USERINFO", STATUS_AUTHED, (*CSession).HandleGetUserInfo}
	OpCodeTable[proto.MSG_GET_BUILDING_DESC] = OpCodeHandler{"MSG_GET_BUILDING_DESC", STATUS_AUTHED, (*CSession).HandleGetBuildingDesc}
	OpCodeTable[proto.MSG_GET_MEMBER_SYS] = OpCodeHandler{"MSG_GET_MEMBER_SYS", STATUS_AUTHED, (*CSession).HandleGetMemberSys}
	OpCodeTable[proto.MSG_GET_USER_LEVEL] = OpCodeHandler{"MSG_GET_USER_LEVEL", STATUS_AUTHED, (*CSession).HandleGetUserLevel}
	OpCodeTable[proto.MSG_GET_ITEMS_LIST] = OpCodeHandler{"MSG_GET_ITEMS_LIST", STATUS_AUTHED, (*CSession).HandleGetItemsList}
	OpCodeTable[proto.MSG_BUY_ITEM] = OpCodeHandler{"MSG_BUY_ITEM", STATUS_AUTHED, (*CSession).HandleBuyItem}
	OpCodeTable[proto.MSG_MODIFY_NICKNAME] = OpCodeHandler{"MSG_MODIFY_NICKNAME", STATUS_AUTHED, (*CSession).HandleModifyNickName}
	OpCodeTable[proto.MSG_SELECT_ROLE] = OpCodeHandler{"MSG_SELECT_ROLE", STATUS_AUTHED, (*CSession).HandleRoleSelect}
	OpCodeTable[proto.MSG_CERTIFICATION] = OpCodeHandler{"MSG_CERTIFICATION", STATUS_AUTHED, (*CSession).HandleCertification}
	OpCodeTable[proto.MSG_BIND_INVITER] = OpCodeHandler{"MSG_BIND_INVITER", STATUS_AUTHED, (*CSession).HandleBindInvitationCode}
	OpCodeTable[proto.MSG_DEPOSIT_REBATE] = OpCodeHandler{"MSG_DEPOSIT_REBATE", STATUS_AUTHED, (*CSession).HandleDepositRebate}
	OpCodeTable[proto.MSG_QUERY_SHOP] = OpCodeHandler{"MSG_QUERY_SHOP", STATUS_AUTHED, (*CSession).HandleQueryShop}
	OpCodeTable[proto.MSG_GET_FRIEND_INFO] = OpCodeHandler{"MSG_GET_FRIEND_INFO", STATUS_AUTHED, (*CSession).HandleGetFriendInfo}
	OpCodeTable[proto.MSG_GET_INVITATION] = OpCodeHandler{"MSG_GET_INVITATION", STATUS_AUTHED, (*CSession).HandleGetInvitation}
	OpCodeTable[proto.MSG_GET_TOEKN_PAY_URL] = OpCodeHandler{"MSG_GET_TOEKN_PAY_URL", STATUS_AUTHED, (*CSession).HandleGetTokenPayUrl} // 已废弃
	OpCodeTable[proto.MSG_USER_QUT_CITY] = OpCodeHandler{"MSG_USER_QUT_CITY", STATUS_AUTHED, (*CSession).HandleQuitCity}
	OpCodeTable[proto.MSG_GET_CITY_USER] = OpCodeHandler{"MSG_GET_CITY_USER", STATUS_AUTHED, (*CSession).HandleGetCityUser}
	OpCodeTable[proto.MSG_GET_SIGNIN_LIST] = OpCodeHandler{"MSG_GET_SIGNIN_LIST", STATUS_AUTHED, (*CSession).HandleGetSignInList}
	OpCodeTable[proto.MSG_SIGN_IN] = OpCodeHandler{"MSG_SIGN_IN", STATUS_AUTHED, (*CSession).HandleSignIn}
	OpCodeTable[proto.MSG_GET_TASK_LIST] = OpCodeHandler{"MSG_GET_TASK_LIST", STATUS_AUTHED, (*CSession).HandleGetTaskList}
	OpCodeTable[proto.MSG_GET_TASK_AWARD] = OpCodeHandler{"MSG_GET_TASK_AWARD", STATUS_AUTHED, (*CSession).HandleGetTaskAward}
	OpCodeTable[proto.MSG_ENTER_SHOP] = OpCodeHandler{"MSG_ENTER_SHOP", STATUS_AUTHED, (*CSession).HandleEnterShop}

	//双旦活动接口
	OpCodeTable[proto.MSG_GET_SWEET_AND_TREE] = OpCodeHandler{"MSG_GET_SWEET_AND_TREE", STATUS_AUTHED, (*CSession).HandleGetSweetAndTree}
	OpCodeTable[proto.MSG_GET_PATCH] = OpCodeHandler{"MSG_GET_PATCH", STATUS_AUTHED, (*CSession).HandleGetPatch}
	OpCodeTable[proto.MSG_SWEET_TREE] = OpCodeHandler{"MSG_SWEET_TREE", STATUS_AUTHED, (*CSession).HandleTradeCdt}
	OpCodeTable[proto.MSG_GET_DOUBLE_YEAR_STATUS] = OpCodeHandler{"MSG_GET_DOUBLE_YEAR_STATUS", STATUS_AUTHED, (*CSession).HandleGetDoubleYearStatus}

	AgentCodeTable[proto.MSG_GET_ALL_ROLE] = AgentHandler{"MSG_GET_ALL_ROLE", STATUS_AUTHED, (*agent).HandleGetAllRole}
	AgentCodeTable[proto.MSG_USER_ADD_ROLE] = AgentHandler{"MSG_USER_ADD_ROLE", STATUS_AUTHED, (*agent).HandleUserAddRole}
	AgentCodeTable[proto.MSG_HEARTBEAT] = AgentHandler{"MSG_HEARTBEAT", STATUS_NOT_AUTHED, (*agent).HandleHEARTBEAT}
	AgentCodeTable[proto.MSG_CREATER_ROLE] = AgentHandler{"MSG_CREATER_ROLE", STATUS_NOT_AUTHED, (*agent).HandleCreateRole}
	AgentCodeTable[proto.MSG_LOGIN] = AgentHandler{"MSG_LOGIN", STATUS_NOT_AUTHED, (*agent).HandleLogin}
	AgentCodeTable[proto.MSG_GET_VERIFICATION_CODE] = AgentHandler{"MSG_GET_VERIFICATION_CODE", STATUS_NOT_AUTHED, (*agent).HandleGetVerificationCode}
	AgentCodeTable[proto.MSG_CHECK_NICK_NAME] = AgentHandler{"MSG_CHECK_NICK_NAME", STATUS_NOT_AUTHED, (*agent).HandleCheckNickName}
	AgentCodeTable[proto.MSG_REGISTER_EMAIL] = AgentHandler{"MSG_REGISTER_EMAIL", STATUS_NOT_AUTHED, (*agent).HandleEmailRegister}
	AgentCodeTable[proto.MSG_REGISTER_PHONE] = AgentHandler{"MSG_REGISTER_PHONE", STATUS_NOT_AUTHED, (*agent).HandleMobileRegister}
	AgentCodeTable[proto.MSG_RESET_PASSWORD] = AgentHandler{"MSG_RESET_PASSWORD", STATUS_NOT_AUTHED, (*agent).HandleResetPwd}
	AgentCodeTable[proto.MSG_ENTER_CITY] = AgentHandler{"MSG_ENTER_CITY", STATUS_NOT_AUTHED, (*agent).HandleEnterCity}
	AgentCodeTable[proto.MSG_REBIND] = AgentHandler{"MSG_REBIND", STATUS_NOT_AUTHED, (*agent).HandleRebind}
	AgentCodeTable[proto.MSG_GET_EMAIL_LIST] = AgentHandler{"MSG_GET_EMAIL_LIST", STATUS_NOT_AUTHED, (*agent).HandleGetEmailList}
	AgentCodeTable[proto.MSG_DEL_EMAIL] = AgentHandler{"MSG_DEL_EMAIL", STATUS_NOT_AUTHED, (*agent).HandleDelEmails}
	AgentCodeTable[proto.MSG_SET_EMAIL_READ] = AgentHandler{"MSG_SET_EMAIL_READ", STATUS_NOT_AUTHED, (*agent).HandleSetEmailRead}
	AgentCodeTable[proto.MSG_COUNT_EMAIL] = AgentHandler{"MSG_COUNT_EMAIL", STATUS_NOT_AUTHED, (*agent).HandleGetEmailCount}
	//AgentCodeTable[proto.MSG_PUSH_USER_MEIAL] = AgentHandler{"MSG_PUSH_USER_MEIAL", STATUS_NOT_AUTHED, (*agent).HandleSendRealNameEmail}  // 已删除
	//AgentCodeTable[proto.MSG_PUSH_NOTICE] = AgentHandler{"MSG_PUSH_NOTICE", STATUS_NOT_AUTHED, (*agent).HandlePushNotice} // 已删除
	AgentCodeTable[proto.MSG_GET_UPGRADE_NOTICE] = AgentHandler{"MSG_GET_UPGRADE_NOTICE", STATUS_NOT_AUTHED, (*agent).HandleGetUpgradeNotice}
	AgentCodeTable[proto.MSG_GET_MERCHANTS_URL] = AgentHandler{"MSG_GET_MERCHANTS_URL", STATUS_NOT_AUTHED, (*agent).HandleGetMerchantUrl} // 已废弃
	AgentCodeTable[proto.MSG_EMAIL_RECEIVE_REWARDS] = AgentHandler{"MSG_EMAIL_RECEIVE_REWARDS", STATUS_NOT_AUTHED, (*agent).EmailReceiveRewards}

	// 本地生活商品展示接口配置.
	AgentCodeTable[proto.MSG_LOCAL_LIFE_INDEX_RECOMMEND] = AgentHandler{"MSG_LOCAL_LIFE_INDEX_RECOMMEND", STATUS_NOT_AUTHED, (*agent).Recommend}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_TOP_SEARCH] = AgentHandler{"MSG_LOCAL_LIFE_TOP_SEARCH", STATUS_NOT_AUTHED, (*agent).TopSearchAndRecordWord}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_HOTEL_SEARCH] = AgentHandler{"MSG_LOCAL_LIFE_HOTEL_SEARCH", STATUS_NOT_AUTHED, (*agent).SearchHotel}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_HOTEL_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_HOTEL_DETAIL", STATUS_NOT_AUTHED, (*agent).GetHotelDetails}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_ROOM_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_ROOM_DETAIL", STATUS_NOT_AUTHED, (*agent).GetHotelRoomDetails}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD] = AgentHandler{"MSG_LOCAL_LIFE_DELETE_SEARCH_RECORD", STATUS_NOT_AUTHED, (*agent).DeleteSearchRecord}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CITY_LIST] = AgentHandler{"MSG_LOCAL_LIFE_CITY_LIST", STATUS_NOT_AUTHED, (*agent).CityList}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_STORE_ClASSIFY] = AgentHandler{"MSG_LOCAL_LIFE_STORE_ClASSIFY", STATUS_NOT_AUTHED, (*agent).StoreClassify}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_INDEX_SELECT_CITY] = AgentHandler{"MSG_LOCAL_LIFE_INDEX_SELECT_CITY", STATUS_NOT_AUTHED, (*agent).IndexSelectCity}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE] = AgentHandler{"MSG_LOCAL_LIFE_CATEGORY_SEARCH_STORE", STATUS_NOT_AUTHED, (*agent).CategoryStoreSearch}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_SEARCH_SUGGEST] = AgentHandler{"MSG_LOCAL_LIFE_SEARCH_SUGGEST", STATUS_NOT_AUTHED, (*agent).SearchSuggest}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_TOP_SEARCH_V2] = AgentHandler{"MSG_LOCAL_LIFE_TOP_SEARCH_V2", STATUS_NOT_AUTHED, (*agent).TopSearchV2}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_STORE_TYPE] = AgentHandler{"MSG_LOCAL_LIFE_STORE_TYPE", STATUS_NOT_AUTHED, (*agent).StoreType}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_STORE_HOTEL_DETAIL", STATUS_NOT_AUTHED, (*agent).StoreHotelDetail}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_STORE_RESTAURANT_DETAIL", STATUS_NOT_AUTHED, (*agent).StoreRestaurantDetail}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_GOODS_HOTEL_DETAIL", STATUS_NOT_AUTHED, (*agent).GoodsHotelDetail}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_DISCOUNT_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_DISCOUNT_DETAIL", STATUS_NOT_AUTHED, (*agent).DiscountDetail}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL] = AgentHandler{"MSG_LOCAL_LIFE_GOODS_RESTAURANT_DETAIL", STATUS_NOT_AUTHED, (*agent).RestaurantDetail}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CITY] = AgentHandler{"MSG_LOCAL_LIFE_CITY", STATUS_NOT_AUTHED, (*agent).CityInfos}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CITY_SUGGEST] = AgentHandler{"MSG_LOCAL_LIFE_CITY_SUGGEST", STATUS_NOT_AUTHED, (*agent).CitySuggest}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CITY_PICK] = AgentHandler{"MSG_LOCAL_LIFE_CITY_PICK", STATUS_NOT_AUTHED, (*agent).CityPick}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE] = AgentHandler{"MSG_LOCAL_LIFE_SEARCH_RECORD_DELETE", STATUS_NOT_AUTHED, (*agent).SearchRecordDeleteV2}
	AgentCodeTable[proto.MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS] = AgentHandler{"MSG_LOCAL_LIFE_CATEGORY_SEARCH_GOODS", STATUS_NOT_AUTHED, (*agent).CategoryGoodsSearch}

	// 统计时间
	AgentCodeTable[proto.MSG_STATISTICS_CITY_ICON] = AgentHandler{"MSG_STATISTICS_CITY_ICON", STATUS_NOT_AUTHED, (*agent).IconStatistics}
	// 统计宝箱有效广告
	AgentCodeTable[proto.MSG_TREASURE_BOX_ADVERTIS] = AgentHandler{"MSG_TREASURE_BOX_ADVERTIS", STATUS_NOT_AUTHED, (*agent).TreasureBoxStatistics}

	// 双旦活动.
	//	AgentCodeTable[proto.MSG_SWEET_TREE] = AgentHandler{"MSG_SWEET_TREE", STATUS_NOT_AUTHED, (*agent).TradeDayCdt}
	AgentCodeTable[proto.MSG_RANK_LIST_UPDATE_PROP] = AgentHandler{"MSG_RANK_LIST_UPDATE_PROP", STATUS_NOT_AUTHED, (*agent).RankListUpdateProp}
	AgentCodeTable[proto.MSG_RANK_LIST_DAY] = AgentHandler{"MSG_RANK_LIST_DAY", STATUS_NOT_AUTHED, (*agent).RankListDay}
	AgentCodeTable[proto.MSG_RANK_LIST_ALL] = AgentHandler{"MSG_RANK_LIST_ALL", STATUS_NOT_AUTHED, (*agent).RankListAll}
	AgentCodeTable[proto.MSG_RANK_LIST_DAY_PROP] = AgentHandler{"MSG_RANK_LIST_DAY_PROP", STATUS_NOT_AUTHED, (*agent).RankListDayProp}
	AgentCodeTable[proto.MSG_RANK_LIST_ALL_PROP] = AgentHandler{"MSG_RANK_LIST_ALL_PROP", STATUS_NOT_AUTHED, (*agent).RankListAllProp}
	AgentCodeTable[proto.MSG_RANK_LIST_INVITE_RECORD] = AgentHandler{"MSG_RANK_LIST_INVITE_RECORD", STATUS_NOT_AUTHED, (*agent).RankListInviteRecord}
	AgentCodeTable[proto.MSG_RANK_LIST_EMAIL_REPAIR_CDT] = AgentHandler{"MSG_RANK_LIST_EMAIL_REPAIR_CDT", STATUS_NOT_AUTHED, (*agent).RankListRepairCdt}

	AgentCodeTable[proto.MSG_DOUBLE_YEAR_DOT] = AgentHandler{"MSG_GET_CONFIGURE_URL", STATUS_NOT_AUTHED, (*agent).RankingDot}

	//宝箱活动
	AgentCodeTable[proto.MSG_OPEN_STREASURE_BOX] = AgentHandler{"MSG_OPEN_STREASURE_BOX", STATUS_NOT_AUTHED, (*agent).OpenStreasureBox}
	AgentCodeTable[proto.MSG_FINISH_STREASURE_BOX] = AgentHandler{"MSG_FINISH_STREASURE_BOX", STATUS_NOT_AUTHED, (*agent).FinishStreasureBox}
	AgentCodeTable[proto.MSG_RECEIVE_STREASURE_BOX] = AgentHandler{"MSG_RECEIVE_STREASURE_BOX", STATUS_NOT_AUTHED, (*agent).ReceiveBoxReward}
	AgentCodeTable[proto.MSG_GET_STREASURE_BOX_RECORD] = AgentHandler{"MSG_GET_STREASURE_BOX_RECORD", STATUS_NOT_AUTHED, (*agent).ReceiveBoxGetRecord}
	AgentCodeTable[proto.MSG_ACTIVITY_STATYS] = AgentHandler{"MSG_ACTIVITY_STATYS", STATUS_NOT_AUTHED, (*agent).ActivitStatus}
	AgentCodeTable[proto.MSG_TREASURE_BOX_DAY_NUM] = AgentHandler{"MSG_ACTIVITY_STATYS", STATUS_NOT_AUTHED, (*agent).GetStreasureBoxDayNum}
	//配置URL
	AgentCodeTable[proto.MSG_GET_CONFIGURE_URL] = AgentHandler{"MSG_GET_CONFIGURE_URL", STATUS_NOT_AUTHED, (*agent).GetConfigureUrl}
}
