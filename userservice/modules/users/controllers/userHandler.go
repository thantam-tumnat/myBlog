package controllers

import (
	"userservice/modules/entities"
	"userservice/modules/logs"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type userHandler struct {
	userService entities.UserService
}

func NewUserHandler(r fiber.Router, userService entities.UserService) {

	controller := &userHandler{userService: userService}

	r.Post("/register", controller.CreateUser)
	r.Post("/login", controller.Login)
	r.Get("/getUsers", controller.GetUsers)
}

func (h *userHandler) Login(c *fiber.Ctx) error {
	var loginRequestBody entities.LoginRequest
	if err := c.BodyParser(&loginRequestBody); err != nil {
		logs.Info("Invalid request to Login :", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	codeStatus, response := h.userService.Login(&loginRequestBody)
	return c.Status(codeStatus).JSON(response)
}

func (h *userHandler) CreateUser(c *fiber.Ctx) error {
	var userRequestBody entities.UserRequest
	if err := c.BodyParser(&userRequestBody); err != nil {
		logs.Info("Invalid request to Create User :", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	codeStatus, response :=  h.userService.UserCreated(&userRequestBody)
	return c.Status(codeStatus).JSON(response)
}

func (h *userHandler) GetUsers(c *fiber.Ctx) error {
	codeStatus, response := h.userService.UserGets()
	return c.Status(codeStatus).JSON(response)
}
