package file_processor

import (
	"2019_2_Covenant/pkg/logger"
	"database/sql"
	"google.golang.org/grpc"
	"net"
)

type FileServer struct {
	grpc     *grpc.Server
	database *sql.DB
	logger   *logger.LogrusLogger
	rootDir  string
}

func NewFileServer(rootDir string, database *sql.DB) *FileServer {
	return &FileServer{
		grpc:     grpc.NewServer(),
		database: database,
		logger:   logger.NewLogrusLogger(),
		rootDir:  rootDir,
	}
}

func (server *FileServer) Start(address string) error {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}

	fileProcessor := NewFileProcessor(server.rootDir, server.database)
	RegisterFilesServer(server.grpc, fileProcessor)

	return server.grpc.Serve(lis)
}

func (server *FileServer) Stop() {
	server.grpc.Stop()
}
