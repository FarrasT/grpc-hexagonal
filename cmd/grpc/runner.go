package main

import (
	"log"

	proto "grpc-hexa/api"
	grpcCtrl "grpc-hexa/internal/controller/grpc"
	"grpc-hexa/internal/core/config"
	"grpc-hexa/internal/core/server/grpc"
	"grpc-hexa/internal/core/service"
	infraConf "grpc-hexa/internal/infra/config"
	"grpc-hexa/internal/infra/repository"

	"go.uber.org/zap"
	googleGrpc "google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

func main() {
	// Initialize logger
	logger, _ := zap.NewProduction()
	undo := zap.ReplaceGlobals(logger)
	defer undo()

	// Initialize the database connection
	db, err := repository.NewDB(
		infraConf.DatabaseConfig{
			Driver:                  "mysql",
			Url:                     "root:password@tcp(127.0.0.1:3306)/grpc_hexa?charset=utf8mb4&parseTime=true&loc=UTC&tls=false&readTimeout=3s&writeTimeout=3s&timeout=3s&clientFoundRows=true",
			ConnMaxLifetimeInMinute: 3,
			MaxOpenConns:            10,
			MaxIdleConns:            1,
		},
	)
	if err != nil {
		log.Fatalf("failed to new database err=%s\n", err.Error())
	}

	// Create the UserRepository
	userRepo := repository.NewUserRepository(db)

	// Create the UserService
	userService := service.NewUserService(userRepo)

	// Create the UserController
	userController := grpcCtrl.NewUserController(userService)

	// Create the gRPC server
	grpcServer, err := grpc.NewGrpcServer(
		config.GrpcServerConfig{
			Port: 9090,
			KeepaliveParams: keepalive.ServerParameters{
				MaxConnectionIdle:     100,
				MaxConnectionAge:      7200,
				MaxConnectionAgeGrace: 60,
				Time:                  10,
				Timeout:               3,
			},
			KeepalivePolicy: keepalive.EnforcementPolicy{
				MinTime:             10,
				PermitWithoutStream: true,
			},
		},
	)
	if err != nil {
		log.Fatalf("failed to new grpc server err=%s\n", err.Error())
	}

	// Start the gRPC server
	go grpcServer.Start(
		func(server *googleGrpc.Server) {
			proto.RegisterUserServiceServer(server, userController)
		},
	)

	// Add shutdown hook to trigger closer resources of service
	grpc.AddShutdownHook(grpcServer, db)
}
