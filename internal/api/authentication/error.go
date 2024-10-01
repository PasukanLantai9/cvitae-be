package authentication

import (
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"net/http"
)

var (
	ErrHashPassword                 = response.NewError(http.StatusInternalServerError, "something went wrong!")
	ErrEmailOrPasswordWreong        = response.NewError(http.StatusBadRequest, "oops!, email or password is wrong")
	ErrUnableToCreateNewAccessToken = response.NewError(http.StatusUnauthorized, "oops!, its look wrong sessionID or refreshToken expired")

	ErrEmailAlreadyExists          = response.NewError(http.StatusConflict, "email already exists")
	ErrUserWithEmailNotFound       = response.NewError(http.StatusNotFound, "user with email not found")
	ErrUserWithIDNotFound          = response.NewError(http.StatusNotFound, "user with ID not found")
	ErrSessionWithIDNotFound       = response.NewError(http.StatusNotFound, "session with ID not found")
	ErrUserWithOauthUserIDNotFound = response.NewError(http.StatusNotFound, "user with oauth user ID not found")

	ErrInvalidOuthProvider = response.NewError(http.StatusBadRequest, "invalid oauth provider specified")
)
