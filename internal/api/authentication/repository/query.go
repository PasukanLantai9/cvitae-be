package authRepository

const (
	queryCreateUser = `
		INSERT INTO Users (id, email, username, password, created_at) 
		VALUES (:id, :email, :username, :password, :created_at)`

	queryGetUserByEmail = `
		SELECT id, email, username, password, created_at 
		FROM Users
		WHERE email = :email`
)

const (
	queryCreateSession = `
		INSERT INTO user_sessions (id, user_id, refresh_token, expires_at, auth_provider, created_at)
		VALUES (:id, :user_id, :refresh_token, :expires_at, :auth_provider, :created_at)`
)
