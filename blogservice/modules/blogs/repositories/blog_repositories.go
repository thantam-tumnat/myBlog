package repositories

import (
	"blogservice/modules/entities"
	"blogservice/modules/logs"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type blogRepositoryRedis struct {
	client *redis.Client
	db     *gorm.DB
}

func NewRepositoryRedis(client *redis.Client, db *gorm.DB) entities.BlogRepository {
	err := db.Table("blog_users").AutoMigrate(&entities.User{})
	if err != nil {
		logs.Error("Error AutoMigrate table blog_users")
		return &blogRepositoryRedis{}
	}
	err = db.AutoMigrate(&entities.Blog{})
	if err != nil {
		logs.Error("Error AutoMigrate table blogs")
		return &blogRepositoryRedis{}
	}
	return &blogRepositoryRedis{
		client: client,
		db:     db,
	}
}

func (r blogRepositoryRedis) CreateUser(user *entities.User) (*entities.User, error) {
	result := r.db.Table("blog_users").Create(&user)
	if result.Error != nil {
		logs.Error("Failed to create User")
		return nil, result.Error
	}

	keys, err := r.client.Keys(context.Background(), "blogrepository*").Result()
	if err != nil {
		logs.Error(fmt.Sprintf("Error retrieving keys: %v", err))
		return nil, err
	}
	for _, key := range keys {
		err := r.client.Del(context.Background(), key).Err()
		if err != nil {
			logs.Error(fmt.Sprintf("Error deleting key: %v", err))
		} else {
			logs.Info(fmt.Sprintf("Deleted key: %v", key))
		}
	}

	return user, nil
}

func (r blogRepositoryRedis) CheckUser(userId int) (uint, error) {
	userUint := uint(userId)
	user := entities.User{}
	err := r.db.Table("blog_users").Where("id=?", userUint).First(&user).Error
	if err != nil {
		logs.Error(fmt.Sprintf("Function CheckUser has error with user_id : %v", userId))
		return userUint, err
	}
	logs.Debug(fmt.Sprintf("Check User Data : %v", userId))
	return userUint, nil
}

func (r blogRepositoryRedis) GetUser(userId uint) (*entities.User, error) {
	user := entities.User{}
	result := r.db.Table("blog_users").Where("id=?", userId).Find(&user)
	if result.Error != nil {
		logs.Error("Failed to get User")
		return nil, result.Error
	}
	logs.Debug(fmt.Sprintf("User Data : %v", user))
	return &user, nil
}

func (r blogRepositoryRedis) CreateBlog(blog *entities.Blog) (*entities.Blog, error) {

	result := r.db.Create(&blog)
	if result.Error != nil {
		logs.Error("Failed to create Blog")
		return nil, result.Error
	}

	keys, err := r.client.Keys(context.Background(), "blogrepository*").Result()
	if err != nil {
		logs.Error(fmt.Sprintf("Error retrieving keys: %v", err))
		return nil, err
	}
	for _, key := range keys {
		err := r.client.Del(context.Background(), key).Err()
		if err != nil {
			logs.Error(fmt.Sprintf("Error deleting key: %v", err))
		} else {
			logs.Info(fmt.Sprintf("Deleted key: %v", key))
		}
	}
	return blog, nil
}

func (r blogRepositoryRedis) GetBlogs() ([]entities.Blog, error) {
	blogs := []entities.Blog{}
	result := r.db.Find(&blogs)
	if result.Error != nil {
		return nil, result.Error
	}

	key := "blogrepository::Getblogs"

	// Check if data is cached in Redis
	productsJson, err := r.client.Get(context.Background(), key).Result()
	if err == nil {
		// Unmarshal cached data if available
		err = json.Unmarshal([]byte(productsJson), &blogs)
		if err == nil {
			logs.Info("blogs Retrieved From Redis")
			return blogs, nil
		}
	}

	// Cache the result in Redis
	data, err := json.Marshal(blogs)
	if err != nil {
		logs.Error(err)
		return blogs, err
	}

	err = r.client.Set(context.Background(), key, string(data), time.Second*10).Err()
	if err != nil {
		logs.Error(err)
		return blogs, err
	}

	return blogs, nil
}
