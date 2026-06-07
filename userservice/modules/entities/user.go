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

type UserRes struct {
	UserID      int    `json:"userId"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserImage   string `json:"userImage"`
}

type UserRequest struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	Name        string `json:"name"`
	Description string `json:"description"`
	UserImage   string `json:"userImage"`
}
