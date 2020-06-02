package driver

import (
	. "immortality/service/data"
	"immortality/service/driver/ceph"
	"immortality/service/driver/three_par"
	"immortality/service/logger"
)

var volumeDrivers = make(map[string]VolumeDriver)

type VolumeDriver interface {
	CreateDisk(req CreateDiskRequest) (resp CreateDiskResponse, err error)
	DeleteDisk(req DeleteDiskRequest) (err error)

	CreateImage(req CreateImageRequest) (err error)
	DeleteImage(req DeleteImageRequest) (err error)

	CreateSnapshot(req CreateSnapshotRequest) (err error)
	DeleteSnapshot(req DeleteSnapshotRequest) (err error)

	ReInitDisk(req ReInitDiskRequest) (resp CreateDiskResponse, err error)
	ResetDisk(req ResetDiskRequest) (err error)
	ResizeDisk(req ResizeDiskRequest) (err error)

	NeedExport(req ExportDiskRequest) (isNeed bool)
	Export(req ExportDiskRequest) (resp ExportDiskResponse, err error)
	CancelExport(req ExportDiskRequest) (err error)

	GetSystemCapacity(req GetSystemCapacityRequest) (result string, err error)
	GetSystemUtilization(req GetSystemUtilizationRequest) (ssd, hdd float64, err error)

	AddDiskQoS(req DiskQoSRequest) (err error)
	RemoveDiskQoS(req DiskQoSRequest) (err error)
	UpdateDiskQoS(req DiskQoSRequest) (err error)
}

func GetDriver(storageType string) (d VolumeDriver, err error) {
	d, ok := volumeDrivers[storageType]
	if ok {
		return d, nil
	}

	switch storageType {
	case CEPH:
		d, err = ceph.CreateCephVolumeDriver()
		if err != nil {
			logger.Log.Error("create ceph volume Driver error [%s]", err)
			return nil, err
		}

	case HPE3PARA:
		d, err = three_par.Create3paraVolumeDriver()
		if err != nil {
			logger.Log.Error("create 3para volume Driver error [%s]", err)
			return nil, err
		}
	default:
		logger.Log.Error("Wrong StorageType :%s", storageType)
		return nil, ErrInvalidDiskType
	}

	volumeDrivers[storageType] = d
	return d, nil
}
