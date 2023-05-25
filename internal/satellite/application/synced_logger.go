package application

import (
	"fmt"
	"github.com/google/uuid"
	"kroseida.org/slixx/pkg/model"
	"kroseida.org/slixx/pkg/utils"
	"time"
)

type SyncedLogger struct {
	Logger      utils.Logger
	CachedLines []*model.SatelliteLogEntry
}

func (l *SyncedLogger) Info(args ...interface{}) {
	l.Logger.Info(args)
	l.appendLine("info", args...)
}

func (l *SyncedLogger) Debug(args ...interface{}) {
	l.Logger.Debug(args)
	// Debug logs are not synced
}

func (l *SyncedLogger) Error(args ...interface{}) {
	l.Logger.Error(args)
	l.appendLine("error", args...)
}

func (l *SyncedLogger) Warn(args ...interface{}) {
	l.Logger.Warn(args)
	l.appendLine("warn", args...)
}

func (l *SyncedLogger) appendLine(level string, args ...interface{}) {
	l.CachedLines = append(l.CachedLines, &model.SatelliteLogEntry{
		Id:          uuid.New(),
		SatelliteId: uuid.UUID{},
		Level:       level,
		Message:     fmt.Sprint(args...),
		LoggedAt:    time.Now(),
	})
}

func (l *SyncedLogger) GetLinesAndClear() []*model.SatelliteLogEntry {
	var cachedLines []*model.SatelliteLogEntry
	cachedLines = append(cachedLines, l.CachedLines...)
	l.CachedLines = []*model.SatelliteLogEntry{}

	return cachedLines
}
