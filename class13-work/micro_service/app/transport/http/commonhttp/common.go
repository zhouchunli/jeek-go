package commonhttp

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"micro_service/app/constant"
	"micro_service/app/deps"
	"micro_service/app/global"
	"micro_service/app/library/utils"
	"micro_service/app/model/api"
	"micro_service/app/transport/http/errorhttp"
	"micro_service/pkg/microerr"
	"net"
	"net/url"
	"strings"
	"time"

	"google.golang.org/grpc/metadata"

	"net/http"

	httptransport "github.com/go-kit/kit/transport/http"
)

// 将http请求信息附加到上下文中
func AddHTTPRequestToContext() httptransport.ServerOption {
	return httptransport.ServerBefore(
		func(ctx context.Context, r *http.Request) context.Context {
			//判断tranceid是否存在，不存在则生成
			if r.Header.Get("x-b3-traceid") == "" {
				ctx = utils.CreateTracingToContext(ctx, global.AppConfig.Common.ServiceName)
			} else {
				//添加tracing信息
				tracingKey := []string{
					"x-request-id",
					"x-b3-traceid",
					"x-b3-parentspanid",
					"x-b3-spanid",
					"x-b3-sampled",
					"x-b3-flags",
				}
				//增加tracing数据
				md := metadata.MD{}
				for _, v := range tracingKey {
					md.Set(v, r.Header.Get(v))
				}
				//单独处理caller字段 (供grpc服务使用，用于定位调用grpc的是哪个服务)
				md.Set("caller", global.AppConfig.Common.ServiceName)
				caller := r.Header.Get("caller")

				ctx = metadata.NewOutgoingContext(ctx, md)
				ctx = context.WithValue(ctx, "traceId", r.Header.Get("x-b3-traceid"))
				ctx = context.WithValue(ctx, "requestTime", time.Now())
				ctx = context.WithValue(ctx, "caller", caller)
			}

			//添加上下文信息
			ctx = context.WithValue(ctx, "httpRequest", r)

			return ctx
		})
}

func ClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if ip != "" {
		return ip
	}

	ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	if ip != "" {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(r.RemoteAddr)); err == nil {
		return ip
	}

	return ""
}

//提取请求参数信息，同时兼容raw和表单提交的形式
func GetRequestString(r *http.Request) []byte {
	var buf []byte
	if r.ParseForm() == nil && r.Form != nil && len(r.Form) > 0 {
		buf = urlValueToStruct(r.Form)
	} else {
		buf, _ = ioutil.ReadAll(r.Body)
		_ = r.Body.Close()
		r.Body = ioutil.NopCloser(bytes.NewBuffer(buf)) //重置，方便打印请求日志时继续读取
	}
	return buf
}

// http.Request.Form转为json字符串
func urlValueToStruct(v url.Values) []byte {
	if len(v) == 0 {
		return []byte{}
	}
	m := make(map[string]interface{})
	// 只取 []string 中第一个
	for k, v := range v {
		m[k] = v[0]
	}
	req, _ := json.Marshal(m)
	return req
}

// 格式化返回的response的消息内容
func EncodeHealthCheckResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json,charset=utf-8")
	type healthRes struct {
		Msg string `json:"msg"`
	}
	//统一返回http的数据
	httpresp := healthRes{
		Msg: "success",
	}

	return json.NewEncoder(w).Encode(httpresp)
}

// 格式化返回的response的消息内容
func EncodeHTTPV1CommonResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	w.Header().Set("Content-type", "application/json,charset=utf-8")

	//统一返回http的数据
	httpresp := api.HttpResponseBody{
		ApiStatus: 0,
		SysStatus: 0,
		Info:      "",
		Timestamp: int(time.Now().Unix()),
		Data:      nil,
	}

	resp := response.(api.EndpointResp)
	if resp.DataError.GetCode() == 0 {
		httpresp.Data = resp.RespData
	} else {
		httpresp.ApiStatus = resp.DataError.GetCode()
		message, ok := microerr.HttpErrMessages[resp.DataError.GetType()]
		if !ok {
			httpresp.ApiStatus = int(microerr.ServerCommonError)
			message = microerr.HttpErrMessages[microerr.ServerCommonError]
		}
		httpresp.Info = message
		if resp.RespData != nil {
			if httpresp.ApiStatus != int(microerr.RedirectLogin) {
				httpresp.Data = api.EmptyRespData{}
			} else {
				httpresp.Data = resp.RespData
			}
		} else {
			httpresp.Data = api.EmptyRespData{}
		}
	}

	//记录日志
	HttpLogging(ctx, resp.ReqData, httpresp, 200)

	return json.NewEncoder(w).Encode(httpresp)
}

func DecodeHttpIotCert(ctx context.Context, r *http.Request, reqInfo *api.IotCertReq) error {
	dataByte := GetRequestString(r)
	if len(dataByte) == 0 {
		return nil
	}
	err := json.Unmarshal(dataByte, reqInfo)
	if err != nil {
		sysErr := microerr.InvalidParameter.New("parseForm or body decode or url decode error").(*microerr.SysErr)
		deps.Logger.ErrorRaw(sysErr, utils.GetContextTraceId(ctx), constant.CommonHttpKey)

		err = &errorhttp.HttpErrorV1{
			Code:    sysErr.GetCode(),
			Message: microerr.HttpErrMessages[sysErr.GetType()],
		}
	}
	return err
}
