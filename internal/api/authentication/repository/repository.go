package authRepository

import (
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

func New(db *sqlx.DB) Repository {
	return &repository{
		DB: db,
	}
}

type repository struct {
	DB *sqlx.DB
}

type Repository interface {
	NewClient(tx bool) (Client, error)
}

func (r *repository) NewClient(tx bool) (Client, error) {
	var db sqlx.ExtContext
	var commitFunc, rollbackFunc func() error

	db = r.DB

	if tx {
		var err error
		txx, err := r.DB.Beginx()
		if err != nil {
			return Client{}, err
		}

		db = txx
		commitFunc = txx.Commit
		rollbackFunc = txx.Rollback
	} else {
		commitFunc = func() error { return nil }
		rollbackFunc = func() error { return nil }
	}

	return Client{
		Users:     &userRepository{q: db},
		Sessions:  &sessionRepository{q: db},
		UserOauth: &userOauthRepository{q: db},
		Commit:    commitFunc,
		Rollback:  rollbackFunc,
	}, nil
}

type Client struct {
	Users interface {
		CreateUser(context.Context, entity.User) error
		GetByEmail(context.Context, string) (entity.User, error)
		GetByID(context.Context, string) (entity.User, error)
	}
	Sessions interface {
		CreateSession(context.Context, entity.Session) error
		GetByID(context.Context, string) (entity.Session, error)
	}
	UserOauth interface {
		CreateUserOauth(context.Context, entity.UserOauth) error
		GetByOauthUserID(context.Context, string, string, entity.AuthProvider) (entity.UserOauth, error)
	}

	Commit   func() error
	Rollback func() error
}

type userRepository struct {
	q sqlx.ExtContext
}

type sessionRepository struct {
	q sqlx.ExtContext
}

type userOauthRepository struct {
	q sqlx.ExtContext
}
