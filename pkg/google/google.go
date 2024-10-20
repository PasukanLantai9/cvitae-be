package google

import (
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"github.com/google/generative-ai-go/genai"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	googleOauthPkg "golang.org/x/oauth2/google"
	"google.golang.org/api/option"
)

type Google struct {
	Outh interface {
		GetUserFromCode(context.Context, string) (entity.User, entity.UserOauth, error)
	}
	Gemini interface {
		GenerateResumeJsonFromPDF(context.Context, []byte) (string, error)
		GenerateExperienceAndSkillsParagrafFromJSON(context.Context, []byte) (string, error)
	}
}

func New() Google {
	googleOauthConf := oauth2.Config{
		RedirectURL:  "http://localhost:5173",
		ClientID:     env.GetString("GOOGLE_CLIENT_ID", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo"},
		Endpoint:     googleOauthPkg.Endpoint,
		ClientSecret: env.GetString("GOOGLE_CLIENT_SECRET", ""),
	}

	key := env.GetString("GEMINI_API_KEY", "")
	typeModel := env.GetString("GEMINI_TYPE_MODEL", "")
	client, err := genai.NewClient(context.Background(), option.WithAPIKey(key))
	if err != nil {
		panic(err)
	}

	model := client.GenerativeModel(typeModel)

	return Google{
		Outh: &googleOauth{
			config: googleOauthConf,
		},
		Gemini: &googleGemini{
			model: model,
		},
	}
}

type googleOauth struct {
	config oauth2.Config
}

type googleGemini struct {
	model *genai.GenerativeModel
}
