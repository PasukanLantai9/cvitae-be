package integration

import (
	"bytes"
	"encoding/json"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	authHandler "github.com/bccfilkom/career-path-service/internal/api/authentication/handler"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	authService "github.com/bccfilkom/career-path-service/internal/api/authentication/service"
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/bccfilkom/career-path-service/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/stretchr/testify/suite"
	"io"
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

	db, err := postgres.NewInstance()
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
		ts.FailNowf("failed to clear data users", err.Error())
	}

	_, err = ts.db.Exec("DELETE FROM user_sessions")
	if err != nil {
		ts.FailNowf("failed to clear data user_sessions", err.Error())
	}

	if err := ts.db.Close(); err != nil {
		ts.FailNowf("database connection failed to close", err.Error())
	}
}

func (ts *AuthTestSuite) TestRegister_NoEmail() {
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

func (ts *AuthTestSuite) TestRegister_PasswordLengthConstraint() {
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

func (ts *AuthTestSuite) TestRegister_Success() {
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

func (ts *AuthTestSuite) TestSingin_WrongEmail() {
	app.Post("/auth/signin", ts.handler.HandleSignin)

	req := map[string]interface{}{
		"email":    "test-salah@example.com",
		"password": "12345678",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusBadRequest, response.StatusCode)
}

func (ts *AuthTestSuite) TestSingin_WrongPassword() {
	app.Post("/auth/signin", ts.handler.HandleSignin)

	req := map[string]interface{}{
		"email":    "test@example.com",
		"password": "123456789",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusBadRequest, response.StatusCode)
}

func (ts *AuthTestSuite) TestSingin_Success() {
	app.Post("/auth/signin", ts.handler.HandleSignin)

	req := map[string]interface{}{
		"email":    "test@example.com",
		"password": "12345678",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusOK, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		ts.FailNowf("failed to read response body: %s", err.Error())
	}

	var responseBody authentication.SinginUserResponse
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("failed to unmarshal response body: %s", err.Error())
	}

	ts.Assert().NotEmpty(responseBody.AccessToken)
	ts.Assert().NotEmpty(responseBody.RefreshToken)
	ts.Assert().NotEmpty(responseBody.SessionID)
	ts.Assert().NotEmpty(responseBody.ExpiresInSeconds)

	ts.Assert().Equal(1800, int(responseBody.ExpiresInSeconds))
}
