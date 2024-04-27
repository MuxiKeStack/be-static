//go:build wireinject

package main

import (
	"github.com/MuxiKeStack/be-static/grpc"
	"github.com/MuxiKeStack/be-static/ioc"
	"github.com/MuxiKeStack/be-static/pkg/grpcx"
	"github.com/MuxiKeStack/be-static/repository"
	"github.com/MuxiKeStack/be-static/repository/cache"
	"github.com/MuxiKeStack/be-static/repository/dao"
	"github.com/MuxiKeStack/be-static/service"
	"github.com/google/wire"
)

func InitGRPCServer() grpcx.Server {
	wire.Build(
		ioc.InitGRPCxKratosServer,
		grpc.NewStaticServiceServer,
		service.NewStaticService,
		repository.NewCachedStaticRepository,
		cache.NewRedisStaticCache,
		dao.NewMongoDBStaticDAO,
		// 第三方
		ioc.InitDB,
		ioc.InitRedis,
		ioc.InitLogger,
		ioc.InitEtcdClient,
	)
	return grpcx.Server(nil)
}
