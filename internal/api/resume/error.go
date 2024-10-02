package resume

import (
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"net/http"
)

var (
	ErrIncorrectObjectID = response.NewError(http.StatusBadRequest, "incorrect object id")
)
