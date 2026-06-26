package controllers

import (
	"blogservice/modules/entities"
	"blogservice/modules/logs"
	"blogservice/modules/middlewares"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type blogHandler struct {
	blogService entities.BlogService
}

func NewBlogHandler(r fiber.Router, blogService entities.BlogService, authMiddleware fiber.Handler) {

	controller := &blogHandler{blogService: blogService}

	r.Post("/createBlog", authMiddleware, controller.CreateBlog)
	r.Get("/getBlogs", controller.GetBlogs)
}

func (h *blogHandler) CreateBlog(c *fiber.Ctx) error {
	// userId comes from the verified JWT, not the URL — callers can't forge it.
	userID, ok := c.Locals(middlewares.UserIDKey).(uint)
	if !ok || userID == 0 {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"message":    "unauthorized",
			"status":     fiber.ErrUnauthorized.Message,
			"statusCode": fiber.ErrUnauthorized.Code,
		})
	}
	userIdInt := int(userID)

	var blogRequestBody entities.BlogRequest
	if err := c.BodyParser(&blogRequestBody); err != nil {
		logs.Info("Invalid request to Create Blog :", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	blog, err := h.blogService.BlogCreated(userIdInt, &blogRequestBody)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":    "Can't create blog",
			"dev":        err.Error(),
			"status":     fiber.ErrBadRequest.Message,
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "blog successfully created.",
		"blogId":  blog.BlogId,
		"data":    blog,
		"status":  fiber.StatusOK,
	})
}

func (h *blogHandler) GetBlogs(c *fiber.Ctx) error {
	response, err := h.blogService.BlogGets()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":    "Can't get blogs",
			"dev": err,
			"status":     fiber.ErrBadRequest.Message,
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}
	return c.JSON(fiber.Map{
		"message": "blogs retrieved successfully",
		"status":  fiber.StatusOK,
		"data":    response,
	})
}
