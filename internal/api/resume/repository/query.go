package resumeRepository

const (
	queryCreateResume = `
		INSERT INTO resumes (id, name, user_id, created_at)
		VALUES (:id, :name, :user_id, :created_at)`

	queryGetResumeByUserID = `
		SELECT id, name, user_id, created_at
		FROM resumes WHERE user_id = :user_id`

	queryDeleteResumeByID = `
		DELETE
		FROM resumes WHERE id = :id`
)
