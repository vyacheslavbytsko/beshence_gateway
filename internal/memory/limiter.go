package memory

import (
	"time"

	"golang.org/x/time/rate"
)

func GetLimiter(bankID string) *rate.Limiter {
	LimitersMutex.Lock()
	defer LimitersMutex.Unlock()

	if l, ok := Limiters[bankID]; ok {
		return l
	}

	l := rate.NewLimiter(rate.Every(6*time.Second), 5)

	Limiters[bankID] = l

	return l
}
