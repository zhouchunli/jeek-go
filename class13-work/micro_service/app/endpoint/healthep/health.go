package healthep

import (
	"context"
	"micro_service/app/model/api"
	"micro_service/pkg/microerr"

	"github.com/go-kit/kit/endpoint"
)

func MakeServerHealthEndpoint() endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		//var reqC = request.(*api.RequestContext)
		respEndpoint := api.EndpointResp{
			DataError: nil,
			ReqData:   request,
			RespData:  nil,
		}

		resp := api.ServerHealthResp{}
		err = microerr.HttpSuccess.New("ok")
		sysErr, ok := err.(*microerr.SysErr)
		if !ok {
			respEndpoint.DataError = microerr.TypeAssertError.New("response error assert fail").(*microerr.SysErr)
		} else {
			respEndpoint.DataError = sysErr
			respEndpoint.RespData = resp
		}

		return respEndpoint, nil
	}
}
