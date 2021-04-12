package proto

//用户获取验证码
type C2SGetVerificationCode struct {
	SysType     int    `json:"sysType"` //业务系统类型，0：Base系统；1：EChain商城系统；2：游戏系统；3：Aries交易所；4：CCMYL交易所
	CountryCode int    `json:"countryCode"`
	Mobile      string `json:"mobile"`
	Email       string `json:"email"`
	UseFor      int    `json:"useFor"`   //1：注册；2：重置密码；3：绑定（手机或邮箱）
	CodeType    int    `json:"codeType"` //1：手机方式；2：邮箱方式
	Language    string `json:"language"` //下发短信的语言，默认为英文：en；简体中文：zh；繁体中文：tc
}
