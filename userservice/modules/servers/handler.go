package servers

import (
	"github.com/gofiber/fiber/v2"
	_userController "userservice/modules/users/controllers"
	_userRepo "userservice/modules/users/repositories"
	_userUse "userservice/modules/users/usecases"

	// _userHandlerConsumer "userservice/modules/consumer/handler"
	_userHandlerProducer "userservice/modules/producer/handler"


)

func (s *server) Handlers() error {

	v1 := s.App.Group("/v1")

	userGroup := v1.Group("/myblogs")

	//kafka produce
	userProduceHandler := _userHandlerProducer.NewEventProducer(s.SyncProducer)

	userRepository := _userRepo.NewRepositoryRedis(s.Redis, userProduceHandler, s.Db)
	userUsecase := _userUse.NewUserService(userRepository,s.Cfg)

	_userController.NewUserHandler(userGroup, userUsecase)
	//
	// End point not found response
	s.App.Use(func(c *fiber.Ctx) error {
		return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
			"status":      fiber.ErrInternalServerError.Message,
			"status_code": fiber.ErrInternalServerError.Code,
			"message":     "error, end point not found",
			"result":      nil,
		})
	})

	return nil
}
