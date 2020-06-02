package model

type CreateDiskParams struct {
	RegionId     string `json:"RegionId"`
	ZoneId       string `json:"ZoneId"`
	SnapshotId   string `json:"SnapshotId"`
	ImageId      string `json:"ImageId"`
	DiskId       string `json:"DiskId"`
	DiskName     string `json:"DiskName"`
	InstanceCode string `json:"InstanceCode"`
	DiskCategory string `json:"-"`
	Size         int64  `json:"Size"`
	Description  string `json:"Description"`
	DiskType     string `json:"DiskType"`
	StorageType  string `json:"StorageType"`
	IsShare      string `json:"IsShare"`
	Qos          string `json:"Qos"`
	PodId        string `json:"PodId"`
}