package logger

import (
	log "github.com/inconshreveable/log15"
)

type (
	// Logger log15
	Logger interface {
		log.Logger
	}

	// Options - logger options
	Options struct {
		Verbose bool
	}
)

// New creates new logger based on logger.Options
func New(opt Options) Logger {
	l := log.New()
	handlers := []log.Handler{}

	lvl := log.LvlInfo
	if opt.Verbose {
		lvl = log.LvlDebug
	}
	handlers = append(handlers, log.LvlFilterHandler(lvl, log.StdoutHandler))

	l.SetHandler(log.MultiHandler(handlers...))
	return l
}
