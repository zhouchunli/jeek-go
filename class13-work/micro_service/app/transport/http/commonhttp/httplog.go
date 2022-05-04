package commonhttp

import (
	"context"
	"encoding/json"
	"micro_service/app/deps"
	"micro_service/app/library/utils"
	"net/http"
	"time"
)

//http请求成功返回
func HttpLogging(ctx context.Context, req interface{}, resp interface{}, code int) {
	var reqStr string
	var respStr string
	var guid string

	//格式化返回数据
	respByte, err := json.Marshal(resp)
	if err != nil {
		return
	}
	respStr = string(respByte)

	//格式化请求数据
	if reqs, ok := req.(string); ok {
		reqStr = reqs
	} else {
		reqByte, err := json.Marshal(req)
		if err != nil {
			return
		}
		reqStr = string(reqByte)
	}
	//获取请求消息
	r := ctx.Value("httpRequest")
	httpReq, ok := r.(*http.Request)
	if !ok {
		return
	}
	//获取原始请求参数
	originalReqByte := GetRequestString(httpReq)
	//取出请求时间
	rquestTime := ctx.Value("requestTime")
	startTime, ok := rquestTime.(time.Time)
	if !ok {
		return
	}
	useTime := time.Now().Sub(startTime)

	if guid == "" {
		guid = "NA"
	}
	traceId := utils.GetContextTraceId(ctx)
	deps.Logger.HttpAccess(httpReq, traceId, "NA", reqStr, string(originalReqByte), guid, respStr, code, useTime)
}

//transportant校验失败需要打印日志
func HTTPDecodeParamsFailLogging(ctx context.Context, r *http.Request, req interface{}) {
	respStr := "NA" //没有返回值

	reqStr, _ := json.Marshal(req)
	//获取请求时间
	requestTime := ctx.Value("requestTime")
	startTime, ok := requestTime.(time.Time)
	if !ok {
		return
	}
	useTime := time.Now().Sub(startTime)

	//获取原始请求参数
	originalReqByte := GetRequestString(r)
	traceId := utils.GetContextTraceId(ctx)

	deps.Logger.HttpAccess(r, traceId, "NA", string(reqStr), string(originalReqByte), "NA", respStr, 500, useTime) //错误码固定为500
}
