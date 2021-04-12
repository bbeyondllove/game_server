package base

import (
	"io/ioutil"
	"time"

	"game_server/core/logger"

	"gopkg.in/yaml.v2"
)

type SettingStruct struct {
	Server struct {
		Debug       bool          `yaml:"debug"`
		LocalHost   string        `yaml:"local_host"`
		LocalPort   string        `yaml:"local_port"`
		HttpPort    string        `yaml:"http_port"`
		MaxConnNum  int           `yaml:"max_conn_num"`
		HttpTimeout time.Duration `yaml:"http_timeout"`
		MaxMsgLen   uint32        `yaml:"max_msg_len"`
		CertFile    string        `yaml:"cert_file"`
		KeyFile     string        `yaml:"key_file"`
	}

	Mysql struct {
		Host        string `yaml:"host"`
		Port        string `yaml:"port"`
		DbName      string `yaml:"dbName"`
		Username    string `yaml:"username"`
		Password    string `yaml:"password"`
		Parameter   string `yaml:"parameter"`
		TablePrefix string `yaml:"tablePrefix"`
		Connects    int    `yaml:"connects"`
		Goroutines  int    `yaml:"goroutines"`
	}

	Redis struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Password string `yaml:"password"`
		DbName   int    `yaml:"db_name"`
		PoolSize int    `yaml:"pool_size"`
	}

	Base struct {
		UserHost            string `yaml:"user_host"`
		UserPort            string `yaml:"user_port"`
		SecurityKey         string `yaml:"security_key"`
		MobileCodeUrl       string `yaml:"mobile_code_url"`
		EmailCodeUrl        string `yaml:"email_code_url"`
		MobileRegisterUrl   string `yaml:"mobile_register_url"`
		EmailRegisterUrl    string `yaml:"email_register_url"`
		MobileLoginUrl      string `yaml:"mobile_login_url"`
		EmailLoginrUrl      string `yaml:"email_login_url"`
		GetUserinfoUrl      string `yaml:"get_userinfo_url"`
		MobileRetsetPwdUrl  string `yaml:"mobile_reset_password_url"`
		EmailRetsetPwdUrl   string `yaml:"email_reset_password_url"`
		GetInviteUsersUrl   string `yaml:"get_invite_users_url"`
		GetGrabComradesUrl  string `yaml:"get_grab_comrades_url"`
		GetMemberSysUrl     string `yaml:"get_member_sys_url"`
		GetUserLevelUrl     string `yaml:"get_user_level_url"`
		GrabComradeUrl      string `yaml:"grab_comrade_url"`
		UpdateUserUrl       string `yaml:"update_user_url"`
		RefreshTokenUrl     string `yaml:"refresh_token_url"`
		BindingInviterUrl   string `yaml:"user_binding_inviter"`
		InviteIncomelistUrl string `yaml:"user_invite_incomelist"`
		GetExchangeListMax  string `yaml:"get_exchange_list_max"`
		UserInvitationUrl   string `yaml:"user_invitation"`
		UserSaveKycUrl      string `yaml:"user_save_kyc"`
	}

	Pay struct {
		TokenPayUrl string `yaml:"token_pay_url"`
	}

	Life struct {
		LocalIp   string `yaml:"http_ip"`
		LocalPort string `yaml:"http_port"`
		Timeout   int    `yaml:"time_out"`
		Secret    string ``
		Debug     bool   `yaml:"debug"`
		HttpUrl   string `yaml:"http_url"`
	}

	Cdt struct {
		CdtUsdRateDefault   float32                  `yaml:"cdt_usd_rate_default"`
		UsdCnyRateDefault   float32                  `yaml:"usd_cny_rate_default"`
		CnyPointRateDefault float32                  `yaml:"cny_point_rate_default"`
		LimitDayCdt         float32                  `yaml:"limit_day_cdt"`
		SpecialLimitDayCdt  float32                  `yaml:"special_limit_day_cdt"`
		EventType           []map[string]interface{} `yaml:"event_type"`
	}
	Statistical struct {
		MakeUpCard int `yaml:"make_up_card"`
	}

	// Doubleyear 双旦活动配置.
	Doubleyear struct {
		TradeCdtDayLimit         float32 `yaml:"trade_cdt_day_limit"`
		PropStartDate            string  `yaml:"prop_start_date"`
		PropEndDate              string  `yaml:"prop_end_date"`
		PropStartTime            string  `yaml:"prop_start_time"`
		PropEndTime              string  `yaml:"prop_end_time"`
		TradeCdtDayTopLimit      float32 `yaml:"trade_cdt_day_top_limit"`
		TradeCdtDayTopTotalLimit float32 `yaml:"trade_cdt_day_top_total_limit"`
		SweetTreeCdtRate         float32 `yaml:"sweet_tree_cdt_rate"`
		FatherChristmasRoleRate  int     `yaml:"father_christmas_role_rate"`

		ChristmasEventInterval       int    `yaml:"christmas_event_interval"`
		ChristmasUnitMaxNum          int    `yaml:"christmas_unit_max_num"`
		ChristmasPatchTradeNum       int    `yaml:"christmas_patch_trade_num"`
		SantaClausExpriteTime        string `yaml:"santa_claus_exprite_time"`
		SantaClausUnlockEffectiveDay int    `yaml:"santa_claus_unlock_effective_day"`
		SantaClausCardId             int    `yaml:"santa_claus_card_id"`
		SantaClausRoleId             int    `yaml:"santa_claus_role_id"`
		SantaEmailTitle              string `yaml:"santa_email_tile"`
		SantaEmailContent            string `yaml:"santa_email_content"`
		SantaBroadContent            string `yaml:"santa_broad_content"`
	}

	// 商户配置
	Merchant struct {
		MerchantEnteringUrl  string `yaml:"merchant_entering_url"`
		ActivityPromotionUrl string `yaml:"activity_promotion_url"`
	}
	// 配置URL
	ConfigureUrl map[string]string `yaml:"configure_url"`
	// 邮箱配置
	Email struct {
		EmailExpireDay int `yaml:"email_expire_day"`
	}
	// 活动角色
	ActivityRoles []map[string]interface{} `yaml:"activity_roles"`
	// 新春活动
	Springfestival struct {
		GoldenCoupleRole      []int  `yaml:"golden_couple_role"`
		RankingListStartDate  string `yaml:"ranking_list_start_date"`
		RankingListEndDate    string `yaml:"ranking_list_end_date"`
		ActivityStartDatetime string `yaml:"activity_start_datetime"`
		ActivityEndDatetime   string `yaml:"activity_end_datetime"`
	}
	// 管理后台
	Admin struct {
		HttpPort string            `yaml:"http_port"`
		Oss      map[string]string `yaml:"oss"`
	}
}

var Setting = new(SettingStruct)

func init() {
	defer logger.Flush()
	yamlFile, err := ioutil.ReadFile("config.yml")
	if err != nil {
		logger.Errorf("Load config error:", err)
		return
	}
	err = yaml.Unmarshal(yamlFile, Setting)
	if err != nil {
		logger.Errorf("Config unmarshal error:", err)
	}
}
