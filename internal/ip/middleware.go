package ip

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

func (mw loggingMiddleware) GetIp() (output []byte, err error) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service",
			zap.String("method", "GetIp"),
			zap.Duration("took", time.Since(begin)),
		)
	}(time.Now())
	output, err = mw.next.GetIp()
	return
}
func (mw instrumentingMiddleware) GetIp() (output []byte, err error) {
	defer func(begin time.Time) {
		lvs := []string{"method", "getip", "error", fmt.Sprint(err != nil)}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	output, err = mw.next.GetIp()
	return
}
