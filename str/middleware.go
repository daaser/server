package str

import (
	"fmt"
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

func (mw loggingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service",
			zap.String("method", "Uppercase"),
			zap.String("input", s),
			zap.String("output", output),
			zap.Duration("took", time.Since(begin)),
			zap.Error(err),
		)
	}(time.Now())
	output, err = mw.next.Uppercase(s)
	return
}

func (mw loggingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service",
			zap.String("method", "Count"),
			zap.String("input", s),
			zap.Int("output", n),
			zap.Duration("took", time.Since(begin)),
		)
	}(time.Now())
	n = mw.next.Count(s)
	return
}

func (mw instrumentingMiddleware) Uppercase(s string) (output string, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "uppercase", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.next.Uppercase(s)
	return
}

func (mw instrumentingMiddleware) Count(s string) (n int) {
	defer func(begin time.Time) {
		lvs := []string{"method", "count", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	n = mw.next.Count(s)
	return
}
