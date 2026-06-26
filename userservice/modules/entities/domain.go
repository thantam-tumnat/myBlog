package entities

import (

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	CreateUser(user *User) (*User, error)
	GetUsers()([]User, error)
	GetUserByUsername(username string) (*User, error)
}

type UserService interface {
	UserCreated(userReq *UserRequest) (int, fiber.Map)
	UserGets()(int,fiber.Map)
	Login(req *LoginRequest) (int, fiber.Map)
}
