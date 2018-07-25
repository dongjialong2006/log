package types

import (
	"time"
)

const (
	DEFAULT_MAX_AGE       = 7 * 24 * time.Hour
	DEFAULT_ROTATION_TIME = 24 * time.Hour
)

const (
	DEFAULT_ROTATION_COUNT        = 6
	DEFAULT_WATCHER_FILES_BY_NUM  = 6
	DEFAULT_WATCHER_FILES_BY_SIZE = 10 * 1024 * 1024
)

const (
	DEFAULT_LOG_NAME = "./log/default.log"
)

const (
	TCP   = "tcp"
	UDP   = "udp"
	HTTP  = "http"
	HTTPS = "https"
)
