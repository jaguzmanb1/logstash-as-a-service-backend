package handlers

import (
	"fmt"
	"logstash-as-a-service-backend/core"
	"logstash-as-a-service-backend/models"
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

func (h *PipelinesHandlers) GetConfiguredPipelinesDetailed(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[GetConfiguredPipelinesDetailed] Handling request to get all configured pipelines")
	pipelines, err := h.PipelineService.GetConfiguredPipelinesDetailed()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ToJSON(&GenericError{Message: err.Error()}, w)
	}
	ToJSON(pipelines[0], w)
}

func (h *PipelinesHandlers) CreatePipeline(w http.ResponseWriter, r *http.Request) {
	h.l.Info("[CreatePipeline] Handling request to create pipeline")
	conf := models.Config{}
	err := FromJSON(&conf, r.Body)
	writeErr(w, err)
	fmt.Printf("%+v\n", conf)
	//h.PipelineService.CreatePipeline(conf)
}

func writeErr(w http.ResponseWriter, err error) {
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		ToJSON(&GenericError{Message: err.Error()}, w)
	}
}
