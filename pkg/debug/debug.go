package debug

import (
	"net/http"
	"net/http/pprof"

	"github.com/ether-echo/user-service/pkg/logger"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	log = logger.Logger().Named("debug_server").Sugar()
)

type ProfilingServer struct {
	*http.ServeMux
}

// Run should be always run in goroutine.
// It runs pprof server with all available pprof handlers
// with it is default paths.
func Run(addr string) {
	server := ProfilingServer{
		http.NewServeMux(),
	}
	server.HandleFunc("/debug/pprof", pprof.Index)
	server.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	server.HandleFunc("/debug/pprof/profile", pprof.Profile)
	server.HandleFunc("/debug/pprof/trace", pprof.Trace)
	server.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	server.Handle("/metrics", promhttp.Handler())

	server.HandleFunc("/healthz", healthCheck)
	server.HandleFunc("/ready", readyCheck)

	defer func() {
		if r := recover(); r != nil {
			log.Errorf("failed to listen and serve: %s", r)
		}
	}()
	//nolint:errcheck
	err := http.ListenAndServe(addr, server)
	if err != nil {
		panic(err)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
