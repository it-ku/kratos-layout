// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"kratos-layout/internal/biz"
	"kratos-layout/internal/conf"
	"kratos-layout/internal/data"
	"kratos-layout/internal/server"
	"kratos-layout/internal/service"
)

import (
	_ "go.uber.org/automaxprocs"
)

// Injectors from wire.go:

// wireApp init kratos application.
func wireApp(env *conf.Env, confServer *conf.Server, confData *conf.Data, bootstrap *conf.Bootstrap, logger log.Logger) (*kratos.App, func(), error) {
	db := data.NewMysqlCmd(bootstrap, logger)
	client := data.NewRedisCmd(confData, logger)
	dataData, cleanup, err := data.NewData(bootstrap, db, client, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUseCase := biz.NewGreeterUseCase(greeterRepo, logger)
	greeterService := service.NewGreeterService(greeterUseCase)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	app := newApp(logger, grpcServer, httpServer)
	return app, func() {
		cleanup()
	}, nil
}