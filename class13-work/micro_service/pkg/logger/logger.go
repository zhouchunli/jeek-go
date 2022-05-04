package logger

import (
	"micro_service/app/global"
	"micro_service/pkg/microerr"
	"net/http"
	"os"
	"time"

	"github.com/rs/zerolog"
)

const (
	Debug = zerolog.DebugLevel
	Info  = zerolog.InfoLevel
	Warn  = zerolog.WarnLevel
	Error = zerolog.ErrorLevel

	codeField          = "code"    // 错误码
	timestampField     = "ts"      // 日志时间-13位时间戳

	respField          = "resp"    // 响应
	callerServiceField = "caller"  // 调用方微服务名称
	pathField          = "path"    // grpc 请求路径
	stackField         = "st"      //堆栈信息

	serviceField = "svc"
	errorFiled   = "err"
	callerField  = "pos"
)

type Level = zerolog.Level
type Logger struct {
	conf Config
	l    zerolog.Logger
}
type Config struct {
	Glv        Level  // log的全局Level
	Svc        string // 微服务项目名称
	StackDepth int    // 堆栈深度
}

/*
 * 实例化一个到标准输出的Logger
 */
func New(conf Config) *Logger {
	// global setting
	zerolog.TimestampFieldName = timestampField
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnixMs
	zerolog.LevelFieldName = levelField
	zerolog.MessageFieldName = messageField
	zerolog.ErrorFieldName = errorFiled
	zerolog.CallerFieldName = callerField
	zerolog.DurationFieldUnit = time.Millisecond
	zerolog.SetGlobalLevel(conf.Glv)
	// instance
	z := zerolog.New(os.Stdout).With().Timestamp().Logger()
	logger := &Logger{l: z}
	return logger
}

func (log *Logger) HttpAccess(r *http.Request, traceId, header, req, reqRaw string, guid string, resp string, retCode int, dur time.Duration) {
	log.l.Log().
		Str(idxTagField, "http").
		Str(traceIdField, traceId).
		Str(guidField, guid).
		Dur(timeCostField, dur).
		Int(codeField, retCode).
		Str(ipField, r.RemoteAddr).
		Str(hostField, r.Host).
		Str(mdField, r.Method).
		Str(uriField, r.RequestURI).
		Str(headerField, header).
		Str(reqField, req).
		Str(reqRawField, reqRaw).
		Str(respField, resp).
		Msg("")
}

func (log *Logger) GrpcAccess(path, req, resp, caller, traceId, ip string, dur time.Duration) {
	log.l.Log().
		Str(idxTagField, "grpc").
		Str(traceIdField, traceId).
		Dur(timeCostField, dur).
		Str(ipField, ip).
		Str(callerServiceField, caller).
		Str(pathField, path).
		Str(reqField, req).
		Str(respField, resp).
		Msg("")
}

func (log *Logger) Event() {

}

func (log *Logger) MqttMsg(traceId, msg, ns, md, cid, ip, topic string) {
	log.l.Log().
		Str(idxTagField, "mqtt").
		Str(rgField, global.AppConfig.Common.ClusterRegion).
		Str(messageField, msg).
		Str(nsField, ns).
		Str(mdField, md).
		Str(cidField, cid).
		Str(traceIdField, traceId).
		Str(peerField, ip).
		Str(topicField, topic).
		Msg("")
}

func (log *Logger) ErrorRaw(err error, traceId, mkey string) {
	if sys, ok := err.(*microerr.SysErr); ok {
		// 系统内部错误
		code := sys.GetCode()
		stack := microerr.GetStackTrace(sys, log.conf.StackDepth)
		msg := sys.Error()
		log.Error(code, msg, stack, traceId, mkey)
	} else {
		// 外部错误
		code := int(microerr.NoTypeErr)
		stack := ""
		msg := err.Error()
		log.Error(code, msg, stack, traceId, mkey)
	}
}

func (log *Logger) Error(code int, msg, stack, traceId, mkey string) {
	log.l.Error().
		Str(idxTagField, "sys").
		Str(traceIdField, traceId).
		Str(rgField, global.AppConfig.Common.ClusterRegion).
		Str(mkeyField, mkey).
		Int(codeField, code).
		Str(stackField, stack).
		Str(messageField, msg).
		Msg("")
}

func (log *Logger) WarnRaw(err error, traceId, mkey string) {
	if sys, ok := err.(*microerr.SysErr); ok {
		// 系统内部错误
		code := sys.GetCode()
		stack := microerr.GetStackTrace(sys, log.conf.StackDepth)
		msg := sys.Error()
		log.Warn(code, msg, stack, traceId, mkey)
	} else {
		// 外部错误
		code := int(microerr.NoTypeErr)
		stack := ""
		msg := err.Error()
		log.Warn(code, msg, stack, traceId, mkey)
	}
}

func (log *Logger) Warn(code int, msg, stack, traceId, mkey string) {
	log.l.Warn().
		Str(idxTagField, "sys").
		Str(traceIdField, traceId).
		Str(rgField, global.AppConfig.Common.ClusterRegion).
		Str(mkeyField, mkey).
		Int(codeField, code).
		Str(stackField, stack).
		Str(messageField, msg).
		Msg("")
}

func (log *Logger) InfoRaw(err error, traceId, mkey string) {
	if sys, ok := err.(*microerr.SysErr); ok {
		// 系统内部错误
		code := sys.GetCode()
		stack := microerr.GetStackTrace(sys, log.conf.StackDepth)
		msg := sys.Error()
		log.Info(code, msg, stack, traceId, mkey)
	} else {
		// 外部错误
		code := int(microerr.NoTypeErr)
		stack := ""
		msg := err.Error()
		log.Info(code, msg, stack, traceId, mkey)
	}
}

func (log *Logger) Info(code int, msg, stack, traceId, mkey string) {
	log.l.Info().
		Str(idxTagField, "sys").
		Str(traceIdField, traceId).
		Str(rgField, global.AppConfig.Common.ClusterRegion).
		Str(mkeyField, mkey).
		Int(codeField, code).
		Str(stackField, stack).
		Str(messageField, msg).
		Msg("")
}

func (log *Logger) DebugRaw(err error, traceId, mkey string) {
	if sys, ok := err.(*microerr.SysErr); ok {
		// 系统内部错误
		code := sys.GetCode()
		stack := microerr.GetStackTrace(sys, log.conf.StackDepth)
		msg := sys.Error()
		log.Debug(code, msg, stack, traceId, mkey)
	} else {
		// 外部错误
		code := int(microerr.NoTypeErr)
		stack := ""
		msg := err.Error()
		log.Debug(code, msg, stack, traceId, mkey)
	}
}

func (log *Logger) Debug(code int, msg, stack, traceId, mkey string) {
	log.l.Debug().
		Str(idxTagField, "sys").
		Str(traceIdField, traceId).
		Str(rgField, global.AppConfig.Common.ClusterRegion).
		Str(mkeyField, mkey).
		Int(codeField, code).
		Str(stackField, stack).
		Str(messageField, msg).
		Msg("")
}

/*
 * TODO: 文件写入 & sample & trace id
 */
