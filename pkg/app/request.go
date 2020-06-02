package app

import (
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

func HttpPost(requestUrl string, requestId string, reqBody interface{}) {
	byteData, err := json.Marshal(reqBody)
	if err != nil {
		ulog.Error("request body error:", err)
		return
	}
	body := bytes.NewReader(byteData)
	//创建post请求
	req, err := http.NewRequest("POST", requestUrl, body)
	if err != nil {
		ulog.Errorf("request uri failed:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json;charset=UTF-8")
	req.Header.Set("Connection", "close")
	if requestId != "" {
		req.Header.Set("RequestId", requestId)
	}

	headers, _ := json.Marshal(req.Header)
	ulog.Debugf("request url:%s, headers:%s, data:%s", requestUrl, headers, byteData)
	//生成client
	client := &http.Client{}
	resp := &http.Response{}
	for i := 0; ; i++ {
		//发送请求
		resp, err = client.Do(req)
		if err != nil {
			ulog.Errorf("request url:%s, body:%s failed.Error msg:%s.", requestUrl, byteData, err)
			if i < 3 && resp != nil && resp.StatusCode == http.StatusInternalServerError {
				time.Sleep(100 * time.Millisecond)
				ulog.Infof("server return 500, try again", i)
				continue
			}
			return
		}
		ulog.Infof("request response status code:%d", resp.StatusCode)
		return
	}
}
