package google

import (
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	googleOauthPkg "golang.org/x/oauth2/google"
)

type Google struct {
	Outh interface {
		GetUserFromCode(context.Context, string) (entity.User, entity.UserOauth, error)
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

	return Google{
		Outh: &googleOauth{
			config: googleOauthConf,
		},
	}
}

type googleOauth struct {
	config oauth2.Config
}
