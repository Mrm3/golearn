package data

import "net/http"

var (
	ErrInvalidPropertyUnit = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-InvalidPropertyUnit",
		Message:            "The specified property unit is invalid",
		DescriptionChinese: "未知的 Unit 值，请查阅文档",
	}

	ErrInvalidPropertyCapacityOrUnit = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-InvalidPropertyCapacityOrUnit",
		Message:            "The specified property capacity or unit is invalid",
		DescriptionChinese: "错误的 Capacity 或 Unit 属性，请查阅文档",
	}

	ErrInvalidPropertyIopsOrBw = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-InvalidPropertyIopsOrBw",
		Message:            "The specified property Iops or Bw is invalid",
		DescriptionChinese: "错误的 Iops 或 Bw 属性，请查阅文档",
	}

	ErrProcessTimeout = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-ProcessTimeout",
		Message:            "Process error",
		DescriptionChinese: "处理超时",
	}

	ErrComputeProcessTimeout = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-ComputeProcessTimeout",
		Message:            "Compute process timeout",
		DescriptionChinese: "计算处理超时",
	}

	ErrComputeProcessError = BusinessLogicError{
		HttpCode:           http.StatusBadRequest,
		ErrorCode:          "UniCloudStorage-ComputeProcessError",
		Message:            "Compute process error",
		DescriptionChinese: "计算处理失败",
	}
)
