package initialize

import (
	"fmt"
	"makedb/global"
	"os"
	"path"

	"makedb/utils"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Zap() *zap.Logger {
	if ok, _ := utils.PathExists(global.MAKEDB_CONFIG.Zap.Director); !ok { // 判断是否有Director文件夹
		// fmt.Printf("create %v directory\n", global.MAKEDB_CONFIG.Zap.Director)
		_ = os.Mkdir(global.MAKEDB_CONFIG.Zap.Director, os.ModePerm)
	}

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:    "time",
		LevelKey:   "level",
		NameKey:    "logger",
		CallerKey:  "caller",
		MessageKey: "msg",
		// StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,    // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.FullCallerEncoder, // 全路径编码器
	}

	// 设置日志级别
	atom := zap.NewAtomicLevelAt(zap.DebugLevel)

	config := zap.Config{
		Level:         atom,                               // 日志级别
		Development:   true,                               // 开发模式，堆栈跟踪
		Encoding:      global.MAKEDB_CONFIG.Zap.LogFormat, // 输出格式 console 或 json
		EncoderConfig: encoderConfig,                      // 编码器配置
		// InitialFields: map[string]interface{}{"serviceName": "makedb"}, // 初始化字段，如：添加一个服务器名称
		OutputPaths: []string{
			"stdout",
			path.Join(global.MAKEDB_CONFIG.Zap.Director, global.MAKEDB_CONFIG.Zap.LogFile),
		}, // 输出到指定文件 stdout（标准输出，正常颜色） stderr（错误输出，红色）
		ErrorOutputPaths: []string{"stderr"},
	}

	// 构建日志
	logger, err := config.Build()
	if err != nil {
		panic(fmt.Sprintf("log init failed: %v", err))
	}
	return logger
}
