package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/api/authentication"
	authHandler "github.com/bccfilkom/career-path-service/internal/api/authentication/handler"
	authRepository "github.com/bccfilkom/career-path-service/internal/api/authentication/repository"
	authService "github.com/bccfilkom/career-path-service/internal/api/authentication/service"
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/bccfilkom/career-path-service/pkg/google"
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
	authServices := authService.New(authRepos, google.Google{})
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

func (ts *AuthTestSuite) signin(email string) authentication.SinginUserResponse {
	app.Post("/auth/signin", ts.handler.HandleSignin)

	dumReq := map[string]interface{}{
		"email":    email,
		"password": "12345678",
	}

	jsonReq, _ := json.Marshal(dumReq)

	request := httptest.NewRequest("POST", "/auth/signin", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	bodyBytes, err := io.ReadAll(response.Body)
	if err != nil {
		ts.FailNowf("read request failed: %s", err.Error())
		return authentication.SinginUserResponse{}
	}

	var responseBody authentication.SinginUserResponse
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("unmarshal request failed: %s", err.Error())
		return authentication.SinginUserResponse{}
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusOK, response.StatusCode)

	return responseBody
}

func (ts *AuthTestSuite) register(email string) {
	app.Post("/auth/register", ts.handler.HandleRegister)

	req := map[string]interface{}{
		"name":     "test user refresh",
		"email":    email,
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

func (ts *AuthTestSuite) TestRefreshToken_WrongSessionID() {
	app.Post("/auth/refresh", ts.handler.HandleRefreshToken)

	ts.register("test-wrong-session-id@gmail.com")
	signinData := ts.signin("test-wrong-session-id@gmail.com")

	req := map[string]interface{}{
		"refreshToken": signinData.RefreshToken,
	}

	jsonReq, _ := json.Marshal(req)

	endpoints := fmt.Sprintf("/auth/refresh?sessionID=%s-wrong", signinData.SessionID)
	request := httptest.NewRequest("POST", endpoints, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusUnauthorized, response.StatusCode)
}

func (ts *AuthTestSuite) TestRefreshToken_WrongSessionRefresh() {
	app.Post("/auth/refresh", ts.handler.HandleRefreshToken)

	ts.register("test-wrong-token@gmail.com")
	signinData := ts.signin("test-wrong-token@gmail.com")

	req := map[string]interface{}{
		"refreshToken": fmt.Sprintf("%s-wrong", signinData.RefreshToken),
	}

	jsonReq, _ := json.Marshal(req)

	endpoints := fmt.Sprintf("/auth/refresh?sessionID=%s", signinData.SessionID)
	request := httptest.NewRequest("POST", endpoints, bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusUnauthorized, response.StatusCode)
}

func (ts *AuthTestSuite) TestRefreshToken_Success() {
	app.Post("/auth/refresh", ts.handler.HandleRefreshToken)

	ts.register("test-refresh-success@gmail.com")
	signinData := ts.signin("test-refresh-success@gmail.com")

	req := map[string]interface{}{
		"refreshToken": fmt.Sprintf("%s", signinData.RefreshToken),
	}

	jsonReq, _ := json.Marshal(req)

	endpoints := fmt.Sprintf("/auth/refresh?sessionID=%s", signinData.SessionID)
	request := httptest.NewRequest("POST", endpoints, bytes.NewBuffer(jsonReq))
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

	var responseBody authentication.RefreshTokenResponse
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("failed to unmarshal response body: %s", err.Error())
	}

	ts.Assert().NotEmpty(responseBody.AccessToken)
	ts.Assert().NotEmpty(responseBody.ExpiresInSeconds)

	ts.Assert().Equal(1800, int(responseBody.ExpiresInSeconds))
}
