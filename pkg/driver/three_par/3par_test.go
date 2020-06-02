package three_par_test

import (
	"immortality/service/driver/three_par"
	"immortality/service/logger"
	"net"
	"net/http"
	"os"
	"testing"
	"time"
)

const RequestId = "requestId"

func getThreePar() (threePar *three_par.ThreeParDriver) {

	stdoutLogger := logger.NewStdoutLogger(os.Stdout, "")
	logger.Log = stdoutLogger

	httpClient := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			MaxIdleConnsPerHost: 1024,
			DialContext: (&net.Dialer{
				Timeout: 3 * time.Second,
			}).DialContext,
			TLSHandshakeTimeout: 10 * time.Second,
		},
	}

	threePar = &three_par.ThreeParDriver{
		HttpClient: httpClient,
		ServerPath: "http://10.254.7.77:8008/api/v1",
		SessionKey: "",
		User:       "3paradm",
		Password:   "3pardata",
	}

	return threePar
}

func TestThreePar_Volume(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	var volumeName string = "Test-volumeName-112"
	err = T.CreateVolume(RequestId, volumeName, "ssd", 1024, map[string]interface{}{"snapCPG": "ssd"})
	if err != nil {
		t.Error("Failed to create volume!", err)
		return
	}
	volume, err := T.GetVolumeByName(RequestId, volumeName)
	if err != nil {
		t.Error("Failed to get volume!", err)
	}
	t.Log(volume)
	//var volumeName string = "Test-volumeName-1572258432"
	err = T.ModifyVolume(RequestId, volumeName, map[string]interface{}{"comment": "test case"})
	if err != nil {
		t.Error("Failed to modify volume!", err)
		return
	}
	//var volumeName string = "Test-volumeName-1572258432"
	err = T.GrowVolume(RequestId, volumeName, 1024)
	if err != nil {
		t.Error("Failed to grow volume!", err)
		return
	}
	err = T.CloneVolume(RequestId, volumeName, "testCopyVolume3", "hdd", map[string]interface{}{"online": true})
	if err != nil {
		t.Error("Failed to modify volume!", err)
		return
	}
	err = T.DeleteVolume(RequestId, volumeName)
	if err != nil {
		t.Errorf("Failed to delete volume!")
		return
	}
	err = T.DeleteVolume(RequestId, "testCopyVolume3")
	if err != nil {
		t.Errorf("Failed to delete volume!")
		return
	}
}

func TestThreePar_Snapshot(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	var volumeName string = "Test-volumeName-113"
	err = T.CreateVolume(RequestId, volumeName, "ssd", 1024, map[string]interface{}{"snapCPG": "ssd"})
	if err != nil {
		t.Error("Failed to create volume!", err)
		return
	}
	var snapshotName string = "Test-snapshotName-001"
	err = T.CreateVolumeSnapshot(RequestId, snapshotName, volumeName, nil)
	if err != nil {
		t.Error("Failed to create snapshot of volume!", err)
		return
	}
	volumeSnaps, err := T.GetSnapshotsOfVolume(RequestId, "ssd", volumeName)
	if err != nil {
		t.Error("Failed to get snapshots of volume!", err)
	}
	t.Log(volumeSnaps)
	//var snapshotName string = "Test-snapshotName-1572260546"
	err = T.PromoteVirtualCopy(RequestId, snapshotName, nil)
	if err != nil {
		t.Error("Failed to create snapshot of volume!", err)
		return
	}
	err = T.RemoveSnapshot(RequestId, snapshotName)
	if err != nil {
		t.Error("Failed to delete snapshot of volume!", err)
		return
	}
	err = T.DeleteVolume(RequestId, volumeName)
	if err != nil {
		t.Errorf("Failed to delete volume!")
		return
	}
}

func TestThreePar_Host(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	err = T.CreateHost(RequestId, "test-case-host-1", []string{"iqn.2019-06.com.example.com:desktop"}, nil, nil)
	if err != nil {
		t.Error("Failed to CreateHost,", err)
		return
	}

	err = T.ModifyHost(RequestId, "test-case-host-1", map[string]interface{}{"pathOperation": 1, "iSCSINames": []string{"iqn.2019-07.com.example.com:desktop"}})
	if err != nil {
		t.Error("Failed to ModifyHost,", err)
		return
	}

	host, err := T.GetHostByName(RequestId, "test-case-host-1")
	if err != nil {
		t.Error("Failed to GetHostByName,", err)
		return
	}
	t.Log(host)
	host2, err := T.GetHostByIqn(RequestId, "iqn.2018-05.com.example.com:desktop")
	if err != nil {
		t.Error("Failed to GetHostByIqn,", err)
		return
	}
	t.Log(host2)

	err = T.DeleteHost(RequestId, "test-case-host-1")
	if err != nil {
		t.Error("Failed to DeleteHost,", err)
		return
	}

}

func TestThreePar_VLUN(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	var volumeName string = "storage_uca_volume_test1"
	err = T.CreateVolume(RequestId, volumeName, "ssd", 1024, map[string]interface{}{"snapCPG": "ssd"})
	if err != nil {
		t.Error("Failed to create volume!", err)
		return
	}

	var hostName string = "storage_uca_host_test1"
	err = T.CreateHost(RequestId, hostName, []string{"iqn.1901-08.org.debian:01:e15349fa9211"}, nil, nil)
	if err != nil {
		t.Error("Failed to CreateHost,", err)
		return
	}

	lunId, err := T.CreateVLUN(RequestId, volumeName, hostName, 1, true, nil)
	if err != nil {
		t.Error("Failed to CreateHost,", err)
		return
	}
	t.Log(lunId)

	vluns, err := T.GetVlunsByHostname(RequestId, hostName)
	if err != nil {
		t.Error("Failed to ModifyHost,", err)
		return
	}
	t.Log(vluns)

	err = T.DeleteVLUN(RequestId, volumeName, hostName, 0)
	if err != nil {
		t.Error("Failed to DeleteHost,", err)
		return
	}

	err = T.DeleteHost(RequestId, hostName)
	if err != nil {
		t.Error("Failed to DeleteHost,", err)
		return
	}

	err = T.DeleteVolume(RequestId, volumeName)
	if err != nil {
		t.Errorf("Failed to delete volume!")
		return
	}

}

func TestThreePar_VolumeSet(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	err = T.CreateVolumeSet(RequestId, "TestCase-volumeSet-1", "", "", nil)
	if err != nil {
		t.Error("Failed to CreateVolumeSet,", err)
		return
	}
	volumeSet, err := T.GetVolumeSet(RequestId, "TestCase-volumeSet-1")
	if err != nil {
		t.Error("Failed to GetVolumeSet,", err)
		return
	}
	t.Log(volumeSet)
	//err = T.ModifyVolumeSet("TestCase-volumeSet-1", map[string]interface{}{"action": 1, "setmembers": []string{"testCopyVolume3"}})
	//if err != nil {
	//	t.Error("Failed to ModifyVolumeSet,", err)
	//	return
	//}
	err = T.DeleteVolumeSet(RequestId, "TestCase-volumeSet-1")
	if err != nil {
		t.Error("Failed to DeleteVolumeSet,", err)
		return
	}
}

func TestThreePar_QoSRules(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	err = T.CreateVolumeSet(RequestId, "TestCase-volumeSet-1", "", "", nil)
	if err != nil {
		t.Error("Failed to CreateVolumeSet,", err)
		return
	}

	err = T.CreateQoSRules(RequestId, "TestCase-volumeSet-1", 1, map[string]interface{}{"bwMinGoalKB": 1024, "bwMaxLimitKB": 1024})
	if err != nil {
		t.Error("Failed to CreateQoSRules,", err)
		return
	}
	result, err := T.GetQoSRule(RequestId, "TestCase-volumeSet-1", "vvset")
	if err != nil {
		t.Error("Failed to GetQoSRule,", err)
		return
	}
	t.Log(result)
	err = T.ModifyQoSRules(RequestId, "TestCase-volumeSet-1", "vvset", map[string]interface{}{"bwMinGoalKB": 2048, "bwMaxLimitKB": 2048})
	if err != nil {
		t.Error("Failed to CreateQoSRules,", err)
		return
	}
	result2, err := T.GetQoSRule(RequestId, "TestCase-volumeSet-1", "vvset")
	if err != nil {
		t.Error("Failed to GetQoSRule,", err)
		return
	}
	t.Log(result2)
	err = T.DeleteQoSRules(RequestId, "TestCase-volumeSet-1", "vvset")
	if err != nil {
		t.Error("Failed to DeleteQoSRules,", err)
		return
	}

	err = T.DeleteVolumeSet(RequestId, "TestCase-volumeSet-1")
	if err != nil {
		t.Error("Failed to DeleteVolumeSet,", err)
		return
	}
}

func TestThreePar_GetSystemCapacity(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	result, err := T.GetSystemCapacity()
	if err != nil {
		t.Error("Failed to GetSystemCapacity,", err)
		return
	}
	t.Log(result)
}

func TestThreeParDriver_CreateVolume(t *testing.T) {
	T := getThreePar()
	err := T.InitSessionKey()
	if err != nil {
		t.Error("Failed to initialize three_par client,", err)
		return
	}

	var volumeName string = "Test-volumeName-119"
	err = T.CreateVolume(RequestId, volumeName, "hybrid-ssd", 1024, map[string]interface{}{"snapCPG": "hybrid-ssd"})
	if err != nil {
		t.Error("Failed to create volume!", err)
		return
	}
	volume, err := T.GetVolumeByName(RequestId, volumeName)
	if err != nil {
		t.Error("Failed to get volume!", err)
	}
	t.Log(volume)

}

//func TestThreeParDriver_GetVlun(t *testing.T) {
//	T := getThreePar()
//	err := T.InitSessionKey()
//	if err != nil {
//		t.Error("Failed to initialize three_par client,", err)
//		return
//	}
//
//	vlun, err := T.GetVlun(RequestId, "cvkDXN", "dxn_2")
//	if err != nil {
//		t.Error("Failed to GetVlun!", err)
//		return
//	}
//	t.Log(vlun)
//
//}
