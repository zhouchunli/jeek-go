// +build wireinject
// The build tag makes sure the stub is not built in the final build.

package di

import (
	"jeek-c4/internal/dao"
	"jeek-c4/internal/server/grpc"
	"jeek-c4/internal/server/http"
	"jeek-c4/internal/service"

	"github.com/google/wire"
)

//go:generate kratos t wire
func InitApp() (*App, func(), error) {
	panic(wire.Build(dao.Provider, service.Provider, http.New, grpc.New, NewApp))
}
