package resumeRepository

import (
	"context"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"go.mongodb.org/mongo-driver/mongo"
)

// New creates a new MongoDB repository
func New(db *mongo.Database, sqlDB *sqlx.DB) Repository {
	return &repository{
		MongoDB: db,
		SqlDB:   sqlDB,
	}
}

type repository struct {
	MongoDB *mongo.Database
	SqlDB   *sqlx.DB
}

type Repository interface {
	NewMongoClient(ctx context.Context, tx bool) (MongoClient, error)
	NewSqlClient(tx bool) (SqlClient, error)
}

func (r *repository) NewMongoClient(ctx context.Context, tx bool) (MongoClient, error) {
	var session mongo.Session
	var commitFunc, rollbackFunc func() error

	if tx {
		var err error
		session, err = r.MongoDB.Client().StartSession()
		if err != nil {
			return MongoClient{}, err
		}

		err = session.StartTransaction()
		if err != nil {
			return MongoClient{}, err
		}

		commitFunc = func() error {
			return session.CommitTransaction(ctx)
		}
		rollbackFunc = func() error {
			return session.AbortTransaction(ctx)
		}
	} else {
		commitFunc = func() error { return nil }
		rollbackFunc = func() error { return nil }
	}

	return MongoClient{
		Resume:   &resumeRepository{db: r.MongoDB, session: session},
		Commit:   commitFunc,
		Rollback: rollbackFunc,
	}, nil
}

func (r *repository) NewSqlClient(tx bool) (SqlClient, error) {
	var db sqlx.ExtContext
	var commitFunc, rollbackFunc func() error

	db = r.SqlDB

	if tx {
		var err error
		txx, err := r.SqlDB.Beginx()
		if err != nil {
			return SqlClient{}, err
		}

		db = txx
		commitFunc = txx.Commit
		rollbackFunc = txx.Rollback
	} else {
		commitFunc = func() error { return nil }
		rollbackFunc = func() error { return nil }
	}

	return SqlClient{
		Resume:   &resumeSQLRepository{q: db},
		Commit:   commitFunc,
		Rollback: rollbackFunc,
	}, nil
}

type MongoClient struct {
	Resume interface {
		CreateResume(ctx context.Context, resume entity.ResumeDetail) (*mongo.InsertOneResult, error)
		GetByIDAndUserID(ctx context.Context, ID string, userID string) (entity.ResumeDetail, error)
	}
	Commit   func() error
	Rollback func() error
}

type resumeRepository struct {
	db      *mongo.Database
	session mongo.Session
}

type SqlClient struct {
	Resume interface {
		CreateResume(ctx context.Context, resume entity.Resume) error
		GetByUserID(ctx context.Context, userID string) ([]entity.Resume, error)
	}
	Commit   func() error
	Rollback func() error
}

type resumeSQLRepository struct {
	q sqlx.ExtContext
}
