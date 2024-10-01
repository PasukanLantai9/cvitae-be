package linkedin

import (
	"github.com/bccfilkom/career-path-service/internal/entity"
	"github.com/bccfilkom/career-path-service/internal/pkg/env"
	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	linkedinOauthPkg "golang.org/x/oauth2/linkedin"
)

type Linkedin struct {
	Outh interface {
		GetUserFromCode(context.Context, string) (entity.User, entity.UserOauth, error)
	}
}

func New() Linkedin {
	linkedinOauthConf := oauth2.Config{
		RedirectURL:  "http://localhost:5173",
		ClientID:     env.GetString("LINKEDIN_CLIENT_ID", ""),
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo"},
		Endpoint:     linkedinOauthPkg.Endpoint,
		ClientSecret: env.GetString("LINKEDIN_CLIENT_SECRET", ""),
	}

	return Linkedin{
		Outh: &linkedinOauth{
			config: linkedinOauthConf,
		},
	}
}

type linkedinOauth struct {
	config oauth2.Config
}
