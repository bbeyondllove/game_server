package proto

//获取聚合支付接口
type C2SGetTokenPayUrl struct {
	Token string `json:"token"`
}

//获取聚合支付接口返回
type S2cGetTokenPayUrl struct {
	S2CCommon
	SharePayUrl string `json:"sharePayUrl"`
}
