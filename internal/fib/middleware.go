package fib

import (
	"math/big"
	"time"

	"github.com/go-kit/kit/metrics"
	"go.uber.org/zap"
)

// Middleware describes a service (as opposed to endpoint) middleware.
type Middleware func(Service) Service

func LoggingMiddleware(logger zap.Logger) Middleware {
	return func(next Service) Service {
		return &loggingMiddleware{
			next:   next,
			logger: logger,
		}
	}
}

func NewInstrumentingMiddleware(
	counter metrics.Counter,
	latency metrics.Histogram,
	s Service,
) Service {
	return &instrumentingMiddleware{
		requestCount:   counter,
		requestLatency: latency,
		next:           s,
	}
}

type loggingMiddleware struct {
	next   Service
	logger zap.Logger
}

type instrumentingMiddleware struct {
	requestCount   metrics.Counter
	requestLatency metrics.Histogram
	next           Service
}

func (mw loggingMiddleware) Fib(n uint64) (b *big.Int) {
	defer func(begin time.Time) {
		mw.logger.Debug(
			"service",
			zap.String("method", "Fib"),
			zap.Uint64("input", n),
			zap.Duration("took", time.Since(begin)),
		)
	}(time.Now())
	b = mw.next.Fib(n)
	return
}

func (mw instrumentingMiddleware) Fib(n uint64) (b *big.Int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "fib", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	b = mw.next.Fib(n)
	return
}
