package service

import (
	"context"
	"github.com/ether-echo/user-service/internal/domain"
)

type IRepository interface {
	RegisterUser(user *domain.User) error
	IsUserRegistered(chatId int64) (bool, error)
	SaveMessage(ctx context.Context, chatID int64, message string) error
	GetTaro(ctx context.Context, chatID int64) (bool, error)
	ChangeAccessTaro(ctx context.Context, chatID int64) error
}

type IRpc interface {
	StartMessage(chatID int64, firstName string, exist bool) error
}

type Service struct {
	repository IRepository
	rpc        IRpc
}

func NewService(repository IRepository, rpc IRpc) *Service {
	return &Service{
		repository: repository,
		rpc:        rpc,
	}
}

func (s *Service) ProcessStart(user *domain.User) error {
	exist, err := s.repository.IsUserRegistered(user.ChatId)
	if err != nil {
		return err
	}

	if !exist {
		err = s.repository.RegisterUser(user)
		if err != nil {
			return err
		}
	}

	return s.rpc.StartMessage(user.ChatId, user.FirstName, exist)
}

func (s *Service) ProcessSave(ctx context.Context, chatId int64, message string) error {
	return s.repository.SaveMessage(ctx, chatId, message)
}

func (s *Service) ProcessChangeAccessTaro(ctx context.Context, chatId int64) (bool, error) {
	IsGotTaro, err := s.repository.GetTaro(ctx, chatId)
	if err != nil {
		return false, err
	}

	if !IsGotTaro {
		err = s.repository.ChangeAccessTaro(ctx, chatId)
		return IsGotTaro, err
	}

	return IsGotTaro, nil
}
