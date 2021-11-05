package handlers

import (
	"logstash-as-a-service-backend/core"
	"net/http"

	"github.com/hashicorp/go-hclog"
)

type PipelinesHandlers struct {
	PipelineService *core.PipelineService
	l               hclog.Logger
}

type GenericError struct {
	Message string `json:"message"`
}

// New crea un handler de usuario con el logger y servicio dado
func New(ps *core.PipelineService, l hclog.Logger) *PipelinesHandlers {
	return &PipelinesHandlers{ps, l}
}

func (h *PipelinesHandlers) GetConfiguredPipelines(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[GetConfiguredPipelines] Handling request to get all configured pipelines")
	pipelines, err := h.PipelineService.GetConfiguredPipelines()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ToJSON(&GenericError{Message: err.Error()}, w)
	}
	ToJSON(&pipelines, w)
}
