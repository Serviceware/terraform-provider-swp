package aipe

import (
	"context"
	"net/http"

	"github.com/hashicorp/errwrap"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

type ApiError struct {
	StatusCode int
	Message    string
}

func (e *ApiError) Error() string {
	return e.Message
}

func ErrorIsNotFound(err error) bool {
	apiError, ok := errwrap.GetType(err, &ApiError{}).(*ApiError)

	tflog.Info(context.TODO(), "ErrorIsNotFound", map[string]interface{}{"apiError": apiError, "ok": ok})

	return ok && apiError != nil && (apiError.StatusCode == http.StatusNotFound || apiError.StatusCode == http.StatusGone)
}
