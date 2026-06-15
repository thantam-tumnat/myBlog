package entities

type BlogRepository interface {
	CreateBlog(blog *Blog) (*Blog, error)
	GetBlogs() ([]Blog, error)

	CreateUser(user *User) (*User, error)
	CheckUser(userId int) (uint, error)
	GetUser(userId uint) (*User, error)
}

type BlogService interface {
	BlogCreated(userId int, blogReq *BlogRequest) (*BlogRes, error)
	BlogGets() (*[]BlogRes, error)
}
