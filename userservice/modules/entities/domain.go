package entities

import (

	"github.com/gofiber/fiber/v2"
)

type UserRepository interface {
	CreateUser(user *User) (*User, error)
	GetUsers()([]User, error)
}

type UserService interface {
	UserCreated(userReq *UserRequest) (int, fiber.Map)
	UserGets()(int,fiber.Map)
}
