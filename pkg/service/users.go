package service

import (
	"context"
	"time"

	"github.com/maxoov1/faq-api/pkg/auth"
	"github.com/maxoov1/faq-api/pkg/hash"
	"github.com/maxoov1/faq-api/pkg/models"
	"github.com/maxoov1/faq-api/pkg/repository"
)

type UsersService struct {
	repo     repository.Users
	hasher   hash.Hasher
	manager  auth.TokenManager
	tokenTTL time.Duration
}

func NewUsersService(repo repository.Users, hasher hash.Hasher, manager auth.TokenManager, ttl time.Duration) *UsersService {
	return &UsersService{
		repo:     repo,
		hasher:   hasher,
		manager:  manager,
		tokenTTL: ttl,
	}
}

func (s *UsersService) SignUp(ctx context.Context, userInput UserInputSignUp) error {
	passwordHash, err := s.hasher.Hash(userInput.Password)
	if err != nil {
		return err
	}

	user := models.User{
		Name:     userInput.Name,
		Email:    userInput.Email,
		Password: passwordHash,
	}

	return s.repo.Create(ctx, user)
}

func (s *UsersService) SignIn(ctx context.Context, userInput UserInputSignIn) (string, error) {
	passwordHash, err := s.hasher.Hash(userInput.Password)
	if err != nil {
		return "", err
	}

	user, err := s.repo.GetByCredentials(ctx, userInput.Email, passwordHash)
	if err != nil {
		return "", err
	}

	return s.manager.NewJWT(user.ID.Hex(), s.tokenTTL)
}
