package utils

import (
	"game_server/game/proto"
	"io/ioutil"
	"net/http"
	"strings"

	"game_server/core/logger"
)

//http的post请求
func HttpPost(reqUrl string, reqData string, encodeType int) ([]byte, error) {
	logger.Debugf("httpPost request:", reqUrl, reqData)
	var resp *http.Response
	var err error
	if encodeType == proto.JSON {
		resp, err = http.Post(reqUrl, "application/json", strings.NewReader(reqData))
	} else {
		resp, err = http.Post(reqUrl, "application/x-www-form-urlencoded", strings.NewReader(reqData))
	}

	if err != nil {
		logger.Error("httpPost failed err:", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debug("httpPost response:", string(body))
	return body, nil
}

//http的Get请求
func HttpGet(reqUrl string, request_data map[string]interface{}) ([]byte, error) {
	data := ""
	if request_data != nil {
		data = "?"
		i := 0
		for k, v := range request_data {
			value := Strval(v)
			data = data + k + "=" + value

			i++
			if i < len(request_data) {
				data = data + "&"
			}
		}
	}
	logger.Debugf("HttpGet request data:", reqUrl+data)
	resp, err := http.Get(reqUrl + data)

	if err != nil {
		logger.Errorf("HttpGet failed err:", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	logger.Debugf("httpPost response:", string(body))
	return body, nil
}

// HttpDelete http中的Delete方法.
func HttpDelete(reqUrl string, param map[string]interface{}) ([]byte, error) {
	req, _ := http.NewRequest("DELETE", reqUrl, nil)

	// 设置参数.
	q := req.URL.Query()
	for k, v := range param {
		q.Add(k, Strval(v))
	}

	logger.Debugf("HttpDelete request data:", reqUrl, param)
	req.URL.RawQuery = q.Encode()
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		logger.Errorf("HttpDelete request fail:", err)
		return nil, err
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		logger.Errorf("HttpDelete read response body fail:", err)
		return nil, err
	}

	return body, nil
}

// 把客户端中的IP端口号去掉
func GetIp(r *http.Request) string {
	var ipAddr string
	if ipAddr = r.Header.Get("X-real-ip"); ipAddr == "" {
		ipAddr = r.RemoteAddr
	}
	return strings.Split(ipAddr, ":")[0]
}
