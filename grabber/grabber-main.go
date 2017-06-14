package grabber

import (
	"fmt"
	"strings"
	"time"

	"github.com/lfkeitel/inca/common"
	"github.com/lfkeitel/inca/targz"
	"github.com/lfkeitel/verbose"
)

var appLogger *verbose.Logger
var stdOutLogger *verbose.Logger
var configGrabRunning = false
var conf common.Config

var totalDevices = 0
var finishedDevices = 0

func init() {
	configGrabRunning = false

	appLogger = verbose.New("grabber")
	stdOutLogger = verbose.New("execStdOut")

	fileLogger, err := verbose.NewFileHandler("logs/main/")
	if err != nil {
		panic("Failed to open logging directory")
	}

	appLogger.AddHandler("file", fileLogger)
	stdOutLogger.AddHandler("file", fileLogger)
}

func LoadConfig(config common.Config) {
	conf = config
	return
}

func PerformConfigGrab() {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() { configGrabRunning = false }()

	// Clean up tftp directory
	removeDir(conf.FullConfDir)

	hosts, err := loadDeviceList(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	dtypes, err := loadDeviceTypes(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = len(hosts)
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dtypes, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	common.UserLogInfo(logText)
	return
}

func PerformSingleRun(name, hostname, brand, method string) {
	if configGrabRunning {
		appLogger.Error("Job already running")
		return
	}

	startTime := time.Now()
	configGrabRunning = true
	defer func() { configGrabRunning = false }()
	name = strings.Replace(name, "-", "_", -1)

	hosts := make([]host, 1)

	hosts[0] = host{
		name:    name,
		address: hostname,
		dtype:   brand,
		method:  method,
	}

	dtypes, err := loadDeviceTypes(conf)
	if err != nil {
		appLogger.Error(err.Error())
		return
	}

	totalDevices = 1
	finishedDevices = 0
	dateSuffix := time.Now().Format("2006012")

	grabConfigs(hosts, dtypes, dateSuffix, conf)
	tarGz.TarGz("archive/"+dateSuffix+".tar.gz", conf.FullConfDir)

	endTime := time.Now()
	logText := fmt.Sprintf("Config grab took %s", endTime.Sub(startTime).String())
	appLogger.Info(logText)
	common.UserLogInfo(logText)
	return
}

func IsRunning() bool {
	return configGrabRunning
}

func Remaining() (total, finished int) {
	if !configGrabRunning {
		if totalDevices == 0 {
			hosts, err := loadDeviceList(conf)
			if err != nil {
				appLogger.Error(err.Error())
				return
			}
			totalDevices = len(hosts)
		}

		if finishedDevices == 0 {
			finishedDevices = -1
		}
	}

	return totalDevices, finishedDevices
}
