package zap

import "GoVoteApi/pkg/logger"

// Error implements logger.Logger
func (l *logger_impl) Error(d logger.LogData) {
	l.z.Errorw(
		d.Message,
		"section", d.Section,
	)
}

// Info implements logger.Logger
func (l *logger_impl) Info(d logger.LogData) {
	l.z.Infow(
		d.Message,
		"section", d.Section,
	)
}

// Warn implements logger.Logger
func (l *logger_impl) Warn(d logger.LogData) {
	l.z.Warnw(
		d.Message,
		"section", d.Section,
	)
}
