package util

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

const (
	XMLHeader  = "xml"
	JSONHeader = "json"
)

const (
	POST = "POST"
	GET  = "GET"
)

type Curl interface {
	SetHeader(key, val string) //设置header头
	Do() ([]byte, error)       //执行请求
}

//请求对象
type ReqParams struct {
	Url    string //地址
	Method string //请求方法
	Header string //请求头 JSON或者XML
	Params []byte //请求参数
}

type reqObj struct {
	req *http.Request
}

//初始请求参数
func (p *ReqParams) InitRequest() (req Curl, err error) {
	var reqParams *bytes.Reader
	obj := new(reqObj)

	if p.Params != nil {
		reqParams = bytes.NewReader(p.Params)
		obj.req, err = http.NewRequest(p.Method, p.Url, reqParams)
	} else {
		obj.req, err = http.NewRequest(p.Method, p.Url, nil)
	}

	if err != nil {
		return nil, err
	}
	if p.Method == POST {
		switch p.Header {
		case JSONHeader:
			obj.req.Header.Set("Content-Type", "application/json;charset=UTF-8")
			break
		case XMLHeader:
			obj.req.Header.Set("Accept", "application/xml")
			obj.req.Header.Set("Content-Status", "application/xml;charset=utf-8")
			break
		default:
			obj.req.Header.Set("Content-Type", "application/json;charset=UTF-8")
		}
	}

	return obj, nil
}

//设置header头
func (obj *reqObj) SetHeader(key, val string) {
	obj.req.Header.Set(key, val)
}

//执行请求
func (obj *reqObj) Do() ([]byte, error) {
	defer func() {
		if er := recover(); er != nil {
			fmt.Print(fmt.Errorf("%v", er))
		}
	}()
	c := http.Client{}
	resp, err := c.Do(obj.req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

// func main() {
// 	// params := map[string]string{
// 	// 	"wd": "测试数据",
// 	// }
// 	// body, _ := json.Marshal(params)
// 	reqParam := &ReqParams{
// 		Url:    "http://localhost:3000/api/test",
// 		Method: "GET",
// 		Header: "json",
// 		//Params: body,
// 	}
// 	req, err := reqParam.InitRequest()
// 	if err != nil {
// 		fmt.Println("request err", err)
// 		return
// 	}
// 	byteData, _ := req.Do()
// 	fmt.Println(string(byteData))
// }
