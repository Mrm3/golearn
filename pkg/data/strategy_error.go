package data

import "net/http"

//lack
var ErrLackOfRequiredFieldStrategyId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrLackOfRequiredFieldStrategyId",
	Message:            "Missing required field StrategyId",
	DescriptionChinese: "缺少必要的字段 StrategyId，请查阅文档。",
}

var ErrLackOfRequiredFieldStrategyName = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrLackOfRequiredFieldStrategyName",
	Message:            "Missing required field StrategyName",
	DescriptionChinese: "缺少必要的字段 StrategyName，请查阅文档。",
}
var ErrLackOfRequiredFieldHours = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrLackOfRequiredFieldHours",
	Message:            "Missing required field Hours",
	DescriptionChinese: "缺少必要的字段 Hours，请查阅文档。",
}

var ErrLackOfRequiredFieldWeeks = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrLackOfRequiredFieldWeeks",
	Message:            "Missing required field Weeks",
	DescriptionChinese: "缺少必要的字段 Weeks，请查阅文档。",
}

var ErrLackOfRequiredFieldDuration = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrLackOfRequiredFieldDuration",
	Message:            "Missing required field Duration",
	DescriptionChinese: "缺少必要的字段 Duration，请查阅文档。",
}

var ErrLackOfRequiredFieldSnapshotQuota = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldSnapshotQuota",
	Message:            "Missing required field SnapshotQuota",
	DescriptionChinese: "缺少必要的字段SnapshotQuota，请查阅文档。",
}

//wrong value
var ErrFieldStrategyIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldStrategyIdWrongValue",
	Message:            "The length of StrategyId is too long!",
	DescriptionChinese: "字段 StrategyId 错误，请查阅文档。",
}

var ErrFieldDurationWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDurationWrongValue",
	Message:            "Invalid value for Duration",
	DescriptionChinese: "字段 Duration 错误，请查阅文档。",
}

//conflict
var ErrConflictStrategyId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictStrategyId",
	Message:            "Repetitive value for StrategyId",
	DescriptionChinese: "StrategyId 已存在，请查阅文档。",
}

var ErrConflictParametersOfStrategy = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ErrConflictParametersOfStrategy",
	Message:            "StrategyId and DiskId can not be set in the same time ",
	DescriptionChinese: "不能同时指定 StrategyId 和 DiskId。",
}

//invalid
var ErrInvalidStrategyId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidStrategyId",
	Message:            "The specified StrategyId is invalid",
	DescriptionChinese: "指定的 StrategyId 不存在",
}

var ErrInvalidStrategyStatus = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidStrategyStatus",
	Message:            "The specified StrategyStatus is invalid",
	DescriptionChinese: "指定的 StrategyStatus 不符合规范",
}

//exist
var ErrExistBindOfStrategy = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ExistBindOfStrategy",
	Message:            "The specified StrategyId has binds",
	DescriptionChinese: "指定的 StrategyId 存在绑定的硬盘",
}

var ErrExistBindOfDisk = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ExistBindOfDisk",
	Message:            "The specified DiskId has binds",
	DescriptionChinese: "指定的 DiskId 存在绑定的策略，请查阅文档",
}

var ErrInExistBindOfStrategy = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InExistBindOfStrategy",
	Message:            "The specified StrategyId has no binds",
	DescriptionChinese: "指定的 StrategyId 不存在绑定的硬盘",
}

var ErrInExistBindOfDisk = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InExistBindOfDisk",
	Message:            "The specified DiskId has no binds",
	DescriptionChinese: "指定的 DiskId 不存在绑定的策略，请查阅文档",
}

var ErrOutOfStrategyDiskQuota = BusinessLogicError{
	HttpCode:           http.StatusBadRequest,
	ErrorCode:          "UniCloudStorage-OutOfStrategyDiskQuota",
	Message:            "The specified strategy has exceeded the quota",
	DescriptionChinese: "指定策略的已绑定磁盘数量已经达到限额",
}
