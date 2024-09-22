package authRepository

const (
	queryCreateUser = `
		INSERT INTO Users (id, email, username, password, created_at) 
		VALUES (:id, :email, :username, :password, :created_at)`
)
