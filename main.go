package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/daaser/server/fib"
	"github.com/daaser/server/header"
	"github.com/daaser/server/str"
	kitprometheus "github.com/go-kit/kit/metrics/prometheus"
	"github.com/gorilla/mux"
	"github.com/heptiolabs/healthcheck"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	defaultPort = "8080"
)

// seq 1 20 | xargs -n1 -P8 bash -c 'i=$0; url="http://localhost:8080/fib/$i"; curl --silent $url'
func main() {
	var (
		addr = envString("PORT", defaultPort)

		httpAddr = flag.String("http.addr", ":"+addr, "HTTP listen address")
		timeout  = flag.Duration(
			"timeout",
			5*time.Second,
			"Time in seconds to wait before forcefully terminating the server.",
		)
	)

	flag.Parse()

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	fieldKeys := []string{"method", "error"}

	var ss str.Service
	{
		ss = str.NewService()
		ss = str.LoggingMiddleware(*logger)(ss)
		ss = str.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "string",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "string",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			ss,
		)
	}

	var fs fib.Service
	{
		fs = fib.NewService()
		fs = fib.LoggingMiddleware(*logger)(fs)
		fs = fib.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "fib",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "fib",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			fs,
		)
	}

	var hs header.Service
	{
		hs = header.NewService()
		hs = header.LoggingMiddleware(*logger)(hs)
		hs = header.NewInstrumentingMiddleware(
			kitprometheus.NewCounterFrom(stdprometheus.CounterOpts{
				Namespace: "api",
				Subsystem: "header",
				Name:      "request_count",
				Help:      "Number of requests received.",
			}, fieldKeys),
			kitprometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
				Namespace: "api",
				Subsystem: "header",
				Name:      "request_latency_microseconds",
				Help:      "Total duration of requests in microseconds.",
			}, fieldKeys),
			hs,
		)
	}

	r := mux.NewRouter()

	// our main API routes
	r.PathPrefix("/string").Handler(str.MakeHandler(ss))
	r.PathPrefix("/fib").Handler(fib.MakeHandler(fs))
	r.Handle("/headers", header.MakeHandler(hs))

	// expose the Promethus metrics we registered above
	r.Handle("/metrics", promhttp.Handler())

	// register readiness and liveness probes
	health := healthcheck.NewHandler()
	health.AddLivenessCheck("goroutine-threshold", healthcheck.GoroutineCountCheck(100))
	health.AddReadinessCheck(
		"check-tcp",
		healthcheck.TCPDialCheck(*httpAddr, 50*time.Millisecond),
	)
	r.HandleFunc("/live", health.LiveEndpoint)
	r.HandleFunc("/ready", health.ReadyEndpoint)

	// register some access control middleware
	r.Use(accessControl)

	err := walkRoute(r)
	if err != nil {
		logger.Error("walkRoute", zap.Error(err))
		os.Exit(1)
	}

	stdOutLogger, _ := zap.NewStdLogAt(logger, zap.ErrorLevel)
	srv := &http.Server{
		Addr:           *httpAddr,
		Handler:        r,
		ErrorLog:       stdOutLogger,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	errs := make(chan error)
	go func() {
		logger.Info(
			"server",
			zap.String("transport", "HTTP"),
			zap.String("addr", *httpAddr),
		)
		errs <- srv.ListenAndServe()
		// errs <- srv.ListenAndServeTLS(
		// 	"./config/cert.pem",
		// 	"./config/key.pem",
		// )
	}()

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errs <- fmt.Errorf("%s", <-c)
	}()

	logger.Warn("terminated", zap.Error(<-errs))
	ctx, cancel := context.WithTimeout(context.Background(), *timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("could not gracefully shutdown https server", zap.Error(err))
	}
}

func accessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

func envString(env, fallback string) string {
	e := os.Getenv(env)
	if e == "" {
		return fallback
	}
	return e
}

func walkRoute(r *mux.Router) error {
	return r.Walk(func(
		route *mux.Route,
		router *mux.Router,
		ancestors []*mux.Route,
	) error {
		pathTemplate, err := route.GetPathTemplate()
		if err == nil {
			fmt.Println("ROUTE:", pathTemplate)
		}
		pathRegexp, err := route.GetPathRegexp()
		if err == nil {
			fmt.Println("Path regexp:", pathRegexp)
		}
		// queriesTemplates, err := route.GetQueriesTemplates()
		// if err == nil {
		// 	fmt.Println("Queries templates:", strings.Join(queriesTemplates, ","))
		// }
		// queriesRegexps, err := route.GetQueriesRegexp()
		// if err == nil {
		// 	fmt.Println("Queries regexps:", strings.Join(queriesRegexps, ","))
		// }
		methods, err := route.GetMethods()
		if err == nil {
			fmt.Println("Methods:", strings.Join(methods, ","))
		}
		fmt.Println()
		return nil
	})
}
