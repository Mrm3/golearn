package data

import "net/http"

var ErrServerInternal = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-InternalServerError",
	Message:            "Internal server error",
	DescriptionChinese: "服务器内部错误.",
}

var ErrServerInternalDB = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-InternalServerError",
	Message:            "Internal database visit error",
	DescriptionChinese: "服务器内访问数据库错误.",
}

var ErrServerInternalMQ = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-InternalServerError",
	Message:            "Internal rabbitmq visit error",
	DescriptionChinese: "服务器内部访问消息队列错误.",
}

var ErrServerInternalDriver = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-InternalDriver",
	Message:            "Internal storage driver error",
	DescriptionChinese: "服务器内部存储驱动错误.",
}

var ErrServerInternalHTTP = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-InternalHTTP",
	Message:            "Internal new request error",
	DescriptionChinese: "服务器内部 HTTP 错误.",
}

var ErrServerExternalHTTP = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "ExternalHTTP",
	Message:            "Visit external Service error",
	DescriptionChinese: "服务器访问外部服务错误.",
}

var ErrServerInternalUUID = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-ServerInternalUUID",
	Message:            "Generate uuid error",
	DescriptionChinese: "服务器内部错误",
}

var ErrServerInternalCache = BusinessLogicError{
	HttpCode:           http.StatusInternalServerError,
	ErrorCode:          "UniCloudStorage-ServerInternalCache",
	Message:            "Internal cache error",
	DescriptionChinese: "服务器内部 Cache 错误.",
}