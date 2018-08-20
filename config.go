package log

const (
	logAge           = "age"
	logName          = "name"
	logTime          = "time"
	logLevel         = "level"
	logCount         = "count"
	logFields        = "fields"
	logContext       = "context"
	logCaller        = "caller"
	logFluent        = "fluent"
	logTerminal      = "terminal"
	logRemoteAddr    = "addr"
	logRemoteProType = "protocol"
	logWatcherEnable = "enable"
	logWatcherByNum  = "watcherByNum"
	logWatcherBySize = "watcherBySize"
)

type option interface {
	Get(key string) interface{}
}

type config struct {
	key   string
	value interface{}
}

func new(key string, value interface{}) *config {
	return &config{
		key:   key,
		value: value,
	}
}

func (c config) Get(key string) interface{} {
	if key == c.key {
		return c.value
	}

	return nil
}
