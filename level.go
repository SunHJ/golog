package golog

type Level int

const (
	Level_Debug Level = iota
	Level_Info
	Level_Warn
	Level_Error
	Level_Fatal
)

func LevelColorStr(l Level) string {
	var sLevel = LevelStr(l)
	var c = ColorFromLevel(l)
	return (logColorPrefix[c] + sLevel + logColorSuffix)
}

func LevelStr(l Level) string {
	var sLevel string
	switch l {
	case Level_Debug:
		sLevel = "[DEBUG]"
	case Level_Info:
		sLevel = "[INFO ]"
	case Level_Warn:
		sLevel = "[WARN ]"
	case Level_Error:
		sLevel = "[ERROR]"
	case Level_Fatal:
		sLevel = "[FATAL]"
	}
	return sLevel
}
