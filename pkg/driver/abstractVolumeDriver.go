package driver

import (
	"database/sql"
	"immortality/service/compute"
	"immortality/service/data"
	"immortality/service/db_model"
	"immortality/service/driver/three_par"
	"immortality/service/handler/model"
	"immortality/service/logger"
	"strconv"
	"strings"
	"time"
)

var AbsDriver AbstractVolumeDriver

type AbstractVolumeDriver struct{}

func (*AbstractVolumeDriver) Init() {}

func (*AbstractVolumeDriver) CreateDisk(req data.CreateDiskRequest) (resp data.CreateDiskResponse, err error) {

	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return resp, err
	}
	resp, err = d.CreateDisk(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "CreateDisk req: %+v err: %s ", req, err)
		return resp, err
	}

	//update database
	err = db_model.UpdateDiskStatusAvailable(data.Db, resp.Wwn, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Update Disk Status Available error:", err)
		return resp, err
	}

	compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     data.DiskStatusMap[2],
		},
	})

	return resp, nil
}

const timeout = 15

//func StartNewDiskMeasurement(req data.CreateDiskRequest, start time.Time) (err error){
//	//start measurement
//	id, _ := uuid.NewV4()
//	if err != nil {
//		logger.Log.Errorf(req.RequestId, "failed to generate new uuid, err: %+v", err)
//		err = data.ErrServerInternalUUID
//
//	}
//
//	m := db_model.DiskMeasurement{
//		ID:      id.String(),
//		DiskID:  req.DiskId,
//		StartAt: sql.NullTime{Time: start, Valid: true},
//		StopAt:  sql.NullTime{Valid: false},
//		Size:    int64(req.Size),
//	}
//
//	err = m.Insert(data.Db)
//	if err != nil {
//		logger.Log.Errorf(req.RequestId, "failed to insert into db, DiskMeasurement: %+v, err: %+v", m, err)
//		return data.ErrServerInternalDB
//	}
//
//	return err
//}

func (h *AbstractVolumeDriver) CreateDisks(req data.CreateDisksRequest) (err error) {

	success := 0
	lastErr := err
	//var diskIds  []string

	for _, diskReq := range req.DisksReq {

		//diskIds = append(diskIds, diskReq.DiskId)

		disk, err := db_model.DiskByDiskID(data.Db, diskReq.DiskId)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Failed to get disk status err:", err)
			lastErr = err
			continue
		}

		//has create successful
		if disk.Status == data.DiskStatusAvailable {
			success++
			continue
		}

		_, err = h.CreateDisk(diskReq)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "Failed to create disk %s, err:%+v", diskReq.DiskId, err)
			lastErr = err
			continue
		}
		success++
	}

	if success != len(req.DisksReq) {
		logger.Log.Errorf(req.RequestId, "Have success:%d, need suucess:%d", success, len(req.DisksReq))
		return lastErr
	}

	////start measurement
	//	//startAt := time.Now()
	//	//for _,diskReq := range req.DisksReq {
	//	//	err = StartNewDiskMeasurement(diskReq, startAt)
	//	//	if err != nil {
	//	//		lastErr = err
	//	//		logger.Log.Errorf(req.RequestId, "Failed to start disk %s measurement, err:%+v", diskReq.DiskId, err)
	//	//		return lastErr
	//	//	}
	//	//}

	compute.DeliveryUnitCallbackRequest{
		DeliveryUnitID: req.DeliveryUnitId,
		Message:        "success",
		Status:         model.DeliverySuccess,
		StartTime:      req.StartAt.UnixNano() / 1e6,
	}.Notify(req.RequestId)

	return nil
}

func (*AbstractVolumeDriver) DeleteDisk(req data.DeleteDiskRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}
	err = d.DeleteDisk(req)
	if err != nil && err != three_par.ErrVolumeDoesNotExist {
		logger.Log.Errorf(req.RequestId, "DeleteDisk req: %+v, err: %s", req, err)
		return err
	}
	if err == three_par.ErrVolumeDoesNotExist {
		logger.Log.Warn1(req.RequestId, "DeleteDisk req: %+v, err: %s", req, err)
		err = nil
	}

	//update disk deleted
	err = db_model.MarkDiskDeleted(data.Db, req.DiskId)
	if err != nil {
		logger.Log.Errorf(req.DiskId, "mark disk is deleted failed! err: %+v", err)
		return err
	}

	//notify ebs core
	compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     "Deleted",
		},
	})

	return nil
}

type DetachForDelete struct {
	requestID string
	disk      *db_model.Disk
	attach    *db_model.Attach
	exports   []*db_model.Export
	threePar  *db_model.ThreePar
	at        time.Time
}

func (d *DetachForDelete) notifyDetach() (err error) {
	optDiskAO := compute.OptDiskAO{
		DiskID:     d.disk.DiskID,
		InstanceID: d.attach.InstanceId,
		UserID:     d.disk.UserID,
	}
	return compute.NotifyDiskDetach(d.requestID, &optDiskAO)
}

func (d *DetachForDelete) notifyUnmap() (err error) {
	for _, export := range d.exports {
		unmapDiskAO := compute.UnmapDiskAO{
			DiskID: d.disk.DiskID,
			IPs:    strings.Split(d.threePar.Ipv4AddrController, ","),
			Lun:    strconv.Itoa(export.CvkLun),
			WWN:    d.disk.ThreeParWWN,
		}

		if err = compute.NotifyDiskUnmap(d.requestID, &unmapDiskAO); err != nil {
			return
		}
	}

	return
}

func (d *DetachForDelete) detach() (err error) {
	d.attach.IsDeleted = 1
	d.attach.UpdateAt = d.at
	d.attach.DeleteAt = sql.NullTime{
		Time:  d.at,
		Valid: true,
	}
	err = d.attach.Update(data.Db)
	if err != nil {
		logger.Log.Errorf(d.requestID, "failed to update attach, attach: %+v, err: %+v", d.attach, err)
		return data.ErrServerInternalDB
	}

	return
}

func (d *DetachForDelete) cancelExport() (err error) {
	for _, export := range d.exports {
		request := data.ExportDiskRequest{
			RequestId:    d.requestID,
			CVKName:      export.CvkName,
			DiskId:       export.DiskId,
			Iqn:          export.Iqn,
			Lun:          export.CvkLun,
			StorageType:  d.disk.StorageType,
			ScheduleInfo: d.disk.ClusterID,
		}
		if err = AbsDriver.CancelExport(request, 1); err != nil {
			return
		}
	}

	return
}

func newDetachForDelete(requestID, diskID string) (*DetachForDelete, error) {
	h := DetachForDelete{
		requestID: requestID,
		//at:        time.Now(),
	}

	disk, err := db_model.DiskByDiskID(data.Db, diskID)
	if err != nil {
		logger.Log.Errorf(requestID, "failed to query disk, disk: %s, err: %+v", diskID, err)
		return nil, data.ErrInvalidDiskId
	}
	h.disk = disk

	if disk.StatusOrig == data.DiskStatusInUse {
		attach, err := db_model.GetAttach(data.Db, disk.DiskID)
		if err == sql.ErrNoRows {
			return nil, data.ErrAttachInformationNotExists
		}
		if err != nil {
			logger.Log.Errorf(requestID, "GetAttach err:", err)
			return nil, data.ErrServerInternalDB
		}

		exports, err := db_model.GetExports(data.Db, disk.DiskID)
		if err != nil {
			logger.Log.Errorf(requestID, "GetExports err:", err)
			return nil, data.ErrServerInternalDB
		}
		if exports == nil {
			return nil, data.ErrAttachInformationNotExists
		}

		h.attach = attach
		h.exports = exports

		threePar, err := db_model.ThreeParByThreeParID(data.Db, disk.ClusterID)
		if err != nil {
			logger.Log.Errorf(requestID, "ThreeParByThreeParID: %s, err: %+v", disk.ClusterID, err)
			return nil, data.ErrServerInternalDB
		}
		h.threePar = threePar
	}

	return &h, nil
}

func (d *DetachForDelete) DoDetach() (err error) {
	if d.disk.StatusOrig == data.DiskStatusInUse {
		if err = d.notifyDetach(); err != nil {
			logger.Log.Errorf(d.requestID,
				"failed to notifyDetach, DetachForDelete: %+v, err: %+v", d, err)
			return
		}

		if err = d.notifyUnmap(); err != nil {
			logger.Log.Errorf(d.requestID,
				"failed to notifyUnmap, DetachForDelete: %+v, err: %+v", d, err)
			return
		}

		if err = d.detach(); err != nil {
			logger.Log.Errorf(d.requestID,
				"failed to detach, DetachForDelete: %+v, err: %+v", d, err)
			return
		}

		if err = d.cancelExport(); err != nil {
			logger.Log.Errorf(d.requestID,
				"failed to cancelExport, DetachForDelete: %+v, err: %+v", d, err)
			return
		}
	}

	return nil
}

func (h *AbstractVolumeDriver) DeleteDisks(req data.DeleteDisksRequest) (err error) {

	//FinishedAt := req.DeleteAt

	for _, diskInfo := range req.DisksInfo {
		d, err := newDetachForDelete(diskInfo.RequestId, diskInfo.DiskId)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "failed to create DetachForDelete err: %+v", err)
			return err
		}

		d.at = req.DeleteAt

		err = d.DoDetach()
		if err != nil {
			logger.Log.Errorf(req.RequestId, "failed to Detach err: %+v", err)
			return err
		}

		err = h.DeleteDisk(diskInfo)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "failed to Delete Disk:%s, err: %+v", diskInfo.DiskId, err)
			return err
		}

	}

	compute.DeliveryUnitCallbackRequest{
		DeliveryUnitID: req.DeliveryUnitID,
		Message:        "success",
		Status:         model.DeliverySuccess,
		StartTime:      req.DeleteAt.UnixNano() / 1e6,
	}.Notify(req.RequestId)

	return nil
}

func (*AbstractVolumeDriver) CreateImage(req data.CreateImageRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.CreateImage(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "Create Image req: %+v, error: %+v", req, err)
		return err
	}

	err = db_model.MarkImageAvailable(data.Db, req.ImageId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "MarkSnapshotAvailable error:", err)
		return err
	}

	return nil
}

func (*AbstractVolumeDriver) DeleteImage(req data.DeleteImageRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.DeleteImage(req)
	if err != nil && err != three_par.ErrVolumeDoesNotExist {
		logger.Log.Errorf(req.RequestId, "DeleteImage  req: %+v, error: %+v", req, err)
		return err
	}
	if err == three_par.ErrVolumeDoesNotExist {
		logger.Log.Warn1(req.RequestId, "DeleteImage req: %+v, err: %s", req, err)
		err = nil
	}
	return nil
}

func (*AbstractVolumeDriver) CreateSnapshot(req data.CreateSnapshotRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.CreateSnapshot(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "CreateSnapShot req: %+v, error: %+v", req, err)
		return err
	}

	err = db_model.MarkSnapshotAvailable(data.Db, req.SnapshotId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "MarkSnapshotAvailable err", err)
		return err
	}

	time.Sleep(1 * time.Second)

	count, err := db_model.DiskCreatingSnapshot(data.Db, req.DiskId)
	if err != nil {
		logger.Log.Errorf(req.RequestId,
			"failed to query snapshot count that creating by disk id, req: %+v, err: %+v", req, err)
		return err
	}

	if count > 0 {
		logger.Log.Infof(req.RequestId, "still has creating snapshot(s): %d, req: %+v, err: %+v", count, req, err)
		return nil
	}

	disk, err := db_model.DiskByDiskIDForUpdate(data.Db, req.DiskId)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "failed to query disk, req: %+v, err: %+v", req, err)
		return err
	}
	disk.StatusOrig, disk.Status = disk.Status, disk.StatusOrig
	disk.UpdatedAt = time.Now()
	err = disk.Save(data.Db)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "failed to save disk, req: %+v, err: %+v", req, err)
		return err
	}

	go compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     data.DiskStatusMap[disk.Status],
		},
	})

	return nil

}

func (*AbstractVolumeDriver) DeleteSnapshot(req data.DeleteSnapshotRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.DeleteSnapshot(req)
	if err != nil && err != three_par.ErrVolumeDoesNotExist {
		logger.Log.Errorf(req.RequestId, "DeleteSnapshot req: %+v, err: %+v", req, err)
		return err
	}
	if err == three_par.ErrVolumeDoesNotExist {
		logger.Log.Warn1(req.RequestId, "DeleteSnapshot req: %+v, err: %s", req, err)
		err = nil
	}
	return nil

}

func (*AbstractVolumeDriver) ReInitDisk(req data.ReInitDiskRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	resp, err := d.ReInitDisk(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "ReInitDisk req: %+v, err: %+v", req, err)
		return err
	}

	//update database
	err = db_model.UpdateDiskStatusAvailable(data.Db, resp.Wwn, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "UpdateDiskStatusAvailable error:", err)
		return err
	}

	go compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     data.DiskStatusMap[data.DiskStatusAvailable],
		},
	})

	return nil

}

func (*AbstractVolumeDriver) ResetDisk(req data.ResetDiskRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Error1(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.ResetDisk(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "ResetDisk req: %+v, err: %+v", req, err)
		return err
	}

	err = db_model.MarkDiskStatus(data.Db, req.OriginalStatus, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "MarkDiskStatus error:", err)
		return err
	}

	go compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     data.DiskStatusMap[req.OriginalStatus],
		},
	})

	return nil

}

func (h *AbstractVolumeDriver) ResizeDisk(req data.ResizeDiskRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}

	err = d.ResizeDisk(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "ResizeDisk req: %+v, err: %+v", req, err)
		return err
	}

	//if req.OriginalStatus == data.DiskStatusInUse {
	//	updateQoSReq := data.DiskQoSRequest{
	//		RequestID:    req.RequestId,
	//		DiskID:       req.DiskId,
	//		DiskCategory: req.DiskCategory,
	//		Size:         uint64(req.NewSize),
	//		StorageType:  req.StorageType,
	//		ScheduleInfo: req.ScheduleInfo,
	//	}
	//	err = h.UpdateDiskQoS(updateQoSReq)
	//	if err != nil {
	//		logger.Log.Errorf(req.RequestId, "failed to update QoS of disk, req: %+v, err: %w", updateQoSReq, err)
	//		return
	//	}
	//}

	if req.OriginalStatus == data.DiskStatusInUse {
		err = compute.NotifyDiskExtended(req.RequestId, req.UserId, req.DiskId, req.NewSize)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "failed to notify compute, disk: %s, size: %d, user: %s, err: %+v",
				req.DiskId, req.NewSize, req.UserId, err)
			return err
		}
	}

	err = db_model.MarkDiskResized(data.Db, req.OriginalStatus, req.NewSize, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "MarkDiskStatus error:", err)
		return err
	}

	compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
		{
			InstanceID: req.DiskId,
			Status:     data.DiskStatusMap[req.OriginalStatus],
		},
	})

	return nil
}

func (h *AbstractVolumeDriver) ResizeDisks(req data.ResizeDisksRequest) (err error) {
	for _, resizeReq := range req.DisksReq {
		err = h.ResizeDisk(resizeReq)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "Resize disk:%s failed, err:%+v", resizeReq.DiskId, err)
			return nil
		}
	}

	compute.DeliveryUnitCallbackRequest{
		DeliveryUnitID: req.DeliveryUnitID,
		Message:        "success",
		Status:         model.DeliverySuccess,
		StartTime:      req.ResizeAt.UnixNano() / 1e6,
	}.Notify(req.RequestId)

	return nil
}

func (*AbstractVolumeDriver) Export(req data.ExportDiskRequest, retryTimes int) (resp data.ExportDiskResponse, err error) {
	export, err := db_model.GetExport3(data.Db, req.DiskId, req.CVKName)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error1(req.RequestId, "check export error.", err)
		return resp, data.ErrServerInternalDB
	}
	if err == sql.ErrNoRows {
		logger.Log.Error1(req.RequestId, "The specified Export is not exist!")
		return resp, data.ErrInvalidExport
	}
	if export != nil && export.Status == data.ExportStatusExported {
		resp.Lun = export.CvkLun
		return resp, nil
	}
	if export != nil && export.Status == data.ExportStatusUnExportFail {
		export.Status = data.ExportStatusExported
		err = export.Save(data.Db)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
			return resp, err
		}
		resp.Lun = export.CvkLun
		return resp, nil
	}

	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return resp, err
	}

	if !d.NeedExport(req) {
		return resp, nil
	}

	resp, err = d.Export(req)
	if err != nil {
		logger.Log.Errorf(req.RequestId, "Export req: %+v, err: %+v", req, err)
		if retryTimes >= 10 {
			export.Status = data.ExportStatusExportFail
			err = export.Save(data.Db)
			if err != nil {
				logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
				return resp, err
			}
			return resp, nil
		}
		return resp, err
	}

	if resp.Lun == data.HEP_3PAR_SPECIAL_LUN {
		a := db_model.Export{
			DiskId:    req.DiskId,
			CvkName:   req.CVKName,
			Iqn:       req.Iqn,
			CvkLun:    resp.Lun,
			Status:    data.ExportStatusExported,
			CreateAt:  time.Now(),
			UpdateAt:  time.Now(),
			IsDeleted: 0,
		}

		err = a.Save(data.Db)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Error saving export to DB:", err)
			return resp, err
		}

		logger.Log.Infof(req.RequestId, "Export lun is 254, exports again")
		resp, err = d.Export(req)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "Export req: %+v, err: %+v", req, err)
			if retryTimes >= 10 {
				export.Status = data.ExportStatusExportFail
				err = export.Save(data.Db)
				if err != nil {
					logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
					return resp, err
				}
				return resp, nil
			}
			return resp, err
		}
	}
	if resp.Lun == data.HEP_3PAR_SPECIAL_LUN {
		logger.Log.Errorf(req.RequestId, "exports again, lun is 254 again")
	}
	logger.Log.Infof(req.RequestId, "Export again, lun is %d", resp.Lun)

	export.Status = data.ExportStatusExported
	export.CvkLun = resp.Lun
	err = export.Save(data.Db)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
		return resp, err
	}

	return resp, nil
}

func (*AbstractVolumeDriver) CancelExport(req data.ExportDiskRequest, retryTimes int) (err error) {
	//check
	export, err := db_model.GetExport3(data.Db, req.DiskId, req.CVKName)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error1(req.RequestId, "check cancel export error.", err)
		return data.ErrServerInternalDB
	}
	if err == sql.ErrNoRows {
		err = nil
	}
	if export != nil {
		req.Lun = export.CvkLun
		d, err := GetDriver(req.StorageType)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
			return err
		}
		err = d.CancelExport(req)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "CancelExport req: %+v, err: %+v", req, err)
			if retryTimes >= 10 {
				export.Status = data.ExportStatusUnExportFail
				err = export.Save(data.Db)
				if err != nil {
					logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
					return err
				}
				return nil
			}
			return err
		}
	}

	//check 3par lun 254
	export2, err := db_model.GetExport4(data.Db, data.HEP_3PAR_SPECIAL_LUN, req.DiskId, req.CVKName)
	if err != nil && err != sql.ErrNoRows {
		logger.Log.Error1(req.RequestId, "check exports error.", err)
		return data.ErrServerInternalDB
	}
	if err == sql.ErrNoRows {
		err = nil
	}

	if export2 != nil {
		req.Lun = data.HEP_3PAR_SPECIAL_LUN
		d, err := GetDriver(req.StorageType)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
			return err
		}
		err = d.CancelExport(req)
		if err != nil {
			logger.Log.Errorf(req.RequestId, "CancelExport req: %+v, err: %+v", req, err)
			if retryTimes >= 10 {
				export2.Status = data.ExportStatusUnExportFail
				err = export2.Save(data.Db)
				if err != nil {
					logger.Log.Error1(req.RequestId, "Error update export to DB:", err)
					return err
				}
				return nil
			}
			return err
		}
	}

	//update database
	err = db_model.UpdateExportIsDeleted(data.Db, req.DiskId, req.CVKName, time.Now(), time.Now())
	if err != nil {
		logger.Log.Error1(req.RequestId, "UpdateAttachIsDeleted error:", err)
		return err
	}

	disk, err := db_model.DiskByDiskID(data.Db, req.DiskId)
	if err == sql.ErrNoRows {
		logger.Log.Error1(req.RequestId, "Cannot find diskId[", req.DiskId, "]")
		return data.ErrInvalidDiskId
	}
	if err != nil {
		logger.Log.Error1(req.RequestId, "DiskByDiskID error:", err)
		return data.ErrServerInternalDB
	}
	if disk.Status == data.DiskStatusDetaching {
		disk.Status = data.DiskStatusAvailable
		err = disk.Save(data.Db)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Error update disk to DB:", err)
			return err
		}

		// notify uco
		go compute.NotifyEbsCoreStatusUpdate(req.RequestId, []compute.EbsCoreStatusUpdateParams{
			{
				InstanceID: disk.DiskID,
				Status:     data.DiskStatusMap[disk.Status],
			},
		})
	}

	return nil

}

func (*AbstractVolumeDriver) AddDiskQoS(req data.DiskQoSRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestID, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}
	if err = d.AddDiskQoS(req); err != nil {
		logger.Log.Errorf(req.RequestID, "failed to add QoS for disk, req: %+v, err: %+v", req, err)
		return err
	}

	return
}

func (*AbstractVolumeDriver) RemoveDiskQoS(req data.DiskQoSRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestID, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}
	if err = d.RemoveDiskQoS(req); err != nil {
		logger.Log.Errorf(req.RequestID, "failed to remove QoS for disk, req: %+v, err: %+v", req, err)
		return err
	}

	return
}

func (*AbstractVolumeDriver) UpdateDiskQoS(req data.DiskQoSRequest) (err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Errorf(req.RequestID, "GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return err
	}
	if err = d.UpdateDiskQoS(req); err != nil {
		logger.Log.Errorf(req.RequestID, "failed to remove QoS for disk, req: %+v, err: %+v", req, err)
		return err
	}

	return
}

func (*AbstractVolumeDriver) GetSystemCapacity(req data.GetSystemCapacityRequest) (result string, err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Error("GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return result, err
	}

	result, err = d.GetSystemCapacity(req)
	if err != nil {
		logger.Log.Error("GetSystemCapacity req: %+v, err: %+v", req, err)
		return result, err
	}
	return result, nil
}

func (*AbstractVolumeDriver) GetSystemUtilization(req data.GetSystemUtilizationRequest) (ssd, hdd float64, err error) {
	d, err := GetDriver(req.StorageType)
	if err != nil {
		logger.Log.Error("GetDriver err: %+v, storageType: %s", err, req.StorageType)
		return ssd, hdd, err
	}

	ssd, hdd, err = d.GetSystemUtilization(req)
	if err != nil {
		logger.Log.Error("GetSystemCapacity req: %+v, err: %+v", req, err)
		return ssd, hdd, err
	}
	return ssd, hdd, nil
}
