package authentication

import (
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"net/http"
)

var (
	ErrHashPassword = response.NewError(http.StatusInternalServerError, "something went wrong!")

	ErrEmailAlreadyExists = response.NewError(http.StatusConflict, "email already exists")
)
