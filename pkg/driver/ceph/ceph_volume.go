package ceph

import (
	"errors"
	"github.com/ceph/go-ceph/rbd"
	"immortality-demo/config"
	"immortality-demo/pkg/data"
	"immortality-demo/pkg/db_model"
	"immortality-demo/pkg/logger"
)

var (
	ErrClusterNotFound = errors.New("ceph cluster not configured")
)

type CephVolumeDriver struct {
	// maps id -> ceph handle
	//idMapping map[string]ceph
	// map category(hdd, hybrid-hdd, ssd) -> ceph handle
	categoryMapping map[string]ceph
	//courierClient   courier.Client
}

func CreateCephVolumeDriver() (v *CephVolumeDriver, err error) {
	//courierClient, err := courier.NewCourierClient()
	//if err != nil {
	//	return nil, errors.New("Failed to init courier:" + err.Error())
	//}

	cephVolumeDriver := &CephVolumeDriver{
		//idMapping:       make(map[string]ceph),
		categoryMapping: make(map[string]ceph),
		//courierClient:   courierClient,
	}

	for _, conf := range config.Config.Ceph {
		ceph, err := newCeph(conf)
		if err != nil {
			return nil, errors.New("Failed to init ceph client: " + err.Error())
		}

		//cephVolumeDriver.idMapping[conf.Id] = ceph
		cephVolumeDriver.categoryMapping[conf.Category] = ceph
	}

	return cephVolumeDriver, nil
}

func (v *CephVolumeDriver) CreateDisk(req data.CreateDiskRequest) (resp data.CreateDiskResponse, err error) {
	//cephConfig, err := config.CephClusterByCategory(req.DiskCategory)
	//if err != nil {
	//	return resp, data.ErrZoneNotAvailable
	//}

	//clusterId := cephConfig.Id

	if req.SnapshotId != "" {
		err = v.createDiskFromSnapshot(req.DiskCategory, req.DiskId, req.SnapshotId)
	} else if req.ImageId != "" {
		err = v.createDiskFromImage(req.DiskCategory, req.DiskId, req.ImageType, req.ImageId, uint64(req.Size))
	} else {
		err = v.createBlankDisk(req.DiskCategory, req.DiskId, uint64(req.Size))
	}

	if err != nil {
		logger.Log.Error1(req.RequestId, "Create disk failed, req:", req, "error:", err)
		return
	}

	defer func() {
		if err != nil {
			ceph, e := v.getClientByCategory(req.DiskCategory)
			if e != nil {
				return // it's a rollback, should not happen
			}
			e = ceph.removeImage(DiskPool, req.DiskId)
			if e != nil {
				logger.Log.Error1(req.RequestId, "Failed to rollback created disk:",
					req.DiskCategory, req.DiskId)
			}
		}
	}()

	return
}

func (v *CephVolumeDriver) createDiskFromSnapshot(category, diskId string, snapshotId string) (err error) {
	var (
		snapshot *db_model.Snapshot
		cluster  ceph
	)

	snapshot, err = db_model.SnapshotBySnapshotID(data.Db, snapshotId)
	if err != nil {
		logger.Log.Error("CreateDisk SnapshotBySnapshotID error:", err)
		return
	}

	cluster, err = v.getClientByCategory(category)
	if err != nil {
		return
	}

	return cluster.createDiskFromSnapshot(snapshot.DiskID, snapshot.SnapshotID, diskId)
}

func (v *CephVolumeDriver) createDiskFromImage(category, diskId, imageType, imageId string, size uint64) (err error) {
	var (
	//dbImage        *image_service.Image
	//imageClusterId string
	)

	//dbImage, err = image_service.GetImageByImageID(imageId)
	//if err != nil {
	//	logger.Log.Error("GetImageByImageID error:", err)
	//	return
	//}

	// customized image only in HDD cluster
	//imageCluster, err := v.getClientByCategory("hdd")
	//if err != nil {
	//	logger.Log.Error("Failed to get HDD cluster")
	//	return err
	//}
	//imageClusterId = imageCluster.ID

	ceph, err := v.getClientByCategory(category)
	if err != nil {
		return
	}
	//if imageClusterId == clusterId {
	//	err = ceph.createDiskFromImage(imageType, imageId, diskId)
	//} else {
	//	var client courierpb.CourierClient
	//	client, err = v.courierClient.PickCourier(imageId)
	//	if err != nil {
	//		logger.Log.Error("PickCourier error:", err)
	//		return
	//	}
	//	req := courierpb.RbdRequest{
	//		NodeID:             config.Config.NodeID,
	//		RequestID:          string(util.GenerateRandomId()),
	//		SourceCluster:      imageClusterId,
	//		SourcePool:         courierpb.PoolType_Image,
	//		SourceRbdId:        imageId,
	//		SourceSnapshotId:   "",
	//		DestinationCluster: clusterId,
	//		DestinationPool:    courierpb.PoolType_Disk,
	//		DestinationRbdId:   diskId,
	//		Size:               uint64(dbImage.Size),
	//	}
	//	reply, err := client.CreateRbd(context.Background(), &req)
	//	if err != nil {
	//		logger.Log.Error(req.NodeID, req.RequestID,
	//			"Courier CreateRbd error:", err)
	//		return err
	//	}
	//	if reply.Error != "" {
	//		return errors.New(req.RequestID + ":" + reply.Error)
	//	}
	//}

	err = ceph.createDiskFromImage(imageType, imageId, diskId)
	if err != nil {
		logger.Log.Error("Create disk from image error:", err)
		return
	}

	defer func() {
		if err != nil {
			e := ceph.removeImage(DiskPool, diskId)
			if e != nil {
				logger.Log.Error("createDiskFromImage rollback error:", e)
			}
		}
	}()

	return ceph.resizeImage(DiskPool, diskId, size)
}

func (v *CephVolumeDriver) createBlankDisk(category, diskId string, size uint64) (err error) {
	var (
		cluster ceph
		image   *rbd.Image
	)

	cluster, err = v.getClientByCategory(category)
	if err != nil {
		return
	}

	image, err = cluster.create(DiskPool, diskId, size)
	if err != nil {
		return
	}

	// FIXME handle this error?
	_ = image.Close()

	return
}

func (v *CephVolumeDriver) DeleteDisk(req data.DeleteDiskRequest) (err error) {
	var (
		cluster ceph
	)

	cluster, err = v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not found for", req.ScheduleInfo, err)
		return
	}

	err = cluster.removeImage(DiskPool, req.DiskId)
	// TODO handle disk not found error
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph delete error:", err)
		return
	}

	return
}

func (v *CephVolumeDriver) CreateImage(req data.CreateImageRequest) (err error) {
	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "hdd ceph cluster not configured", err)
		return
	}

	logger.Log.Info1(req.RequestId, "cluster", cluster.ID, "sourceCluster", req.ScheduleInfo)

	//if cluster.ID == req.ScheduleInfo {
	//	// source and image are of same ceph cluster
	//	err = cluster.createImageSameCluster(req.SourceDiskId, req.ImageId)
	//} else {
	//	// source and image are in different ceph clusters
	//	client, err := v.courierClient.PickCourier(req.SourceDiskId)
	//	if err != nil {
	//		logger.Log.Error("PickCourier error:", err)
	//		return err
	//	}
	//	req := courierpb.RbdRequest{
	//		NodeID:             config.Config.NodeID,
	//		RequestID:          string(util.GenerateRandomId()),
	//		SourceCluster:      req.ScheduleInfo,
	//		SourcePool:         courierpb.PoolType_Disk,
	//		SourceRbdId:        req.SourceDiskId,
	//		SourceSnapshotId:   "",
	//		DestinationCluster: cluster.ID,
	//		DestinationPool:    courierpb.PoolType_Image,
	//		DestinationRbdId:   req.ImageId,
	//		Size:               req.Size,
	//	}
	//	reply, err := client.CreateRbd(context.Background(), &req)
	//	if err != nil {
	//		logger.Log.Error(req.NodeID, req.RequestID,
	//			"Courier CreateRbd error:", err)
	//		return err
	//	}
	//	if reply.Error != "" {
	//		return errors.New(req.RequestID + ":" + reply.Error)
	//	}
	//}

	err = cluster.createImageSameCluster(req.SourceDiskId, req.ImageId)
	if err != nil {
		logger.Log.Error("Create image error, req:", req, "error:", err)
		return
	}

	defer func() {
		if err != nil {
			logger.Log.Info1(req.RequestId, "Rollback: CreateImage -> RemoveImage",
				"image ID:", req.ImageId)
			_ = cluster.removeImage(ImagePool, req.ImageId)
		}
	}()

	return
}

func (v *CephVolumeDriver) DeleteImage(req data.DeleteImageRequest) (err error) {
	cluster, err := v.getClientByCategory("hdd")
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph HDD cluster not configured", err)
		return
	}

	err = cluster.removeImage(ImagePool, req.ImageId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph delete error:", err)
		return
	}

	return
}

func (v *CephVolumeDriver) CreateSnapshot(req data.CreateSnapshotRequest) (err error) {
	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not configured for", req.ScheduleInfo, err)
		return
	}

	_, err = cluster.createSnapshot(req.DiskId, req.SnapshotId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "CreateSnapshot failed:", err)
		return
	}

	defer func() {
		if err != nil {
			logger.Log.Info1(req.RequestId, "Rollback: CreateSnapshot -> RemoveSnapshot",
				"disk:", req.DiskId, "snapshot:", req.SnapshotId)
			_ = cluster.removeSnapshot(req.DiskId, req.SnapshotId)
		}
	}()

	return
}

func (v *CephVolumeDriver) DeleteSnapshot(req data.DeleteSnapshotRequest) (err error) {
	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not configured for", req.ScheduleInfo, err)
		return
	}

	err = cluster.removeSnapshot(req.DiskId, req.SnapshotId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph RemoveSnapshot error:", err)
		return
	}

	return
}

func (v *CephVolumeDriver) ReInitDisk(req data.ReInitDiskRequest) (resp data.CreateDiskResponse, err error) {
	// 1. delete original disk
	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not configured for", req.ScheduleInfo, err)
		return
	}

	err = cluster.removeImage(DiskPool, req.DiskId)
	// TODO handle disk not found error
	if err != nil {
		logger.Log.Error1(req.RequestId, "RemoveImage error:", err)
		return
	}

	// 2. create a new one
	if req.ImageId != "" {
		err = v.createDiskFromImage(req.ScheduleInfo, req.DiskId, req.ImageType, req.ImageId, req.Size)
	} else if req.SnapshotId != "" {
		err = v.createDiskFromSnapshot(req.ScheduleInfo, req.DiskId,
			req.SnapshotId)
	} else {
		err = v.createBlankDisk(req.ScheduleInfo, req.DiskId, req.Size)
	}
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to recreate disk:", err)
		return
	}

	return
}

func (v *CephVolumeDriver) ResetDisk(req data.ResetDiskRequest) (err error) {
	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not configured for", req.ScheduleInfo, err)
		return
	}

	err = cluster.rollbackSnapshot(req.DiskId, req.SnapshotId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph RollbackSnapshot error:", err)
		return
	}

	//Check the size between snap and disk,and if the snap size is smaller, make disk to resize.
	if req.SnapSize < req.DiskSize {
		err = cluster.resizeImage(DiskPool, req.DiskId, req.DiskSize)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Ceph ResizeImage error:", err)
			//update database
			err = db_model.ResizeDisk(data.Db, req.SnapSize, req.DiskId)
			if err != nil {
				logger.Log.Error1(req.RequestId, "DB update disk size error:", err)
				return
			}

		}
	}

	return
}

func (v *CephVolumeDriver) ResizeDisk(req data.ResizeDiskRequest) (err error) {
	//defer func() {
	//	err = db_model.MarkDiskStatus(data.Db, req.OriginalStatus, req.DiskId)
	//}()

	cluster, err := v.getClientByCategory(req.DiskCategory)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph cluster not configured for", req.ScheduleInfo, err)
		return
	}

	err = cluster.resizeImage(DiskPool, req.DiskId, uint64(req.NewSize))
	if err != nil {
		logger.Log.Error1(req.RequestId, "Ceph ResizeImage diskid[%d],size[%d] error:", err, req.DiskId, req.NewSize)
		return
	}

	return
}

//func (v *CephVolumeDriver) getClientById(id string) (ceph, error) {
//	cluster, ok := v.idMapping[id]
//	if !ok {
//		return ceph{}, ErrClusterNotFound
//	}
//
//	return cluster, nil
//}

func (v *CephVolumeDriver) getClientByCategory(category string) (ceph, error) {
	cluster, ok := v.categoryMapping[category]
	if !ok {
		return ceph{}, ErrClusterNotFound
	}

	return cluster, nil
}

func (v *CephVolumeDriver) NeedExport(req data.ExportDiskRequest) (isNeed bool) {
	return false
}

func (v *CephVolumeDriver) Export(req data.ExportDiskRequest) (resp data.ExportDiskResponse, err error) {
	return resp, nil
}

func (v *CephVolumeDriver) CancelExport(req data.ExportDiskRequest) (err error) {
	return nil
}

func (v *CephVolumeDriver) GetSystemCapacity(req data.GetSystemCapacityRequest) (result string, err error) {
	return result, nil
}

func (v *CephVolumeDriver) GetSystemUtilization(req data.GetSystemUtilizationRequest) (ssd, hdd float64, err error) {
	return ssd, hdd, nil
}

func (v *CephVolumeDriver) AddDiskQoS(req data.DiskQoSRequest) (err error) {
	logger.Log.Warnf(req.RequestID, "CephVolumeDriver.AddDiskQoS does not been implemented")
	return nil
}

func (v *CephVolumeDriver) RemoveDiskQoS(req data.DiskQoSRequest) (err error) {
	logger.Log.Warnf(req.RequestID, "CephVolumeDriver.RemoveDiskQoS does not been implemented")
	return nil
}

func (v *CephVolumeDriver) UpdateDiskQoS(req data.DiskQoSRequest) (err error) {
	logger.Log.Warnf(req.RequestID, "CephVolumeDriver.UpdateDiskQoS does not been implemented")
	return nil
}
