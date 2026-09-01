package repository

import (
	"Blog_project_with_Go/internal/models"
	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post *models.Post) error
	GetAll(page, limit int) ([]models.Post, int64, error)
	GetByID(id string) (*models.Post, error)
	Update(post *models.Post) error
	Delete(post *models.Post) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db: db}
}

func (r *postRepository) Create(post *models.Post) error {
	return r.db.Create(post).Error
}

func (r *postRepository) GetAll(page, limit int) ([]models.Post, int64, error) {
	var posts []models.Post
	var total int64

	r.db.Model(&models.Post{}).Count(&total)

	offset := (page - 1) * limit
	err := r.db.
		Limit(limit).
		Offset(offset).
		Order("created_at desc").
		Find(&posts).Error

	return posts, total, err
}

func (r *postRepository) GetByID(id string) (*models.Post, error) {
	var post models.Post
	err := r.db.First(&post, id).Error
	if err != nil {
		return nil, err
	}
	return &post, nil
}

func (r *postRepository) Update(post *models.Post) error {
	return r.db.Save(post).Error
}

func (r *postRepository) Delete(post *models.Post) error {
	return r.db.Delete(post).Error
}