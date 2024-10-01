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

	if user.Password == "" {
		return authentication.SinginUserResponse{}, authentication.ErrEmailOrPasswordWreong
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

func (s *authService) GenerateNewAccessToken(ctx context.Context, req entity.Session) (authentication.RefreshTokenResponse, error) {
	repo, err := s.authRepository.NewClient(false)
	if err != nil {
		return authentication.RefreshTokenResponse{}, err
	}

	originData, err := repo.Sessions.GetByID(ctx, req.ID)
	if err != nil {
		return authentication.RefreshTokenResponse{}, authentication.ErrUnableToCreateNewAccessToken
	}

	if originData.RefreshToken != req.RefreshToken {
		return authentication.RefreshTokenResponse{}, authentication.ErrUnableToCreateNewAccessToken
	}

	if time.Now().After(originData.ExpiresAt) {
		return authentication.RefreshTokenResponse{}, authentication.ErrUnableToCreateNewAccessToken
	}

	user, err := repo.Users.GetByID(ctx, originData.UserID)
	if err != nil {
		return authentication.RefreshTokenResponse{}, err
	}

	accessToken, err := generateAccessToken(user, entity.AuthProviderNative)
	if err != nil {
		return authentication.RefreshTokenResponse{}, err
	}

	return authentication.RefreshTokenResponse{
		AccessToken:      accessToken,
		ExpiresInSeconds: int64((time.Minute * 30).Seconds()),
	}, nil
}

func (s *authService) CallbackGoogleOauth(ctx context.Context, code string) (authentication.SinginUserResponse, error) {
	repo, err := s.authRepository.NewClient(true)

	defer func() {
		if err != nil {
			if errTx := repo.Rollback(); errTx != nil {
				err = errTx
				return
			}
		}
	}()

	user, userOauth, err := s.google.Outh.GetUserFromCode(ctx, code)
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	originDataUser, err := repo.Users.GetByEmail(ctx, user.Email)
	if err != nil {
		if errors.Is(err, authentication.ErrUserWithEmailNotFound) {
			ulid, err := helper.NewUlidFromTimestamp(time.Now())
			if err != nil {
				return authentication.SinginUserResponse{}, err
			}

			user := entity.User{
				ID:       ulid,
				Username: user.Username,
				Email:    user.Email,
			}

			if err := repo.Users.CreateUser(ctx, user); err != nil {
				return authentication.SinginUserResponse{}, err
			}

			originDataUser.ID = user.ID
		} else {
			return authentication.SinginUserResponse{}, err
		}
	}

	originDataUserOauth, err := repo.UserOauth.GetByOauthUserID(ctx, userOauth.OAuthUserID, originDataUser.ID, entity.AuthProviderGoogle)
	if err != nil {
		if errors.Is(err, authentication.ErrUserWithOauthUserIDNotFound) {
			ulid, err := helper.NewUlidFromTimestamp(time.Now())
			if err != nil {
				return authentication.SinginUserResponse{}, err
			}

			userOauth.ID = ulid
			userOauth.UserID = originDataUser.ID
			userOauth.CreatedAt = time.Now()

			if err := repo.UserOauth.CreateUserOauth(ctx, userOauth); err != nil {
				return authentication.SinginUserResponse{}, err
			}

			originDataUserOauth.ID = userOauth.ID
			originDataUserOauth.CreatedAt = userOauth.CreatedAt
			originDataUserOauth.UserID = userOauth.UserID
			originDataUserOauth.OAuthUserID = userOauth.OAuthUserID
			originDataUserOauth.AccessToken = userOauth.AccessToken
			originDataUserOauth.RefreshToken = userOauth.RefreshToken
			originDataUserOauth.Provider = userOauth.Provider
		} else {
			return authentication.SinginUserResponse{}, err
		}
	}

	accessToken, refreshToken, err := generateFreshToken(user, entity.AuthProviderGoogle)
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	sessionID, err := helper.NewUlidFromTimestamp(time.Now())
	if err != nil {
		return authentication.SinginUserResponse{}, err
	}

	session := entity.Session{
		ID:           sessionID,
		UserID:       originDataUser.ID,
		AuthProvider: entity.AuthProviderGoogle,
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
		RefreshToken:     refreshToken,
		AccessToken:      accessToken,
		ExpiresInSeconds: int64((time.Minute * 30).Seconds()),
		SessionID:        sessionID,
	}, nil
}
