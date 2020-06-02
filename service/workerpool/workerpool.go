package workerpool

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/coreos/pkg/capnslog"
	"immortality-demo/pkg/driver"
	"immortality-demo/pkg/gredis"
	"time"
)

var ulog = capnslog.NewPackageLogger("immortality", "workpool")

//New 注册工作池，并设置最大并发数
//new workpool and set the max number of concurrencies
func NewWorkPool(max int) *WorkerPool {
	if max < 1 {
		max = 1
	}

	p := &WorkerPool{
		Job: make(chan Job, 2*max),
	}

	go p.loop(max)
	return p
}

//SetTimeout 设置超时时间
func (p *WorkerPool) SetTimeout(timeout time.Duration) {
	p.timeout = timeout
}

func (p *WorkerPool) startQueue() {
	for {
		var job Job
		tmp, err := gredis.TryPop("queue")
		if err != nil {
			fmt.Println(err)
		}
		if tmp != nil {
			err = json.Unmarshal(tmp.([]byte), &job)
			if err != nil {
				fmt.Println(err)
			}
			p.Job <- job
		} else {
			time.Sleep(time.Second)
			fmt.Println("tmp is nil ")
		}
		//time.Sleep(time.Second)
	}
}

func testJob(job Job) error {
	fmt.Println(job.Action)
	if job.Action == "CreatDisk0" {
		return errors.New("error test")
	}
	return nil
}

func (p *WorkerPool) loop(maxWorkersCount int) {
	go p.startQueue() //Startup queue , 启动队列

	p.wg.Add(maxWorkersCount) // Maximum number of work cycles,最大的工作协程数
	//Start Max workers, 启动max个worker
	for i := 0; i < maxWorkersCount; i++ {
		go func() {
			defer p.wg.Done()
			// worker 开始干活
			for job := range p.Job {

				closed := make(chan struct{}, 1)
				// Set timeout, priority task timeout.有设置超时,优先task 的超时
				if p.timeout > 0 {
					ct, cancel := context.WithTimeout(context.Background(), p.timeout)
					go func() {
						select {
						case <-ct.Done():
							ulog.Error(ct.Err())
							cancel()
						case <-closed:
						}
					}()
				}
				//Points of Execution.真正执行的点
				//err := testJob(job)
				data, err := json.Marshal(job.Detail)
				if err != nil {
					fmt.Println(err)
				}
				err = driver.Dispatch(job.Action, string(data))
				if err != nil {
					fmt.Println(err)
					ulog.Error(err)
				}
				//time.Sleep(time.Second)
				close(closed)
				if err != nil {
					fmt.Println(err)
					ulog.Error(err)
				}
			}
		}()
	}
}
