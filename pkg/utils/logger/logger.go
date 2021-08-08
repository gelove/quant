package logger

import (
	"fmt"
	"time"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var sugar *zap.SugaredLogger

func init() {
	LogConf()
}

/**
 * 获取日志
 * filePath 日志文件路径
 * level 日志级别
 * maxSize 每个日志文件保存的最大尺寸 单位：M
 * maxBackups 日志文件最多保存多少个备份
 * maxAge 文件最多保存多少天
 * compress 是否压缩
 * serviceName 服务名
 */
func LogConf() {
	now := time.Now()
	hook := &lumberjack.Logger{
		Filename:   fmt.Sprintf("logs/%04d-%02d-%02d.log", now.Year(), now.Month(), now.Day()), // filePath
		MaxSize:    256,                                                                        // megabytes
		MaxBackups: 30,
		MaxAge:     30,    // days
		Compress:   false, // disabled by default
	}
	defer hook.Close()
	/*zap 的 Config 非常的繁琐也非常强大，可以控制打印 log 的所有细节，因此对于我们开发者是友好的，有利于二次封装。
	  但是对于初学者则是噩梦。因此 zap 提供了一整套的易用配置，大部分的姿势都可以通过一句代码生成需要的配置。
	*/
	enConfig := zap.NewProductionEncoderConfig() //生成配置

	// 时间格式
	enConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	w := zapcore.AddSync(hook)
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(enConfig), //编码器配置
		w,                                   //打印到控制台和文件
		zap.InfoLevel,                       //日志等级
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar = logger.Sugar()
	Info("logger start...")
}

func Debug(args ...interface{}) {
	sugar.Debug(args)
}

func Info(args ...interface{}) {
	sugar.Info(args)
}

func Warn(args ...interface{}) {
	sugar.Warn(args)
}

func Error(args ...interface{}) {
	sugar.Error(args)
}

func DPanic(args ...interface{}) {
	sugar.DPanic(args)
}

func Panic(args ...interface{}) {
	sugar.Panic(args)
}

func Fatal(args ...interface{}) {
	sugar.Fatal(args)
}

func Debugf(template string, args ...interface{}) {
	sugar.Debugf(template, args)
}

func Infof(template string, args ...interface{}) {
	sugar.Infof(template, args)
}

func Warnf(template string, args ...interface{}) {
	sugar.Warnf(template, args)
}

func Errorf(template string, args ...interface{}) {
	sugar.Errorf(template, args)
}

func DPanicf(template string, args ...interface{}) {
	sugar.DPanicf(template, args)
}

func Panicf(template string, args ...interface{}) {
	sugar.Panicf(template, args)
}

func Fatalf(template string, args ...interface{}) {
	sugar.Fatalf(template, args)
}
