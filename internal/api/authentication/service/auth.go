package authService

import (
	"errors"
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

func (s *authService) SinginUser(ctx context.Context, req authentication.SigninUserRequest) (authentication.SinginUserResponse, error) {
	repo, err := s.authRepository.NewClient(true)
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	defer func() {
		if err != nil {
			if errTx := repo.Rollback(); errTx != nil {
				err = errTx
				return
			}
		}
	}()

	user, err := repo.Users.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, authentication.ErrUserWithEmailNotFound) {
			return authentication.SinginUserResponse{}, authentication.ErrEmailOrPasswordWreong
		}
		return authentication.SinginUserResponse{}, err
	}

	if err := helper.ComparePassword(user.Password, req.Password); err != nil {
		return authentication.SinginUserResponse{}, authentication.ErrEmailOrPasswordWreong
	}

	accessToken, refreshToken, err := generateFreshToken(user, entity.AuthProviderNative)
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	sessionID, err := helper.NewUlidFromTimestamp(time.Now())
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	session := entity.Session{
		ID:           sessionID,
		UserID:       user.ID,
		AuthProvider: entity.AuthProviderNative,
		RefreshToken: refreshToken,
		ExpiresAt:    time.Now().Add(time.Hour * 48),
	}

	if err := repo.Sessions.CreateSession(ctx, session); err != nil {
		return authentication.SinginUserResponse{}, err
	}

	if err := repo.Commit(); err != nil {
		return authentication.SinginUserResponse{}, err
	}

	return authentication.SinginUserResponse{
		AccessToken:      accessToken,
		SessionID:        session.ID,
		RefreshToken:     refreshToken,
		ExpiresInSeconds: int64((time.Minute * 30).Seconds()),
	}, nil
}
