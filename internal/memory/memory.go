package memory

import (
	"sync"

	"golang.org/x/time/rate"
)

var (
	Banks         = map[string]Bank{}
	Challenges    = map[string]Challenge{}
	APIURLss      = map[string]APIURLs{}
	Mutex         sync.Mutex
	Limiters      = map[string]*rate.Limiter{}
	LimitersMutex sync.Mutex
)
