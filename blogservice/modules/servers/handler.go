package servers

import (
	"github.com/gofiber/fiber/v2"
	_blogController "blogservice/modules/blogs/controllers"
	_blogRepo "blogservice/modules/blogs/repositories"
	_blogUse "blogservice/modules/blogs/usecases"

	_blogHandlerConsumer "blogservice/modules/consumer/handler"
	_blogUsecaseConsumer "blogservice/modules/consumer/usecases"

)

func (s *server) Handlers() error {

	v1 := s.App.Group("/v1")

	blogGroup := v1.Group("/myblogs")

	blogRepository := _blogRepo.NewRepositoryRedis(s.Redis,s.Db)
	blogUsecases := _blogUse.NewUserService(blogRepository,s.Cfg)

	//kafka consume
	blogUsecasesConsumer := _blogUsecaseConsumer.NewBlogEventHanlder(blogUsecases,blogRepository)
	blogConsumerHandler := _blogHandlerConsumer.NewConsumerHandler(blogUsecasesConsumer)
	s.consumerGroupHandler = blogConsumerHandler

	_blogController.NewBlogHandler(blogGroup, blogUsecases)

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
