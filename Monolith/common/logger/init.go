package logger

import (
	"io/ioutil"
	"os"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v2"
)

type LoggerConfig struct {
	Level         int    `yaml:"log_level"`
	Format        string `yaml:"log_format"`
	TimeKey       string `yaml:"log_time_key"`
	LevelKey      string `yaml:"log_level_key"`
	NameKey       string `yaml:"log_name_key"`
	CallerKey     string `yaml:"log_caller_key"`
	MessageKey    string `yaml:"log_message_key"`
	StacktraceKey string `yaml:"log_stacktrace_key"`
	Path          string `yaml:"log_path"`
	MaxSize       int    `yaml:"log_max_size"`
	MaxBackUps    int    `yaml:"log_max_backups"`
	MaxAge        int    `yaml:"log_max_age"`
}

func readLoggerConfig() *LoggerConfig {
	config := new(LoggerConfig)
	file, _ := ioutil.ReadFile("config.yaml")
	yaml.Unmarshal(file, config)
	return config
}

// 设置日志级别、输出格式和日志文件的路径
func init() {
	config := readLoggerConfig()
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        config.TimeKey,
		LevelKey:       config.LevelKey,
		NameKey:        config.NameKey,
		CallerKey:      config.CallerKey,
		MessageKey:     config.MessageKey,
		StacktraceKey:  config.StacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,    // 大写编码器
		EncodeTime:     zapcore.ISO8601TimeEncoder,     // ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     // 短路径编码器(相对路径+行号)
		EncodeName:     zapcore.FullNameEncoder,
	}

	// 设置日志输出格式
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 添加日志切割归档功能
	hook := lumberjack.Logger{
		Filename:   config.Path,       // 日志文件路径
		MaxSize:    config.MaxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: config.MaxBackUps, // 日志文件最多保存多少个备份
		MaxAge:     config.MaxAge,     // 文件最多保存多少天
		Compress:   true,              // 是否压缩
	}

	core := zapcore.NewCore(
		encoder, // 编码器配置
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stderr), zapcore.AddSync(&hook)), // 打印到控制台和文件
		zap.NewAtomicLevelAt(zapcore.Level(config.Level)),                               // 日志级别
	)

	// 开启文件及行号
	caller := zap.AddCaller()
	stackTrace := zap.AddStacktrace(zapcore.ErrorLevel)

	logger := zap.New(core, caller, stackTrace)

	// 将自定义的logger替换为全局的logger
	zap.ReplaceGlobals(logger)
}
