package limiter

import "time"

type Limiter struct {
	sem chan struct{}
}

func NewLimiter(n int) *Limiter {
	return &Limiter{
		sem: make(chan struct{}, n),
	}
}

func (l *Limiter) TryAcquire(timeout time.Duration) bool {
	select {
	case l.sem <- struct{}{}:
		return true
	case <-time.After(timeout):
		return false
	}
}

func (l *Limiter) Release() {
	<-l.sem
}
