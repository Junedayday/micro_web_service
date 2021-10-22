package zlog

import (
	"context"
	"fmt"

	"github.com/natefinch/lumberjack"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	// Logger为zap提供的原始日志，但使用起来比较烦，有强类型约束
	logger *zap.Logger
	// SugaredLogger为zap提供的一个通用性更好的日志组件，作为本项目的核心日志组件
	Sugar *zap.SugaredLogger
)

func Init(logPath string) {
	// 日志暂时只开放一个配置 - 配置文件路径，有需要可以后续开放
	hook := lumberjack.Logger{
		Filename: logPath,
	}
	w := zapcore.AddSync(&hook)

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		zap.InfoLevel,
	)

	logger = zap.New(core, zap.AddCaller())
	Sugar = logger.Sugar()
	return
}

// 命名和原生的Zap Log尽量一致，方便理解
func Sync() {
	logger.Sync()
}

func WithTrace(ctx context.Context) *zap.SugaredLogger {
	var jTraceId jaeger.TraceID
	if parent := opentracing.SpanFromContext(ctx); parent != nil {
		parentCtx := parent.Context()
		if tracer := opentracing.GlobalTracer(); tracer != nil {
			mySpan := tracer.StartSpan("my info", opentracing.ChildOf(parentCtx))
			if sc, ok := mySpan.Context().(jaeger.SpanContext); ok {
				jTraceId = sc.TraceID()
			}
			defer mySpan.Finish()
		}
	}

	return Sugar.With(zap.String(jaeger.TraceContextHeaderName, fmt.Sprint(jTraceId)))
}
