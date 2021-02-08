package header

import (
	"net/http"
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

func (mw loggingMiddleware) Headers(req *http.Request) (b []byte) {
	defer func(begin time.Time) {
		mw.logger.Info(
			"service",
			zap.String("method", "Headers"),
			zap.Duration("took", time.Since(begin)),
		)
	}(time.Now())
	b = mw.next.Headers(req)
	return
}

func (mw instrumentingMiddleware) Headers(req *http.Request) (b []byte) {
	defer func(begin time.Time) {
		lvs := []string{"method", "headers", "error", "false"}
		mw.requestCount.With(lvs...).Add(1)
		mw.requestLatency.With(lvs...).Observe(time.Since(begin).Seconds())
	}(time.Now())
	b = mw.next.Headers(req)
	return
}
