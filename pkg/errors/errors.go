package errors

import "net/http"

type AppError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func (e *AppError) Error() string {
	return e.Message
}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

func NewNotFoundError(message string) *AppError {
	if message == "" {
		message = "リソースが見つかりませんでした"
	}
	return NewAppError(http.StatusNotFound, message)
}

func NewBadRequestError(message string) *AppError {
	if message == "" {
		message = "不正なリクエストです"
	}
	return NewAppError(http.StatusBadRequest, message)
}

func NewInternalServerError(message string) *AppError {
	if message == "" {
		message = "サーバー内部でエラーが発生しました"
	}
	return NewAppError(http.StatusInternalServerError, message)
}
