package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/Guldana11/gophermart/database"
	"github.com/Guldana11/gophermart/models"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo database.UserRepository
}

func NewUserService(repo database.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(ctx context.Context, login, password string) (*models.User, error) {
	if login == "" || password == "" {
		return nil, errors.New("login and password required")
	}
	user, err := s.repo.CreateUser(ctx, login, password)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) Authenticate(ctx context.Context, login, password string) (*models.User, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, errors.New("invalid login/password")
	}
	return user, nil
}

func (s *UserService) Login(ctx context.Context, login, password string) (*models.User, error) {
	user, err := s.repo.GetUserByLogin(ctx, login)
	if err != nil {
		return nil, fmt.Errorf("invalid credentials")
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return nil, fmt.Errorf("invalid credentials")
	}

	return user, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
