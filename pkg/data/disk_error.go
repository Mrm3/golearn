package data

import "net/http"

// lack
var ErrLackOfRequiredField = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredField",
	Message:            "Missing required field SnapshotId or ImageId",
	DescriptionChinese: "系统盘需要指定 SnapshotId 或 ImageId,请查阅文档。",
}

var ErrLackOfRequiredFieldDiskId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldDiskId",
	Message:            "Missing required field DiskId",
	DescriptionChinese: "缺少必要的字段 DiskId，请查阅文档。",
}

var ErrLackOfRequiredFieldDiskIds = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldDiskIds",
	Message:            "Missing required field DiskIds",
	DescriptionChinese: "缺少必要的字段 DiskIds，请查阅文档。",
}

var ErrLackOfRequiredFieldDiskName = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldDiskName",
	Message:            "Missing required field DiskName",
	DescriptionChinese: "缺少必要的字段 DiskName，请查阅文档。",
}

var ErrLackOfRequiredFieldInstanceCodeOrDiskCategory = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldInstanceCodeOrDiskCategory",
	Message:            "Missing required field InstanceCode or DiskCategory",
	DescriptionChinese: "缺少必要的字段 InstanceCode 或 DiskCategory，请查阅文档。",
}

var ErrLackOfRequiredFieldDiskCategory = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldDiskCategory",
	Message:            "Missing required field DiskCategory",
	DescriptionChinese: "缺少必要的字段 DiskCategory，请查阅文档。",
}

var ErrLackOfRequiredFieldDiskType = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldDiskType",
	Message:            "Missing required field DiskType",
	DescriptionChinese: "缺少必要的字段 DiskType，请查阅文档。",
}

var ErrLackOfRequiredFieldStorageType = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldStorageType",
	Message:            "Missing required field StorageType",
	DescriptionChinese: "缺少必要的字段StorageType，请查阅文档。",
}

var ErrLackOfRequiredFieldPageNumber = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldPageNumber",
	Message:            "Missing required field PageNumber",
	DescriptionChinese: "缺少必要的字段 PageNumber，请查阅文档。",
}

var ErrLackOfRequiredFieldPageSize = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldPageSize",
	Message:            "Missing required field PageSize ",
	DescriptionChinese: "缺少必要的字段 PageSize，请查阅文档。",
}

// wrong value
var ErrFieldDiskIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskIdWrongValue",
	Message:            "Invalid value for DiskId",
	DescriptionChinese: "字段 DiskId 错误，请查阅文档。",
}

var ErrFieldDiskNameWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskNameWrongValue",
	Message:            "Invalid value for DiskName",
	DescriptionChinese: "字段DiskName错误，请查阅文档。",
}

var ErrFieldStorageTypeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldStorageTypeWrongValue",
	Message:            "Invalid value for StorageType",
	DescriptionChinese: "字段 StorageType 错误，请查阅文档。",
}

var ErrFieldInstanceCodeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldInstanceCodeWrongValue",
	Message:            "Invalid value for InstanceCode",
	DescriptionChinese: "字段 InstanceCode 错误，请查阅文档。",
}

var ErrFieldDiskCategoryWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskCategoryWrongValue",
	Message:            "Invalid value for DiskCategory",
	DescriptionChinese: "字段 DiskCategory 错误，请查阅文档。",
}

var ErrFieldDiskTypeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskTypeWrongValue",
	Message:            "Invalid value for DiskType",
	DescriptionChinese: "字段 DiskType 错误，请查阅文档。",
}

var ErrFieldDiskIsShareWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskIsShareWrongValue",
	Message:            "Invalid value for IsShare",
	DescriptionChinese: "字段 IsShare 错误，请查阅文档。",
}

var ErrFieldSizeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldSizeWrongValue",
	Message:            "Invalid value for Size",
	DescriptionChinese: "字段 Size 错误，请查阅文档。",
}

var ErrFieldZoneIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldZoneIdWrongValue",
	Message:            "Invalid value for ZoneId",
	DescriptionChinese: "字段 ZoneId 错误，请查阅文档。",
}

var ErrFieldRegionIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldRegionIdWrongValue",
	Message:            "Invalid value for RegionId",
	DescriptionChinese: "字段 RegionId 错误，请查阅文档。",
}

var ErrFieldPodIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldPodIdWrongValue",
	Message:            "Invalid value for PodId",
	DescriptionChinese: "字段 PodId 错误，请查阅文档。",
}

var ErrFieldPageNumberWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldPageNumberWrongValue",
	Message:            "Invalid value for PageNumber",
	DescriptionChinese: "字段PageNumber错误，请查阅文档。",
}

var ErrFieldPageSizeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldPageSizeWrongValue",
	Message:            "Invalid value for PageSize",
	DescriptionChinese: "字段PageSize错误，请查阅文档。",
}

// conflict
var ErrConflictDiskId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictDiskId",
	Message:            "Repetitive value for DiskId",
	DescriptionChinese: "DiskId 已存在，请查阅文档。",
}

var ErrConflictDiskIds = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictDiskIds",
	Message:            "Repetitive value for DiskIds",
	DescriptionChinese: "参数中包含重复 DiskId，请查阅文档。",
}

var ErrConflictParameters = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictParameters",
	Message:            "SnapshotId and ImageId cannot be set at the same time.",
	DescriptionChinese: "不能同时指定 SnapshotId 和 ImageId。",
}

var ErrConflictSnapshotId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictSnapshotId",
	Message:            "Repetitive value for SnapshotId",
	DescriptionChinese: "SnapshotId 已存在，请查阅文档。",
}

var ErrConflictIsShare = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ConflictIsShare",
	Message:            "SystemDisk can't been set yes for IsShare",
	DescriptionChinese: "系统盘不能被共享。",
}

//status err
var ErrBadDiskDeletionStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-BadDiskDeletionStatus",
	Message:            "Disk can only be deleted when available",
	DescriptionChinese: "只能在云盘处于Available或使用中状态时可以删除",
}

var ErrDiskAttachBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskAttachBadStatus",
	Message:            "Disk can only be attached when Available",
	DescriptionChinese: "硬盘只有在可用状态才可以挂载",
}

var ErrDiskDetachBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskDetachBadStatus",
	Message:            "Disk can only be detached when In-use",
	DescriptionChinese: "硬盘只有在使用中状态才可以卸载",
}

var ErrNotifyDiskAttachBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-NotifyDiskAttachBadStatus",
	Message:            "Disk can only be notified attached when attaching",
	DescriptionChinese: "硬盘只有在正在挂载状态下才可以通知挂载成功",
}

var ErrNotifyDiskDetachBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-NotifyDiskDetachBadStatus",
	Message:            "Disk can only be notified detached when detaching",
	DescriptionChinese: "硬盘只有在正在卸载状态下才可以通知卸载成功",
}

//other
var ErrInvalidDiskType = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidDiskType",
	Message:            "Cannot detach system disk",
	DescriptionChinese: "不能卸载系统盘",
}

var ErrInvalidDiskId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidDiskId",
	Message:            "The specified DiskId is invalid",
	DescriptionChinese: "指定的 DiskId 不存在",
}

var ErrInvalidQos = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidQos",
	Message:            "The specified Qos is invalid",
	DescriptionChinese: "指定的云硬盘 Qos 不合法。",
}

var ErrInvalidTimeString = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidTimeString",
	Message:            "The specified time string is invalid",
	DescriptionChinese: "指定的时间无效",
}

var ErrInvalidStartTimeOrEndTime = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidStartTimeOrEndTime",
	Message:            "The specified start time or end time is invalid",
	DescriptionChinese: "指定的计量时间无效",
}

var ErrIsInTheMeasurement = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-IsInTheMeasurement",
	Message:            "The specified disk is in the measurement",
	DescriptionChinese: "指定磁盘已经在计量中",
}

var ErrIsNotInTheMeasurement = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-IsNotInTheMeasurement",
	Message:            "The specified disk is not in the measurement",
	DescriptionChinese: "指定磁盘未计量",
}

var ErrInvalidInstanceCode = BusinessLogicError{
	HttpCode:           http.StatusBadRequest,
	ErrorCode:          "UniCloudStorage-InvalidInstanceCode",
	Message:            "The specified InstanceCode is invalid",
	DescriptionChinese: "错误的 InstanceCode 属性，请查阅文档",
}

var ErrAttachInformationNotExists = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-AttachInformationNotExists",
	Message:            "The specified disk and instance attachInformation does not exists",
	DescriptionChinese: "指定磁盘和主机的挂载信息不存在",
}
