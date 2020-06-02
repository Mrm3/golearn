package disk_service

import (
	"database/sql"
	"encoding/json"
	"immortality-demo/pkg/data"
	"immortality-demo/pkg/db_model"
	"immortality-demo/pkg/gredis"
	"immortality-demo/pkg/logger"
	"immortality-demo/service/model"
	"time"
)

type CreateDiskHandler struct {
}

func (h *CreateDiskHandler) Handle(param model.CreateDiskParams, userId, requestId string) (result map[string]interface{}, err error) {

	//param, err = verifyCreateDisk(requestId, param)
	if err != nil {
		logger.Log.Error1(requestId, "CreateDisk param verify err:", err, "CreateDiskParams:", param)
		return nil, err
	}

	var imageType string
	var clusterId string
	if param.SnapshotId != "" {
		snapshot, err := db_model.SnapshotBySnapshotID(data.Db, param.SnapshotId)
		if err == sql.ErrNoRows {
			logger.Log.Error1(requestId, "cannot find snapshot:", param.SnapshotId)
			return nil, data.ErrInvalidSnapshotId
		}
		if err != nil {
			logger.Log.Error1(requestId, "CreateDisk SnapshotBySnapshotID error:", err)
			return nil, err
		}
		if snapshot.UserID != userId {
			logger.Log.Error1(requestId, "CreateDisk snapshot userId doesn't match")
			return nil, data.ErrInvalidSnapshotId
		}
		param.Size = snapshot.Size >> 30
		clusterId = snapshot.ClusterID
	} else if param.ImageId != "" {
		image, err := image_service.GetImageByImageID(param.ImageId)
		if err != nil {
			logger.Log.Error1(requestId, "CreateDisk getImageByImageID error:", err)
			return nil, err
		}
		if param.Size < image.Size {
			return nil, data.ErrInvalidSize
		}
		imageType = data.IMAGE_TYPE_CUSTOM
		if image.UserId == data.SpatialUserPublic {
			imageType = data.IMAGE_TYPE_PUBLIC
		}
		if imageType == data.IMAGE_TYPE_CUSTOM {
			if image.UserId != userId {
				logger.Log.Error1(requestId, "CreateDisk image userId doesn't match")
				return nil, data.ErrInvalidImageId
			}
		}
		clusterId = image.ClusterId
	}

	storageType := param.StorageType
	if storageType == data.HPE3PARA && param.SnapshotId == "" && imageType != data.IMAGE_TYPE_CUSTOM {
		scheduler := schedule.GetScheduler()
		info := data.ScheduleParam{
			DiskCategory: param.DiskCategory,
			PodId:        param.PodId,
		}
		var scheduleResult data.ScheduleResult
		scheduleResult, err = scheduler.Schedule(info)
		if err != nil {
			logger.Log.Error1(requestId, "Schedule error, schedule info:%s", info)
			return nil, err
		}
		clusterId = scheduleResult.Id
	} else if storageType == data.CEPH {
		//cephCluster, err := config.CephClusterByCategory(param.DiskCategory)
		//if err != nil {
		//	return nil, data.ErrZoneNotAvailable
		//}
		//clusterId = cephCluster.Id
	}

	disk := db_model.Disk{
		DiskID:       param.DiskId,
		Status:       data.DiskStatusCreating,
		DiskType:     data.DiskTypeMap2[param.DiskType],
		Name:         param.DiskName,
		Region:       param.RegionId,
		Zone:         param.ZoneId,
		Category:     data.CategoryMap[param.DiskCategory],
		Size:         param.Size << 30,
		Description:  param.Description,
		FromSnapshot: param.SnapshotId,
		FromImage:    param.ImageId,
		UserID:       userId,
		ClusterID:    clusterId,
		StorageType:  param.StorageType,
		IsShare:      data.DiskIsShareMap2[param.IsShare],
		Qos:          param.Qos,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	err = disk.Save(data.Db)
	if err != nil {
		logger.Log.Error1(requestId, "Error saving disk to DB:", err)
		return nil, data.ErrServerInternalDB
	}

	asyncReq := data.CreateDiskRequest{
		RequestId:    requestId,
		DiskId:       param.DiskId,
		StorageType:  storageType,
		DiskCategory: param.DiskCategory,
		Size:         uint64(param.Size << 30),
		UserId:       userId,
		ScheduleInfo: clusterId,
		Qos:          param.Qos,
		SnapshotId:   param.SnapshotId,
		ImageId:      param.ImageId,
		ImageType:    imageType,
	}

	jsonBytes, err := json.Marshal(asyncReq)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return nil, err
	}
	err = gredis.Push("queue", jsonBytes)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to publish:", err)
		return nil, data.ErrServerInternalMQ
	}

	result = map[string]interface{}{
		"DiskId":    param.DiskId,
		"RequestId": requestId,
	}

	return result, nil
}
