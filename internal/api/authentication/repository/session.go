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

type SessionDB struct {
	ID           string              `db:"id"`
	UserID       string              `db:"user_id"`
	RefreshToken string              `db:"refresh_token"`
	CreatedAt    string              `db:"created_at"`
	ExpiresAt    time.Time           `db:"expires_at"`
	AuthProvider entity.AuthProvider `db:"auth_provider"`
}

func (data *SessionDB) format() entity.Session {
	return entity.Session{
		ID:           data.ID,
		UserID:       data.UserID,
		RefreshToken: data.RefreshToken,
		CreatedAt:    data.CreatedAt,
		ExpiresAt:    data.ExpiresAt,
		AuthProvider: data.AuthProvider,
	}
}

func (r *sessionRepository) CreateSession(ctx context.Context, session entity.Session) error {
	argsKV := map[string]interface{}{
		"id":            session.ID,
		"user_id":       session.UserID,
		"refresh_token": session.RefreshToken,
		"expires_at":    session.ExpiresAt,
		"auth_provider": session.AuthProvider,
		"created_at":    time.Now(),
	}

	query, args, err := sqlx.Named(queryCreateSession, argsKV)
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

func (r *sessionRepository) GetByID(ctx context.Context, id string) (entity.Session, error) {
	argsKV := map[string]interface{}{
		"id": id,
	}

	query, args, err := sqlx.Named(queryGetSessionByID, argsKV)
	if err != nil {
		return entity.Session{}, err
	}
	query = r.q.Rebind(query)

	var session SessionDB
	if err := r.q.QueryRowxContext(ctx, query, args...).StructScan(&session); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.Session{}, authentication.ErrSessionWithIDNotFound
		}
		return entity.Session{}, err
	}

	return session.format(), nil
}
