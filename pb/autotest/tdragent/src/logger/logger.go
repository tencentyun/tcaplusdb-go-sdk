package logger

import (
	"encoding/xml"
	"fmt"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io/ioutil"
	"os"
	"time"
)

type LogInterface interface {
	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type logCfg struct {
	XMLName       xml.Name `xml:"LogConfig"` // 指定最外层的标签为LogConfig
	LogPath       string   `xml:"LogPath"`
	LogName       string   `xml:"LogName"`
	LogLevel      string   `xml:"LogLevel"`      // DEBUG INFO WARN ERROR, default: INFO
	LogType       string   `xml:"LogType"`       // json console, default: console
	LogMaxSize    int      `xml:"LogMaxSize"`    // MB default:100 MB
	LogMaxBackups int      `xml:"LogMaxBackups"` // default 100
}

var Logger LogInterface
var LogConf *logCfg

func Init(cfgPath string) error {
	if Logger != nil {
		return nil
	}
	if len(cfgPath) > 0 {
		if err := setLogCfg(cfgPath); err != nil {
			return err
		}
	}

	//默认打印到./log
	if LogConf == nil {
		LogConf = &logCfg{
			LogPath:       "./log",
			LogName:       "TDRAgent",
			LogLevel:      "DEBUG",
			LogType:       "console",
			LogMaxSize:    100,
			LogMaxBackups: 100,
		}
	}

	errCore := zapcore.NewCore(getEncoder(), getErrLogWriter(), zapcore.ErrorLevel)
	switch LogConf.LogLevel {
	case "DEBUG":
		core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.DebugLevel)
		Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	case "INFO":
		core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.InfoLevel)
		Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	case "WARN":
		core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.WarnLevel)
		Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	case "ERROR":
		core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.ErrorLevel)
		Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	default:
		core := zapcore.NewCore(getEncoder(), getLogWriter(), zapcore.InfoLevel)
		Logger = zap.New(zapcore.NewTee(core, errCore), zap.AddCaller(), zap.AddCallerSkip(1)).Sugar()
	}
	return nil
}

func setLogCfg(cfgPath string) error {
	data, err := ioutil.ReadFile(cfgPath)
	if err != nil {
		fmt.Println("ReadFile " + cfgPath + " err:" + err.Error())
		return err
	}

	//读取的数据为json格式，需要进行解码
	LogConf = new(logCfg)
	err = xml.Unmarshal(data, LogConf)
	if err != nil {
		fmt.Println("json Unmarshal " + cfgPath + " err:" + err.Error())
		return err
	}

	//校验日志路径
	if len(LogConf.LogPath) > 0 {
		_, err := os.Stat(LogConf.LogPath)
		if err != nil {
			if os.IsNotExist(err) {
				//创建目录
				if err := os.MkdirAll(LogConf.LogPath, os.ModePerm); err != nil {
					fmt.Println("MkdirAll " + LogConf.LogPath + " err:" + err.Error())
					return err
				}
			} else {
				fmt.Println("Stat " + LogConf.LogPath + " err:" + err.Error())
				return err
			}
		}
	}

	if len(LogConf.LogName) == 0 {
		LogConf.LogName = "TcaplueApi"
	}

	if LogConf.LogLevel != "DEBUG" && LogConf.LogLevel != "INFO" &&
		LogConf.LogLevel != "WARN" && LogConf.LogLevel != "ERROR" {
		LogConf.LogLevel = "INFO"
	}

	if LogConf.LogType != "json" && LogConf.LogType != "console" {
		LogConf.LogType = "console"
	}

	//最大一个G日志文件
	if LogConf.LogMaxSize < 1 || LogConf.LogMaxSize > 1024 {
		LogConf.LogMaxSize = 100
	}

	if LogConf.LogMaxBackups <= 0 {
		LogConf.LogMaxBackups = 100
	}

	fmt.Println(*LogConf)
	return nil
}

func formatEncodeTime(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(fmt.Sprintf("%v", t.Format("2006-01-02T15:04:05.000000Z")))
}

func getEncoder() zapcore.Encoder {
	if LogConf.LogType == "json" {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     formatEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		return zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "t",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "trace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     formatEncodeTime,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		return zapcore.NewConsoleEncoder(encoderConfig)
	}
}

//Filename: 日志文件的位置
//MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
//MaxBackups：保留旧文件的最大个数
//MaxAges：保留旧文件的最大天数
//Compress：是否压缩/归档旧文件
func getLogWriter() zapcore.WriteSyncer {
	if len(LogConf.LogPath) == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   LogConf.LogPath + "/" + LogConf.LogName + ".log",
		MaxSize:    LogConf.LogMaxSize,
		MaxBackups: LogConf.LogMaxBackups,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

//Filename: 日志文件的位置
//MaxSize：在进行切割之前，日志文件的最大大小（以MB为单位）
//MaxBackups：保留旧文件的最大个数
//MaxAges：保留旧文件的最大天数
//Compress：是否压缩/归档旧文件
func getErrLogWriter() zapcore.WriteSyncer {
	if len(LogConf.LogPath) == 0 {
		return zapcore.AddSync(os.Stdout)
	}

	lumberJackLogger := &lumberjack.Logger{
		Filename:   LogConf.LogPath + "/" + LogConf.LogName + ".error",
		MaxSize:    LogConf.LogMaxSize,
		MaxBackups: LogConf.LogMaxBackups,
		Compress:   false,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func DEBUG(s string, args ...interface{}) {
	Logger.Debugf(s, args...)
}

func INFO(s string, args ...interface{}) {
	Logger.Infof(s, args...)
}

func WARN(s string, args ...interface{}) {
	Logger.Warnf(s, args...)
}

func ERR(s string, args ...interface{}) {
	Logger.Errorf(s, args...)
}
