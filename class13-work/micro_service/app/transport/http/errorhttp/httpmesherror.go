package errorhttp

import (
	"micro_service/app/model/api"

	json "github.com/json-iterator/go"
)

type HttpMeshError struct {
	Code    int
	Message string
}

func (herr *HttpMeshError) Error() string {
	return herr.Message
}

//http v1版本服务在transport层直接返回错误时使用的代码 （特性化代码 不具有通用性）
func (herr *HttpMeshError) MarshalJSON() ([]byte, error) {
	body := struct {
		Code    int         `json:"code"`
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
	}{
		Code:    herr.Code,
		Message: herr.Message,
		Data:    api.EmptyRespData{},
	}
	bodyByte, errs := json.Marshal(body)
	return bodyByte, errs
}

//http v1版本在返回错误的时候必须指定为200
func (herr *HttpMeshError) StatusCode() int {
	return 200
}
