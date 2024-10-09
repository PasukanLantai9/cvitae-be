package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bccfilkom/career-path-service/internal/api/resume"
	resumeHandler "github.com/bccfilkom/career-path-service/internal/api/resume/handler"
	resumeRepository "github.com/bccfilkom/career-path-service/internal/api/resume/repository"
	resumeService "github.com/bccfilkom/career-path-service/internal/api/resume/service"
	"github.com/bccfilkom/career-path-service/internal/config"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/middleware"
	"github.com/bccfilkom/career-path-service/internal/pkg/cronjob"
	"github.com/bccfilkom/career-path-service/pkg/mongo"
	"github.com/bccfilkom/career-path-service/pkg/postgres"
	redisdb "github.com/bccfilkom/career-path-service/pkg/redis"
	jobMatching "github.com/bccfilkom/career-path-service/pkg/rpc/job_matching"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
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
	redis       *redis.Client
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

	redisClient, err := redisdb.NewInstance()
	if err != nil {
		ts.FailNowf("database connection failed ", err.Error())
	}

	mlRPCClient := jobMatching.NewRpcClient(logrus.New())

	validator := config.NewValidator()
	cronjobClient := cronjob.NewInstance()
	resumeRepos := resumeRepository.New(mongos, db, redisClient)
	resumeServices := resumeService.New(resumeRepos, mlRPCClient)
	resumeHandlers := resumeHandler.New(resumeServices, logrus.New(), validator)

	ts.db = db
	ts.mongo = mongos
	ts.redis = redisClient
	ts.handler = resumeHandlers
	ts.accessToken = token

	cronjobClient.SetupCronJob(time.Second*5, func() error {
		return resumeServices.SyncResumesFromRedisToMongo(redisClient, mongos.Collection("resume"))
	})
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

	err = ts.redis.FlushAll(context.Background()).Err()
	if err != nil {
		ts.FailNowf("failed to clear data in Redis", err.Error())
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

func (ts *ResumeTestSuite) TestUpdateResume_Success() {
	id := ts.CreateData()

	app.Put("/resume/:id", middleware.JWTAccessToken(), ts.handler.HandleUpdateResume)

	updateData := resume.ResumeDetailDTO{
		PersonalDetails: resume.PersonalDetails{
			FullName:      "Updated Name",
			PhoneNumber:   "123-456-7890",
			Email:         "updated.email@example.com",
			Linkedin:      "https://linkedin.com/in/updated",
			PortfolioURL:  "https://portfolio.com/updated",
			Description:   "Updated description",
			AddressString: "Updated Address",
		},
		ProfessionalExperience: []resume.Experience{
			{
				StartDate:   resume.Date{Month: "January", Year: 2021},
				EndDate:     resume.Date{Month: "Present", Year: 2024},
				RoleTitle:   "Updated Role",
				CompanyName: "Updated Company",
				Location:    "Updated Location",
				Current:     true,
			},
		},
		Education: []resume.Education{
			{
				StartDate:   resume.Date{Month: "August", Year: 2017},
				EndDate:     resume.Date{Month: "May", Year: 2021},
				School:      "Updated University",
				Location:    "Updated Location",
				DegreeLevel: "Bachelor's",
				Major:       "Updated Major",
				GPA:         3.8,
				MaxGPA:      4.0,
			},
		},
		LeadershipExperience: []resume.Leadership{
			{
				StartDate:        resume.Date{Month: "June", Year: 2019},
				EndDate:          resume.Date{Month: "May", Year: 2020},
				RoleTitle:        "Team Leader",
				OrganisationName: "Updated Organization",
				Location:         "Updated Location",
				Current:          false,
			},
		},
		Others: []resume.Achievement{
			{
				Name:        "Updated Achievement",
				Date:        resume.Date{Month: "December", Year: 2020},
				Category:    "Award",
				Elaboration: resume.Elaboration{Text: "Updated description of the achievement."},
			},
		},
	}

	jsonData, _ := json.Marshal(updateData)

	endpoint := fmt.Sprintf("/resume/%s", id)
	request := httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusNoContent, response.StatusCode)
}

func (ts *ResumeTestSuite) TestUpdateResume_SuccessWithCache() {
	id := ts.CreateData()

	app.Put("/resume/:id", middleware.JWTAccessToken(), ts.handler.HandleUpdateResume)

	updateData := resume.ResumeDetailDTO{
		PersonalDetails: resume.PersonalDetails{
			FullName:      "Updated Name",
			PhoneNumber:   "123-456-7890",
			Email:         "updated.email@example.com",
			Linkedin:      "https://linkedin.com/in/updated",
			PortfolioURL:  "https://portfolio.com/updated",
			Description:   "Updated description",
			AddressString: "Updated Address",
		},
		ProfessionalExperience: []resume.Experience{
			{
				StartDate:   resume.Date{Month: "January", Year: 2021},
				EndDate:     resume.Date{Month: "Present", Year: 2024},
				RoleTitle:   "Updated Role",
				CompanyName: "Updated Company",
				Location:    "Updated Location",
				Current:     true,
			},
		},
		Education: []resume.Education{
			{
				StartDate:   resume.Date{Month: "August", Year: 2017},
				EndDate:     resume.Date{Month: "May", Year: 2021},
				School:      "Updated University",
				Location:    "Updated Location",
				DegreeLevel: "Bachelor's",
				Major:       "Updated Major",
				GPA:         3.8,
				MaxGPA:      4.0,
			},
		},
		LeadershipExperience: []resume.Leadership{
			{
				StartDate:        resume.Date{Month: "June", Year: 2019},
				EndDate:          resume.Date{Month: "May", Year: 2020},
				RoleTitle:        "Team Leader",
				OrganisationName: "Updated Organization",
				Location:         "Updated Location",
				Current:          false,
			},
		},
		Others: []resume.Achievement{
			{
				Name:        "Updated Achievement",
				Date:        resume.Date{Month: "December", Year: 2020},
				Category:    "Award",
				Elaboration: resume.Elaboration{Text: "Updated description of the achievement."},
			},
		},
	}

	jsonData, _ := json.Marshal(updateData)

	endpoint := fmt.Sprintf("/resume/%s", id)
	request := httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusNoContent, response.StatusCode)

	updateData.PersonalDetails.FullName = "Updated Name ke 2"

	jsonData, _ = json.Marshal(updateData)

	endpoint = fmt.Sprintf("/resume/%s", id)
	request = httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err = app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	time.Sleep(10 * time.Second)

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusNoContent, response.StatusCode)
}

func (ts *ResumeTestSuite) TestScoringResume_Success() {
	id := ts.CreateData()

	app.Put("/resume/:id", middleware.JWTAccessToken(), ts.handler.HandleUpdateResume)

	updateData := resume.ResumeDetailDTO{
		PersonalDetails: resume.PersonalDetails{
			FullName:      "Updated Name",
			PhoneNumber:   "123-456-7890",
			Email:         "updated.email@example.com",
			Linkedin:      "https://linkedin.com/in/updated",
			PortfolioURL:  "https://portfolio.com/updated",
			Description:   "Updated description",
			AddressString: "Updated Address",
		},
		ProfessionalExperience: []resume.Experience{
			{
				StartDate:   resume.Date{Month: "Jul", Year: 2024},
				EndDate:     resume.Date{},
				RoleTitle:   "Backend Developer",
				CompanyName: "Dakana Techworks",
				Location:    "Malang",
				Current:     true,
				Elaboration: []resume.Elaboration{
					{Text: "Designed database schema with 15++ tables."},
					{Text: "Maintained CI/CD pipeline using GitHub Actions."},
				},
			},
		},
		Education: []resume.Education{
			{
				StartDate:   resume.Date{Month: "Aug", Year: 2023},
				EndDate:     resume.Date{Month: "Aug", Year: 2027},
				School:      "Universitas Brawijaya",
				Location:    "Kota Malang, Indonesia",
				DegreeLevel: "Degree",
				Major:       "Informatics Engineering",
				GPA:         3.89,
				MaxGPA:      4.00,
				Elaboration: []resume.Elaboration{
					{Text: "Participated in hackathon competitions."},
					{Text: "Fullstack developer for student orientation committee."},
				},
			},
		},
		LeadershipExperience: []resume.Leadership{
			{
				StartDate:        resume.Date{Month: "Jan", Year: 2024},
				EndDate:          resume.Date{},
				RoleTitle:        "Fullstack Developer",
				OrganisationName: "Student Orientation Committee",
				Location:         "Brawijaya",
				Current:          true,
				Elaboration: []resume.Elaboration{
					{Text: "Built software solution for over 900+ users."},
				},
			},
		},
		Others: []resume.Achievement{
			{
				Name:        "2nd Place at Social Startup Games Hackfest 2024",
				Date:        resume.Date{Year: 2024},
				Category:    "Achievements",
				Elaboration: resume.Elaboration{Text: "2nd Place at Social Startup Games Hackfest 2024."},
			},
			{
				Name:        "Certified Kubernetes Administrator",
				Date:        resume.Date{Year: 2023},
				Category:    "Certificates",
				Elaboration: resume.Elaboration{Text: "Certified Kubernetes Administrator (CKA)."},
			},
			{
				Name:        "AWS Solutions Architect",
				Date:        resume.Date{Year: 2024},
				Category:    "Certificates",
				Elaboration: resume.Elaboration{Text: "AWS Certified Solutions Architect – Associate."},
			},
			{
				Name:        "Technical Skills",
				Category:    "Skills",
				Elaboration: resume.Elaboration{Text: "Proficient in Go, JavaScript, and Python. Strong skills in RESTful API design, microservices architecture, and containerization (Docker, Kubernetes)."},
			},
			{
				Name:        "Soft Skills",
				Category:    "Skills",
				Elaboration: resume.Elaboration{Text: "Effective communication, problem-solving, and time management."},
			},
		},
	}

	jsonData, _ := json.Marshal(updateData)

	endpoint := fmt.Sprintf("/resume/%s", id)
	request := httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err := app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusNoContent, response.StatusCode)

	updateData.PersonalDetails.FullName = "Updated Name ke 2"

	jsonData, _ = json.Marshal(updateData)

	endpoint = fmt.Sprintf("/resume/%s", id)
	request = httptest.NewRequest("PUT", endpoint, bytes.NewBuffer(jsonData))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)
	request.Header.Set("Content-Type", "application/json")

	response, err = app.Test(request)
	if err != nil {
		ts.FailNowf("request failed: %s", err.Error())
	}

	time.Sleep(10 * time.Second)

	ts.Assert().NotNil(response)
	ts.Assert().Equal(http.StatusNoContent, response.StatusCode)

	app.Get("/scoring/:id", middleware.JWTAccessToken(), ts.handler.HandleScoringResume)

	endpoint = fmt.Sprintf("/scoring/%s", id)
	request = httptest.NewRequest("GET", endpoint, bytes.NewBuffer(nil))
	request.Header.Set("Authorization", "Bearer "+ts.accessToken)

	response, err = app.Test(request, 2000)
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
