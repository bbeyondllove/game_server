package message

import (
	"io"
	"net/http"
	"net/http/httptest"
	// "strings"
	"testing"

	"github.com/julienschmidt/httprouter"
)

//充值
func TestChargeDeduce(t *testing.T) {

	// var (
	// 	param    io.Reader
	// 	expected string
	// 	rr       *httptest.ResponseRecorder
	// 	err      error
	// )
	// //----------------请求参数为空时 start -----------------------
	// param = strings.NewReader("")
	// expected = `{"code":9002,"result":"","text":"Invalid arguments"}`

	// rr, err = newRequestRecorder(http.MethodPost, "/api/ChargeDeduce", param, ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }
	// // 检测返回的状态码
	// if status := rr.Code; status != http.StatusOK {
	// 	t.Errorf("handler returned wrong status code: got %v want %v",
	// 		status, http.StatusOK)
	// }

	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
	//----------------请求参数为空时 end -----------------------

	// //代理登录名不存在
	// param = strings.NewReader("agent=dlebo011")
	// expected = `{"code":9004,"result":"","text":"Invalid Agent"}`
	// rr,err = newRequestRecorder("POST", "/deposit", param,ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

	// //金额为负数
	// param = strings.NewReader("agent=dlebo01&username=test1002&amount=-10")

	// rr,err = newRequestRecorder("POST", "/deposit", param,ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// //验证ip，ip有限制，下面断言就不执行
	// expected = `"code":9003`
	// if strings.Index(rr.Body.String(),expected) != -1 {
	// 	return
	// }
	// expected = `{"code":9002,"result":"","text":"Invalid arguments"}`
	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

	// //token不正确
	// param = strings.NewReader("agent=dlebo01&username=test1002&amount=10&token=sdsdsdsdsdsd")
	// expected = `{"code":9002,"result":"","text":"Invalid arguments"}`
	// rr,err = newRequestRecorder("POST", "/deposit", param,ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if rr.Body.String() != expected {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

	// //不存在的玩家，会先注册，然后再充值，返回成功信息
	// param = strings.NewReader("agent=dlebo01&username=mytest1002&amount=10&token=de14e6bc2c21fdbe6d6c1fe9487c370aebcc8dab")
	// expected = `"code":0`
	// rr,err = newRequestRecorder("POST", "/deposit", param,ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if strings.Index(rr.Body.String(),expected) == -1 {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }

	// //玩家存在，充值
	// param = strings.NewReader("agent=dlebo01&username=test1002&amount=10&token=a0482cc5d8b11f6f6d1ed8065dd48b41f3ce02a7")
	// expected = `"code":0`
	// rr,err = newRequestRecorder("POST", "/deposit", param,ChargeDeduce)
	// if err != nil {
	// 	t.Fatal(err)
	// }

	// if strings.Index(rr.Body.String(),expected) == -1 {
	// 	t.Errorf("handler returned unexpected body: got %v want %v",
	// 		rr.Body.String(), expected)
	// }
}

/**
发起请求公共函数
*/
func newRequestRecorder(method string, strPath string, params io.Reader, fnHandler func(w http.ResponseWriter, r *http.Request, param httprouter.Params)) (*httptest.ResponseRecorder, error) {
	req, err := http.NewRequest(method, strPath, params)
	if err != nil {
		return nil, err
	}
	//设置头部信息
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("X-real-ip", "127.0.0.1") //构造一个ip地址
	//创建一个 ResponseRecorder 来记录响应
	rr := httptest.NewRecorder()
	fnHandler(rr, req, nil)
	return rr, nil
}
