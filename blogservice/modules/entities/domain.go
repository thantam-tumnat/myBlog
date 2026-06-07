package entities

import (
	"github.com/gofiber/fiber/v2"
)

type BlogRepository interface {
	CreateBlog(blog *Blog) (*Blog, error)
	GetBlogs() ([]Blog, error)

	CreateUser(user *User) (*User, error)
	CheckUser(userId int) (uint, error)
	GetUser(userId uint) (*User, error)
}

type BlogService interface {
	BlogCreated(userId int, blogReq *BlogRequest) (int, fiber.Map)
	BlogGets() (*[]BlogRes, error)
}
