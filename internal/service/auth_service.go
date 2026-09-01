package service

import (
	"errors"

	"Blog_project_with_Go/internal/dto"
	"Blog_project_with_Go/internal/pkg/utils"
	"Blog_project_with_Go/internal/models"
	"Blog_project_with_Go/internal/repository"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.AuthResponse, error)
	Login(req dto.LoginRequest) (*dto.AuthResponse, error)
}


type authService struct {
	userRepo repository.UserRepository
}


func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &authService{userRepo: userRepo}
}


func (s *authService) Register(req dto.RegisterRequest) (*dto.AuthResponse, error) {
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	err = s.userRepo.Create(&user)
	if err != nil {
		return nil, errors.New("bu email yoki username allaqachon mavjud")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{Token: token}, nil
}


func (s *authService) Login(req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := s.userRepo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email yoki parol noto'g'ri")
	}

	if !utils.CheckPassword(req.Password, user.Password) {
		return nil, errors.New("email yoki parol noto'g'ri")
	}

	token, err := utils.GenerateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	return &dto.AuthResponse{Token: token}, nil
}