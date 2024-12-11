package logutils

import (
	"fmt"
	"os"
	"path"
	"time"
	"wallet_sdk/models"
	"wallet_sdk/utils/dir"

	zaprotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	level            zapcore.Level
	timeFormatPrefix string
)

func Log(archive string, zapConf models.Zap) (logger *zap.Logger) {
	if ok, _ := dir.PathExists(zapConf.Director); !ok { // 判断是否有Director文件夹
		fmt.Printf("create %v directory\n", zapConf.Director)
		_ = os.Mkdir(zapConf.Director, os.ModePerm)
	}

	switch zapConf.Level { // 初始化配置文件的Level
	case "debug":
		level = zap.DebugLevel
	case "info":
		level = zap.InfoLevel
	case "warn":
		level = zap.WarnLevel
	case "error":
		level = zap.ErrorLevel
	case "dpanic":
		level = zap.DPanicLevel
	case "panic":
		level = zap.PanicLevel
	case "fatal":
		level = zap.FatalLevel
	default:
		level = zap.InfoLevel
	}

	timeFormatPrefix = zapConf.Prefix

	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(archive, zapConf), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore(archive, zapConf))
	}
	if zapConf.ShowLine {
		logger = logger.WithOptions(zap.AddCaller())
	}
	logger = logger.WithOptions(zap.Hooks(ZapHook(zap.ErrorLevel)))

	return logger
}

// getEncoderConfig 获取zapcore.EncoderConfig
func getEncoderConfig(zapConf models.Zap) (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:     "message",
		LevelKey:       "level",
		TimeKey:        "time",
		NameKey:        "logger",
		CallerKey:      "caller",
		StacktraceKey:  zapConf.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     CustomTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}
	switch {
	case zapConf.EncodeLevel == "LowercaseLevelEncoder": // 小写编码器(默认)
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	case zapConf.EncodeLevel == "LowercaseColorLevelEncoder": // 小写编码器带颜色
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder
	case zapConf.EncodeLevel == "CapitalLevelEncoder": // 大写编码器
		config.EncodeLevel = zapcore.CapitalLevelEncoder
	case zapConf.EncodeLevel == "CapitalColorLevelEncoder": // 大写编码器带颜色
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

// getEncoder 获取zapcore.Encoder
func getEncoder(zapConf models.Zap) zapcore.Encoder {
	if zapConf.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig(zapConf))
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig(zapConf))
}

// getEncoderCore 获取Encoder的zapcore.Core
func getEncoderCore(archive string, zapConf models.Zap) (core zapcore.Core) {
	writer, err := GetWriteSyncer(archive, zapConf) // 使用file-rotatelogs进行日志分割
	if err != nil {
		fmt.Printf("Get Write Syncer Failed err:%v", err.Error())
		return
	}
	return zapcore.NewCore(getEncoder(zapConf), writer, level)
}

// 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format(timeFormatPrefix + "2006-01-02 15:04:05.000"))
}

// zap logger中加入file-rotatelogs
func GetWriteSyncer(archive string, zapConf models.Zap) (zapcore.WriteSyncer, error) {
	var filePath string
	if archive == "" {
		filePath = path.Join(zapConf.Director, "%Y-%m-%d.log")
	} else {
		filePath = path.Join(zapConf.Director, archive+"-%Y-%m-%d.log")
	}
	var linkName zaprotatelogs.Option
	if archive == "" {
		linkName = zaprotatelogs.WithLinkName(zapConf.LinkName)
	} else {
		linkName = zaprotatelogs.WithLinkName(zapConf.LinkName + "_" + archive)
	}

	fileWriter, err := zaprotatelogs.New(
		filePath,
		linkName,
		zaprotatelogs.WithMaxAge(7*24*time.Hour),
		zaprotatelogs.WithRotationTime(24*time.Hour),
	)
	if zapConf.LogInConsole {
		return zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(fileWriter)), err
	}
	return zapcore.AddSync(fileWriter), err
}

func ZapHook(level zapcore.Level) func(zapcore.Entry) error {
	return func(entry zapcore.Entry) error {
		return nil
	}
}

func LogInfof(logCli *zap.Logger, format string, a ...any) {
	logStr := fmt.Sprintf(format, a...)
	logCli.Info(logStr)
}

func LogErrorf(logCli *zap.Logger, format string, a ...any) {
	logStr := fmt.Sprintf(format, a...)
	logCli.Error(logStr)
}
