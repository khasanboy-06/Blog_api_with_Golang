package service

import (
	"errors"
	"Blog_project_with_Go/internal/dto"
	"Blog_project_with_Go/internal/models"
	"Blog_project_with_Go/internal/repository"
)

type PostService interface {
	CreatePost(userID uint, req dto.CreatePostRequest) (*models.Post, error)
	GetPosts(page, limit int) ([]models.Post, int64, int64, error)
	GetPost(id string) (*models.Post, error)
	UpdatePost(id string, userID uint, req dto.UpdatePostRequest) (*models.Post, error)
	DeletePost(id string, userID uint) error
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) CreatePost(userID uint, req dto.CreatePostRequest) (*models.Post, error) {
	post := models.Post{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	err := s.repo.Create(&post)
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (s *postService) GetPosts(page, limit int) ([]models.Post, int64, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 10 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	posts, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)
	return posts, total, totalPages, nil
}

func (s *postService) GetPost(id string) (*models.Post, error) {
	return s.repo.GetByID(id)
}

func (s *postService) UpdatePost(id string, userID uint, req dto.UpdatePostRequest) (*models.Post, error) {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return nil, errors.New("Post topilmadi")
	}

	if post.UserID != userID {
		return nil, errors.New("bu postni o'zgartirishga huquqingiz yo'q")
	}

	if req.Title != "" {
		post.Title = req.Title
	}
	if req.Content != "" {
		post.Content = req.Content
	}

	err = s.repo.Update(post)
	if err != nil {
		return nil, err
	}

	return post, nil
}

func (s *postService) DeletePost(id string, userID uint) error {
	post, err := s.repo.GetByID(id)
	if err != nil {
		return errors.New("Post topilmadi")
	}

	if post.UserID != userID {
		return errors.New("bu postni o'chirishga huquqingiz yo'q")
	}

	return s.repo.Delete(post)
}