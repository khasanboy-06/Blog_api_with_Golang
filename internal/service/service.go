package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"Blog_project_with_Go/internal/database"
	"Blog_project_with_Go/internal/dto"
	"Blog_project_with_Go/internal/models"
	"Blog_project_with_Go/internal/repository"
	"github.com/redis/go-redis/v9"
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


type PaginatedPostsResponse struct{
	Posts []models.Post `json:"posts"`
	Total int64 `json:"total"`
	TotalPages int64 `json:"total_pages"`
}


func clearPostsCache() {
	ctx := context.Background()
	iter := database.RedisClient.Scan(ctx, 0, "posts:page:*", 0).Iterator()
	for iter.Next(ctx) {
		database.RedisClient.Del(ctx, iter.Val())
	}
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

	clearPostsCache()

	return &post, nil
}

func (s *postService) GetPosts(page, limit int) ([]models.Post, int64, int64, error) {
	if page < 1 {
		page = 1
	}
	if limit < 10 {
		limit = 10
	}
	
	cacheKey := fmt.Sprintf("posts:page:%d:limit:%d", page, limit)
	ctx := context.Background()

	cachedData, err := database.RedisClient.Get(ctx, cacheKey).Result()
	if err == nil {
		var res PaginatedPostsResponse
		if jsonErr := json.Unmarshal([]byte(cachedData), &res); jsonErr == nil {
			return res.Posts, res.Total, res.TotalPages, nil
		}
	} else if err != redis.Nil {
		fmt.Printf("Redis error: %v\n", err)
	}

	posts, total, err := s.repo.GetAll(page, limit)
	if err != nil {
		return nil, 0, 0, err
	}

	totalPages := (total + int64(limit) - 1) / int64(limit)

	responseObj := PaginatedPostsResponse{
		Posts:      posts,
		Total:      total,
		TotalPages: totalPages,
	}
	bytes, err := json.Marshal(responseObj)
	if err == nil {
		database.RedisClient.Set(ctx, cacheKey, bytes, 10*time.Minute)
	}

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

	clearPostsCache()

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

	clearPostsCache()

	return s.repo.Delete(post)
}