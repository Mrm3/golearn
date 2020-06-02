package workerpool

import (
	"fmt"
	"immortality-demo/config"
	"immortality-demo/pkg/gredis"
	"testing"
)

func TestWorkerPoolError(t *testing.T) {

	config.LoadConfig()
	config.SetupLogging()
	gredis.Setup()
	NewWorkPool(3) //Set the maximum number of threads，设置最大线程数
	//
	//
	var jobs []Job

	for i := 0; i < 10; i++ {
		job := Job{
			Action: fmt.Sprintf("CreatDisk%d", i),
			Detail: nil,
		}
		jobs = append(jobs, job)
	}
	for _, v := range jobs {
		gredis.Push("queue", v)
	}

	var res chan bool
	<-res

	fmt.Println("down")
}
