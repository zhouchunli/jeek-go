package httpcli

import (
	"bytes"
	"io/ioutil"
	"micro_service/app/model/api"
	"micro_service/pkg/microerr"
	"net/http"
	"time"

	json "github.com/json-iterator/go"
)

// 发送GET请求
// url：         请求地址
func Get(url string) (string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "*/*")
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, _ := ioutil.ReadAll(resp.Body)

	return string(result), nil
}

// 发送POST请求
// url：         请求地址
// data：        POST请求提交的数据
// contentType： 请求体格式，如：application/json
func Post(url string, data interface{}, contentType string) (string, error) {
	if contentType == "" {
		contentType = "application/json"
	}

	// 超时时间：5秒
	client := &http.Client{Timeout: 5 * time.Second}
	jsonStr, _ := json.Marshal(data)
	resp, err := client.Post(url, contentType, bytes.NewBuffer(jsonStr))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	result, _ := ioutil.ReadAll(resp.Body)
	return string(result), nil
}

//发送http请求
func sendHttp(url string, body interface{}) (sysErr error) {
	reshttp := api.HttpRes{}
	str, err := Post(url, body, "")
	if err != nil {
		sysErr = microerr.ServerCommonError.Wrap(err, "Post err")
		return sysErr
	}
	err = json.Unmarshal([]byte(str), &reshttp)
	if err != nil {
		sysErr = microerr.SerializeError.Wrap(err, "Unmarshal err")
		return sysErr
	}
	if reshttp.Code != int(microerr.Success) {
		sysErr = microerr.ErrType(reshttp.Code).New(reshttp.Msg)
		return sysErr
	}
	return nil
}
