package entities

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Username    string `db:"username" json:"username"`
	Password    string `db:"password" json:"password"`
	Name        string `db:"name" json:"name"`
	Description string `db:"description" json:"description"`
	UserImage   string `db:"user_image" json:"userImage"`
}

type Blog struct {
	gorm.Model
	Title       string `db:"title" json:"title"`
	BlogDesc string `db:"blog_desc" json:"blogDesc"`
	Content     string `db:"content" json:"content"`
	CoverImage  string `db:"cover_image" json:"coverImage"`
	UserId      uint `db:"user_id" json:"userId"`
	UserName    string `db:"user_name" json:"userName"`
	UserDesc    string `db:"user_desc" json:"userDesc"`
	UserImage   string `db:"user_image" json:"userImage"`
}

type BlogRes struct {
	BlogId     uint `json:"blogId"`
	Title      string `json:"title"`
	BlogDesc   string `json:"blogDesc"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
	UserName   string `json:"userName"`
	UserDesc   string `json:"userDesc"`
	UserImage  string `json:"userImage"`
}

type BlogRequest struct {
	Title      string `json:"title"`
	BlogDesc   string `json:"blogDesc"`
	Content    string `json:"content"`
	CoverImage string `json:"coverImage"`
}
