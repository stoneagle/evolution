package req

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"time"
)

type ReqMethod string

const (
	GetMethod  ReqMethod = "GET"
	PostMethod ReqMethod = "POST"
)

func Exec(method ReqMethod, url string, params []byte, timeout time.Duration) (body []byte, err error) {
	var req *http.Request
	switch method {
	case GetMethod:
		req, err = http.NewRequest(string(method), url, nil)
	case PostMethod:
		req, err = http.NewRequest(string(method), url, bytes.NewBuffer(params))
		req.Header.Set("Content-Type", "application/json")
	}
	if err != nil {
		return
	}
	client := &http.Client{
		Timeout: time.Duration(timeout),
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)
	return
}
