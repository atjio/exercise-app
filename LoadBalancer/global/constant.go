package global

import (
	"time"
)

const (
	HEALTHCHECK_URL string = "/healthcheck"
	HEALTHCHECK_DELAY_IN_MS time.Duration = 1000 * time.Millisecond
	HTTP_PROTOCOL string = "http://"
	NODE_MAX_TIMEOUT time.Duration = 2000 * time.Millisecond
)