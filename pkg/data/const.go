package data

import (
	"time"
)

const DiskSnapshotQuota = 64
const ProjectDiskQuota = 500
const ProjectImageQuota = 50
const ProjectDiskTotalSizeQuota = 32768 // GB
const VmInstanceDataDiskQuota = 15

var MaxDate = time.Unix(253402214400, 0) // UTC 9999-12-31 00:00:00

const (
	DiskStatusCreating         int8 = 1
	DiskStatusAvailable        int8 = 2
	DiskStatusDeleting         int8 = 3
	DiskStatusError            int8 = 4
	DiskStatusAttaching        int8 = 5
	DiskStatusDetaching        int8 = 6
	DiskStatusInUse            int8 = 7
	DiskStatusResetting        int8 = 8
	DiskStatusResizing         int8 = 9
	DiskStatusCreatingSnapshot int8 = 10
)

var DiskStatusMap = map[int8]string{
	1:  "Creating",
	2:  "Available",
	3:  "Deleting",
	4:  "Error",
	5:  "Attaching",
	6:  "Detaching",
	7:  "In-use",
	8:  "Resetting",
	9:  "Resizing",
	10: "CreatingSnapshot",
}

var DiskStatusMap2 = map[string]int8{
	"Creating":         1,
	"Available":        2,
	"Deleting":         3,
	"Error":            4,
	"Attaching":        5,
	"Detaching":        6,
	"In-use":           7,
	"Resetting":        8,
	"Resizing":         9,
	"CreatingSnapshot": 10,
}

var DiskStatusChineseMap = map[int8]string{
	1:  "创建中",
	2:  "可用",
	3:  "删除中",
	4:  "错误",
	5:  "挂载中",
	6:  "卸载中",
	7:  "使用中",
	8:  "恢复中",
	9:  "扩容中",
	10: "创建快照中",
}

const (
	SnapshotStatusCreating  int8 = 1
	SnapshotStatusAvailable int8 = 2
	SnapshotStatusDeleting  int8 = 3
	SnapshotStatusError     int8 = 4
)

var SnapshotStatusMap = map[int8]string{
	1: "Creating",
	2: "Available",
	3: "Deleting",
	4: "Error",
}

var SnapshotStatusMap2 = map[string]int8{
	"Creating":  1,
	"Available": 2,
	"Deleting":  3,
	"Error":     4,
}

var CategoryMap = map[string]int8{
	"hdd":        1,
	"hybrid-hdd": 2,
	"ssd":        3,
}

var CategoryMap2 = map[int8]string{
	1: "hdd",
	2: "hybrid-hdd",
	3: "ssd",
}

const (
	DiskTypeSystem int8 = 1
	DiskTypeData   int8 = 2
)

var DiskTypeMap = map[int8]string{
	1: "system",
	2: "data",
}

var DiskTypeMap2 = map[string]int8{
	"system": 1,
	"data":   2,
}

const (
	QuotaTypeProject    int8 = 1
	QuotaTypeUser       int8 = 2
	QuotaTypeVmInstance int8 = 3
	QuotaTypeDisk       int8 = 4
	QuotaTypeSnapshot   int8 = 5
	QuotaTypeImage      int8 = 6
	QuotaTypeTotalSize  int8 = 7
)

const (
	ChargeTypePrepayed   int8 = 1
	ChargeTypePayAsYouGo int8 = 2
)

var ChargeTypeMap = map[int8]string{
	1: "prepayed",
	2: "pay-as-you-go",
}

var ChargeTypeMap2 = map[string]int8{
	"prepayed":      1,
	"pay-as-you-go": 2,
}

const (
	ImageStatusCreating  int8 = 1
	ImageStatusAvailable int8 = 2
	ImageStatusError     int8 = 3
)

var ImageStatusMap = map[int8]string{
	1: "Creating",
	2: "Available",
	3: "Error",
}

var ImageStatusMap2 = map[string]int8{
	"Creating":  1,
	"Available": 2,
	"Error":     3,
}

var AvailableImageStatus = []string{
	"all",
	"Creating",
	"Available",
	"Error",
}

const (
	ImageArchitectureI386  int8 = 1
	ImageArchitectureX8664 int8 = 2
)

var ImageArchitectureMap = map[int8]string{
	1: "i386",
	2: "x86_64",
}

var ImageArchitectureMap2 = map[string]int8{
	"i386":   1,
	"x86_64": 2,
}

const (
	ImageFormatQcow2 int8 = 1
	ImageFormatRaw   int8 = 2
)

var ImageFormatMap = map[int8]string{
	1: "qcow2",
	2: "raw",
}

const (
	OSTypeLinux   int8 = 1
	OSTypeWindows int8 = 1
)

var OSTypeMap = map[int8]string{
	1: "Linux",
	2: "Windows",
}

var OSTypeMap2 = map[string]int8{
	"Linux":   1,
	"Windows": 2,
}

const (
	ResizeNotifyRetryTimes int = 3
)

const (
	DiskCategoryHDD       string = "hdd"
	DiskCategorySSD       string = "ssd"
	DiskCategoryHybridHDD string = "hybrid-hdd"
)

const (
	AttachStatusSuccess string = "success"
	AttachStatusFail    string = "fail"
)

const (
	DetachStatusSuccess string = "success"
	DetachStatusFail    string = "fail"
)

const (
	HPE3PARA              = "3par"
	HPE3PAR_MAX_DISK_SIZE = 65536 // 64T
	HPE3PAR_MAX_ID_LEN    = 27
	CEPH                  = "ceph"
)

const (
	Hpe3parBwMinGoalKB = 1
	Hpe3parIoMinGoal   = 1
)

const (
	THREE_PAR_STATUS_EXPORTED  string = "exported"
	THREE_PAR_STATUS_EXPORTING string = "exporting"
)

const (
	IMAGE_TYPE_PUBLIC = "public"
	IMAGE_TYPE_CUSTOM = "custom"
)

var InstanceCodeToDiskCategory = map[string]string{
	"ebs.highIO.ssd": DiskCategorySSD,
	"ebs.hybrid.hdd": DiskCategoryHybridHDD,
}

var DiskCategoryToInstanceCode = map[string]string{
	DiskCategorySSD:       "ebs.highIO.ssd",
	DiskCategoryHybridHDD: "ebs.hybrid.hdd",
}

const (
	DatetimeFormatString = "2006-01-02 15:04:05"
)

const (
	JobDispatchQueue = "IMMORTALITY_ASYNCHRONOUS_JOBS"
)

const (
	DISK_IS_SHARE_YES = 1
	DISK_IS_SHARE_NO  = 2
)

var DiskIsShareMap = map[int8]string{
	1: "yes",
	2: "no",
}

var DiskIsShareMap2 = map[string]int8{
	"yes": 1,
	"no":  2,
}

const (
	ExportStatusExporting    = 1
	ExportStatusExported     = 2
	ExportStatusExportFail   = 3
	ExportStatusUnExportFail = 4
	ExportStatusUnExported   = 5
)

var ExportStatusMap = map[int]string{
	1: "exporting",
	2: "exported",
	3: "exportFailed",
	4: "unExportFailed",
	5: "unExported",
}

var ExportStatusMap2 = map[string]int8{
	"exporting":      1,
	"exported":       2,
	"exportFailed":   3,
	"UnExportFailed": 4,
	"UnExported":     5,
}

const (
	HEP_3PAR_SPECIAL_LUN = 254
)

const (
	SpatialUserPublic = "public"
)

const (
	MaxAttempts  = 3
	AttemptDelay = time.Duration(time.Second)
)

const (
	KafkaImmortalityTopic = "unicloud_monitor_immortality"
)
