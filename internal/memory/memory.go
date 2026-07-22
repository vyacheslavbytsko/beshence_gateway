package memory

import "sync"

var (
	Banks      = map[string]Bank{}
	Challenges = map[string]Challenge{}
	APIURLss   = map[string]APIURLs{}
	Mutex      sync.Mutex
)
