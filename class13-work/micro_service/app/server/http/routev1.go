package http

import (
	"micro_service/app/transport/http/healthhttp"
	"github.com/julienschmidt/httprouter"
)

func route(router *httprouter.Router) {
	router.Handler(Get, "/healthCheck", healthhttp.MakeServerHealthHTTPTransport())
}
