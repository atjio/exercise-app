package global

import (
	"time"
)

type AppState struct {
	DelayInMS time.Duration
}

var State AppState
