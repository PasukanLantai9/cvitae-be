package resumeRepository

import (
	"time"

	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/jmoiron/sqlx"
	"golang.org/x/net/context"
)

type resumeDBSQL struct {
	ID        string    `db:"id"`
	Name      string    `db:"name"`
	UserID    string    `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

func (data *resumeDBSQL) format() entity.Resume {
	return entity.Resume{
		ID:        data.ID,
		Name:      data.Name,
		UserID:    data.UserID,
		CreatedAt: data.CreatedAt,
	}
}

func (r *resumeSQLRepository) CreateResume(ctx context.Context, resume entity.Resume) error {
	argsKV := map[string]interface{}{
		"id":         resume.ID,
		"name":       resume.Name,
		"user_id":    resume.UserID,
		"created_at": resume.CreatedAt,
	}

	query, args, err := sqlx.Named(queryCreateResume, argsKV)
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

func (r *resumeSQLRepository) GetByUserID(ctx context.Context, userID string) ([]entity.Resume, error) {
	argsKV := map[string]interface{}{
		"user_id": userID,
	}

	query, args, err := sqlx.Named(queryGetResumeByUserID, argsKV)
	if err != nil {
		return nil, err
	}
	query = r.q.Rebind(query)

	resumes := make([]entity.Resume, 0)
	rows, err := r.q.QueryxContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		resume := resumeDBSQL{}
		if err := rows.StructScan(&resume); err != nil {
			return nil, err
		}
		resumes = append(resumes, resume.format())
	}

	return resumes, nil
}

func (r *resumeSQLRepository) DeleteById(ctx context.Context, ID string) error {
	argsKV := map[string]interface{}{
		"id": ID,
	}

	query, args, err := sqlx.Named(queryDeleteResumeByID, argsKV)
	if err != nil {
		return err
	}
	query = r.q.Rebind(query)

	_, err = r.q.ExecContext(ctx, query, args...)
	return err
}
