package data

import (
	"time"
)

type CreateDiskRequest struct {
	RequestId    string `json:"request_id"`
	DiskId       string `json:"disk_id"`
	DiskCategory string `json:"disk_category"`
	SnapshotId   string `json:"snapshot_id"`
	ImageId      string `json:"image_id"`
	Size         uint64 `json:"size"` // in byte
	StorageType  string `json:"storage_type"`
	Qos          string `json:"qos"`
	UserId       string `json:"user_id"`
	ScheduleInfo string `json:"schedule_info"`
	ImageType    string
}

type CreateDisksRequest struct {
	//DeliveryVersion string
	RequestId      string
	DeliveryID     string
	DeliveryUnitId string
	StartAt        time.Time
	//UserId          string
	//ServiceType     string
	//Region          string
	DisksReq []CreateDiskRequest
}

type CreateDiskResponse struct {
	Wwn string
}

type DeleteDiskRequest struct {
	RequestId    string `json:"request_id"`
	DiskId       string `json:"disk_id"`
	DiskCategory string `json:"disk_category"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type DeleteDisksRequest struct {
	RequestId      string `json:"request_id"`
	DeliveryID     string `json:"delivery_id"`
	DeliveryUnitID string `json:"delivery_unit_id"`
	DeleteAt       time.Time
	DisksInfo      []DeleteDiskRequest
}

type CreateImageRequest struct {
	RequestId    string `json:"request_id"`
	SourceDiskId string `json:"source_disk_id"`
	DiskCategory string `json:"disk_category"`
	ImageId      string `json:"image_id"`
	Size         uint64 `json:"size"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"` // for ceph, it is ceph cluster_id ;for 3par , it is 3par id
}

type DeleteImageRequest struct {
	RequestId    string `json:"request_id"`
	ImageId      string `json:"image_id"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type CreateSnapshotRequest struct {
	RequestId    string `json:"request_id"`
	DiskCategory string `json:"disk_category"`
	DiskId       string `json:"disk_id"`
	SnapshotId   string `json:"snapshot_id"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type DeleteSnapshotRequest struct {
	RequestId    string `json:"request_id"`
	DiskCategory string `json:"disk_category"`
	DiskId       string `json:"disk_id"`
	SnapshotId   string `json:"snapshot_id"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type ReInitDiskRequest struct {
	RequestId      string `json:"request_id"`
	DiskCategory   string `json:"disk_category"`
	DiskId         string `json:"disk_id"`
	ImageId        string `json:"image_id"`
	SnapshotId     string `json:"snapshot_id"`
	Size           uint64 `json:"size"`
	OriginalStatus int8   `json:"original_status"`
	StorageType    string `json:"storage_type"`
	ScheduleInfo   string `json:"schedule_info"`
	ImageType      string `json:"image_type"`
}

type ResetDiskRequest struct {
	RequestId      string `json:"request_id"`
	DiskCategory   string `json:"disk_category"`
	DiskId         string `json:"disk_id"`
	DiskType       int8   `json:"disk_type"`
	SnapshotId     string `json:"snapshot_id"`
	OriginalStatus int8   `json:"original_status"`
	UserId         string `json:"UserId"`
	SnapSize       uint64 `json:"snap_size"`
	DiskSize       uint64 `json:"disk_size"`
	StorageType    string `json:"storage_type"`
	ScheduleInfo   string `json:"schedule_info"`
}

type ResizeDiskRequest struct {
	RequestId      string `json:"request_id"`
	DiskCategory   string `json:"disk_category"`
	DiskId         string `json:"disk_id"`
	OldSize        int64  `json:"old_size"` // in byte
	NewSize        int64  `json:"new_size"` // in byte
	DiskType       int8   `json:"disk_type"`
	OriginalStatus int8   `json:"original_status"`
	StorageType    string `json:"storage_type"`
	ScheduleInfo   string `json:"schedule_info"`
	UserId         string `json:"UserId"`
	QoS            string `json:"qos"`
}

type ResizeDisksRequest struct {
	RequestId      string `json:"request_id"`
	DeliveryUnitID string `json:"delivery_unit_id"`
	ResizeAt       time.Time
	DisksReq       []ResizeDiskRequest
}

type ExportDiskRequest struct {
	RequestId    string `json:"request_id"`
	CVKName      string `json:"cvk_name"`
	DiskId       string `json:"disk_id"`
	Iqn          string `json:"iqn"`
	Lun          int    `json:"lun"`
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type ExportDiskResponse struct {
	Lun            int
	ThreeParWWN    string
	ThreeParDataIP string
}

type GetSystemCapacityRequest struct {
	StorageType  string `json:"storage_type"`
	ScheduleInfo string `json:"schedule_info"`
}

type GetSystemUtilizationRequest struct {
	RequestID    string
	StorageType  string
	ScheduleInfo string
}

type DiskQoSRequest struct {
	RequestID    string
	DiskID       string
	DiskCategory string
	Size         uint64
	StorageType  string
	ScheduleInfo string
}

//type DeliveryRequest struct {
//	RequestId string
//	Delivery  model.Delivery
//}
