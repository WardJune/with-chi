package limiter

import (
	"log/slog"
	"sync/atomic"
	"time"
)

type AdaptiveLimiter struct {
	maxInFlight atomic.Int64
	inFlight    atomic.Int64

	targetLatency time.Duration

	emaLatency atomic.Int64
	alpha      float64
}

func NewAdaptiveLimiter(
	initialMax int64,
	target time.Duration,
	alpha float64,
) *AdaptiveLimiter {
	l := &AdaptiveLimiter{
		targetLatency: target,
		alpha:         alpha,
	}

	l.maxInFlight.Store(initialMax)
	l.emaLatency.Store(target.Nanoseconds())
	return l
}

func (l *AdaptiveLimiter) TryAcquire() bool {
	for {
		cur := l.inFlight.Load()
		max := l.maxInFlight.Load()

		if cur >= max {
			return false
		}

		if l.inFlight.CompareAndSwap(cur, cur+1) {
			return true
		}
	}
}

func (l *AdaptiveLimiter) Release() {
	l.inFlight.Add(-1)
}

func (l *AdaptiveLimiter) Obeserve(d time.Duration) {
	old := time.Duration(l.emaLatency.Load())
	newEma := time.Duration(
		l.alpha*float64(d) + (1-l.alpha)*float64(old),
	)

	// slog.Info("ema update",
	// 	"old_ms", old.Milliseconds(),
	// 	"new_ms", newEma.Milliseconds(),
	// 	"raw_ms", d.Milliseconds(),
	// )

	l.emaLatency.Store(newEma.Nanoseconds())
}

func (l *AdaptiveLimiter) Adjust() {
	ema := time.Duration(l.emaLatency.Load())
	current := l.maxInFlight.Load()

	// slog.Info("adjust check",
	// 	"ema_ms", ema.Milliseconds(),
	// 	"target_ms", l.targetLatency.Milliseconds(),
	// 	"current_max", current,
	// )

	switch {
	case ema > l.targetLatency && current > 1:
		l.maxInFlight.Store(current - 1)
		slog.Info("decrease max", "new", current-1)
	case ema < l.targetLatency/2:
		l.maxInFlight.Store(current + 1)
		slog.Info("increase max", "new", current+1)
	}

	l.ExportMetrics()
}
