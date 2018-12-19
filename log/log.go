package log

import (
	"fmt"
	"github.com/astaxie/beego/logs"
)

const (
	OutputConsole = "console"
	OutputFile    = "file"
)

func SetLogger(outputWay string, configs ...string) error {
	if outputWay == OutputConsole {
		return logs.SetLogger(outputWay, "")
	}
	if outputWay == OutputFile {
		if len(configs) == 0 {
			return fmt.Errorf("lack of filename for log")
		}
		return logs.SetLogger(OutputFile, fmt.Sprintf(`{"filename":"%v"}`, configs[0]))
	}

	return fmt.Errorf("unknown output way")
}

func Info(format string, v ...interface{}) {
	logs.Info(format, v...)
}

func Debug(format string, v ...interface{}) {
	logs.Debug(format, v...)
}

func Warn(format string, v ...interface{}) {
	logs.Warn(format, v...)
}

func Error(format string, v ...interface{}) {
	logs.Error(format, v...)
}

func Trace(format string, v ...interface{}) {
	logs.Trace(format, v...)
}
