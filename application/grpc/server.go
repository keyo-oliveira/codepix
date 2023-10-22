package grpc

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/keyo-oliveira/codepix/application/grpc/pb"
	"github.com/keyo-oliveira/codepix/application/usecases"
	"github.com/keyo-oliveira/codepix/infra/repository"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func StartGrpcServer(database *gorm.DB, port int) {
	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)

	pixRepository := repository.PixKeyRepositoryDb{Db: database}
	pixUseCase := usecases.PixUseCase{PixKeyRepository: pixRepository}
	PixGrpcService := NewPixGrpcService(pixUseCase)
	pb.RegisterPixServiceServer(grpcServer, PixGrpcService)

	address := fmt.Sprintf("0.0.0.0:%d", port)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal("cannot start Grpc server", err)
	}
	log.Printf("gRPC server has been started on port %d", port)

	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal("cannot start grpc server", err)
	}
}
