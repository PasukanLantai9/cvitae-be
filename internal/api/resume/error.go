package resume

import (
	"github.com/bccfilkom/career-path-service/internal/pkg/response"
	"net/http"
)

var (
	ErrObjectIDNotProvided = response.NewError(http.StatusUnprocessableEntity, "object ID not provided")
	ErrIncorrectObjectID   = response.NewError(http.StatusBadRequest, "incorrect object ID")

	ErrPDFFileNotProvided = response.NewError(http.StatusUnprocessableEntity, "pdf file not provided")
	ErrInvalidResumeFile  = response.NewError(http.StatusBadRequest, "invalid resume file")
)
