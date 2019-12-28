package usecase

import (
	"2019_2_Covenant/pkg/models"
	"2019_2_Covenant/pkg/session/repository"
	"context"
	"github.com/golang/protobuf/ptypes"
)

type SessionUsecase struct {
	sessionRepo repository.SessionsClient
}

func toModel(session *repository.Session) *models.Session {
	expires, _ := ptypes.Timestamp(session.Expires)
	return &models.Session{
		ID: session.Id,
		UserID: session.UserId,
		Expires: expires,
		Data: session.Data,
	}
}

func fromModel(session *models.Session) *repository.Session {
	expires, _ := ptypes.TimestampProto(session.Expires)
	return &repository.Session{
		Id: session.ID,
		UserId: session.UserID,
		Expires: expires,
		Data: session.Data,
	}
}

func NewSessionUsecase(sr repository.SessionsClient) *SessionUsecase {
	return &SessionUsecase{
		sessionRepo: sr,
	}
}

func (sUC *SessionUsecase) Get(ctx context.Context, value string) (*models.Session, error) {
	sess, err := sUC.sessionRepo.Get(ctx, &repository.GetRequest{Value:value})

	if err != nil {
		return nil, err
	}

	return toModel(sess), nil
}

func (sUC *SessionUsecase) Store(ctx context.Context, newSession *models.Session) error {
	session, err := sUC.sessionRepo.Store(ctx, fromModel(newSession))
	if err != nil {
		return err
	}

	*newSession = *toModel(session)
	return nil
}

func (sUC *SessionUsecase) DeleteByID(ctx context.Context, id uint64) error {
	_, err := sUC.sessionRepo.DeleteByID(ctx, &repository.DeleteByIDRequest{Id:id})
	return err
}
