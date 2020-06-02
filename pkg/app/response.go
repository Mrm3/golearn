package app

import (
	"github.com/gin-gonic/gin"
	e "immortality-demo/pkg/error"
)

type Response struct {
	Status  string      `json:"Status"`
	Message string      `json:"Message"`
	Detail  interface{} `json:"Detail"`
}

type Callback struct {
	CallbackUrl string `json:"CallbackUrl" validate:"url"`
	RequestId   string `json:"RequestId"`
}

type HeaderInfo struct {
	UserId    string `json:"X-User-Id" validate:"X-User-Id"`
	RequestId string `json:"RequestId"`
}

type Gin struct {
	C *gin.Context
}

// Response setting gin.JSON
func (g *Gin) Response(httpCode int, status string, message string, data interface{}) {
	g.C.JSON(httpCode, Response{
		Status:  status,
		Message: message,
		Detail:  data,
	})
}

type CommonError struct {
	Errors map[string]interface{} `json:"errors"`
}

func CallbackWithError(status string, err error, callback Callback) {
	response := Response{}
	response.Message = err.Error()
	response.Status = status
	HttpPost(callback.CallbackUrl, callback.RequestId, response)
}

func CallbackWithSuccess(detail interface{}, callback Callback) {
	response := Response{}
	response.Status = e.SUCCESS
	response.Detail = detail
	HttpPost(callback.CallbackUrl, callback.RequestId, response)
}
