package healthhttp

import (
	"context"
	"micro_service/app/endpoint/healthep"
	"micro_service/app/model/api"
	"micro_service/app/transport/http/commonhttp"
	"net/http"

	"github.com/go-kit/kit/endpoint"

	httptransport "github.com/go-kit/kit/transport/http"
)

func MakeServerHealthHTTPTransport() http.Handler {
	return httptransport.NewServer(
		MakeWrappedServerHealthEndpoint(),
		DecodeHTTPRequestServerHealthV1,
		commonhttp.EncodeHealthCheckResponse,
		commonhttp.AddHTTPRequestToContext(),
	)
}

func MakeWrappedServerHealthEndpoint() endpoint.Endpoint {
	ep := healthep.MakeServerHealthEndpoint()
	//ep = middleware.NotTokenSignWrapV1(ep)
	return ep
}

func DecodeHTTPRequestServerHealthV1(ctx context.Context, r *http.Request) (interface{}, error) {
	var (
		requestCtx = api.ServerHealthResp{}
		//reqInfo    api.ServerHealthReq
		err error
	)
	//requestCtx.ReqInfo = reqInfo
	return &requestCtx, err
}
