package users

import (
	"errors"

	"Server/internal/utils"

	"gorm.io/gorm"
)

type Service interface {
	GetByID(id uint) (*UserResponse, error)
	Register(req RegisterRequest) (*LoginResponse, error)
	Login(req LoginRequest) (*LoginResponse, error)
	Update(id uint, req UpdateUserRequest) (*UserResponse, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetByID(id uint) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *service) Register(req RegisterRequest) (*LoginResponse, error) {
	// Check if user already exists
	_, err := s.repo.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user already exists")
	}
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &User{
		Email:     req.Email,
		Password:  hashedPassword,
		FirstName: req.FirstName,
		LastName:  req.LastName,
		IsActive:  true,
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

func (s *service) Login(req LoginRequest) (*LoginResponse, error) {
	user, err := s.repo.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return nil, errors.New("invalid credentials")
	}

	if !user.IsActive {
		return nil, errors.New("account is inactive")
	}

	// Generate JWT
	token, err := utils.GenerateJWT(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		Token: token,
		User:  *s.toUserResponse(user),
	}, nil
}

func (s *service) Update(id uint, req UpdateUserRequest) (*UserResponse, error) {
	user, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Update fields if provided
	if req.FirstName != "" {
		user.FirstName = req.FirstName
	}
	if req.LastName != "" {
		user.LastName = req.LastName
	}
	if req.Email != "" {
		user.Email = req.Email
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}

	return s.toUserResponse(user), nil
}

func (s *service) toUserResponse(user *User) *UserResponse {
	return &UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		IsActive:  user.IsActive,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
