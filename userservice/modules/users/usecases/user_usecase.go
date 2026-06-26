package usecases

import (
	"fmt"
	"userservice/configs"
	"userservice/modules/entities"
	"userservice/modules/logs"
	"userservice/pkg/utils"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
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
	if user.Username == "" || user.Password == "" {
		return fiber.ErrBadRequest.Code, fiber.Map{
			"message": "username and password are required",
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		logs.Error("Can't hash password", zap.String("error", err.Error()))
		return fiber.ErrInternalServerError.Code, fiber.Map{
			"message": "Internal server error",
		}
	}

	userCreate := entities.User{
		Username:    user.Username,
		Password:    string(hashedPassword),
		Role:        "USER",
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

func (s *userService) Login(req *entities.LoginRequest) (int, fiber.Map) {
	if req.Username == "" || req.Password == "" {
		return fiber.ErrBadRequest.Code, fiber.Map{
			"message": "username and password are required",
		}
	}

	user, err := s.useRepo.GetUserByUsername(req.Username)
	if err != nil {
		// Same response for unknown user vs wrong password (avoid user enumeration).
		logs.Info("Login failed: user not found", zap.String("username", req.Username))
		return fiber.ErrUnauthorized.Code, fiber.Map{
			"message": "invalid username or password",
		}
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		logs.Info("Login failed: wrong password", zap.String("username", req.Username))
		return fiber.ErrUnauthorized.Code, fiber.Map{
			"message": "invalid username or password",
		}
	}

	token, err := utils.GenerateAccessToken(s.cfg.Jwt.Secret, user.ID, user.Username, user.Role)
	if err != nil {
		logs.Error("Can't generate token", zap.String("error", err.Error()))
		return fiber.ErrInternalServerError.Code, fiber.Map{
			"message": "Internal server error",
		}
	}

	return fiber.StatusOK, fiber.Map{
		"message":     "login successful",
		"accessToken": token,
		"tokenType":   "Bearer",
		"status":      "OK",
		"statusCode":  200,
	}
}
