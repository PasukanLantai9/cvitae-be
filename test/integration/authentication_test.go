package integration

import (
	"bytes"
	"encoding/json"
	authHandler "github.com/bccfilkom/career-path-service/internal/api/authentication/handler"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	authService "github.com/bccfilkom/career-path-service/internal/api/authentication/service"
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"testing"
)

type AuthTestSuite struct {
	suite.Suite
	db      *sqlx.DB
	handler *authHandler.AuthHandler
}

func TestAuthSuite(t *testing.T) {
	suite.Run(t, new(AuthTestSuite))
}

func (ts *AuthTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		ts.FailNowf("failed load environment", err.Error())
	}

	db, err := sqlx.Open("postgres", GenerateDSN())
	if err != nil {
		ts.FailNowf("database connection failed ", err.Error())
	}

	if err := db.Ping(); err != nil {
		ts.FailNowf("database ping failed ", err.Error())
	}

	validator := config.NewValidator()
	authRepos := authRepository.New(db)
	authServices := authService.New(authRepos)
	authHandlers := authHandler.New(authServices, validator)

	ts.db = db
	ts.handler = authHandlers
}

func (ts *AuthTestSuite) TearDownSuite() {
	_, err := ts.db.Exec("DELETE FROM users")
	if err != nil {
		ts.FailNowf("failed to clear data", err.Error())
	}

	if err := ts.db.Close(); err != nil {
		ts.FailNowf("database connection failed to close", err.Error())
	}
}

func (ts *AuthTestSuite) TestCreateUser_NoEmail() {
	app.Post("/auth/register", ts.handler.HandleRegister)

	req := map[string]interface{}{
		"name":     "test user",
		"password": "12345678",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusUnprocessableEntity, response.StatusCode)
}

func (ts *AuthTestSuite) TestCreateUser_PasswordLengthConstraint() {
	app.Post("/auth/register", ts.handler.HandleRegister)

	req := map[string]interface{}{
		"name":     "test user",
		"email":    "test@example.com",
		"password": "123456",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusUnprocessableEntity, response.StatusCode)
}

func (ts *AuthTestSuite) TestCreateUser_Success() {
	app.Post("/auth/register", ts.handler.HandleRegister)

	req := map[string]interface{}{
		"name":     "test user",
		"email":    "test@example.com",
		"password": "12345678",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/register", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusCreated, response.StatusCode)
}
