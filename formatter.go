// Package easy allows to easily format output of Logrus logger
package easy

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

const (
	// Default log format will output [INFO]: 2006-01-02T15:04:05Z07:00 - Log message
	defaultLogFormat       = "[%lvl%]: %time% - %msg%"
	defaultTimestampFormat = time.RFC3339
)

// Formatter implements logrus.Formatter interface.
type Formatter struct {
	// Timestamp format
	TimestampFormat string
	// Available standard keys: time, msg, lvl
	// Also can include custom fields but limited to strings.
	// All of fields need to be wrapped inside %% i.e %time% %msg%
	LogFormat string
}

// Format building log message.
func (f *Formatter) Format(entry *logrus.Entry) ([]byte, error) {
	output := f.LogFormat
	if output == "" {
		output = defaultLogFormat
	}

	timestampFormat := f.TimestampFormat
	if timestampFormat == "" {
		timestampFormat = defaultTimestampFormat
	}

	output = strings.Replace(output, "%time%", entry.Time.Format(timestampFormat), 1)

	output = strings.Replace(output, "%msg%", entry.Message, 1)

	level := strings.ToUpper(entry.Level.String())
	output = strings.Replace(output, "%lvl%", level, 1)

	for k, val := range entry.Data {
		switch v := val.(type) {
		case string:
			if strings.Contains(output, "%"+k+"%") {
				output = strings.Replace(output, "%"+k+"%", v, 1)
			} else {
				output = output + fmt.Sprintf(" | %s=%s", k, v)
			}
		case int:
			s := strconv.Itoa(v)
			if strings.Contains(output, "%"+k+"%") {
				output = strings.Replace(output, "%"+k+"%", s, 1)
			} else {
				output = output + fmt.Sprintf(" | %s=%s", k, s)
			}
		case bool:
			s := strconv.FormatBool(v)
			if strings.Contains(output, "%"+k+"%") {
				output = strings.Replace(output, "%"+k+"%", s, 1)
			} else {
				output = output + fmt.Sprintf(" | %s=%s", k, s)
			}
		}
	}

	output = output + "\n"

	return []byte(output), nil
}
