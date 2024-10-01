package authRepository

import (
	"database/sql"
	"errors"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
	"time"
)

type UserOauthDB struct {
	ID           string              `db:"id"`
	UserID       string              `db:"user_id"`
	Provider     entity.AuthProvider `db:"provider"`
	OAuthUserID  string              `db:"oauth_user_id"`
	AccessToken  string              `db:"access_token"`
	RefreshToken string              `db:"refresh_token"`
	CreatedAt    time.Time           `db:"created_at"`
}

func (data *UserOauthDB) format() entity.UserOauth {
	return entity.UserOauth{
		ID:           data.ID,
		UserID:       data.UserID,
		Provider:     data.Provider,
		OAuthUserID:  data.OAuthUserID,
		AccessToken:  data.AccessToken,
		RefreshToken: data.RefreshToken,
		CreatedAt:    data.CreatedAt,
	}
}

func (r *userOauthRepository) CreateUserOauth(ctx context.Context, oauth entity.UserOauth) error {
	argsKV := map[string]interface{}{
		"id":            oauth.ID,
		"provider":      oauth.Provider,
		"user_id":       oauth.UserID,
		"oauth_user_id": oauth.OAuthUserID,
		"access_token":  oauth.AccessToken,
		"refresh_token": oauth.RefreshToken,
		"created_at":    oauth.CreatedAt,
	}

	query, args, err := sqlx.Named(queryCreateOauthUser, argsKV)
	if err != nil {
		return err
	}
	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	if err != nil {
		return err
	}

	return nil
}

func (r *userOauthRepository) GetByOauthUserID(ctx context.Context, OauthUserID string, userID string, provider entity.AuthProvider) (entity.UserOauth, error) {
	argsKV := map[string]interface{}{
		"user_id":       userID,
		"oauth_user_id": OauthUserID,
		"provider":      provider,
	}

	query, args, err := sqlx.Named(queryGetByOauthUserID, argsKV)
	if err != nil {
		return entity.UserOauth{}, err
	}
	query = r.q.Rebind(query)

	var oauthUserDB UserOauthDB
	if err := r.q.QueryRowxContext(ctx, query, args...).StructScan(&oauthUserDB); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.UserOauth{}, authentication.ErrUserWithOauthUserIDNotFound
		}
		return entity.UserOauth{}, err
	}

	return oauthUserDB.format(), nil
}
