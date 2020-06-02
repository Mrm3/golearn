package gredis

import (
	"encoding/json"
	"fmt"
	"immortality-demo/config"
	"reflect"
	"testing"
)

type Job struct {
	Action string
	Detail interface{}
}

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

func TestRedisSet(t *testing.T) {

	config.LoadConfig()
	config.SetupLogging()
	Setup()
	test := CreateDiskRequest{
		RequestId:    "1",
		DiskId:       "disk1",
		DiskCategory: "ceph",
		SnapshotId:   "snap1",
		ImageId:      "img1",
		Size:         0,
		StorageType:  "ceph",
		Qos:          "100",
		UserId:       "zq",
		ScheduleInfo: "100",
		ImageType:    "linux",
	}

	v := Job{
		Action: "createDisk",
		Detail: test,
	}
	Push("queue", v)
	i, err := Len("queue")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("the len of the creatdisk is ", i)
	var job Job
	tmp, _ := TryPop("queue")
	if tmp != nil {
		err = json.Unmarshal(tmp.([]byte), &job)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(reflect.TypeOf(job.Detail))

		data,err := json.Marshal(job.Detail)
		if err != nil {
			fmt.Println(err)
		}
		data2:= string(data)
		//fmt.Println(data)
		var req CreateDiskRequest
		err = json.Unmarshal([]byte(data2), &req)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(req)

	}else{
		fmt.Println("tmp is nil ")
	}

	//if job.Action != "" {
	//	fmt.Println(tmp)
	//}

}
