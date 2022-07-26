package logstream

import "fmt"

// LogBlock represents a log message with a channel id that logs are coming from.
type LogBlock struct {
	ChanID string
	Msg    string
}

func (l *LogBlock) String() string {
	return fmt.Sprintf("[%s] %s", l.ChanID, l.Msg)
}
