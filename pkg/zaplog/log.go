package zaplog

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"kratos-layout/internal/conf"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"os"
	"path"
	"strings"
	"time"
)

var level zapcore.Level
var zapLogConfig *conf.Zap

func Zap(zapConf *conf.Zap) (logger *zap.Logger) {
	zapLogConfig = zapConf
	level = getLevel(zapLogConfig.Level)
	// 如果是 debug 或者是 error, 输出 stack trace
	if level == zap.DebugLevel || level == zap.ErrorLevel {
		logger = zap.New(getEncoderCore(), zap.AddStacktrace(level))
	} else {
		logger = zap.New(getEncoderCore())
	}
	return logger
}

func getEncoderCore() (core zapcore.Core) {
	encoder := getEncoder()
	infoWrite := getLogWriter("info")
	errorWrite := getLogWriter("error")

	infoLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})
	warnLevel := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	core = zapcore.NewTee(
		zapcore.NewCore(encoder, zapcore.AddSync(os.Stdout), zapcore.DebugLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(infoWrite), infoLevel),
		zapcore.NewCore(encoder, zapcore.AddSync(errorWrite), warnLevel),
	)
	return
}

func getEncoder() zapcore.Encoder {
	if zapLogConfig.Format == "json" {
		return zapcore.NewJSONEncoder(getEncoderConfig())
	}
	return zapcore.NewConsoleEncoder(getEncoderConfig())
}

func getLogWriter(level string) io.Writer {
	linkFile := path.Join(zapLogConfig.Director, zapLogConfig.LinkName)
	hook, err := rotatelogs.New(
		path.Join(zapLogConfig.Director, "%Y-%m-%d_"+level+".log"),
		rotatelogs.WithLinkName(linkFile),         // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(7*24*time.Hour),     // 文件最多保存期限，这里是7天
		rotatelogs.WithRotationTime(24*time.Hour), // 日志切割时间间隔
	)
	if err != nil {
		panic(err)
	}
	return hook
}

func getEncoderConfig() (config zapcore.EncoderConfig) {
	config = zapcore.EncoderConfig{
		MessageKey:    "message",
		LevelKey:      "level",
		TimeKey:       "time",
		NameKey:       "logger",
		CallerKey:     "caller",
		StacktraceKey: "stacktrace",
		LineEnding:    zapcore.DefaultLineEnding,
		//EncodeLevel:   zapcore.LowercaseLevelEncoder,
		EncodeTime: CustomTimeEncoder,
		//EncodeTime:     zapcore.RFC3339TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder,
	}

	switch zapLogConfig.EncodeLevel {
	case "LowercaseLevelEncoder":
		config.EncodeLevel = zapcore.LowercaseLevelEncoder // 小写编码器(默认)
	case "LowercaseColorLevelEncoder":
		config.EncodeLevel = zapcore.LowercaseColorLevelEncoder // 小写编码器带颜色
	case "CapitalLevelEncoder":
		config.EncodeLevel = zapcore.CapitalLevelEncoder // 大写编码器
	case "CapitalColorLevelEncoder":
		config.EncodeLevel = zapcore.CapitalColorLevelEncoder // 大写编码器带颜色
	default:
		config.EncodeLevel = zapcore.LowercaseLevelEncoder
	}
	return config
}

func getLevel(level string) (lv zapcore.Level) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		lv = zap.DebugLevel
	case "info":
		lv = zap.InfoLevel
	case "warn":
		lv = zap.WarnLevel
	case "error":
		lv = zap.ErrorLevel
	case "dpanic":
		lv = zap.DPanicLevel
	case "panic":
		lv = zap.PanicLevel
	case "fatal":
		lv = zap.FatalLevel
	default:
		lv = zap.InfoLevel
	}
	return
}

// CustomTimeEncoder 自定义日志输出时间格式
func CustomTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006/01/02 - 15:04:05"))
}
