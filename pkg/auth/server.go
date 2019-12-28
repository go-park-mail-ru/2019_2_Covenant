package auth

import (
	"2019_2_Covenant/pkg/logger"
	session "2019_2_Covenant/pkg/session/repository"
	user "2019_2_Covenant/pkg/user/repository"
	"database/sql"
	"google.golang.org/grpc"
	"net"
)

type AuthServer struct {
	grpc     *grpc.Server
	database *sql.DB
	logger   *logger.LogrusLogger
}

func NewAuthServer(database *sql.DB) *AuthServer {
	return &AuthServer{
		grpc:     grpc.NewServer(),
		database: database,
		logger:   logger.NewLogrusLogger(),
	}
}

func (server *AuthServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	userRepository := user.NewUserRepository(server.database)
	sessionRepository := session.NewSessionRepository(server.database)

	user.RegisterUsersServer(server.grpc, userRepository)
	session.RegisterSessionsServer(server.grpc, sessionRepository)

	return server.grpc.Serve(lis)
}

func (server *AuthServer) Stop() {
	server.grpc.Stop()
}
