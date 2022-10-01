package logger

type (
	Logger interface {
		Error(LogData)
		Warn(LogData)
		Info(LogData)
	}

	LogData struct {
		Section string
		Message string
	}
)
