package linkedin

import (
	"encoding/json"
	"github.com/bccfilkom/career-path-service/internal/entity"
	"golang.org/x/net/context"
	"io"
)

type linkedinOauthUser struct {
	Email         string `json:"email"`
	FamilyName    string `json:"family_name"`
	GivenName     string `json:"given_name"`
	HD            string `json:"hd"`
	ID            string `json:"id"`
	Name          string `json:"name"`
	Picture       string `json:"picture"`
	VerifiedEmail bool   `json:"verified_email"`
}

func (data *linkedinOauthUser) formatUser() entity.User {
	return entity.User{
		Email:    data.Email,
		Username: data.Name,
	}
}

func (data *linkedinOauthUser) formatOauth(accessToken string, refreshToken string) entity.UserOauth {
	return entity.UserOauth{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		OAuthUserID:  data.ID,
		Provider:     entity.AuthProviderGoogle,
	}
}

func (g *linkedinOauth) GetUserFromCode(ctx context.Context, code string) (entity.User, entity.UserOauth, error) {
	token, err := g.config.Exchange(ctx, code)
	if err != nil {
		return entity.User{}, entity.UserOauth{}, err
	}

	client := g.config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return entity.User{}, entity.UserOauth{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return entity.User{}, entity.UserOauth{}, err
	}

	var userInfo linkedinOauthUser
	err = json.Unmarshal(body, &userInfo)
	if err != nil {
		return entity.User{}, entity.UserOauth{}, err
	}

	return userInfo.formatUser(), userInfo.formatOauth(token.AccessToken, token.RefreshToken), err
}
