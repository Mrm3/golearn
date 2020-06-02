package three_par

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"immortality-demo/pkg/logger"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	HYBRID_SSD_CPG string = "hybrid-ssd"
	HYBRID_HDD_CPG string = "hybrid-hdd"
	SSD_CPG        string = "ssd"
	HDD_CPG        string = "hdd"
	IMAGE_CPG      string = "image"
)

const sessionCookieName = "X-Hp3Par-Wsapi-Sessionkey"

const virtualCopy float64 = 3

const (
	GrowVolume         float64 = 3
	PromoteVirtualCopy float64 = 4
)

const (
	QoS_TargetType_VVSET  int8 = 1
	QoS_TargetType_SYS    int8 = 2
	QoS_TargetType_DOMAIN int8 = 4
)

const (
	TASKSTATUS_DONE      float64 = 1 //The task has finished.
	TASKSTATUS_ACTIVE    float64 = 2 //The task in progress.
	TASKSTATUS_CANCELLED float64 = 3 //The task was canceled.
	TASKSTATUS_FAILED    float64 = 4 //The task failed.
)

const (
	HOST_EDIT_OPERATION_ADD    int = 1
	HOST_EDIT_OPERATION_REMOVE int = 2
)

var (
	QoSTargetType = map[int8]string{
		QoS_TargetType_VVSET: "vvset",
		QoS_TargetType_SYS:   "sys",
	}

	ErrQosRuleDoesNotExist = errors.New("QoS rule does not exist")
	ErrQoSRuleExistent     = errors.New("the QoS rule exists")
	ErrSetDoesNotExist     = errors.New("the set does not exist")
	ErrVolumeNotInSet      = errors.New("volume is not part of the set")
	ErrVolumeHasInSet      = errors.New("object is already part of the set")
)

var (
	ErrVolumeDoesNotExist = errors.New("the volume does not exist")
)

type AuthenticateBody struct {
	User     string `json:"user"`
	Password string `json:"password"`
}

type AuthenticateResponse struct {
	Key string `json:"key"`
}

type HttpResponseBody struct {
	Code ErrCode `json:"code"`
	Desc string  `json:"desc"`
	Ref  string  `json:"ref"`
}

func (h HttpResponseBody) String() string {
	return fmt.Sprintf("code: %d, desc: %s, ref: %s", h.Code, h.Desc, h.Ref)
}

type ErrCode int

const (
	NonExistentVolume  ErrCode = 23
	NonExistentQoSRule ErrCode = 100
	HaveExistentVolume ErrCode = 104
	ExistentQoSRule    ErrCode = 114
)

const (
	HPE3ParSystemIsBusy ErrCode = 270
)

type VolumeInfos struct {
	Total   int          `json:"total"`
	Members []VolumeInfo `json:"members"`
}

type VolumeInfo struct {
	Id           float64 `json:"id"`
	Name         string  `json:"name"`
	TotalUsedMiB float64 `json:"totalUsedMiB"`
	SizeMiB      float64 `json:"sizeMiB"`
	UserCPG      string  `json:"userCPG"`
	SnapCPG      string  `json:"snapCPG"`
	Wwn          string  `json:"wwn"`
}

type VolumeSnapInfo struct {
	Id       float64 `json:"id"`
	Name     string  `json:"name"`
	SizeMiB  float64 `json:"sizeMiB"`
	CopyOf   string  `json:"copyOf"`
	CopyType float64 `json:"copyType"`
	Wwn      string  `json:"wwn"`
}

type HostInfo struct {
	Id         float64     `json:"id"`
	Name       string      `json:"name"`
	FCPaths    interface{} `json:"FCPaths"`
	ISCSIPaths []IscsiPath `json:"iSCSIPaths"`
	Persona    float64     `json:"persona"`
}

type IscsiPath struct {
	Name      string      `json:"name"`
	PortPos   interface{} `json:"portPos"`
	IPAddr    string      `json:"IPAddr"`
	HostSpeed float64     `json:"hostSpeed"`
}

type PortPos struct {
	Node int `json:"node"`
	Slot int `json:"slot"`
	Port int `json:"port"`
}

type VlunIfo struct {
	Lun        float64 `json:"lun"`
	VolumeName string  `json:"volumeName"`
	Hostname   string  `json:"hostname"`
	VolumeWWN  string  `json:"volumeWWN"`
}

type VolumeSetInfo struct {
	Id               float64  `json:"id"`
	Uuid             string   `json:"uuid"`
	Name             string   `json:"name"`
	Comment          string   `json:"comment"`
	Setmembers       []string `json:"setmembers"`
	FlashCachePolicy float64  `json:"flashCachePolicy"`
	QosEnabled       bool     `json:"qosEnabled"`
}

type TaskId struct {
	TaskId float64 `json:"taskid"`
}

type TaskInfo struct {
	Id         float64 `json:"id"`
	Type       float64 `json:"type"`
	Name       string  `json:"name"`
	Status     float64 `json:"status"`
	StartTime  string  `json:"startTime"`
	FinishTime string  `json:"finishTime"`
	User       string  `json:"user"`
}

type QoSRules struct {
	BwMinGoalKB  uint64 `json:"bwMinGoalKB"`
	BwMaxLimitKB uint64 `json:"bwMaxLimitKB"`
	IoMinGoal    uint32 `json:"ioMinGoal"`
	IoMaxLimit   uint32 `json:"ioMaxLimit"`
}

func (q QoSRules) String() string {
	return fmt.Sprintf("{\"bwMinGoalKB\": %d, \"bwMaxLimitKB\": %d, \"ioMinGoal\": %d, \"ioMaxLimit\": %d}",
		q.BwMinGoalKB, q.BwMaxLimitKB, q.IoMinGoal, q.IoMaxLimit)
}

func (q QoSRules) FromString(qos string) (rules QoSRules, err error) {
	err = json.Unmarshal([]byte(qos), &rules)
	return
}

type RawReservedSpaceData struct {
	UsedMiB  float64 `json:"usedMiB"`
	SnapMiB  float64 `json:"snapMiB"`
	AdminMiB float64 `json:"adminMiB"`
	TotalMiB float64 `json:"totalMiB"`
}

type UserSpaceData struct {
	UsedMiB     float64 `json:"usedMiB"`
	FreeMiB     float64 `json:"freeMiB"`
	ReservedMiB float64 `json:"reservedMiB"`
}

type SnapAdminData struct {
	UsedMiB     float64 `json:"usedMiB"`
	FreeMiB     float64 `json:"freeMiB"`
	ReservedMiB float64 `json:"reservedMiB"`
	VCopyMiB    float64 `json:"vcopyMiB"`
}

type TotalSpaceData struct {
	UsedMiB        float64 `json:"usedMiB"`
	VirtualSizeMiB float64 `json:"virtualSizeMiB"`
	ReservedMiB    float64 `json:"reservedMiB"`
	VCopyMiB       float64 `json:"vcopyMiB"`
	HostWriteMiB   float64 `json:"hostWriteMiB"`
}

type CapacityEfficiency struct {
	Compaction       float64 `json:"compaction"`
	Compression      float64 `json:"compression"`
	DataReduction    float64 `json:"dataReduction"`
	OverProvisioning float64 `json:"overProvisioning"`
	Deduplication    float64 `json:"deduplication"`
}

type AtTimeVolumeSpaceData struct {
	Domain             string               `json:"domain"`
	ID                 uint32               `json:"id"`
	Name               string               `json:"name"`
	BaseID             uint32               `json:"baseId"`
	WWN                string               `json:"wwn"`
	SnapCPG            string               `json:"snapCPG"`
	UserCPG            string               `json:"userCPG"`
	ProvisioningType   int                  `json:"provisioningType"`
	CopyType           int                  `json:"copyType"`
	CompressionState   int                  `json:"compressionState"`
	VvsetName          string               `json:"vvsetName"`
	RawReserved        RawReservedSpaceData `json:"rawReserved"`
	UserSpace          UserSpaceData        `json:"userSpace"`
	SnapSpace          SnapAdminData        `json:"snapSpace"`
	AdminSpace         SnapAdminData        `json:"adminSpace"`
	TotalSpace         TotalSpaceData       `json:"totalSpace"`
	CapacityEfficiency CapacityEfficiency   `json:"capacityEfficiency"`
	CompressionGcKBPS  float64              `json:"compressionGcKBPS"`
}

type AtTimeVolumeSpaceResponse struct {
	SampleTime    string                  `json:"sampleTime"`
	SampleTimeSec int32                   `json:"sampleTimeSec"`
	Total         int32                   `json:"total"`
	Members       []AtTimeVolumeSpaceData `json:"members"`
	//Links       []interface{}           `json:"links"`
}

func (T *ThreeParDriver) InitSessionKey() (err error) {
	request := AuthenticateBody{
		User:     T.User,
		Password: T.Password}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error("Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/credentials", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error("NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	logger.Log.Info("InitSessionKey Request to three_par:", requestUrl, string(jsonBytes))
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error("Failed to connect to three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		var response AuthenticateResponse
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error("ReadJsonBody error:", err)
			return
		}
		if response.Key != "" {
			T.SessionKey = response.Key
		}
	} else {
		var response HttpResponseBody
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error("ReadJsonBody error:", err)
			return
		}

		logger.Log.Info("Failed to authenticate,three_par response:", response.String())
		err = errors.New("Failed to authenticate,three_par response:" + response.String())
		return err
	}
	return
}

//basic methods
func (T *ThreeParDriver) CreateVolume(requestId, volumeName, cpgName string, sizeMiB int64, optional map[string]interface{}) (err error) {
	request := map[string]interface{}{"name": volumeName, "cpg": cpgName, "sizeMiB": sizeMiB}
	if optional != nil {
		for k, v := range optional {
			request[k] = v
		}
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateVolume Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CreateVolume(requestId, volumeName, cpgName, sizeMiB, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to createVolume:", err)
				return err
			}
			return nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			if response.Code == HPE3ParSystemIsBusy {
				sleepTime := rand.Intn(5-1) + 1
				logger.Log.Debugf(requestId, "Waiting %d seconds to CreateVolume again!", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)
				err = T.CreateVolume(requestId, volumeName, cpgName, sizeMiB, optional)
				if err != nil {
					logger.Log.Error1(requestId, "Failed to createVolume:", err)
					return err
				}
				return nil
			}

			if response.Code == 22 && response.Desc == "volume exists" {
				logger.Log.Warnf(requestId, "Volume [%s] has exists", volumeName)
				return nil
			}

			logger.Log.Error1(requestId, "Failed to createVolume,three_par response:", response.String())
			err = errors.New("Failed to createVolume,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) createVolumeFromImage(requestId, srcName, destName string, sizeMiB int64, cpgName string) (err error) {
	image, err := T.GetVolumeByName(requestId, srcName)
	if err != nil && err != ErrVolumeDoesNotExist {
		logger.Log.Error1(requestId, "Failed to GetVolumeByName:", err)
		return
	}
	if err == ErrVolumeDoesNotExist {
		logger.Log.Error1(requestId, "The source image is not exist")
		return
	}

	err = T.CreateVolume(requestId, destName, cpgName, int64(image.SizeMiB), map[string]interface{}{"snapCPG": cpgName, "tpvv": true})
	if err != nil {
		logger.Log.Error1(requestId, "Failed to CreateVolume:", err)
		return
	}

	optional := map[string]interface{}{"online": false, "priority": 1}
	err = T.CloneVolume(requestId, srcName, destName, cpgName, optional)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to CloneVolume:", err)
		return
	}
	//resize to newSize
	amount := sizeMiB - int64(image.SizeMiB)
	if amount > 0 {
		err = T.GrowVolume(requestId, destName, amount)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to GrowVolume:", err)
			return err
		}
	}
	return err
}

func (T *ThreeParDriver) createVolumeFromSnapshot(requestId, srcName, destName, cpgName string) (err error) {
	snapshot, err := T.GetVolumeByName(requestId, srcName)
	if err != nil && err != ErrVolumeDoesNotExist {
		logger.Log.Error1(requestId, "Failed to GetVolumeByName:", err)
		return
	}
	if err == ErrVolumeDoesNotExist {
		logger.Log.Error1(requestId, "The source snapshot is not exist")
		return
	}

	err = T.CreateVolume(requestId, destName, cpgName, int64(snapshot.SizeMiB), map[string]interface{}{"snapCPG": cpgName, "tpvv": true})
	if err != nil {
		logger.Log.Error1(requestId, "Failed to CreateVolume:", err)
		return
	}

	optional := map[string]interface{}{"online": false, "priority": 1}
	err = T.CloneVolume(requestId, srcName, destName, cpgName, optional)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to CloneVolume:", err)
		return
	}
	return
}

func (T *ThreeParDriver) UnAuthenticate() (err error) {
	requestUrl := fmt.Sprintf("%s/credentials/%s", T.ServerPath, T.SessionKey)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error("NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info("UnAuthenticate Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error("Failed to request three_par:", err)
		return
	}
	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error("Failed to init three_par:", err)
				return err
			}
			err = T.UnAuthenticate()
			if err != nil {
				logger.Log.Error("Failed to UnAuthenticate:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error("ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info("Failed to unAuthenticate,three_par response:", response.String())
			err = errors.New("Failed to unAuthenticate,three_par response:" + response.String())
			return err
		}
	} else {
		T.SessionKey = ""
	}
	return
}

func (T *ThreeParDriver) GetVolumeByName(requestId, volumeName string) (volume VolumeInfo, err error) {
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, volumeName)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetVolumeByName Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if resp.Body != nil {
			err = readJsonBody(resp.Body, &volume)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return
			}
			jsonBytes2, err := json.Marshal(volume)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to marshal json:", err)
				return volume, err
			}
			logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
		}
		return volume, nil
	case http.StatusNotFound:
		logger.Log.Error1(requestId, "GetVolumeByName response: The volume does not exist.")
		return volume, ErrVolumeDoesNotExist
	default:
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return volume, err
			}
			volume, err = T.GetVolumeByName(requestId, volumeName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetVolumeByName:", err)
				return volume, err
			}
			return volume, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return volume, err
			}
			if response.Code == HPE3ParSystemIsBusy {
				sleepTime := rand.Intn(5-1) + 1
				logger.Log.Debugf(requestId, "Waiting %d seconds to GetVolumeByName again!", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)
				volume, err = T.GetVolumeByName(requestId, volumeName)
				if err != nil {
					logger.Log.Error1(requestId, "Failed to GetVolumeByName:", err)
					return volume, err
				}
				return volume, nil
			} else {
				logger.Log.Info1(requestId, "Failed to get volume,three_par response:", response.String())
				err = errors.New("Failed to get volume,three_par response:" + response.String())
				return volume, err
			}
		}
	}
}

func (T *ThreeParDriver) ModifyVolume(requestId, name string, parameters map[string]interface{}) (err error) {
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, name)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "ModifyVolume Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.ModifyVolume(requestId, name, parameters)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyVolume:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to modify volume,three_par response:", response.String())
			err = errors.New("Failed to modify volume,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) GrowVolume(requestId, name string, sizeMiB int64) (err error) {
	request := map[string]interface{}{"action": GrowVolume, "sizeMiB": sizeMiB}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, name)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GrowVolume Request to three_par:", requestUrl, string(jsonBytes))
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.GrowVolume(requestId, name, sizeMiB)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CloneVolume:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to grow volume,three_par response:", response.String())
			err = errors.New("Failed to grow volume,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) CloneVolume(requestId, srcName, destName, destCpg string, optional map[string]interface{}) (err error) {
	parameters := map[string]interface{}{"destVolume": destName, "destCPG": destCpg}
	if optional != nil {
		for k, v := range optional {
			parameters[k] = v
		}
	}
	if v, ok := parameters["online"]; ok {
		value, ok := v.(bool)
		if !ok {
			logger.Log.Error1(requestId, "Convert value error")
		}
		if !value {
			delete(parameters, "destCPG")
		}
	} else {
		delete(parameters, "destCPG")
	}
	request := map[string]interface{}{"action": "createPhysicalCopy", "parameters": parameters}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, srcName)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CloneVolume Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CloneVolume(requestId, srcName, destName, destCpg, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CloneVolume:", err)
				return err
			}
			return nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to clone volume,three_par response:", response.String())
			err = errors.New("Failed to clone volume,three_par response:" + response.String())
			return err
		}
	}
	var taskId TaskId
	if resp.Body != nil {
		err = readJsonBody(resp.Body, &taskId)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return err
		}
		jsonBytes2, err := json.Marshal(taskId)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	if taskId.TaskId != 0 {
		var (
			//	retryTime int = 1
			sleepTime int = 1
		)

		for {
			//sleepTime =  1 << retryTime

			logger.Log.Debugf(requestId, "Waiting %d seconds to getTask", sleepTime)

			time.Sleep(time.Duration(sleepTime) * time.Second)

			task, err := T.getTask(requestId, taskId.TaskId)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to getTask,", err)
				return err
			}
			if task.Status == TASKSTATUS_DONE || task.Status == TASKSTATUS_CANCELLED {
				break
			} else if task.Status == TASKSTATUS_FAILED {
				err = errors.New("Failed to clone volume,taskId is:" + strconv.FormatFloat(task.Id, 'f', 0, 64))
				return err
			}
			if task.Status != TASKSTATUS_ACTIVE {
				logger.Log.Warn1(requestId, "Clone Task current status:[%d],not in expect status:[TASKSTATUS_ACTIVE]", task.Status)
			}

			sleepTime = rand.Intn(10-1) + 1

			//if retryTime > TaskMaxRetryTimes {
			//	logger.Log.Warnf(requestId, "Have retried %d times to getTask, break", retryTime)
			//
			//	break
			//}

		}
	}
	return nil
}

func (T *ThreeParDriver) DeleteVolume(requestId, name string) (err error) {
	//check for snapshots
	volume, err := T.GetVolumeByName(requestId, name)
	if err != nil {
		logger.Log.Error1(requestId, "GetVolumeByName error:", err)
		return
	}
	snapshots, err := T.GetSnapshotsOfVolume(requestId, volume.SnapCPG, name)
	if snapshots != nil {
		for _, volume := range snapshots {
			err = T.RemoveSnapshot(requestId, volume.Name)
			if err != nil {
				logger.Log.Error1(requestId, "RemoveSnapshot error:", err)
				return
			}
		}
	}

	//delete volume
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, name)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "DeleteVolume Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.DeleteVolume(requestId, name)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to DeleteVolume:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to DeleteVolume,three_par response:", response.String())
			err = errors.New("Failed to DeleteVolume,three_par response:" + response.String())
			return err
		}
	}
	return
}

//snapshot methods
func (T *ThreeParDriver) CreateVolumeSnapshot(requestId, name, copyOfName string, optional map[string]interface{}) (err error) {
	parameters := map[string]interface{}{"name": name}
	if optional != nil {
		for k, v := range optional {
			parameters[k] = v
		}
	}
	request := map[string]interface{}{"action": "createSnapshot", "parameters": parameters}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, copyOfName)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateVolumeSnapshot Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CreateVolumeSnapshot(requestId, name, copyOfName, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CreateSnapshot:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to create snapshot,three_par response:", response.String())
			err = errors.New("Failed to create snapshot,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) PromoteVirtualCopy(requestId, snapshot string, optional map[string]interface{}) (err error) {
	request := map[string]interface{}{"action": PromoteVirtualCopy}
	if optional != nil {
		for k, v := range optional {
			request[k] = v
		}
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, snapshot)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "PromoteVirtualCopy Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.PromoteVirtualCopy(requestId, snapshot, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to PromoteVirtualCopy:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to promoteVirtualCopy,three_par response:", response.String())
			err = errors.New("Failed to promoteVirtualCopy,three_par response:" + response.String())
			return err
		}
	}
	var taskId TaskId
	if resp.Body != nil {
		err = readJsonBody(resp.Body, &taskId)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return err
		}
		jsonBytes2, err := json.Marshal(taskId)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	if taskId.TaskId != 0 {
		for {
			time.Sleep(1 * time.Second)
			task, err := T.getTask(requestId, taskId.TaskId)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to getTask,", err)
				return err
			}
			if task.Status == TASKSTATUS_DONE || task.Status == TASKSTATUS_CANCELLED {
				break
			} else if task.Status == TASKSTATUS_FAILED {
				err = errors.New("Failed to PromoteVirtualCopy,taskId is:" + strconv.FormatFloat(task.Id, 'f', 0, 64))
				return err
			}
		}
	}
	return
}

func (T *ThreeParDriver) RemoveSnapshot(requestId, snapshotName string) (err error) {
	requestUrl := fmt.Sprintf("%s/volumes/%s", T.ServerPath, snapshotName)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "RemoveSnapshot Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.RemoveSnapshot(requestId, snapshotName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to RemoveSnapshot:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to remove snapshot,three_par response:", response.String())
			err = errors.New("Failed to remove snapshot,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) GetSnapshotsOfVolume(requestId, snapCPG, volName string) (snaps []VolumeSnapInfo, err error) {
	requestUrl := fmt.Sprintf("%s/volumes", T.ServerPath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	q := req.URL.Query()
	q.Add("query", "\"copyOf\tEQ\t"+volName+"\"")
	req.URL.RawQuery = q.Encode()
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetSnapshotsOfVolume Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return nil, err
			}
			snaps, err = T.GetSnapshotsOfVolume(requestId, snapCPG, volName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetSnapshotsOfVolume:", err)
				return nil, err
			}
			return snaps, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return nil, err
			}

			logger.Log.Info1(requestId, "Failed to get snapshots of volume,three_par response:", response.String())
			err = errors.New("Failed to get snapshots of volume,three_par response:" + response.String())
			return nil, err
		}
	}
	if resp.Body != nil {
		var response map[string]interface{}
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return
		}
		if _, ok := response["members"]; ok {
			list, ok := response["members"].([]interface{})
			if ok {
				for _, v := range list {
					var snap VolumeSnapInfo
					jsonBytesTemp, err := json.Marshal(v)
					if err != nil {
						logger.Log.Error1(requestId, "Failed to marshal json:", err)
						return nil, err
					}
					err = json.Unmarshal(jsonBytesTemp, &snap)
					if err != nil {
						return nil, err
					}
					if snap.CopyType == virtualCopy {
						snaps = append(snaps, snap)
					}
				}
			}
		}
		//jsonBytes2, err := json.Marshal(response)
		//if err != nil {
		//	logger.Log.Error1(requestId,"Failed to marshal json:", err)
		//	return nil, err
		//}
		//logger.Log.Info1(requestId,"three_par response:", string(jsonBytes2))
	}
	return snaps, nil
}

//HOST methods
func (T *ThreeParDriver) CreateHost(requestId, hostName string, iscsiNames, FCWwns []string, optional map[string]interface{}) (err error) {
	request := map[string]interface{}{"name": hostName}
	if iscsiNames != nil {
		request["iSCSINames"] = iscsiNames
	}
	if FCWwns != nil {
		request["FCWWNs"] = FCWwns
	}
	if optional != nil {
		for k, v := range optional {
			request[k] = v
		}
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/hosts", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateHost Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CreateHost(requestId, hostName, iscsiNames, FCWwns, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CreateHost:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to create host,three_par response:", response.String())
			err = errors.New("Failed to create host,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) ModifyHost(requestId, hostName string, parameters map[string]interface{}) (err error) {
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/hosts/%s", T.ServerPath, hostName)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "ModifyHost Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.ModifyHost(requestId, hostName, parameters)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyHost:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to modify host,three_par response:", response.String())
			err = errors.New("Failed to modify host,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) DeleteHost(requestId, hostName string) (err error) {
	//check vluns
	host, err := T.GetHostByName(requestId, hostName)
	if err != nil {
		logger.Log.Error1(requestId, "GetHostByName error:", err)
		return
	}
	if host.Name == "" {
		logger.Log.Error1(requestId, "the host is not exist")
		return
	}
	vluns, err := T.GetVlunsByHostname(requestId, hostName)
	if vluns != nil {
		for _, host := range vluns {
			err = T.DeleteVLUN(requestId, host.VolumeName, host.Hostname, host.Lun)
			if err != nil {
				logger.Log.Error1(requestId, "DeleteVLUN error:", err)
				return
			}
		}
	}
	//delete host
	requestUrl := fmt.Sprintf("%s/hosts/%s", T.ServerPath, hostName)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "DeleteHost Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.DeleteHost(requestId, hostName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to DeleteHost:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to remove host,three_par response:", response.String())
			err = errors.New("Failed to remove host,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) GetHostByName(requestId, hostName string) (host HostInfo, err error) {
	requestUrl := fmt.Sprintf("%s/hosts/%s", T.ServerPath, hostName)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetHostByName Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return host, err
			}
			host, err = T.GetHostByName(requestId, hostName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetHostByName:", err)
				return host, err
			}
			return host, nil
		} else {
			result2, err := checkHostExist(resp)
			if !result2 {
				return host, nil
			}
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return host, err
			}
			if response.Code == HPE3ParSystemIsBusy {
				sleepTime := rand.Intn(5-1) + 1
				logger.Log.Debugf(requestId, "Waiting %d seconds to GetHostByName again!", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)
				host, err = T.GetHostByName(requestId, hostName)
				if err != nil {
					logger.Log.Error1(requestId, "Failed to GetHostByName:", err)
					return host, err
				}
				return host, nil
			} else {
				logger.Log.Info1(requestId, "Failed to get host,three_par response:", response.String())
				err = errors.New("Failed to get host,three_par response:" + response.String())
				return host, err
			}
		}
	}
	if resp.Body != nil {
		err = readJsonBody(resp.Body, &host)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return
		}
		jsonBytes2, err := json.Marshal(host)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return host, err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	return host, nil
}

func (T *ThreeParDriver) GetHostByIqn(requestId, iqn string) (host HostInfo, err error) {
	requestUrl := fmt.Sprintf("%s/hosts", T.ServerPath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	q := req.URL.Query()
	q.Add("query", "\"iSCSIPaths[name\tEQ\t"+iqn+"]\"")
	req.URL.RawQuery = q.Encode()
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetHostByIqn Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return host, err
			}
			host, err = T.GetHostByIqn(requestId, iqn)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetHostByIqn:", err)
				return host, err
			}
			return host, nil
		} else {
			result2, err := checkHttpResponseOfHostExist(resp)
			if result2 {
				return host, nil
			}
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return host, err
			}

			logger.Log.Info1(requestId, "Failed to get GetHostByIqn of host,three_par response:", response.String())
			err = errors.New("Failed to get GetHostByIqn of host,three_par response:" + response.String())
			return host, err
		}
	}
	if resp.Body != nil {
		var response map[string]interface{}
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return
		}
		if _, ok := response["members"]; ok {
			list, ok := response["members"].([]interface{})
			if ok {
				if len(list) > 0 {
					value := list[0]
					jsonBytesTemp, err := json.Marshal(value)
					if err != nil {
						logger.Log.Error1(requestId, "Failed to marshal json:", err)
						return host, err
					}
					err = json.Unmarshal(jsonBytesTemp, &host)
					if err != nil {
						return host, err
					}
				}
			}
		}
		jsonBytes2, err := json.Marshal(response)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return host, err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	return host, nil
}

//VLUN methods
func (T *ThreeParDriver) CreateVLUN(requestId, volumeName, hostName string, lun int8, auto bool, optional map[string]interface{}) (lunId int, err error) {
	request := map[string]interface{}{"volumeName": volumeName, "hostname": hostName, "lun": lun}
	if optional != nil {
		for k, v := range optional {
			request[k] = v
		}
	}
	if auto {
		request["autoLun"] = true
		request["maxAutoLun"] = 0
		request["lun"] = 0
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/vluns", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateVLUN Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == 201 {
		location := resp.Header["Location"][0]
		lunStr := strings.Replace(location, "/api/v1/vluns/", "", -1)
		lunSplice := strings.Split(lunStr, ",")
		lunId, err = strconv.Atoi(lunSplice[1])
		return lunId, nil
	}
	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return lunId, err
			}
			lunId, err = T.CreateVLUN(requestId, volumeName, hostName, lun, auto, optional)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CreateVLUN:", err)
				return lunId, err
			}
			return lunId, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return lunId, err
			}

			if response.Code == HPE3ParSystemIsBusy {
				sleepTime := rand.Intn(5-1) + 1
				logger.Log.Debugf(requestId, "Waiting %d seconds to CreateVLUN again!", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)
				lunId, err = T.CreateVLUN(requestId, volumeName, hostName, lun, auto, optional)
				if err != nil {
					logger.Log.Error1(requestId, "Failed to CreateVLUN:", err)
					return lunId, err
				}
				return lunId, nil
			}

			logger.Log.Info1(requestId, "Failed to create vlun,three_par response:", response.String())
			err = errors.New("Failed to create vlun,three_par response:" + response.String())
			return lunId, err
		}
	}
	return
}

func (T *ThreeParDriver) DeleteVLUN(requestId, volumeName, hostName string, lunId float64) (err error) {
	requestUrl := fmt.Sprintf("%s/vluns/%s,%d,%s", T.ServerPath, volumeName, int64(lunId), hostName)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "DeleteVLUN Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.DeleteVLUN(requestId, volumeName, hostName, lunId)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to DeleteVLUN:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}
			if response.Code == 19 && response.Desc == "VLUN does not exist" {
				logger.Log.Warn("VLUN does not exist: ", response.String())
				return err
			}

			if response.Code == HPE3ParSystemIsBusy {
				sleepTime := rand.Intn(5-1) + 1
				logger.Log.Debugf(requestId, "Waiting %d seconds to DeleteVLUN again!", sleepTime)
				time.Sleep(time.Duration(sleepTime) * time.Second)
				err = T.DeleteVLUN(requestId, volumeName, hostName, lunId)
				if err != nil {
					logger.Log.Error1(requestId, "Failed to DeleteVLUN:", err)
					return err
				}
				return err
			}

			logger.Log.Info1(requestId, "Failed to remove VLUN,three_par response:", response.String())
			err = errors.New("Failed to remove VLUN,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) GetVlunsByHostname(requestId, hostName string) (vluns []VlunIfo, err error) {
	requestUrl := fmt.Sprintf("%s/vluns", T.ServerPath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	q := req.URL.Query()
	q.Add("query", "\"hostname\tEQ\t"+hostName+"\"")
	req.URL.RawQuery = q.Encode()
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetVlunsByHostname Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return nil, err
			}
			vluns, err = T.GetVlunsByHostname(requestId, hostName)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetVlunsByHostname:", err)
				return nil, err
			}
			return vluns, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return nil, err
			}

			logger.Log.Info1(requestId, "Failed to get vluns of host,three_par response:", response.String())
			err = errors.New("Failed to get vluns of host,three_par response:" + response.String())
			return nil, err
		}
	}
	if resp.Body != nil {
		var response map[string]interface{}
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return
		}
		if _, ok := response["members"]; ok {
			list, ok := response["members"].([]interface{})
			if ok {
				for _, v := range list {
					var vlun VlunIfo
					jsonBytesTemp, err := json.Marshal(v)
					if err != nil {
						logger.Log.Error1(requestId, "Failed to marshal json:", err)
						return nil, err
					}
					err = json.Unmarshal(jsonBytesTemp, &vlun)
					if err != nil {
						return nil, err
					}
					vluns = append(vluns, vlun)
				}
			}
		}
		//jsonBytes2, err := json.Marshal(response)
		//if err != nil {
		//	logger.Log.Error1(requestId,"Failed to marshal json:", err)
		//	return nil, err
		//}
		//logger.Log.Info1(requestId,"three_par response:", string(jsonBytes2))
	}
	return vluns, nil
}

//func (T *ThreeParDriver) GetVlun(requestId, hostName, volumeName string) (vlun VlunIfo, err error) {
//	vluns, err := T.GetVlunsByHostname(requestId, hostName)
//	if err != nil {
//		logger.Log.Error1(requestId, "Failed to GetVlunsByHostname,", err)
//		return vlun, err
//	}
//	if vluns != nil {
//		for _, value := range vluns {
//			if value.VolumeName == volumeName {
//				vlun = value
//				break
//			}
//		}
//	}
//	return vlun, nil
//}

//volumeSet methods
func (T *ThreeParDriver) CreateVolumeSet(requestId, name, domain, comment string, setMembers []string) (err error) {
	request := map[string]interface{}{"name": name}
	if domain != "" {
		request["domain"] = domain
	}
	if comment != "" {
		request["comment"] = comment
	}
	if setMembers != nil {
		request["setmembers"] = setMembers
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumesets", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateVolumeSet Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CreateVolumeSet(requestId, name, domain, comment, setMembers)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CreateVolumeSet:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}
			if response.Code == 101 && response.Desc == "Set exists" {
				logger.Log.Warnf(requestId, "Volume set %s has exits", name)
				return nil
			}

			logger.Log.Info1(requestId, "Failed to create volumeSet,three_par response:", response.String())
			err = errors.New("Failed to create volumeSet,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) DeleteVolumeSet(requestId, name string) (err error) {
	requestUrl := fmt.Sprintf("%s/volumesets/%s", T.ServerPath, name)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "DeleteVolumeSet Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return
	case http.StatusNotFound:
		return ErrSetDoesNotExist
	case http.StatusForbidden:
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.DeleteVolumeSet(requestId, name)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to DeleteVolumeSet:", err)
				return err
			}
		}
	}

	var response HttpResponseBody
	err = readJsonBody(resp.Body, &response)
	if err != nil {
		logger.Log.Error1(requestId, "ReadJsonBody error:", err)
		return err
	}

	logger.Log.Info1(requestId, "Failed to remove volumeset,three_par response:", response.String())
	err = errors.New("Failed to remove volumeset,three_par response:" + response.String())
	return err
}

func (T *ThreeParDriver) ModifyVolumeSet(requestId, name string, parameters map[string]interface{}) (err error) {
	jsonBytes, err := json.Marshal(parameters)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/volumesets/%s", T.ServerPath, name)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "ModifyVolumeSet Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return
	case http.StatusForbidden:
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.ModifyVolumeSet(requestId, name, parameters)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyVolumeSet:", err)
				return err
			}
			return err
		}
	default:
		var response HttpResponseBody

		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return err
		}

		if response.Code == HPE3ParSystemIsBusy {
			sleepTime := rand.Intn(5-1) + 1
			logger.Log.Debugf(requestId, "Waiting %d seconds to ModifyVolumeSet again!", sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			err = T.ModifyVolumeSet(requestId, name, parameters)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyVolumeSet:", err)
				return err
			}
			return err
		}

		switch response.Code {
		case NonExistentVolume:
			return ErrVolumeNotInSet
		case HaveExistentVolume:
			return ErrVolumeHasInSet
		default:
			logger.Log.Error1(requestId, "Failed to modify volumeset, three_par response:", response.String())
			return errors.New("Failed to modify volumeset,three_par response:" + response.String())
		}
	}

	return
}

func (T *ThreeParDriver) GetVolumeSet(requestId, name string) (volumeSet VolumeSetInfo, err error) {
	requestUrl := fmt.Sprintf("%s/volumesets/%s", T.ServerPath, name)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetVolumeSet Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return volumeSet, err
			}
			volumeSet, err = T.GetVolumeSet(requestId, name)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetVolumeSet:", err)
				return volumeSet, err
			}
			return volumeSet, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return volumeSet, err
			}

			logger.Log.Info1(requestId, "Failed to get volumeset,three_par response:", response.String())
			err = errors.New("Failed to get volumeset,three_par response:" + response.String())
			return volumeSet, err
		}
	}
	if resp.Body != nil {
		err = readJsonBody(resp.Body, &volumeSet)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return
		}
		jsonBytes2, err := json.Marshal(volumeSet)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return volumeSet, err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	return volumeSet, nil
}

//QoS Priority Optimization methods
func (T *ThreeParDriver) CreateQoSRules(requestId, targetName string, targetType int8, qosRules map[string]interface{}) (err error) {
	request := map[string]interface{}{"name": targetName, "type": targetType}
	if qosRules != nil {
		for k, v := range qosRules {
			request[k] = v
		}
	}
	jsonBytes, err := json.Marshal(request)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/qos", T.ServerPath)
	req, err := http.NewRequest("POST", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "CreateQoSRules Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 201 {
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.CreateQoSRules(requestId, targetName, targetType, qosRules)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to CreateQoSRules:", err)
				return err
			}
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return err
			}

			logger.Log.Info1(requestId, "Failed to create QoSRules,three_par response:", response.String())

			if response.Code == ExistentQoSRule {
				return ErrQoSRuleExistent
			}

			err = errors.New("Failed to create QoSRules,three_par response:" + response.String())
			return err
		}
	}
	return
}

func (T *ThreeParDriver) ModifyQoSRules(requestId, targetName string, targetType string, qosRules map[string]interface{}) (err error) {
	jsonBytes, err := json.Marshal(qosRules)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to marshal json:", err)
		return
	}
	requestUrl := fmt.Sprintf("%s/qos/%s:%s", T.ServerPath, targetType, targetName)
	req, err := http.NewRequest("PUT", requestUrl, bytes.NewReader(jsonBytes))
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "ModifyQoSRules Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return
	case http.StatusForbidden:
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.ModifyQoSRules(requestId, targetName, targetType, qosRules)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyQoSRules:", err)
				return err
			}
			return err
		}
	default:
		var response HttpResponseBody

		err = readJsonBody(resp.Body, &response)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return err
		}

		switch response.Code {
		case HPE3ParSystemIsBusy:
			sleepTime := rand.Intn(5-1) + 1
			logger.Log.Debugf(requestId, "Waiting %d seconds to ModifyQoSRules again!", sleepTime)
			time.Sleep(time.Duration(sleepTime) * time.Second)
			err = T.ModifyQoSRules(requestId, targetName, targetType, qosRules)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to ModifyQoSRules:", err)
				return err
			}
			return err
		case NonExistentQoSRule:
			logger.Log.Warnf(requestId, "QoS rule does not exist, response: %s", response.String())
			return ErrQosRuleDoesNotExist
		default:
			logger.Log.Error1(requestId, "Failed to modify QoS rule, three_par response:", response.String())
			return errors.New("Failed to modify QoS rule,three_par response:" + response.String())
		}
	}

	return
}

func (T *ThreeParDriver) DeleteQoSRules(requestId, targetName string, targetType string) (err error) {
	requestUrl := fmt.Sprintf("%s/qos/%s:%s", T.ServerPath, targetType, targetName)
	req, err := http.NewRequest("DELETE", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "DeleteQoSRules Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return
	case http.StatusForbidden:
		result, err := checkHttpResponseOfSessionKey(resp)
		if result {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return err
			}
			err = T.DeleteQoSRules(requestId, targetName, targetType)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to DeleteQoSRules:", err)
				return err
			}
		}
	}

	var response HttpResponseBody
	err = readJsonBody(resp.Body, &response)
	if err != nil {
		logger.Log.Error1(requestId, "ReadJsonBody error:", err)
		return err
	}

	logger.Log.Info1(requestId, "Failed to remove QoSRules,three_par response:", response.String())
	err = errors.New("Failed to remove QoSRules,three_par response:" + response.String())
	return err
}

func (T *ThreeParDriver) GetQoSRule(requestId, targetName string, targetType string) (result string, err error) {
	requestUrl := fmt.Sprintf("%s/qos/%s:%s", T.ServerPath, targetType, targetName)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "GetQoSRule Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		if resp.Body != nil {
			jsonBytes, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				return result, err
			}
			result = string(jsonBytes)
			logger.Log.Info1(requestId, "three_par response:", result)
		}
		return result, nil
	case http.StatusForbidden:
		result2, err := checkHttpResponseOfSessionKey(resp)
		if result2 {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return result, err
			}
			result, err = T.GetQoSRule(requestId, targetName, targetType)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to GetQoSRule:", err)
				return result, err
			}
			return result, nil
		}
	case http.StatusNotFound:
		return result, ErrQosRuleDoesNotExist
	}

	var response HttpResponseBody
	err = readJsonBody(resp.Body, &response)
	if err != nil {
		logger.Log.Error1(requestId, "ReadJsonBody error:", err)
		return result, err
	}

	if response.Code == HPE3ParSystemIsBusy {
		sleepTime := rand.Intn(5-1) + 1
		logger.Log.Debugf(requestId, "Waiting %d seconds to GetQoSRule again!", sleepTime)
		time.Sleep(time.Duration(sleepTime) * time.Second)
		result, err = T.GetQoSRule(requestId, targetName, targetType)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to DeleteVLUN:", err)
			return result, err
		}
		return result, err
	}

	logger.Log.Info1(requestId, "Failed to get Qos rule response:", response.String())
	err = errors.New("Failed to get Qos rule response:" + response.String())
	return result, err
}

//system information methods
func (T *ThreeParDriver) GetSystemCapacity() (result string, err error) {
	requestUrl := fmt.Sprintf("%s/capacity", T.ServerPath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error("NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info("GetSystemCapacity Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error("Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result2, err := checkHttpResponseOfSessionKey(resp)
		if result2 {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error("Failed to init three_par:", err)
				return result, err
			}
			result, err = T.GetSystemCapacity()
			if err != nil {
				logger.Log.Error("Failed to GetSystemCapacity:", err)
				return result, err
			}
			return result, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error("ReadJsonBody error:", err)
				return result, err
			}

			logger.Log.Info("Failed to get system capacity,three_par response:", response.String())
			err = errors.New("Failed to get system capacity,three_par response:" + response.String())
			return result, err
		}
	}
	if resp.Body != nil {
		jsonBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return result, err
		}
		result = string(jsonBytes)
		//logger.Log.Info1(requestId,"three_par response:", result)
	}
	return result, nil
}

// GetSystemUtilization 计算虚拟卷利用率
func (T *ThreeParDriver) GetSystemUtilization() (ssd, hdd float64, err error) {
	requestUrl := fmt.Sprintf("%s/systemreporter/attime/volumespacedata/hires;groupby:userCPG", T.ServerPath)
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error("NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info("GetSystemUtilization Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error("Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result2, err := checkHttpResponseOfSessionKey(resp)
		if result2 {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error("Failed to init three_par:", err)
				return ssd, hdd, err
			}
			ssd, hdd, err = T.GetSystemUtilization()
			if err != nil {
				logger.Log.Error("Failed to GetSystemUtilization:", err)
				return ssd, hdd, err
			}
			return ssd, hdd, nil
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error("ReadJsonBody error:", err)
				return ssd, hdd, err
			}

			logger.Log.Info("Failed to get volumes, three_par response:", response.String())
			err = errors.New("Failed to get volumes, three_par response:" + response.String())
			return ssd, hdd, err
		}
	}

	var volumeSpaceResponse AtTimeVolumeSpaceResponse
	err = readJsonBody(resp.Body, &volumeSpaceResponse)
	if err != nil {
		logger.Log.Error("ReadJsonBody error:", err)
		return ssd, hdd, err
	}

	var ssdTotalUsedMiB, ssdSizeMiB, hddTotalUsedMiB, hddSizeMiB float64
	for _, m := range volumeSpaceResponse.Members {
		switch m.UserCPG {
		case IMAGE_CPG, SSD_CPG, HYBRID_SSD_CPG:
			ssdTotalUsedMiB += m.TotalSpace.UsedMiB
			ssdSizeMiB += m.TotalSpace.VirtualSizeMiB
		case HYBRID_HDD_CPG:
			hddTotalUsedMiB += m.TotalSpace.UsedMiB
			hddSizeMiB += m.TotalSpace.VirtualSizeMiB
		default:
		}
	}

	if ssdSizeMiB > 0 {
		ssd = ssdTotalUsedMiB / ssdSizeMiB
	}

	if hddSizeMiB > 0 {
		hdd = hddTotalUsedMiB / hddSizeMiB
	}

	return ssd, hdd, nil
}

func (T *ThreeParDriver) getTask(requestId string, taskId float64) (task TaskInfo, err error) {
	requestUrl := fmt.Sprintf("%s/tasks/%d", T.ServerPath, int64(taskId))
	req, err := http.NewRequest("GET", requestUrl, nil)
	if err != nil {
		logger.Log.Error1(requestId, "NewRequest error:", err)
		return
	}
	req.Header.Set(sessionCookieName, T.SessionKey)
	logger.Log.Info1(requestId, "getTask Request to three_par:", requestUrl)
	resp, err := T.HttpClient.Do(req)
	if err != nil {
		logger.Log.Error1(requestId, "Failed to request three_par:", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		result2, err := checkHttpResponseOfSessionKey(resp)
		if result2 {
			err = T.InitSessionKey()
			if err != nil {
				logger.Log.Error1(requestId, "Failed to init three_par:", err)
				return task, err
			}
			task, err = T.getTask(requestId, taskId)
			if err != nil {
				logger.Log.Error1(requestId, "Failed to getTask:", err)
				return task, err
			}
			return task, err
		} else {
			var response HttpResponseBody
			err = readJsonBody(resp.Body, &response)
			if err != nil {
				logger.Log.Error1(requestId, "ReadJsonBody error:", err)
				return task, err
			}
			if response.Code == HPE3ParSystemIsBusy {
				logger.Log.Info1(requestId, "getTask,three_par response:", response.String())
				task := TaskInfo{
					Id:         0,
					Type:       0,
					Name:       "",
					Status:     2,
					StartTime:  "",
					FinishTime: "",
					User:       "",
				}
				return task, nil
			}
			logger.Log.Info1(requestId, "Failed to getTask,three_par response:", response.String())
			err = errors.New("Failed to getTask,three_par response:" + response.String())
			return task, err
		}
	}
	if resp.Body != nil {
		err = readJsonBody(resp.Body, &task)
		if err != nil {
			logger.Log.Error1(requestId, "ReadJsonBody error:", err)
			return task, err
		}
		jsonBytes2, err := json.Marshal(task)
		if err != nil {
			logger.Log.Error1(requestId, "Failed to marshal json:", err)
			return task, err
		}
		logger.Log.Info1(requestId, "three_par response:", string(jsonBytes2))
	}
	return task, nil
}

func checkHttpResponseOfSessionKey(resp *http.Response) (result bool, err error) {
	result = false
	if resp.StatusCode == 403 {
		var response HttpResponseBody
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			return
		}
		if response.Code == 6 {
			result = true
		}
	}
	return result, nil
}

func checkHostExist(resp *http.Response) (result bool, err error) {
	result = true
	if resp.StatusCode == 404 {
		var response HttpResponseBody
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			return
		}
		if response.Code == 17 {
			result = false
		}
	}
	return result, nil
}

func checkHttpResponseOfHostExist(resp *http.Response) (result bool, err error) {
	result = false
	if resp.StatusCode == 404 {
		var response HttpResponseBody
		err = readJsonBody(resp.Body, &response)
		if err != nil {
			return
		}
		if response.Code == 17 {
			result = true
		}
	}
	return result, nil
}

func readJsonBody(body io.ReadCloser, out interface{}) (err error) {
	defer func() {
		_ = body.Close()
	}()

	jsonBytes, err := ioutil.ReadAll(body)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return err
	}
	return nil
}
