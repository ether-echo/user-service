package adapter

import (
	"context"
	"github.com/ether-echo/user-service/internal/domain"
	"github.com/ether-echo/user-service/pkg/logger"

	up "github.com/ether-echo/protos/userProcessor"
)

var (
	log = logger.Logger().Named("grpc_server").Sugar()
)

type IRepository interface {
	ProcessSave(ctx context.Context, chatId int64, message string) error
	ProcessChangeAccessTaro(ctx context.Context, chatId int64) (bool, error)
	ProcessGetAllUsers(ctx context.Context) ([]domain.User, error)
	ProcessGetAllChatId(ctx context.Context) ([]int64, error)
}

type UserService struct {
	up.UnimplementedUserServiceServer
	Repository IRepository
}

func (u *UserService) SaveMessage(ctx context.Context, req *up.MessageRequest) (*up.MessageResponse, error) {
	err := u.Repository.ProcessSave(ctx, req.ChatId, req.Message)
	if err != nil {
		log.Error(err)
	}

	return &up.MessageResponse{
		Success: true,
	}, nil
}

func (u *UserService) SetTaro(ctx context.Context, req *up.SetTaroRequest) (*up.SetTaroResponse, error) {
	IsGotTaro, err := u.Repository.ProcessChangeAccessTaro(ctx, req.ChatId)
	if err != nil {
		log.Error(err)
	}

	return &up.SetTaroResponse{
		TaroIsGot: IsGotTaro,
		Success:   true,
	}, nil
}

func (u *UserService) GetAllID(ctx context.Context, _ *up.Empty) (*up.IdList, error) {
	usersId, err := u.Repository.ProcessGetAllChatId(ctx)
	if err != nil {
		log.Error(err)
	}

	return &up.IdList{
		Ids: usersId,
	}, nil
}

func (u *UserService) GetAllUsers(ctx context.Context, _ *up.Empty) (*up.UserList, error) {
	users, err := u.Repository.ProcessGetAllUsers(ctx)
	if err != nil {
		log.Error(err)
	}

	var userList []*up.User

	for _, user := range users {
		userList = append(userList, &up.User{
			ChatId:    user.ChatId,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Username:  user.Username,
		})
	}

	return &up.UserList{
		Users: userList,
	}, nil

}
