package controllers

import (
	"blogservice/modules/entities"
	"blogservice/modules/logs"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type blogHandler struct {
	blogService entities.BlogService
}

func NewBlogHandler(r fiber.Router, blogService entities.BlogService) {

	controller := &blogHandler{blogService: blogService}

	r.Post("/createBlog/:userId", controller.CreateBlog)
	r.Get("/getBlogs", controller.GetBlogs)
}

func (h *blogHandler) CreateBlog(c *fiber.Ctx) error {
	//param
	userID := c.Params("userId")

	if userID == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":    "userId could't empty",
			"status":     fiber.ErrBadRequest.Message,
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	userIdInt, err := strconv.Atoi(userID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message":    "Can't convert string to int",
			"status":     fiber.ErrBadRequest.Message,
			"statusCode": fiber.ErrBadRequest.Code,
		})
	}

	var blogRequestBody entities.BlogRequest
	if err := c.BodyParser(&blogRequestBody); err != nil {
		logs.Info("Invalid request to Create Blog :", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Invalid request",
		})
	}

	codeStatus, response := h.blogService.BlogCreated(userIdInt, &blogRequestBody)
	return c.Status(codeStatus).JSON(response)
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
