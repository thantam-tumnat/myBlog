package repositories

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"userservice/modules/entities"
	"userservice/modules/logs"

	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

type userRepositoryRedis struct {
	client  *redis.Client
	db      *gorm.DB
	produce entities.EventProducer
}

func NewRepositoryRedis(client *redis.Client, produce entities.EventProducer, db *gorm.DB) entities.UserRepository {
	err := db.AutoMigrate(&entities.User{})
	if err != nil {
		logs.Error("Error AutoMigrate table user")
		return &userRepositoryRedis{}
	}
	return &userRepositoryRedis{
		client:  client,
		db:      db,
		produce: produce,
	}
}

func (r userRepositoryRedis) CreateUser(user *entities.User) (*entities.User, error) {
	result := r.db.Create(&user)
	if result.Error != nil {
		logs.Error("Failed to create User")
		return nil, result.Error
	}

	keys, err := r.client.Keys(context.Background(), "userrepository*").Result()
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

	//produce
	userPayload := entities.UserCreated{
		UserID:      int(user.ID),
		Name:        user.Name,
		Description: user.Description,
		UserImage:   user.UserImage,
	}
	errProduce := r.produce.Produce(&userPayload)
	if errProduce != nil {
		logs.Error(err)
		return nil, err
	}

	return user, nil
}

func (r userRepositoryRedis) GetUsers() ([]entities.User, error) {
	users := []entities.User{}
	result := r.db.Find(&users)
	if result.Error != nil {
		return nil, result.Error
	}

	key := "userrepository::GetUsers"

	// Check if data is cached in Redis
	productsJson, err := r.client.Get(context.Background(), key).Result()
	if err == nil {
		// Unmarshal cached data if available
		err = json.Unmarshal([]byte(productsJson), &users)
		if err == nil {
			logs.Info("Tags Retrieved From Redis")
			return users, nil
		}
	}

	// Cache the result in Redis
	data, err := json.Marshal(users)
	if err != nil {
		logs.Error(err)
		return users, err
	}

	err = r.client.Set(context.Background(), key, string(data), time.Second*10).Err()
	if err != nil {
		logs.Error(err)
		return users, err
	}

	return users, nil
}
