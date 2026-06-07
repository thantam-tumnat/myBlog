package usecase

import (
	"blogservice/modules/entities"
	"blogservice/modules/logs"
	"encoding/json"
	"fmt"

	"gorm.io/gorm"
)

type BlogEventHandler struct {
	blogSrv       entities.BlogService
	blogRepo      entities.BlogRepository
}

func NewBlogEventHanlder(blogSrv entities.BlogService, blogRepo entities.BlogRepository) *BlogEventHandler {
	return &BlogEventHandler{
		blogSrv:       blogSrv,
		blogRepo:      blogRepo,
	}
}

func (c *BlogEventHandler) Handle(topic string, message []byte) error {
	switch topic {
	case entities.UserCreated{}.TopicName():
		logs.Info(fmt.Sprintf("Received Topic: %s", string(topic)))
		logs.Info(fmt.Sprintf("Received message: %s", string(message)))
		event := &entities.UserCreated{}
		err := json.Unmarshal(message, event)
		if err != nil {
			logs.Error(fmt.Sprintf("Consume %v has error : %v", event.TopicName(), err))
			return err
		}
		data := entities.User{
			Model: gorm.Model{ID: uint(event.UserID)},
			Name: event.Name,
			Description: event.Description,
			UserImage: event.UserImage,
		}
		_ , err = c.blogRepo.CreateUser(&data)
		if err != nil {
			logs.Error(err)
			return err
		}
		logs.Debug(fmt.Sprintf("[%v] Created event : %v", topic, event.TopicName()))
		logs.Debug(fmt.Sprintf("Data : %v", string(message)))

	}
	return nil
}
