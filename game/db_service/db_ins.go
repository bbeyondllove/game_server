package db_service

var (
	GameWalletIns                           *GameWallet
	ExchangeRecordIns                       *ExchangeRecord
	UserIns                                 *User
	HouseIns                                *House
	WorldMapIns                             *WorldMap
	UserLevelDescIns                        *UserLevelConfig
	UserKnapsackIns                         *UserKnapsack
	ItemIns                                 *Items
	BuildTypeIns                            *BuildType
	RoleInfoIns                             *RoleInfo
	CertificationIns                        *Certification
	UserLevelConfigIns                      *UserLevelConfig
	TasksIns                                *Tasks
	TaskAwardIns                            *TaskAward
	AwardsIns                               *Awards
	UserTaskIns                             *UserTask
	EmailIns                                *EmailInfo
	NoticeIns                               *NoticeInfo
	TimeNoticeIns                           *TimeNoticeInfo
	ActivityDataIns                         *ActivityDateInfo
	DoubleYearIns                           *DoubleYearUserItem
	ActivityInvitationIns                   *ActivityInvite
	EmailPrizeIns                           *EmailPrizeInfo
	ActivityRolesIns                        *ActivityRoles
	TreasureBoxCdtConfigIns                 *TreasureBoxCdtConfig
	UserTreasureBoxRecordIns                *UserTreasureBoxRecord
	CertificationRecordIns                  *CertificationRecord
	AdminUserIns                            *AdminUser
	AdminRoleIns                            *AdminRole
	AdminLimitIns                           *AdminLimit
	ActivityRedEnvelopeIns                  *ActivityRedEnvelope
	ActivityConfigIns                       *ActivityConfig
	StatisticsDayIns                        *StatisticsDay
	StatisticsRetainedIns                   *StatisticsRetained
	StatisticsActiveCountIns                *StatisticsActiveCount
	StatisticsTreasureBoxIns                *StatisticsTreasureBox
	StatisticsRealTreasureBoxIns            *StatisticsRealTreasureBox
	StatisticsSignInIns                     *StatisticsSignIn
	StatisticsDoubleYearCdtIns              *StatisticsDoubleYearCdt
	StatisticsDoubleYearUserDayCdtIns       *StatisticsDoubleYearUserDayCdt
	StatisticsDoubleYearFragmentIns         *StatisticsDoubleYearFragment
	StatisticsDoubleYearUserFragmentIns     *StatisticsDoubleYearUserFragment
	StatisticsDoubleYearDailyRankingIns     *StatisticsDoubleYearDailyRanking
	StatisticsDoubleYearUserDailyRankingIns *StatisticsDoubleYearUserDailyRanking
	StatisticsDoubleYearTotalRankingIns     *StatisticsDoubleYearTotalRanking
	StatisticsDoubleYearUserTotalRankingIns *StatisticsDoubleYearUserTotalRanking
	UserNoticeIns                           *UserNotice
)

func DbInit() {
	GameWalletIns = &GameWallet{}
	UserIns = &User{}
	HouseIns = &House{}
	WorldMapIns = &WorldMap{}
	ExchangeRecordIns = &ExchangeRecord{}
	UserLevelConfigIns = &UserLevelConfig{}
	UserKnapsackIns = &UserKnapsack{}
	ItemIns = &Items{}
	BuildTypeIns = &BuildType{}
	CertificationIns = &Certification{}
	TasksIns = &Tasks{}
	TaskAwardIns = &TaskAward{}
	AwardsIns = &Awards{}
	UserTaskIns = &UserTask{}
	EmailIns = &EmailInfo{}
	NoticeIns = &NoticeInfo{}
	TimeNoticeIns = &TimeNoticeInfo{}
	ActivityDataIns = &ActivityDateInfo{}
	DoubleYearIns = &DoubleYearUserItem{}
	ActivityInvitationIns = &ActivityInvite{}
	EmailPrizeIns = &EmailPrizeInfo{}
	ActivityRolesIns = &ActivityRoles{}
	TreasureBoxCdtConfigIns = &TreasureBoxCdtConfig{}
	UserTreasureBoxRecordIns = &UserTreasureBoxRecord{}
	CertificationRecordIns = &CertificationRecord{}
	AdminUserIns = &AdminUser{}
	AdminRoleIns = &AdminRole{}
	AdminLimitIns = &AdminLimit{}
	ActivityRedEnvelopeIns = &ActivityRedEnvelope{}
	ActivityConfigIns = &ActivityConfig{}
	StatisticsDayIns = &StatisticsDay{}
	StatisticsRetainedIns = &StatisticsRetained{}
	StatisticsActiveCountIns = &StatisticsActiveCount{}
	StatisticsTreasureBoxIns = &StatisticsTreasureBox{}
	StatisticsRealTreasureBoxIns = &StatisticsRealTreasureBox{}
	StatisticsSignInIns = &StatisticsSignIn{}
	StatisticsDoubleYearCdtIns = &StatisticsDoubleYearCdt{}
	StatisticsDoubleYearUserDayCdtIns = &StatisticsDoubleYearUserDayCdt{}
	StatisticsDoubleYearFragmentIns = &StatisticsDoubleYearFragment{}
	StatisticsDoubleYearUserFragmentIns = &StatisticsDoubleYearUserFragment{}
	StatisticsDoubleYearDailyRankingIns = &StatisticsDoubleYearDailyRanking{}
	StatisticsDoubleYearUserDailyRankingIns = &StatisticsDoubleYearUserDailyRanking{}
	StatisticsDoubleYearTotalRankingIns = &StatisticsDoubleYearTotalRanking{}
	StatisticsDoubleYearUserTotalRankingIns = &StatisticsDoubleYearUserTotalRanking{}
	UserNoticeIns = &UserNotice{}
}
