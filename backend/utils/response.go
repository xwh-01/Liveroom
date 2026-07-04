package utils

import "net/http"

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Success(data interface{}) Response {
	return Response{Code: 0, Msg: "ok", Data: data}
}

func Error(code int, msg string) Response {
	return Response{Code: code, Msg: msg}
}

var (
	ErrBadRequest = Error(http.StatusBadRequest, "bad request")
	ErrInternal   = Error(http.StatusInternalServerError, "internal server error")
)
