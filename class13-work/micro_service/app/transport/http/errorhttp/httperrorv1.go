package errorhttp

import (
	"time"

	json "github.com/json-iterator/go"
)

type HttpErrorV1 struct {
	Code    int
	Message string
}

func (herr *HttpErrorV1) Error() string {
	return herr.Message
}

//http v1版本服务在transport层直接返回错误时使用的代码 （特性化代码 不具有通用性）
func (herr *HttpErrorV1) MarshalJSON() ([]byte, error) {
	body := struct {
		ApiStatus int         `json:"apiStatus"`
		SysStatus int         `json:"sysStatus"`
		Info      string      `json:"info"`
		Timestamp int         `json:"timestamp"`
		Data      interface{} `json:"data"`
	}{
		ApiStatus: herr.Code,
		SysStatus: 0,
		Info:      herr.Message,
		Timestamp: int(time.Now().Unix()),
		Data:      nil,
	}
	bodyByte, errs := json.Marshal(body)
	return bodyByte, errs
}

//http v1版本在返回错误的时候必须指定为200
func (herr *HttpErrorV1) StatusCode() int {
	return 200
}
