package logit

type Level int

type LevelConversion map[Level]string

const (
	LevelTrace Level = iota
	LevelDebug
	LevelInfo
	LevelWarn
	LevelError
	LevelCritical
	LevelEmergency
	LevelFatal
)

var defaultLevelString LevelConversion = map[Level]string{
	LevelTrace:     "TRACE",
	LevelDebug:     "DEBUG",
	LevelInfo:      "INFO",
	LevelWarn:      "WARN",
	LevelError:     "ERROR",
	LevelCritical:  "CRITICAL",
	LevelEmergency: "EMGERGENCY",
	LevelFatal:     "FATAL",
}
