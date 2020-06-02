package driver

import (
	"encoding/json"
	. "immortality-demo/pkg/data"

	"immortality-demo/pkg/logger"
)

func Dispatch(action string, payload string) (err error) {
	switch action {
	case "CreateDisk":
		var req CreateDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		logger.Log.Infof(req.RequestId, "CreateDisk:%+v", req)
		_, err = AbsDriver.CreateDisk(req)
		return err

	case "CreateDisks":
		var req CreateDisksRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		logger.Log.Infof(req.RequestId, "CreateDisks:%+v", req)
		return AbsDriver.CreateDisks(req)

	case "DeleteDisk":
		var req DeleteDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.DeleteDisk(req)

	case "DeleteDisks":
		var req DeleteDisksRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		logger.Log.Infof(req.RequestId, "DeleteDisks:%+v", req)
		return AbsDriver.DeleteDisks(req)

	case "CreateSnapshot":
		var req CreateSnapshotRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}

		return AbsDriver.CreateSnapshot(req)

	case "DeleteSnapshot":
		var req DeleteSnapshotRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}

		return AbsDriver.DeleteSnapshot(req)

	case "CreateImage":
		var req CreateImageRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}

		return AbsDriver.CreateImage(req)

	case "DeleteImage":
		var req DeleteImageRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}

		return AbsDriver.DeleteImage(req)

	case "ResetDisk":
		var req ResetDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}

		return AbsDriver.ResetDisk(req)

	case "ReInitDisk":

		var req ReInitDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.ReInitDisk(req)

	case "ResizeDisk":
		var req ResizeDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.ResizeDisk(req)

	case "ResizeDisks":
		var req ResizeDisksRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.ResizeDisks(req)

	case "Export":
		var req ExportDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		_, err = AbsDriver.Export(req, retryTimes)
		return err

	case "CancelExport":
		var req ExportDiskRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.CancelExport(req, retryTimes)

	case "AddDiskQoS":
		var req DiskQoSRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.AddDiskQoS(req)

	case "RemoveDiskQoS":
		var req DiskQoSRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.RemoveDiskQoS(req)

	case "UpdateDiskQoS":
		var req DiskQoSRequest
		err := json.Unmarshal([]byte(payload), &req)
		if err != nil {
			logger.Log.Error("json Unmarshal failed! json:%s", payload)
			return err
		}
		return AbsDriver.UpdateDiskQoS(req)

	case "ResizeDelivery":
		//	var req DeliveryRequest
		//	if err := json.Unmarshal([]byte(payload), &req); err != nil {
		//		logger.Log.Error("json Unmarshal failed! json:%s", payload)
		//		return err
		//	}
		//	return businesshandler.ResizeDeliveryHandler{}.Handle(&req)
		//
		//default:
		//	logger.Log.Error("Wrong volume action [%s]", action)
	}
	return err
}
