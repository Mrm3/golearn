package config

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/coreos/pkg/capnslog"
	rotatelogs "github.com/lestrrat/go-file-rotatelogs"
	"github.com/pkg/errors"
	"immortality-demo/pkg/file"
	"immortality-demo/pkg/util"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	// Copyright.
	Copyright         = "copyright 2020 unicloud.com"
	environmentPrefix = "IMMORTALITY_"
	logDirName        = "/var/log/immortality/"
	immortalityLog    = "/var/log/immortality/immortality.log"
)

var ulog = capnslog.NewPackageLogger("immortality", "config")

var configPath = flag.String("conf",
	"/etc/immortality/conf.toml", "config file path")

var Config *Configuration

type CephConfig struct {
	ConfPath string
	//	Id       string
	Category string // hdd/hybrid/ssd
	User     string
}

type Configuration struct {
	NodeID                 string
	DBPath                 string
	Listen                 string
	Ceph                   []CephConfig
	Amqp                   string
	ComputeServiceEndPoint string
	DeliveryCenter         string
	EbsCore                string
	RegionId               string
	ZoneId                 string
	KafkaAddresses         []string
	CourierAddresses       []string
	ThreeParJobPeriod      string //unit: second
	MonitorJobPeriod       string //unit: second
	ImageServiceEndPoint   string
	AmqPrefetchCount       int
	AmqPrefetchSize        int
	RedisHost              string
	RedisPassword          string
	RedisMaxIdle           int
	RedisMaxActive         int
	RedisIdleTimeout       time.Duration
}

var categoryMap = map[string]string{
	"hdd":    "hdd",
	"hybrid": "hybrid-hdd",
	"ssd":    "ssd",
}

var ErrClusterNotConfigured = errors.New(
	"ceph cluster for specified category not configured")

func CephClusterByCategory(requestCategory string) (cluster CephConfig, err error) {
	for _, c := range Config.Ceph {
		if requestCategory == categoryMap[c.Category] {
			return c, nil
		}
	}
	return CephConfig{}, ErrClusterNotConfigured
}

func DefaultConfiguration() *Configuration {
	cfg := &Configuration{
		// TODO: hostname is always `immortality` in k8s
		NodeID: util.GetHostName(),
		DBPath: "mysql://root:cloudos@10.0.47.235:3306/immortality2?parseTime=true&loc=Local",
		Listen: "0.0.0.0:10010",
		Ceph: []CephConfig{
			{
				ConfPath: "/etc/ceph/hdd.conf",
				Category: "hdd",
				User:     "admin",
			},
			{
				ConfPath: "/etc/ceph/ssd.conf",
				Category: "ssd",
				User:     "admin",
			},
			{
				ConfPath: "/etc/ceph/hybrid.conf",
				Category: "hybrid",
				User:     "admin",
			},
		},
		//Amqp:              "amqp://ct:123456@10.0.47.160:5672",
		Amqp:           "amqp://lyk:123456@10.0.47.162:5672",
		RegionId:       "cn-beijing",
		ZoneId:         "cn-beijing-a",
		KafkaAddresses: []string{"10.254.7.245:9092", "10.254.7.243:9092", "10.254.7.242:9092"},
		//KafkaAddresses:       []string{"10.0.47.165:9092"},
		CourierAddresses:     []string{},
		ThreeParJobPeriod:    "3600",
		MonitorJobPeriod:     "180",
		ImageServiceEndPoint: "http://10.254.7.230:10011",
		AmqPrefetchCount:     10,
		AmqPrefetchSize:      0,
		RedisHost:            "127.0.0.1:6379",
		RedisPassword:        "123456",
		RedisMaxIdle:         30,
		RedisMaxActive:       30,
		RedisIdleTimeout:     200,
	}
	return cfg
}

/*
e.g.
export 'IMMORTALITY_DB_PATH=mysql://root@127.0.0.1:3306/immortality?parseTime=true&loc=Local'
export 'IMMORTALITY_LISTEN=0.0.0.0:10010'
export 'IMMORTALITY_AMQP=amqp://guest:guest@localhost:5672'
export 'IMMORTALITY_CEPH_1=/etc/ceph/hdd.conf,1e082bd2-9faa-4438-8037-a03d37604b78,hdd,cas'
*/
func loadFromEnvironmentVariable() {
	for _, env := range os.Environ() {
		pair := strings.SplitN(env, "=", 2)
		if len(pair) != 2 {
			continue
		}
		name, value := pair[0], pair[1]
		if !strings.HasPrefix(name, environmentPrefix) {
			continue
		}
		name = strings.TrimPrefix(name, environmentPrefix)
		switch name {
		case "DB_PATH":
			Config.DBPath = value
		case "LISTEN":
			Config.Listen = value
		case "AMQP":
			Config.Amqp = value
		case "REGIONID":
			Config.RegionId = value
		case "ZONEID":
			Config.ZoneId = value
		case "COURIER_ADDRESSES":
			addresses := strings.Split(value, ",")
			Config.CourierAddresses = addresses
		case "KAFKA_ADDRESSES":
			addresses := strings.Split(value, ",")
			Config.KafkaAddresses = addresses
		case "ThreeParJobPeriod":
			Config.ThreeParJobPeriod = value
		case "MonitorJobPeriod":
			Config.MonitorJobPeriod = value
		case "ComputeServiceEndPoint":
			Config.ComputeServiceEndPoint = value
		case "DELIVERY_CENTER":
			Config.DeliveryCenter = value
		case "EBS_CORE":
			Config.EbsCore = value
		case "IMAGE_SERVICE_ENDPOINT":
			Config.ImageServiceEndPoint = value
		case "AMQ_PREFETCH_COUNT":
			Config.AmqPrefetchCount, _ = strconv.Atoi(value)
		case "AMQ_PREFETCH_SIZE":
			Config.AmqPrefetchSize, _ = strconv.Atoi(value)
		default:
			if strings.HasPrefix(name, "CEPH") {
				parts := strings.Split(value, ",")
				if len(parts) != 4 {
					_, _ = fmt.Fprintln(os.Stderr,
						"malformed config for env", name)
					continue
				}
				Config.Ceph = append(Config.Ceph, CephConfig{
					ConfPath: parts[0],
					//Id:       parts[1],
					Category: parts[2],
					User:     parts[3],
				})
			}
		}
	}
}

func LoadConfig() {
	Config = DefaultConfiguration()
	if _, err := os.Stat(*configPath); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "config file does not exist, skip config file")
		Config.Ceph = make([]CephConfig, 0)
		loadFromEnvironmentVariable()
		return
	}
	_, err := toml.DecodeFile(*configPath, &Config)
	if err != nil {
		fmt.Println(err)
		panic("Failed to decode config file,Please check\n" + err.Error())
	}
}

func SetupLogging() {
	capnslog.SetGlobalLogLevel(capnslog.DEBUG)
	file.IsNotExistMkDir(logDirName)
	writer := GetWriter(immortalityLog)
	capnslog.SetFormatter(capnslog.NewPrettyFormatter(writer, true))
}

func GetWriter(filename string) io.Writer {
	writer, err := rotatelogs.New(
		filename+".%Y-%m-%d",
		rotatelogs.WithLinkName(filename),         // 生成软链，指向最新日志文
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最大保存时间
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)

	if err != nil {
		ulog.Fatalf("config local file system logger error. %+v", errors.WithStack(err))
	}
	return writer
}
