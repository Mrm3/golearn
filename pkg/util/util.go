package util

import (
	"fmt"
	"github.com/coreos/pkg/capnslog"
	"immortality-demo/pkg/data"
	"immortality-demo/pkg/db_model"
	"math"
	"math/rand"
	"os"
	"os/exec"
	"time"
)

var ulog = capnslog.NewPackageLogger("immortality", "util")

func GetHostName() (name string) {
	name, err := os.Hostname()
	if err != nil {
		panic(err)
	}
	return name
}

var alphaNumericTable = []byte("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ")

func GenerateRandomId() string {
	alpha := make([]byte, 16, 16)
	for i := 0; i < 16; i++ {
		n := rand.Intn(len(alphaNumericTable))
		alpha[i] = alphaNumericTable[n]
	}
	return string(alpha)
}

var numericTable = []byte("0123456789")

func CatchPanic(status string) {
	if err := recover(); err != nil {
		errMsg := fmt.Sprintf("occur panic:%v", err)
		ulog.Error(errMsg)
		return
	}
}

// QosByCapacity calculate the max bandwidth and iops by capacity(GiB) and specification
func QosByCapacity(capacity int64, spec string) (bwMaxLimitKB uint64, ioMaxLimit uint32, err error) {
	if capacity > data.HPE3PAR_MAX_DISK_SIZE {
		err = data.ErrInvalidSize
		return
	}

	item, ok := data.Cache.Load(spec)
	if !ok {
		err = data.ErrInvalidInstanceCode
		return
	}

	diskSpec := item.(db_model.DiskSpec)
	ioBase := float64(diskSpec.IOPSBase)
	ioFactor := diskSpec.IOPSFactor
	ioMax := float64(diskSpec.IOPSMax)
	bwBase := float64(diskSpec.BandwidthBase)
	bwFactor := diskSpec.BandwidthFactor
	bwMax := float64(diskSpec.BandwidthMax)

	ioMaxLimit = uint32(math.Min(ioBase+ioFactor*float64(capacity), ioMax))
	bwMaxLimitKB = uint64(math.Min(bwBase+bwFactor*float64(capacity), bwMax) * 1024)

	return
}

type stop struct {
	error
}

func Retry(attempts int, sleep time.Duration, fn func() error) error {
	if err := fn(); err != nil {
		if e, ok := err.(stop); ok {
			return e.error
		}
		if attempts--; attempts > 0 {
			ulog.Warningf("retry func error: %s. attempts #%d after %s.", err.Error(), attempts, sleep)
			time.Sleep(sleep)
			return Retry(attempts, sleep, fn)
		}
		return err
	}
	return nil
}

func ExecCmd(cmd string) (string, error) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	flag := r.Intn(100000000)
	ulog.Infof("starting execute cmd:[%s], flag:%d\n", cmd, flag)
	output, err := exec.Command("bash", "-c", cmd).Output()
	if err != nil {
		return "", err
	}
	context := ""
	if len(output) > 0 {
		ulog.Infof("flag %d execute output:%s", flag, output)
		context = fmt.Sprintf("%s", output)
	}
	return context, nil
}
