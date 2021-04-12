package local_life

const (
	// LocalLifeVerifyOk 验证公共参数OK编码.
	LocalLifeVerifyOk = iota
)

// 状态码.
const (
	// LocalLifeSuccess 请求成功.
	LocalLifeSuccess = 200
	// LocalLifeRequireToken 请求参数必需要有token.
	LocalLifeRequireToken = 1000 + iota
	// LocalLifeRequireSignature 请求参数必需要有sign签名.
	LocalLifeRequireSignature
	// LocalLifeSignatureError sign签名错误.
	LocalLifeSignatureError
	// LocalLifeRequireArgs 缺少相关参数，{args}为对应的实际具体参数.
	LocalLifeRequireArgs
	// LocalLifeServerBusy 调用本地生活接口失败或超时.
	LocalLifeServerBusy
	// LocalLifeTokenError token不合法.
	LocalLifeTokenError
	// LocalLifeArgumentError 解析客户端参数错误.
	LocalLifeArgumentError
)

// statusCodeMessage 状态码对应信息.
var statusCodeMessage = map[int]string{
	LocalLifeSuccess:          "OK",
	LocalLifeRequireToken:     "require token",
	LocalLifeRequireSignature: "require sign",
	LocalLifeSignatureError:   "signature error",
	LocalLifeRequireArgs:      "require args [%v]", // 缺少必要参数，{%v}代表具体参数名.
	LocalLifeServerBusy:       "server is busy",
	LocalLifeTokenError:       "token error",
	LocalLifeArgumentError:    "parse argument error",
}
