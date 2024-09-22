package authService

import (
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/helper"
	"golang.org/x/net/context"
	"time"
)

func (s *authService) RegisterUser(ctx context.Context, req authentication.CreateUserRequest) error {
	repo, err := s.authRepository.NewClient(false)
	if err != nil {
		return err
	}

	hashedPassword, err := helper.HashPassword(req.Paassword)
	if err != nil {
		return err
	}

	ulid, err := helper.NewUlidFromTimestamp(time.Now())
	if err != nil {
		return err
	}

	user := entity.User{
		ID:       ulid,
		Username: req.Name,
		Password: hashedPassword,
		Email:    req.Email,
	}

	if err := repo.Users.CreateUser(ctx, user); err != nil {
		return err
	}

	return nil
}
