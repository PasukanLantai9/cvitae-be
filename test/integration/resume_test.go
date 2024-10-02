package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	resumeHandler "github.com/bccfilkom/career-path-service/internal/api/resume/handler"
	resumeRepository "github.com/bccfilkom/career-path-service/internal/api/resume/repository"
	resumeService "github.com/bccfilkom/career-path-service/internal/api/resume/service"
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/middleware"
	"github.com/bccfilkom/career-path-service/pkg/mongo"
	"github.com/bccfilkom/career-path-service/pkg/postgres"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson"
	mongo2 "go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/net/context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type ResumeTestSuite struct {
	suite.Suite
	db          *sqlx.DB
	mongo       *mongo2.Database
	handler     *resumeHandler.ResumeHandler
	accessToken string
}

func TestResumeSuite(t *testing.T) {
	suite.Run(t, new(ResumeTestSuite))
}

func (ts *ResumeTestSuite) CreateData() string {
	app.Post("/resume", middleware.JWTAccessToken(), ts.handler.HandleCreateResume)

	req := map[string]interface{}{
		"name": "test-resume-create",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/resume", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusCreated, response.StatusCode)

	bodyBytes, err := io.ReadAll(response.Body)

	var responseBody struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("request failed: %s", err.Error())
		return "0"
	} else {
		return responseBody.ID
	}
}

func (ts *ResumeTestSuite) SetupSuite() {
	if err := godotenv.Load("../../.env"); err != nil {
		ts.FailNowf("failed load environment", err.Error())
	}

	db, err := postgres.NewInstance()
	if err != nil {
		ts.FailNowf("database connection failed ", err.Error())
	}

	user, hashPass, err := generateUser("resume-user@example.com", "resume-user")
	if err != nil {
		ts.FailNowf("generate user failed ", err.Error())
	}

	_, err = db.Exec("INSERT INTO users (id, username, email, password, created_at) VALUES ($1, $2, $3, $4, $5)", user.ID, user.Username, user.Email, hashPass, time.Now())
	if err != nil {
		ts.FailNowf("generate user failed ", err.Error())
	}

	token, err := generateAccessToken(user, entity.AuthProviderNative)
	if err != nil {
		ts.FailNowf("generate token failed ", err.Error())
	}

	if err := db.Ping(); err != nil {
		ts.FailNowf("database ping failed ", err.Error())
	}

	mongos, err := mongo.NewMongoInstance()
	if err != nil {
		ts.FailNowf("database connection failed ", err.Error())
	}

	validator := config.NewValidator()
	resumeRepos := resumeRepository.New(mongos, db)
	resumeServices := resumeService.New(resumeRepos)
	resumeHandlers := resumeHandler.New(resumeServices, logrus.New(), validator)

	ts.db = db
	ts.mongo = mongos
	ts.handler = resumeHandlers
	ts.accessToken = token
}

func (ts *ResumeTestSuite) TearDownSuite() {
	_, err := ts.db.Exec("DELETE FROM resumes")
	if err != nil {
		ts.FailNowf("failed to clear data resumes", err.Error())
	}

	_, err = ts.db.Exec("DELETE FROM users")
	if err != nil {
		ts.FailNowf("failed to clear data users", err.Error())
	}

	_, err = ts.mongo.Collection("resume").DeleteMany(context.TODO(), bson.M{})
	if err != nil {
		ts.FailNowf("failed to clear data resume in mongo", err.Error())
	}

	if err := ts.db.Close(); err != nil {
		ts.FailNowf("database connection failed to close", err.Error())
	}
}

func (ts *ResumeTestSuite) TestCreateResume_Success() {
	app.Post("/resume", middleware.JWTAccessToken(), ts.handler.HandleCreateResume)

	req := map[string]interface{}{
		"name": "test-resume-create",
	}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/resume", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusCreated, response.StatusCode)
}

func (ts *ResumeTestSuite) TestCreateResume_FailedNoName() {
	app.Post("/resume", middleware.JWTAccessToken(), ts.handler.HandleCreateResume)

	req := map[string]interface{}{}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("POST", "/resume", bytes.NewBuffer(jsonReq))
	request.Header.Set("Content-Type", "application/json")
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusUnprocessableEntity, response.StatusCode)
}

func (ts *ResumeTestSuite) TestGetResume_Success() {
	ts.CreateData()

	app.Get("/resume", middleware.JWTAccessToken(), ts.handler.HandleGetUserResume)

	req := map[string]interface{}{}

	jsonReq, _ := json.Marshal(req)

	request := httptest.NewRequest("GET", "/resume", bytes.NewBuffer(jsonReq))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	bodyBytes, err := io.ReadAll(response.Body)

	var responseBody any
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	} else {
		ts.T().Log(responseBody)
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusOK, response.StatusCode)
}

func (ts *ResumeTestSuite) TestGetResumeDetail_Success() {
	id := ts.CreateData()

	app.Get("/resume/:id", middleware.JWTAccessToken(), ts.handler.HandleGetResumeByID)

	req := map[string]interface{}{}

	jsonReq, _ := json.Marshal(req)

	endpoint := fmt.Sprintf("%s/%s", "/resume", id)
	request := httptest.NewRequest("GET", endpoint, bytes.NewBuffer(jsonReq))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	bodyBytes, err := io.ReadAll(response.Body)

	var responseBody any
	if err := json.Unmarshal(bodyBytes, &responseBody); err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	} else {
		ts.T().Log(responseBody)
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusOK, response.StatusCode)
}
