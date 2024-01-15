package errors

import "net/http"

var (
	ErrAuth       = &APIError{status: http.StatusUnauthorized, msg: "invalid token"}
	ErrNotFound   = &APIError{status: http.StatusNotFound, msg: "not found"}
	ErrTechnical  = &APIError{status: http.StatusInternalServerError, msg: "technical error"}
	ErrBadRequest = &APIError{status: http.StatusBadRequest, msg: "bad request"}
	ErrBusiness   = &APIError{status: http.StatusConflict, msg: "business error"}
)

type APIError struct {
	status int    `json:"status"`
	msg    string `json:"msg"`
}

func (e APIError) Error() string {
	return e.msg
}

func (e APIError) APIError() (int, string) {
	return e.status, e.msg
}

func (e APIError) GetCode() int {
	return e.status
}
