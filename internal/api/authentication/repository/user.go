package authRepository

import (
	"database/sql"
	"errors"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"golang.org/x/net/context"
	"time"
)

type userDB struct {
	ID        string    `db:"id"`
	Username  string    `db:"username"`
	Password  string    `db:"password"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
}

func (data *userDB) format() entity.User {
	return entity.User{
		ID:       data.ID,
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
	}
}

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
		if errors.As(err, &pqErr) {
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

func (r *userRepository) GetByEmail(ctx context.Context, email string) (entity.User, error) {
	argsKV := map[string]interface{}{
		"email": email,
	}

	query, args, err := sqlx.Named(queryGetUserByEmail, argsKV)
	if err != nil {
		return entity.User{}, err
	}
	query = r.q.Rebind(query)

	var user userDB
	if err := r.q.QueryRowxContext(ctx, query, args...).StructScan(&user); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return entity.User{}, authentication.ErrUserWithEmailNotFound
		}
		return entity.User{}, err
	}

	return user.format(), err
}
