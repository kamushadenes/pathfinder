package pathfinder

import (
	"io"
	"os"
	"runtime"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/x-cray/logrus-prefixed-formatter"
)

var Logger *logrus.Logger

var PanicLevel = logrus.PanicLevel
var FatalLevel = logrus.FatalLevel
var ErrorLevel = logrus.ErrorLevel
var WarnLevel = logrus.WarnLevel
var InfoLevel = logrus.InfoLevel
var DebugLevel = logrus.DebugLevel
var AllLevels = logrus.AllLevels

// WriterHook is a hook that writes logs of specified LogLevels to specified Writer
type LogHook struct {
	Writer    io.Writer
	LogLevels []logrus.Level
	Formatter logrus.Formatter
}

// Fire will be called when some logging function is called with current hook
// It will format log entry to string and write it to appropriate writer
func (hook *LogHook) Fire(entry *logrus.Entry) error {
	line, err := hook.Formatter.Format(entry)
	if err != nil {
		return err
	}
	_, err = hook.Writer.Write(line)
	return err
}

// Levels define on which log levels this hook would trigger
func (hook *LogHook) Levels() []logrus.Level {
	return hook.LogLevels
}

func GetLogLevels(subset string) []logrus.Level {
	switch subset {
	case "all":
		return AllLevels
	case "error":
		return []logrus.Level{ErrorLevel, FatalLevel, PanicLevel}
	case "warn":
		return []logrus.Level{WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	case "info":
		return []logrus.Level{InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	case "debug":
		return []logrus.Level{DebugLevel, InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
	case "infoOnly":
		return []logrus.Level{InfoLevel}
	case "infoWarn":
		return []logrus.Level{InfoLevel, WarnLevel}
	}

	return []logrus.Level{InfoLevel, WarnLevel, ErrorLevel, FatalLevel, PanicLevel}
}

func GetLogHook(writer io.Writer, logLevels []logrus.Level) *LogHook {
	return &LogHook{Writer: writer, LogLevels: logLevels, Formatter: &logrus.JSONFormatter{TimestampFormat: time.RFC3339}}
}

func GetColoredFormatter() *prefixed.TextFormatter {
	formatter := new(prefixed.TextFormatter)
	formatter.FullTimestamp = true
	formatter.SetColorScheme(&prefixed.ColorScheme{
		PrefixStyle:    "blue+b",
		TimestampStyle: "white+h",
	})

	return formatter
}

func GetLogger() *logrus.Logger {
	logger := logrus.New()

	logger.Formatter = GetColoredFormatter()
	logger.Level = logrus.InfoLevel

	return logger
}

func AddLogHook(hook *LogHook) {
	Logger.AddHook(hook)
}

func GetJsonLogger() *logrus.Logger {
	logger := logrus.New()

	formatter := new(logrus.JSONFormatter)

	logger.Formatter = formatter
	logger.Level = logrus.InfoLevel

	return logger
}

func GetPosition(frameSkip int) (string, string, int) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	var frame runtime.Frame
	for i := 0; i <= frameSkip; i++ {
		frame, _ = frames.Next()
	}
	return frame.Function, frame.File, frame.Line
}

func GetLogFields(m map[string]interface{}) map[string]interface{} {
	//function, file, line := GetPosition(1)

	nm := make(map[string]interface{})

	//nm["function"] = function
	//nm["position"] = fmt.Sprintf("%s:%d", file, line)

	for k, v := range m {
		nm[k] = v
	}

	return nm
}

func init() {
	Logger = GetLogger()

	if os.Getenv("DEBUG") != "" {
		Logger.SetLevel(DebugLevel)
	}
}
