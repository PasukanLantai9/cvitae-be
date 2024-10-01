package authRepository

const (
	queryCreateUser = `
		INSERT INTO Users (id, email, username, password, created_at) 
		VALUES (:id, :email, :username, :password, :created_at)`

	queryGetUserByEmail = `
		SELECT id, email, username, password, created_at 
		FROM Users
		WHERE email = :email`

	queryGetUserByID = `
		SELECT id, email, username, password, created_at 
		FROM Users
		WHERE id = :id`
)

const (
	queryCreateSession = `
		INSERT INTO user_sessions (id, user_id, refresh_token, expires_at, auth_provider, created_at)
		VALUES (:id, :user_id, :refresh_token, :expires_at, :auth_provider, :created_at)`

	queryGetSessionByID = `
		SELECT id, user_id, refresh_token, expires_at, auth_provider, created_at
		FROM user_sessions
		WHERE id = :id`
)

const (
	queryCreateOauthUser = `
		INSERT INTO user_oauth (id, user_id, provider, oauth_user_id, access_token, refresh_token, created_at)
		VALUES (:id, :user_id, :provider, :oauth_user_id, :access_token, :refresh_token, :created_at)`

	queryGetByOauthUserID = `
		SELECT id, user_id, provider, oauth_user_id, access_token, refresh_token, created_at
		FROM user_oauth
		WHERE oauth_user_id = :oauth_user_id AND user_id = :user_id AND provider = :provider`
)
