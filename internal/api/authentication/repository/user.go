package authRepository

import (
	"errors"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"time"
)

func (r *userRepository) CreateUser(ctx context.Context, user entity.User) error {
	argsKV := map[string]interface{}{
		"id":         user.ID,
		"email":      user.Email,
		"username":   user.Username,
		"password":   user.Password,
		"created_at": time.Now(),
	}

	query, args, err := sqlx.Named(queryCreateUser, argsKV)
	if err != nil {
		return err
	}
	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	if err != nil {
		var pqErr *pq.Error
		if ok := errors.As(err, &pqErr); ok {
			switch pqErr.Code {
			case "23505":
				if pqErr.Constraint == "users_email_key" {
					return authentication.ErrEmailAlreadyExists
				}
			}
		}
		return err
	}

	return nil
}
