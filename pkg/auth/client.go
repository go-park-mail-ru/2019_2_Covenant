package auth

import (
	sessions "2019_2_Covenant/pkg/session/repository"
	users "2019_2_Covenant/pkg/user/repository"
	"google.golang.org/grpc"
)

type AuthClient struct {
	usersClient    users.UsersClient
	sessionsClient sessions.SessionsClient
}

func NewAuthClient(conn *grpc.ClientConn) *AuthClient {
	return &AuthClient{
		usersClient:    users.NewUsersClient(conn),
		sessionsClient: sessions.NewSessionsClient(conn),
	}
}

func (client *AuthClient) User() users.UsersClient {
	return client.usersClient
}

func (client *AuthClient) Session() sessions.SessionsClient {
	return client.sessionsClient
}
