package data

import (
	"encoding/json"
	"net/http"
)

var ErrCanNotDownSize = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-CanNotDownSize",
	Message:            "You could only enlarge your disk",
	DescriptionChinese: "只能增大硬盘，以防数据丢失",
}

var ErrDiskExpired = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskExpired",
	Message:            "The specified disk is expired",
	DescriptionChinese: "指定的云硬盘已过期",
}

var ErrDiskFromDiskSnapshotExists = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskFromDiskSnapshotExists",
	Message:            "There're disks created from this disk's snapshots",
	DescriptionChinese: "存在由此硬盘的快照创建的云硬盘",
}

var ErrDiskFromImageExists = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskFromImageExists",
	Message:            "There're disks created from this image",
	DescriptionChinese: "存在由此镜像创建的云硬盘",
}

var ErrDiskFromSnapshotExists = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskFromSnapshotExists",
	Message:            "There're disks created from this snapshot",
	DescriptionChinese: "存在由此快照创建的云硬盘",
}

var ErrDiskNotSystem = BusinessLogicError{
	HttpCode:           502,
	ErrorCode:          "UniCloudStorage-DiskNotSystem",
	Message:            "Cannot create image from data disk",
	DescriptionChinese: "指定的 DiskId 不是系统盘",
}

var ErrDiskQuotaExceeded = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskQuotaExceeded",
	Message:            "Quota of disk per project(500) exceeded",
	DescriptionChinese: "每用户最多创建500个硬盘，如有特殊需求请提交工单",
}

var ErrDiskReInitBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskReInitBadStatus",
	Message:            "The status of the disk should be In-use when re-init",
	DescriptionChinese: "恢复操作时硬盘状态必须为In-use",
}

var ErrDiskResettingBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskResettingBadStatus",
	Message:            "The status of the disk should be In-use/Available when resetting",
	DescriptionChinese: "恢复操作时硬盘状态必须为In-use/Available",
}

var ErrDiskDataResettingBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskDataResettingBadStatus",
	Message:            "The status of the data disk should be Available when resetting",
	DescriptionChinese: "恢复操作时数据盘状态必须为Available",
}

var ErrDiskExportAndCancelExportBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskExportAndCancelExportBadStatus",
	Message:            "Disk can only be exported when it is 3par",
	DescriptionChinese: "只有3par设备才可以被导出和取消导出",
}

var ErrDiskResizingBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskResizingBadStatus",
	Message:            "Disk can only be resized when Available or In-use",
	DescriptionChinese: "云盘只能在 Available 或 In-use 状态才可以扩容",
}

var ErrDiskExportBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskExportBadStatus",
	Message:            "Disk has been exported",
	DescriptionChinese: "此云盘在此cvk上已经被导出",
}

var ErrDiskBad3ParInfo = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskBad3parInfo",
	Message:            "Disk has no 3par info",
	DescriptionChinese: "该磁盘没有3par信息",
}

var ErrDiskExportingBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskExportingBadStatus",
	Message:            "Disk can only be exported when Available or In-use",
	DescriptionChinese: "云盘只能在 Available 或 In-use 状态才可以被导出",
}

var ErrDiskTotalSizeQuotaExceeded = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskTotalSizeQuotaExceeded",
	Message:            "Quota of total disk size per project(32768 GB) exceeded",
	DescriptionChinese: "用户硬盘配额为32768 GB，如有特殊需求请提交工单",
}

var ErrFieldActionWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldActionWrongValue",
	Message:            "Invalid value for Action",
	DescriptionChinese: "字段Action错误，请查阅文档。",
}

var ErrDiskCancelExportBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskCancelExportBadStatus",
	Message:            "Disk has not been exported",
	DescriptionChinese: "此云盘在此cvk上不是导出状态",
}
var ErrInvalidUser = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-DiskWrongUser",
	Message:            "user invalid",
	DescriptionChinese: "请求的磁盘不属于该用户",
}

var ErrFieldArchitectureWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldArchitectureWrongValue",
	Message:            "Invalid value for Architecture",
	DescriptionChinese: "字段Architecture错误，请查阅文档。",
}

var ErrFieldCategoryWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldCategoryWrongValue",
	Message:            "Invalid value for Category",
	DescriptionChinese: "字段Category错误，请查阅文档。",
}

var ErrFieldClusterIDWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldClusterIDWrongValue",
	Message:            "Invalid value for ClusterID",
	DescriptionChinese: "字段ClusterID错误，请查阅文档。",
}

var ErrFieldConfigWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldConfigWrongValue",
	Message:            "Invalid value for Config",
	DescriptionChinese: "字段Config错误，请查阅文档。",
}

var ErrFieldCountWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldCountWrongValue",
	Message:            "Invalid value for Count",
	DescriptionChinese: "字段Count错误，请查阅文档。",
}

var ErrFieldDescriptionWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDescriptionWrongValue",
	Message:            "Invalid value for Description",
	DescriptionChinese: "字段Description错误，请查阅文档。",
}

var ErrFieldDiskCreatedWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskCreatedWrongValue",
	Message:            "Invalid value for DiskCreated",
	DescriptionChinese: "字段DiskCreated错误，请查阅文档。",
}

var ErrFieldDiskIdsWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskIdsWrongValue",
	Message:            "Invalid value for DiskIds",
	DescriptionChinese: "字段DiskIds错误，请查阅文档。",
}

var ErrFieldDiskQuotaWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskQuotaWrongValue",
	Message:            "Invalid value for DiskQuota",
	DescriptionChinese: "字段DiskQuota错误，请查阅文档。",
}

var ErrFieldDiskSinglePurchaseQuotaWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskSinglePurchaseQuotaWrongValue",
	Message:            "Invalid value for DiskSinglePurchaseQuota",
	DescriptionChinese: "字段DiskSinglePurchaseQuota错误，请查阅文档。",
}

var ErrFieldDiskTotalSizeQuotaWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDiskTotalSizeQuotaWrongValue",
	Message:            "Invalid value for DiskTotalSizeQuota",
	DescriptionChinese: "字段DiskTotalSizeQuota错误，请查阅文档。",
}

var ErrFieldDisksWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldDisksWrongValue",
	Message:            "Invalid value for DisksReq",
	DescriptionChinese: "字段Disks错误，请查阅文档。",
}

var ErrFieldEndTimeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldEndTimeWrongValue",
	Message:            "Invalid value for EndTime",
	DescriptionChinese: "字段EndTime错误，请查阅文档。",
}

var ErrFieldFilterWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldFilterWrongValue",
	Message:            "Invalid value for Filter",
	DescriptionChinese: "字段Filter错误，请查阅文档。",
}

var ErrFieldImageIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldImageIdWrongValue",
	Message:            "Invalid value for ImageId",
	DescriptionChinese: "字段ImageId错误，请查阅文档。",
}

var ErrFieldImageNameWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldImageNameWrongValue",
	Message:            "Invalid value for ImageName",
	DescriptionChinese: "字段ImageName错误，请查阅文档。",
}

var ErrFieldImagesWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldImagesWrongValue",
	Message:            "Invalid value for Images",
	DescriptionChinese: "字段Images错误，请查阅文档。",
}

var ErrFieldInstanceIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldInstanceIdWrongValue",
	Message:            "Invalid value for InstanceId",
	DescriptionChinese: "字段InstanceId错误，请查阅文档。",
}

var ErrFieldInstancesWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldInstancesWrongValue",
	Message:            "Invalid value for Instances",
	DescriptionChinese: "字段Instances错误，请查阅文档。",
}

var ErrFieldIqnWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldIqnWrongValue",
	Message:            "Invalid value for Iqn",
	DescriptionChinese: "字段Iqn错误，请查阅文档。",
}

var ErrFieldNameWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldNameWrongValue",
	Message:            "Invalid value for Name",
	DescriptionChinese: "字段Name错误，请查阅文档。",
}

var ErrFieldNewSizeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldNewSizeWrongValue",
	Message:            "Invalid value for NewSize",
	DescriptionChinese: "字段NewSize错误，请查阅文档。",
}

var ErrFieldOSTypeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldOSTypeWrongValue",
	Message:            "Invalid value for OSType",
	DescriptionChinese: "字段OSType错误，请查阅文档。",
}

var ErrFieldOperatingSystemWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldOperatingSystemWrongValue",
	Message:            "Invalid value for OperatingSystem",
	DescriptionChinese: "字段OperatingSystem错误，请查阅文档。",
}

var ErrFieldPayTypeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldPayTypeWrongValue",
	Message:            "Invalid value for PayType",
	DescriptionChinese: "字段PayType错误，请查阅文档。",
}

var ErrFieldQosWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldQosWrongValue",
	Message:            "Invalid value for Qos",
	DescriptionChinese: "字段Qos错误，请查阅文档。",
}

var ErrFieldRequestIdWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldRequestIdWrongValue",
	Message:            "Invalid value for RequestId",
	DescriptionChinese: "字段RequestId错误，请查阅文档。",
}

var ErrFieldResultWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldResultWrongValue",
	Message:            "Invalid value for Result",
	DescriptionChinese: "字段Result错误，请查阅文档。",
}

var ErrFieldStartTimeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldStartTimeWrongValue",
	Message:            "Invalid value for StartTime",
	DescriptionChinese: "字段StartTime错误，请查阅文档。",
}

var ErrFieldStatusWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldStatusWrongValue",
	Message:            "Invalid value for Status",
	DescriptionChinese: "字段Status错误，请查阅文档。",
}

var ErrFieldStepWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldStepWrongValue",
	Message:            "Invalid value for Step",
	DescriptionChinese: "字段Step错误，请查阅文档。",
}

var ErrFieldThreeParDataIPWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldThreeParDataIPWrongValue",
	Message:            "Invalid value for ThreeParDataIP",
	DescriptionChinese: "字段ThreeParDataIP错误，请查阅文档。",
}

var ErrFieldThreeParLunWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldThreeParLunWrongValue",
	Message:            "Invalid value for ThreeParLun",
	DescriptionChinese: "字段ThreeParLun错误，请查阅文档。",
}

var ErrFieldThreeParWWNWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldThreeParWWNWrongValue",
	Message:            "Invalid value for ThreeParWWN",
	DescriptionChinese: "字段ThreeParWWN错误，请查阅文档。",
}

var ErrFieldTotalCountWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldTotalCountWrongValue",
	Message:            "Invalid value for TotalCount",
	DescriptionChinese: "字段TotalCount错误，请查阅文档。",
}

var ErrFieldVolumeIDWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldVolumeIDWrongValue",
	Message:            "Invalid value for VolumeID",
	DescriptionChinese: "字段VolumeID错误，请查阅文档。",
}

var ErrFieldVolumeSizeWrongValue = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-FieldVolumeSizeWrongValue",
	Message:            "Invalid value for VolumeSize",
	DescriptionChinese: "字段VolumeSize错误，请查阅文档。",
}

var ErrImageQuotaExceeded = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-ImageQuotaExceeded",
	Message:            "Quota of image per project(50) exceeded",
	DescriptionChinese: "每个用户最多创建50个镜像，如有特殊需求请提交工单。",
}

var ErrInitSourceNotFound = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InitSourceNotFound",
	Message:            "The image or snapshot used to create the disk has been removed",
	DescriptionChinese: "创建硬盘所使用的镜像或快照已删除",
}

var ErrInstanceDiskAttachmentQuotaExceeded = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-InstanceDiskAttachmentQuotaExceeded",
	Message:            "Quota of data disk per VM instance(15) exceeded",
	DescriptionChinese: "每个虚拟机实例最多挂载15个硬盘",
}

var ErrInstanceReInitBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-InstanceReInitBadStatus",
	Message:            "The status of instance should be stopped",
	DescriptionChinese: "恢复操作时虚机实例必须关机",
}

var ErrInstanceResetBadStatus = BusinessLogicError{
	HttpCode:           403,
	ErrorCode:          "UniCloudStorage-InstanceResetBadStatus",
	Message:            "The status of instance should be stopped",
	DescriptionChinese: "回滚操作时虚机实例必须关机",
}

var ErrInvalidExport = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidExport",
	Message:            "The specified Export is not exist",
	DescriptionChinese: "指定的 export 不存在",
}

var ErrInvalidImageId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidImageId",
	Message:            "The specified ImageId is invalid",
	DescriptionChinese: "指定的 ImageId 不存在",
}

var ErrInvalidInstanceId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidInstanceId",
	Message:            "The specified InstanceId is invalid",
	DescriptionChinese: "指定的 InstanceId 不存在",
}

var ErrInvalidIQN = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidIQN",
	Message:            "The specified IQN is invalid",
	DescriptionChinese: "指定的 IQN 不存在",
}

var ErrInvalidSize = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidSize",
	Message:            "The specified size is invalid",
	DescriptionChinese: "指定的云硬盘 Size 不合法。",
}

var ErrLackOfRequiredFieldAction = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldAction",
	Message:            "Missing required field Action",
	DescriptionChinese: "缺少必要的字段Action，请查阅文档。",
}

var ErrLackOfRequiredFieldImageId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldImageId",
	Message:            "Missing required field ImageId",
	DescriptionChinese: "缺少必要的字段ImageId，请查阅文档。",
}

var ErrLackOfRequiredFieldInstanceId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldInstanceId",
	Message:            "Missing required field InstanceId",
	DescriptionChinese: "缺少必要的字段InstanceId，请查阅文档。",
}

var ErrLackOfRequiredFieldIqn = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldIqn",
	Message:            "Missing required field Iqn",
	DescriptionChinese: "缺少必要的字段Iqn，请查阅文档。",
}

var ErrLackOfRequiredFieldCvkName = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldCvkName",
	Message:            "Missing required field CvkName",
	DescriptionChinese: "缺少必要的字段CvkName，请查阅文档。",
}

var ErrLackOfRequiredFieldName = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldName",
	Message:            "Missing required field Name",
	DescriptionChinese: "缺少必要的字段Name，请查阅文档。",
}

var ErrLackOfRequiredFieldNewSize = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldNewSize",
	Message:            "Missing required field NewSize",
	DescriptionChinese: "缺少必要的字段NewSize，请查阅文档。",
}

var ErrLackOfRequiredFieldPodId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldPodId",
	Message:            "Missing required field PodId",
	DescriptionChinese: "缺少必要的字段PodId，请查阅文档。",
}

var ErrLackOfRequiredFieldStatus = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-LackOfRequiredFieldStatus",
	Message:            "Missing required field Status",
	DescriptionChinese: "缺少必要的字段Status，请查阅文档。",
}

var ErrZoneNotAvailable = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-ZoneNotAvailable",
	Message:            "The specified zone does't exist or not providing requested service",
	DescriptionChinese: "指定的可用区不存在或者不提供此服务",
}

var ErrJsonDecode = BusinessLogicError{
	HttpCode:           500,
	ErrorCode:          "UniCloudStorage-DecodeJsonError",
	Message:            "Decode Json Error",
	DescriptionChinese: "Json Decode 出错",
}

var ErrJsonEncode = BusinessLogicError{
	HttpCode:           500,
	ErrorCode:          "UniCloudStorage-EncodeJsonError",
	Message:            "Encode Json Error",
	DescriptionChinese: "Json Encode 出错",
}

var ErrCvkRequired = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-export cvk nil",
	Message:            "Need cvk name for export 3para volume",
	DescriptionChinese: "3par volume 导出需要CVK的名字",
}

var ErrInvalidLengthOfDiskId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidLengthOfDiskId",
	Message:            "The specified DiskId is invalid",
	DescriptionChinese: "指定的 DiskId 长度不符合规范，请查阅文档",
}

var ErrInvalidLengthOfImageId = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-InvalidLengthOfImageId",
	Message:            "The specified ImageId is invalid",
	DescriptionChinese: "指定的 ImageId 长度不符合规范，请查阅文档",
}

var ErrUnauthorized = BusinessLogicError{
	HttpCode:           http.StatusUnauthorized,
	ErrorCode:          "UniCloudStorage-Unauthorized",
	Message:            "Unauthorized user",
	DescriptionChinese: "用户认证失败",
}

var ErrOpNotSupport = BusinessLogicError{
	HttpCode:           400,
	ErrorCode:          "UniCloudStorage-OperatonNotSupport",
	Message:            "Operatoin not support",
	DescriptionChinese: "操作不支持",
}

type BusinessLogicError struct {
	HttpCode           int
	ErrorCode          string
	Message            string
	DescriptionChinese string
}

type BusinessLogicErrorJson struct {
	ErrorCode          string `json:"ErrorCode"`
	Message            string `json:"Message"`
	RequestID          string `json:"RequestId"`
	DescriptionChinese string `json:"DescriptionChinese"`
}

func (e BusinessLogicError) Error() string {
	return e.Message
}

func (e BusinessLogicError) WriteResponse(w http.ResponseWriter,
	requestID string) (err error) {

	w.WriteHeader(e.HttpCode)
	logicError := BusinessLogicErrorJson{
		ErrorCode:          e.ErrorCode,
		Message:            e.Message,
		RequestID:          requestID,
		DescriptionChinese: e.DescriptionChinese,
	}
	bytes, err := json.Marshal(logicError)
	if err != nil {
		return err
	}
	_, err = w.Write(bytes)
	return err
}
