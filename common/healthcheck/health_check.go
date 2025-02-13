package healthcheck

import (
	"encoding/json"
	"net/http"

	"github.com/cloudwego/kitex/pkg/klog"
)

type HealthResponse struct {
	Status      string `json:"status"`
	ServiceName string `json:"service"`
}

func StartHealthCheck(addr string, serviceName string) {
	handler := newHealthCheckHandler(serviceName)
	http.HandleFunc("/health", handler.healthCheckHandler)

	// 在后台启动健康检查服务
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			klog.Errorf("Health check server failed to start: %v", err)
		}
	}()
}

type healthCheckHandler struct {
	serviceName string
}

func newHealthCheckHandler(serviceName string) *healthCheckHandler {
	return &healthCheckHandler{
		serviceName: serviceName,
	}
}

func (h *healthCheckHandler) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:      "ok",
		ServiceName: h.serviceName,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(response)
}
