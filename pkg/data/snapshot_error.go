package data

var (
	ErrFieldSnapshotIdWrongValue = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-FieldSnapshotIdWrongValue",
		Message:            "Invalid value for SnapshotId",
		DescriptionChinese: "字段SnapshotId错误，请查阅文档。",
	}

	ErrFieldSnapshotIdsWrongValue = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-FieldSnapshotIdsWrongValue",
		Message:            "Invalid value for SnapshotIds",
		DescriptionChinese: "字段SnapshotIds错误，请查阅文档。",
	}

	ErrFieldSnapshotNameWrongValue = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-FieldSnapshotNameWrongValue",
		Message:            "Invalid value for SnapshotName",
		DescriptionChinese: "字段SnapshotName错误，请查阅文档。",
	}

	ErrFieldSnapshotsWrongValue = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-FieldSnapshotsWrongValue",
		Message:            "Invalid value for Snapshots",
		DescriptionChinese: "字段Snapshots错误，请查阅文档。",
	}

	ErrInvalidSnapshotId = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-InvalidSnapshotId",
		Message:            "The specified SnapshotId is invalid",
		DescriptionChinese: "指定的 SnapshotId 不存在",
	}

	ErrLackOfRequiredFieldSnapshotId = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-LackOfRequiredFieldSnapshotId",
		Message:            "Missing required field SnapshotId",
		DescriptionChinese: "缺少必要的字段SnapshotId，请查阅文档。",
	}

	ErrLackOfRequiredFieldSnapshotIds = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-LackOfRequiredFieldSnapshotIds",
		Message:            "Missing required field SnapshotIds",
		DescriptionChinese: "缺少必要的字段SnapshotIds，请查阅文档。",
	}

	ErrSnapshotNotBelongToDisk = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-SnapshotNotBelongToDisk",
		Message:            "The specified snapshot is incompatible with disk",
		DescriptionChinese: "快照与硬盘不匹配",
	}

	ErrSnapshotQuotaExceeded = BusinessLogicError{
		HttpCode:           403,
		ErrorCode:          "UniCloudStorage-SnapshotQuotaExceeded",
		Message:            "Quota of snapshot per disk(64) exceeded",
		DescriptionChinese: "一块硬盘最多创建64份快照",
	}

	ErrInvalidLengthOfSnapshotId = BusinessLogicError{
		HttpCode:           400,
		ErrorCode:          "UniCloudStorage-InvalidLengthOfSnapshotId",
		Message:            "The specified SnapshotId is invalid",
		DescriptionChinese: "指定的SnapshotId不符合规范，请查阅文档",
	}

	ErrNoSnapshotsBelongToUser = BusinessLogicError{
		HttpCode:           403,
		ErrorCode:          "UniCloudStorage-NoSnapshotsBelongToUser",
		Message:            "No snapshots belong to specified user",
		DescriptionChinese: "没有快照属于指定的用户",
	}

	ErrNoSnapshotsBelongToDisk = BusinessLogicError{
		HttpCode:           403,
		ErrorCode:          "UniCloudStorage-NoSnapshotsBelongToDisk",
		Message:            "No snapshots belong to specified disk",
		DescriptionChinese: "没有快照属于指定的磁盘",
	}

	ErrSnapshotDoesNotBelongToUser = BusinessLogicError{
		HttpCode:           403,
		ErrorCode:          "UniCloudStorage-SnapshotDoesNotBelongToUser",
		Message:            "Snapshot does not belong to specified user",
		DescriptionChinese: "快照不属于指定的用户",
	}
)
