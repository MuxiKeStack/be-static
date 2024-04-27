package grpc

import (
	"context"
	staticv1 "github.com/MuxiKeStack/be-api/gen/proto/static/v1"
	"github.com/MuxiKeStack/be-static/domain"
	"github.com/MuxiKeStack/be-static/service"
	"github.com/go-kratos/kratos/v2/transport/grpc"
)

type StaticServiceServer struct {
	staticv1.UnimplementedStaticServiceServer
	svc service.StaticService
}

func NewStaticServiceServer(svc service.StaticService) *StaticServiceServer {
	return &StaticServiceServer{svc: svc}
}

func (s *StaticServiceServer) GetStaticByName(ctx context.Context, request *staticv1.GetStaticByNameRequest) (*staticv1.GetStaticByNameResponse, error) {
	static, err := s.svc.GetStaticByName(ctx, request.GetName())
	return &staticv1.GetStaticByNameResponse{
		Static: &staticv1.Static{
			Name:    static.Name,
			Content: static.Content,
		},
	}, err
}

func (s *StaticServiceServer) SaveStatic(ctx context.Context, request *staticv1.SaveStaticRequest) (*staticv1.SaveStaticResponse, error) {
	err := s.svc.SaveStatic(ctx, domain.Static{
		Name:    request.GetStatic().GetName(),
		Content: request.GetStatic().GetContent(),
	})
	return &staticv1.SaveStaticResponse{}, err
}

func (s *StaticServiceServer) Register(server *grpc.Server) {
	staticv1.RegisterStaticServiceServer(server, s)
}
