package proto

//获取商户URL
type C2SGetMerchantUrl struct {
	Token string `json:"token"`
}

//获取聚合支付接口返回
type S2CGetMerchantUrl struct {
	S2CCommon
	MerchantEnteringUrl  string `json:"merchantEnteringUrl"`
	ActivityPromotionUrl string `json:"activityPromotionUrl"`
}

//获取商户URL
type C2SConfigureUrl struct {
	Token string `json:"token"`
}

//获取聚合支付接口返回
type S2CConfigureUrl struct {
	S2CCommon
	Data map[string]string `json:"data"`
}
