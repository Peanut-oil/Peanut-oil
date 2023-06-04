package pkg

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"os"
	"time"
)

func InitLogger(logPath string) {
	exist, _ := pathExists(logPath)
	if !exist {
		err := os.MkdirAll(logPath, os.ModePerm)
		if err != nil {
			log.Fatal("init logger mkdir failed!", err)
		}
	}
	// 禁用logrus自带日志输出
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.InfoLevel)
	logInfo, err := rotatelogs.New(
		logPath+"/info.%Y%m%d_%H.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Fatal("初始化info日志失败")
	}
	logError, err := rotatelogs.New(
		logPath+"/error.%Y%m%d_%H.log",
		rotatelogs.WithMaxAge(7*24*time.Hour),
		rotatelogs.WithRotationTime(time.Hour),
	)
	if err != nil {
		log.Fatal("初始化error日志失败")
	}
	// 设置不同级别文件的输出目录
	lfHook := NewHook(
		WriterMap{
			logrus.DebugLevel: logInfo,
			logrus.InfoLevel:  logInfo,
			logrus.WarnLevel:  logInfo,
			logrus.ErrorLevel: logError,
			logrus.FatalLevel: logError,
			logrus.PanicLevel: logError,
		},
		&logrus.JSONFormatter{},
	)
	logrus.AddHook(lfHook)
}

func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}
