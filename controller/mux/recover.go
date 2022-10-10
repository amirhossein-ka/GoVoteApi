package mux

import (
	dto "GoVoteApi/DTO"
	"GoVoteApi/pkg/logger"
	"log"
	"net/http"
)

type recovery struct {
	logger logger.Logger
}

func (r *recovery) recoverHttp(w http.ResponseWriter) {
	if err := recover(); err != nil {
		switch e := err.(type) {
		case error:
			writeJson(w, http.StatusInternalServerError, dto.Error{Status: dto.StatusError, Error: e.Error()})
			//r.logger.Error(logger.LogData{Message: e.Error(), Section: "unknown"})
		case string:
			writeJson(w, http.StatusInternalServerError, dto.Error{Status: dto.StatusError, Error: e})
			r.logger.Error(logger.LogData{Message: e, Section: "unknown"})
		case *Error:
			e.handle(r.logger, w)
		case Error:
			(&e).handle(r.logger, w)
		default:
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("err: %#v\n", e)
		}
	}
}

type LogLevel uint8

const (
	InfoLevel LogLevel = iota
	WarnLevel
	ErrorLevel
)

type Error struct {
	Err        error
	Section    string
	StatusCode int
	LogLevel   LogLevel
	Log        bool
}

func (e *Error) handle(l logger.Logger, w http.ResponseWriter) {
	switch e.LogLevel {
	case InfoLevel:
		writeJson(w, e.StatusCode, dto.Error{Status: dto.StatusError, Error: e.Err.Error()})
		if e.Log {
			l.Info(logger.LogData{Section: e.Section, Message: e.Err.Error()})
		}
	case WarnLevel:
		writeJson(w, e.StatusCode, dto.Error{Status: dto.StatusError, Error: e.Err.Error()})
		if e.Log {
			l.Warn(logger.LogData{Section: e.Section, Message: e.Err.Error()})
		}
	case ErrorLevel:
		writeJson(w, e.StatusCode, dto.Error{Status: dto.StatusError, Error: e.Err.Error()})
		if e.Log {
			l.Error(logger.LogData{Section: e.Section, Message: e.Err.Error()})
		}
	default:
		writeJson(w, e.StatusCode, dto.Error{Status: dto.StatusError, Error: e.Err.Error()})
		if e.Log {
			l.Info(logger.LogData{Section: e.Section, Message: e.Err.Error()})
		}
	}
}
