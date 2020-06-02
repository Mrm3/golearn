package three_par

import (
	"database/sql"
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/gofrs/uuid"

	. "immortality-demo/pkg/data"
	"immortality-demo/pkg/db_model"
	"immortality-demo/pkg/logger"
	"immortality-demo/pkg/util"
)

type ThreeParDriver struct {
	HttpClient *http.Client
	ServerPath string
	SessionKey string
	User       string
	Password   string
}

type ThreeParVolumeDriver struct {
	ThreePars map[string]*ThreeParDriver
}

func Create3paraVolumeDriver() (v *ThreeParVolumeDriver, err error) {
	threePars := map[string]*ThreeParDriver{}

	httpClient := &http.Client{
		Timeout: 120 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			DialContext: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	//Get 3par information
	pars, err := db_model.ThreePars(Db)
	if err != nil {
		logger.Log.Error("GetThreePars error:", err)
		return nil, err
	}

	if pars != nil {
		for _, par := range pars {
			threeParId := par.ThreeParId
			threePars[threeParId] = &ThreeParDriver{
				HttpClient: httpClient,
				ServerPath: "http://" + par.Ipv4AddrManagement + "/api/v1",
				SessionKey: "",
				User:       par.Username,
				Password:   par.Password,
			}
		}
	}
	for _, threePar := range threePars {
		err = threePar.InitSessionKey()
		if err != nil {
			logger.Log.Error("Failed to initSessionKey:", err)
			return nil, err
		}
	}

	//package threeParVolumeDriver
	threeParVolumeDriver := &ThreeParVolumeDriver{
		ThreePars: threePars,
	}
	return threeParVolumeDriver, nil
}

func (p *ThreeParVolumeDriver) CreateDisk(req CreateDiskRequest) (resp CreateDiskResponse, err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return resp, err
	}

	req.Size = req.Size >> 20
	var cpgName string
	if req.DiskCategory == DiskCategorySSD {
		cpgName = SSD_CPG
	} else if req.DiskCategory == DiskCategoryHDD {
		cpgName = HDD_CPG
	} else {
		cpgName = HYBRID_SSD_CPG
	}

	if req.ImageId != "" {
		err = T.createVolumeFromImage(req.RequestId, req.ImageId, req.DiskId, int64(req.Size), cpgName)
	} else if req.SnapshotId != "" {
		err = T.createVolumeFromSnapshot(req.RequestId, req.SnapshotId, req.DiskId, cpgName)
	} else {
		optional := map[string]interface{}{"snapCPG": cpgName, "tpvv": true}
		err = T.CreateVolume(req.RequestId, req.DiskId, cpgName, int64(req.Size), optional)
	}

	//defer func() {
	//	if err != nil {
	//		e := T.DeleteVolume(req.RequestId, req.DiskId)
	//		if e != nil {
	//			logger.Log.Error1(req.RequestId, "Failed to rollback created disk:", req.ScheduleInfo, req.DiskId)
	//		}
	//	}
	//}()

	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to create disk:", req, err)
		return resp, err
	}

	//get wwn of disk
	disk, err := T.GetVolumeByName(req.RequestId, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to GetVolumeByName:", req.DiskId, err)
		return
	}

	resp.Wwn = disk.Wwn

	return resp, nil
}

func (p *ThreeParVolumeDriver) DeleteDisk(req DeleteDiskRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	err = T.DeleteVolume(req.RequestId, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to deleteDisk:", req.DiskId, err)
		return err
	}
	return nil
}

func (p *ThreeParVolumeDriver) CreateImage(req CreateImageRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	var destCpg string
	disk, err := T.GetVolumeByName(req.RequestId, req.SourceDiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to GetVolumeByName,", req.SourceDiskId, err)
		return
	}
	if disk.Name != "" {
		destCpg = disk.UserCPG
	}

	err = T.CreateVolume(req.RequestId, req.ImageId, destCpg, int64(disk.SizeMiB), map[string]interface{}{"snapCPG": destCpg, "tpvv": true})
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to CreateVolume:", err)
		return
	}

	optional := map[string]interface{}{"online": false, "priority": 1}
	err = T.CloneVolume(req.RequestId, disk.Name, req.ImageId, destCpg, optional)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to createImage,", err)
		return err
	}

	return nil
}

func (p *ThreeParVolumeDriver) DeleteImage(req DeleteImageRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	err = T.DeleteVolume(req.RequestId, req.ImageId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to deleteImage,ImageId:", req.ImageId, err)
		return err
	}
	return nil
}

func (p *ThreeParVolumeDriver) CreateSnapshot(req CreateSnapshotRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	optional := map[string]interface{}{}
	err = T.CreateVolumeSnapshot(req.RequestId, req.SnapshotId, req.DiskId, optional)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to createSnapshot:", req.DiskId, err)
		return err
	}

	return nil
}

func (p *ThreeParVolumeDriver) DeleteSnapshot(req DeleteSnapshotRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	err = T.DeleteVolume(req.RequestId, req.SnapshotId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to deleteSnapshot:", req.SnapshotId, err)
		return err
	}
	return nil
}

func (p *ThreeParVolumeDriver) ReInitDisk(req ReInitDiskRequest) (resp CreateDiskResponse, err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return
	}

	//delete original disk
	err = T.DeleteVolume(req.RequestId, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to deleteVolume:", req.DiskId, err)
		return
	}

	//create a new one
	req.Size = req.Size >> 20
	var cpgName string
	if req.DiskCategory == DiskCategorySSD {
		cpgName = SSD_CPG
	} else if req.DiskCategory == DiskCategoryHDD {
		cpgName = HDD_CPG
	} else {
		cpgName = HYBRID_SSD_CPG
	}
	if req.ImageId != "" {
		err = T.createVolumeFromImage(req.RequestId, req.ImageId, req.DiskId, int64(req.Size), cpgName)
	} else if req.SnapshotId != "" {
		err = T.createVolumeFromSnapshot(req.RequestId, req.SnapshotId, req.DiskId, cpgName)
	} else {
		optional := map[string]interface{}{"snapCPG": cpgName, "tpvv": true}
		err = T.CreateVolume(req.RequestId, req.DiskId, cpgName, int64(req.Size), optional)
	}
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to create new disk:", req, err)
		return
	}

	//get wwn of disk
	disk, err := T.GetVolumeByName(req.RequestId, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to GetVolumeByName:", req.DiskId, err)
		return
	}

	resp.Wwn = disk.Wwn

	return
}

func (p *ThreeParVolumeDriver) ResetDisk(req ResetDiskRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	optional := map[string]interface{}{}
	err = T.PromoteVirtualCopy(req.RequestId, req.SnapshotId, optional)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to resetDisk:", req.SnapshotId, err)
		return err
	}

	return nil
}

func (p *ThreeParVolumeDriver) ResizeDisk(req ResizeDiskRequest) (err error) {

	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	amount := (req.NewSize - req.OldSize) >> 20
	if amount > 0 {
		err = T.GrowVolume(req.RequestId, req.DiskId, amount)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Failed to GrowVolume:", req.DiskId, err)
			return err
		}
	}

	return nil
}

func (p *ThreeParVolumeDriver) Export(req ExportDiskRequest) (resp ExportDiskResponse, err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return resp, err
	}

	var hostName string

	hostName = req.CVKName
	hostName = strings.ReplaceAll(hostName, "(", "")
	hostName = strings.ReplaceAll(hostName, ")", "")
	host, err := T.GetHostByName(req.RequestId, hostName)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to GetHostByName:", hostName, err)
		return resp, err
	}
	if host.Name == "" {
		//create host
		err = T.CreateHost(req.RequestId, hostName, []string{req.Iqn}, nil, nil)
		if err != nil {
			logger.Log.Error1(req.RequestId, "Failed to CreateHost:", req, err)
			return resp, err
		}
	} else {
		//check the iqn of host
		not := true
		for _, IscsiPath := range host.ISCSIPaths {
			if IscsiPath.Name == req.Iqn {
				not = false
				break
			}
		}
		if not {
			//modify the iqn of host
			parameters := map[string]interface{}{"pathOperation": HOST_EDIT_OPERATION_ADD,
				"iSCSINames": []string{req.Iqn}}
			err = T.ModifyHost(req.RequestId, hostName, parameters)
			if err != nil {
				logger.Log.Error1(req.RequestId, "Failed to ModifyHost:", req, err)
				return resp, err
			}
		}
	}

	//create VLUN
	lunId, err := T.CreateVLUN(req.RequestId, req.DiskId, req.CVKName, 0, true, nil)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to CreateVLUN:", req, err)
		return resp, err
	}
	//get disk
	disk, err := T.GetVolumeByName(req.RequestId, req.DiskId)
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to GetVolumeByName:", err)
		return resp, err
	}

	resp.Lun = lunId
	resp.ThreeParWWN = disk.Wwn
	return resp, nil
}

func (p *ThreeParVolumeDriver) CancelExport(req ExportDiskRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestId, err)
		return err
	}

	var hostName = req.CVKName
	//hostName = strings.ReplaceAll(hostName, "(", "")
	//hostName = strings.ReplaceAll(hostName, ")", "")
	//host, err := T.GetHostByName(req.RequestId, hostName)
	//if err != nil {
	//	logger.Log.Error1(req.RequestId, "Failed to GetHostByName:", hostName, err)
	//	return err
	//}
	//if host.Name == "" {
	//	err = errors.New("The host does not exist,hostName:" + hostName)
	//	logger.Log.Error1(req.RequestId, err)
	//	return err
	//}

	//delete VLUN
	err = T.DeleteVLUN(req.RequestId, req.DiskId, hostName, float64(req.Lun))
	if err != nil {
		logger.Log.Error1(req.RequestId, "Failed to DeleteVLUN:", req, err)
		return err
	}

	return nil
}

func (p *ThreeParVolumeDriver) GetSystemCapacity(req GetSystemCapacityRequest) (result string, err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error(err)
		return result, err
	}

	result, err = T.GetSystemCapacity()
	if err != nil {
		logger.Log.Error("Failed to GetSystemCapacity:", err)
		return result, err
	}

	return result, nil
}

func (p *ThreeParVolumeDriver) GetSystemUtilization(req GetSystemUtilizationRequest) (ssd, hdd float64, err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Errorf(req.RequestID, "%s", err)
		return ssd, hdd, err
	}

	ssd, hdd, err = T.GetSystemUtilization()
	if err != nil {
		logger.Log.Errorf(req.RequestID, "failed to GetSystemUtilization:", err)
		return ssd, hdd, err
	}

	return ssd, hdd, err
}

func (p *ThreeParVolumeDriver) NeedExport(req ExportDiskRequest) (isNeed bool) {
	return true
}

// AddDiskQoS 设定磁盘限速规则
func (p *ThreeParVolumeDriver) AddDiskQoS(req DiskQoSRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestID, err)
		return err
	}

	capacity := int64(req.Size >> 30)

	switch req.DiskCategory {
	case DiskCategorySSD:
		err = T.CreateVolumeSet(req.RequestID, req.DiskID, "", "", []string{req.DiskID})
		if err != nil {
			logger.Log.Errorf(req.RequestID, "failed to CreateVolumeSet, req: %+v, err: %+v", req, err)
			return
		}

		// 计算 QoS
		var bwMaxLimitKB uint64
		var ioMaxLimit uint32
		bwMaxLimitKB, ioMaxLimit, err = util.QosByCapacity(capacity, DiskCategoryToInstanceCode[req.DiskCategory])
		if err != nil {
			logger.Log.Errorf(req.RequestID, "failed to QosByCapacity, req: %+v, err: %+v", req, err)
			return
		}
		qosRules := map[string]interface{}{
			"bwMinGoalKB":  Hpe3parBwMinGoalKB,
			"bwMaxLimitKB": bwMaxLimitKB,
			"ioMinGoal":    Hpe3parIoMinGoal,
			"ioMaxLimit":   ioMaxLimit,
		}
		err = T.CreateQoSRules(req.RequestID, req.DiskID, QoS_TargetType_VVSET, qosRules)
		if err != nil && err != ErrQoSRuleExistent {
			logger.Log.Error1(req.RequestID, "failed to CreateQoSRules:", req.DiskID, qosRules, err)
			return fmt.Errorf("failed to create QoS level, err: %w", err)
		}
		return nil
	case DiskCategoryHybridHDD, DiskCategoryHDD:
		var tx *sql.Tx
		if tx, err = Db.Begin(); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to start db transaction: %+v", err)
			return ErrServerInternalDB
		}
		defer func() {
			if err != nil {
				if e := tx.Rollback(); e != nil {
					logger.Log.Errorf(req.RequestID, "failed to rollback, err: %+v", e)
				}
			}
		}()

		var level *db_model.QoSLevel
		if level, err = db_model.QoSLevelByCapForUpdate(tx, capacity); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to QoSLevelByCapForUpdate, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}

		group := fmt.Sprintf("%s-%d", level.ID, level.NextGroupID)
		params := map[string]interface{}{
			"action":     1, // Adds a member to the VV set
			"setmembers": []string{req.DiskID},
		}
		err = T.ModifyVolumeSet(req.RequestID, group, params)
		if err == ErrVolumeHasInSet {
			if err = tx.Commit(); err != nil {
				logger.Log.Errorf(req.RequestID, "failed to commit tx, req: %+v, err: %+v", req, err)
				return ErrServerInternalDB
			}
		}
		if err != nil {
			logger.Log.Errorf(req.RequestID, "failed to ModifyVolumeSet, req: %+v, err: %+v", req, err)
			return
		}

		var id uuid.UUID
		if id, err = uuid.NewV4(); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to generate new uuid, err: %+v", err)
			return ErrServerInternalUUID
		}
		at := time.Now()
		diskQoS := db_model.DiskQoS{
			ID:        id.String(),
			DiskID:    req.DiskID,
			LevelID:   level.ID,
			GroupID:   level.NextGroupID,
			CreatedAt: at,
			UpdatedAt: at,
			DeletedAt: sql.NullTime{Valid: false},
		}
		if err = diskQoS.Insert(tx); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to insert db, record: %+v, err: %+v", diskQoS, err)
			return
		}

		level.NextGroupID = level.NextGroupID%level.GroupCount + 1
		if err = level.Update(tx); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to Update, level: %+v, err: %+v", level, err)
			return ErrServerInternalDB
		}

		if err = tx.Commit(); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to commit tx, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}
		return
	default:
		logger.Log.Warnf(req.RequestID, "unexpected disk category, req: %+v", req)
		return
	}
}

// RemoveDiskQoS 移除磁盘限速规则
func (p *ThreeParVolumeDriver) RemoveDiskQoS(req DiskQoSRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestID, err)
		return err
	}

	switch req.DiskCategory {
	case DiskCategorySSD:
		// 3par will automatically delete QoS
		err = T.DeleteVolumeSet(req.RequestID, req.DiskID)
		if err != nil && err != ErrSetDoesNotExist {
			logger.Log.Errorf(req.RequestID, "failed to remove vvset, err: %+v", err)
			return
		}
		return nil
	case DiskCategoryHybridHDD, DiskCategoryHDD:
		var tx *sql.Tx
		if tx, err = Db.Begin(); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to start db transaction: %+v", err)
			return ErrServerInternalDB
		}
		defer func() {
			if err != nil {
				if e := tx.Rollback(); e != nil {
					logger.Log.Errorf(req.RequestID, "failed to rollback, err: %+v", e)
				}
			}
		}()

		var diskQoS *db_model.DiskQoS
		if diskQoS, err = db_model.DiskQoSByDiskIDForUpdate(tx, req.DiskID); err != nil && err != sql.ErrNoRows {
			logger.Log.Errorf(req.RequestID, "failed to DiskQoSByDiskIDForUpdate, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}
		if diskQoS != nil {
			if err = diskQoS.Delete(tx); err != nil {
				logger.Log.Errorf(req.RequestID, "failed to Delete, record: %+v, err: %+v", diskQoS, err)
				return ErrServerInternalDB
			}
		}

		group := fmt.Sprintf("%s-%d", diskQoS.LevelID, diskQoS.GroupID)
		params := map[string]interface{}{
			"action":     2, // Removes a member from the VV set
			"setmembers": []string{req.DiskID},
		}
		if err = T.ModifyVolumeSet(req.RequestID, group, params); err != nil && err != ErrVolumeNotInSet {
			logger.Log.Errorf(req.RequestID, "failed to ModifyVolumeSet, vvset: %s, params: %+v, err: %+v", group, params, err)
			return
		}

		if err = tx.Commit(); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to commit tx, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}

		return nil
	default:
		logger.Log.Warnf(req.RequestID, "unexpected disk category, req: %+v", req)
		return
	}
}

// UpdateDiskQoS 更新磁盘限速规则
func (p *ThreeParVolumeDriver) UpdateDiskQoS(req DiskQoSRequest) (err error) {
	T := p.ThreePars[req.ScheduleInfo]
	if T == nil {
		err = errors.New("3par information is empty,clusterId:" + req.ScheduleInfo)
		logger.Log.Error1(req.RequestID, err)
		return
	}

	capacity := int64(req.Size >> 30)

	switch req.DiskCategory {
	case DiskCategorySSD:
		// 计算 QoS
		var bwMaxLimitKB uint64
		var ioMaxLimit uint32
		bwMaxLimitKB, ioMaxLimit, err = util.QosByCapacity(capacity, DiskCategoryToInstanceCode[req.DiskCategory])
		if err != nil {
			logger.Log.Errorf(req.RequestID, "failed to QosByCapacity, req: %+v, err: %+v", req, err)
			return
		}
		qosRules := map[string]interface{}{
			"bwMinGoalKB":  Hpe3parBwMinGoalKB,
			"bwMaxLimitKB": bwMaxLimitKB,
			"ioMinGoal":    Hpe3parIoMinGoal,
			"ioMaxLimit":   ioMaxLimit,
		}
		err = T.ModifyQoSRules(req.RequestID, req.DiskID, QoSTargetType[QoS_TargetType_VVSET], qosRules)
		if err != nil && err != ErrQosRuleDoesNotExist {
			logger.Log.Error1(req.RequestID, "failed to ModifyQoSRules:", req.DiskID, qosRules, err)
			return fmt.Errorf("failed to modify QoS rule, err: %w", err)
		}
		return nil
	case DiskCategoryHybridHDD, DiskCategoryHDD:
		var current string
		if current, err = db_model.DiskQoSIDByDiskID(Db, req.DiskID); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to query record, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}
		var shouldBe string
		if shouldBe, err = db_model.QoSLevelByCap(Db, capacity); err != nil {
			logger.Log.Errorf(req.RequestID, "failed to query record, req: %+v, err: %+v", req, err)
			return ErrServerInternalDB
		}

		if current != shouldBe {
			logger.Log.Infof(req.RequestID, "updates needed, current: %s, should be: %s", current, shouldBe)

			if err = p.RemoveDiskQoS(req); err != nil {
				return fmt.Errorf("failed to remove QoS, req: %+v, err: %w", req, err)
			}
			if err = p.AddDiskQoS(req); err != nil {
				return fmt.Errorf("failed to add QoS, req: %+v, err: %w", req, err)
			}
		}

		logger.Log.Infof(req.RequestID, "no updates needed, current: %s, should be: %s", current, shouldBe)

		return
	default:
		logger.Log.Warnf(req.RequestID, "unexpected disk category, req: %+v", req)
		return
	}
}
