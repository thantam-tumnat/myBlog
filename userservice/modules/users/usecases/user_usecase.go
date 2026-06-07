package usecases

import (
	"fmt"
	"userservice/configs"
	"userservice/modules/entities"
	"userservice/modules/logs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type userService struct {
	useRepo entities.UserRepository
	cfg     *configs.Configs
}

func NewUserService(useRepo entities.UserRepository, cfg *configs.Configs) entities.UserService {
	return &userService{
		useRepo: useRepo,
		cfg:     cfg,
	}
}

func (s *userService) UserCreated(user *entities.UserRequest) (int, fiber.Map) {
	userCreate := entities.User{
		Username:    user.Username,
		Password:    user.Password,
		Name:        user.Name,
		Description: user.Description,
		UserImage:   user.UserImage,
	}

	created, err := s.useRepo.CreateUser(&userCreate)
	if err != nil {
		logs.Info("Can't create user", zap.String("error", err.Error()))

		return fiber.ErrBadRequest.Code, fiber.Map{
			"message": fmt.Sprintf("BadRequest : %v", err),
		}
	}

	userIdResponse := created.ID

	return fiber.StatusOK, fiber.Map{
		"message":    "user successfully created.",
		"TaskId":     userIdResponse,
		"data":       created,
		"status":     "OK",
		"statusCode": 200,
	}

}

func (s *userService) UserGets() (int, fiber.Map) {
	users, err := s.useRepo.GetUsers()
	if err != nil {
		logs.Info("Can't get user", zap.String("error", err.Error()))

		return fiber.ErrBadRequest.Code, fiber.Map{
			"message": fmt.Sprintf("BadRequest : %v", err),
		}
	}
	usersResponse := []entities.UserRes{}
	for _, item := range users{
		Response := entities.UserRes{
			UserID: int(item.ID),
			Username: item.Username,
			Password: item.Password,
			Name: item.Name,
			Description: item.Description,
			UserImage: item.UserImage,
		}
		usersResponse = append(usersResponse, Response)
	}

	return fiber.StatusOK, fiber.Map{
		"message":    "Get User successfully",
		"data":       usersResponse,
		"status":     "OK",
		"statusCode": 200,
	}
}
