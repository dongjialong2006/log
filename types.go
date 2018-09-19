package log

import (
	"time"
)

const (
	DEFAULT_ROTATION_TIME = 24 * time.Hour
)

const (
	DEFAULT_ROTATION_COUNT        = 6
	DEFAULT_WATCHER_FILES_BY_NUM  = 6
	DEFAULT_WATCHER_FILES_BY_SIZE = 5 * 1024 * 1024
)

const (
	DEFAULT_LOG_NAME  = "./log/default.log"
	DEFAULT_LOG_LEVEL = "debug"
)

const (
	LOG_DEBUG_LEVEL = "debug"
	LOG_INFO_LEVEL  = "info"
	LOG_ERROR_LEVEL = "error"
	LOG_FATAL_LEVEL = "fatal"
	LOG_WARN_LEVEL  = "warn"
	LOG_PANIC_LEVEL = "panic"
)

const (
	TCP   = "tcp"
	UDP   = "udp"
	HTTP  = "http"
	HTTPS = "https"
)
