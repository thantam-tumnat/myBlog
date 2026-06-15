package usecases

import (
	"blogservice/configs"
	"blogservice/modules/entities"
	"blogservice/modules/logs"

	"go.uber.org/zap"
)

type blogService struct {
	blogRepo entities.BlogRepository
	cfg      *configs.Configs
}

func NewUserService(blogRepo entities.BlogRepository, cfg *configs.Configs) entities.BlogService {
	return &blogService{
		blogRepo: blogRepo,
		cfg:      cfg,
	}
}

func (s *blogService) BlogCreated(userId int, blogReq *entities.BlogRequest) (*entities.BlogRes, error) {
	check, err := s.blogRepo.CheckUser(userId)
	if err != nil {
		logs.Info("There are no users in the system.", zap.String("error", err.Error()))
		return nil, err
	}
	user, err := s.blogRepo.GetUser(check)
	if err != nil {
		logs.Info("Can't get user.", zap.String("error", err.Error()))
		return nil, err
	}
	blog := entities.Blog{
		Title:      blogReq.Title,
		BlogDesc:   blogReq.BlogDesc,
		Content:    blogReq.Content,
		CoverImage: blogReq.CoverImage,
		UserId:     user.ID,
		UserName:   user.Name,
		UserDesc:   user.Description,
		UserImage:  user.UserImage,
	}

	created, err := s.blogRepo.CreateBlog(&blog)
	if err != nil {
		logs.Info("Can't create blog", zap.String("error", err.Error()))
		return nil, err
	}

	return &entities.BlogRes{
		BlogId:     created.ID,
		Title:      created.Title,
		BlogDesc:   created.BlogDesc,
		Content:    created.Content,
		CoverImage: created.CoverImage,
		UserName:   created.UserName,
		UserDesc:   created.UserDesc,
		UserImage:  created.UserImage,
	}, nil
}

func (s *blogService) BlogGets() (*[]entities.BlogRes, error) {
	blogs, err := s.blogRepo.GetBlogs()
	if err != nil {
		logs.Info("Can't get blogs", zap.String("error", err.Error()))

		return nil, err
	}
	blogsResponse := []entities.BlogRes{}
	for _, item := range blogs {
		Response := entities.BlogRes{
			BlogId:     item.ID,
			Title:      item.Title,
			BlogDesc:   item.BlogDesc,
			Content:    item.Content,
			CoverImage: item.CoverImage,
			UserName:   item.UserName,
			UserDesc:   item.UserDesc,
			UserImage:  item.UserImage,
		}
		blogsResponse = append(blogsResponse, Response)
	}
	return &blogsResponse,nil
}
