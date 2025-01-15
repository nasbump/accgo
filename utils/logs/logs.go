package logs

import (
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

var logHnd zerolog.Logger

func init() {
	curProcID := strconv.Itoa(os.Getpid())

	zerolog.LevelFieldName = "L"
	zerolog.CallerFieldName = "F"
	zerolog.MessageFieldName = "Msg"
	zerolog.TimestampFieldName = "T"
	zerolog.FloatingPointPrecision = 2
	zerolog.TimeFieldFormat = "20060102-150405.000"
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
		var b strings.Builder
		b.WriteString(curProcID + "/")
		b.WriteString(filepath.Base(file))
		b.WriteString(":" + strconv.Itoa(line))
		return b.String()
	}

	// 默认日志对象
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	consoleWriter := zerolog.ConsoleWriter{Out: zerolog.SyncWriter(os.Stdout), NoColor: true, TimeFormat: zerolog.TimeFieldFormat}
	logHnd = zerolog.New(consoleWriter).With().Timestamp().Caller().Logger()
}

func LogsInit(logLev int) {
	zerolog.SetGlobalLevel(zerolog.Level(logLev)) //zerolog.DebugLevel)
}

type zlogger struct {
	*zerolog.Event
}

func Trace() *zlogger {
	z := &zlogger{logHnd.Trace()}
	return z
}
func Debug() *zlogger {
	z := &zlogger{logHnd.Debug()}
	return z
}
func Info() *zlogger {
	z := &zlogger{logHnd.Info()}
	return z
}
func Warn(err error) *zlogger {
	z := &zlogger{logHnd.Warn()}
	return z.Fail(err)
}
func Error(err error) *zlogger {
	z := &zlogger{logHnd.Error()}
	return z.Fail(err)
}
func Fatal(err error) *zlogger {
	z := &zlogger{logHnd.Fatal()}
	return z.Fail(err)
}
func Panic(err error) *zlogger {
	z := &zlogger{logHnd.Panic()}
	return z.Fail(err)
}
func (z *zlogger) Fail(err error) *zlogger {
	if err != nil {
		z.AnErr("err", err)
	}

	return z
}
func Catch(err error) *zlogger {
	const size = 16 << 10
	buf := make([]byte, size)
	buf = buf[:runtime.Stack(buf, false)]
	var sb strings.Builder
	sb.Write(buf)
	if err != nil {
		log.Printf("err: %v\n", err)
	}
	log.Printf("stack: %s\n", sb.String())

	z := &zlogger{logHnd.Fatal()}
	z.Str("stack", sb.String())
	return z
}
