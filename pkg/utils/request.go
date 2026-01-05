package utils

import (
	"emm/internal/models"
	//"github.com/gin-gonic/gin"
)

func Ok[T any](value T) models.RequestResult[T] {
	return models.RequestResult[T]{
		Result: true,
		Value:  value,
	}
}

func OkWithMessage[T any](value T, msg string) models.RequestResult[T] {
	return models.RequestResult[T]{
		Result:  true,
		Value:   value,
		Message: msg,
	}
}

func SimpleOkMessage(msg string) models.RequestResult[models.EmptyRequestResult] {
	return models.RequestResult[models.EmptyRequestResult]{
		Result:  true,
		Message: msg,
	}
}

func Err(err error, msg string) models.ErrorResult {
	resp := models.ErrorResult{
		Result: false,
		Error:  msg,
	}

	/*if gin.Mode() == gin.DebugMode && err != nil {
		resp.Details = err.Error()
	}*/

	return resp
}

func ErrWithSimpleMessage(msg string) models.ErrorResult {
	resp := models.ErrorResult{
		Result: false,
		Error:  msg,
	}

	/*if gin.Mode() == gin.DebugMode {
		resp.Details = msg
	}*/

	return resp
}
