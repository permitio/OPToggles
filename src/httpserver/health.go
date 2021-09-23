package httpserver

import (
	"net/http"
)

const HealthEndpoint string = "/health"
const LiveEndpoint string = "/health/live"
const StartedEndpoint string = "/health/started"

type HealthServer struct {
	policyTrackHealthy bool
	targetSyncHealthy  bool
	started            bool
}

func (hs *HealthServer) SetStarted() {
	hs.started = true
}

func (hs *HealthServer) SetTargetSyncHealth(healthy bool) {
	hs.targetSyncHealthy = healthy
}

func (hs *HealthServer) SetPolicyTrackHealth(healthy bool) {
	hs.policyTrackHealthy = healthy
}

func (hs *HealthServer) RegisterHandlers() {
	http.HandleFunc(LiveEndpoint, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	http.HandleFunc(StartedEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if hs.started {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	})
	http.HandleFunc(HealthEndpoint, func(w http.ResponseWriter, r *http.Request) {
		if !hs.started {
			w.WriteHeader(http.StatusServiceUnavailable)
		} else if hs.targetSyncHealthy && hs.policyTrackHealthy {
			w.WriteHeader(http.StatusOK)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
	})
}

func NewHealthServer() *HealthServer {
	return &HealthServer{
		policyTrackHealthy: true,
		targetSyncHealthy:  true,
		started:            false,
	}
}

var HealthServ *HealthServer

func init() {
	HealthServ = NewHealthServer()
	HealthServ.RegisterHandlers()
}
