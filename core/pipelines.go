package core

import (
	"logstash-as-a-service-backend/data"
	"logstash-as-a-service-backend/models"

	"github.com/hashicorp/go-hclog"
)

type PipelineService struct {
	l  hclog.Logger
	fs *data.FileService
}

func NewPipelineService(l hclog.Logger, fs *data.FileService) *PipelineService {
	return &PipelineService{l, fs}
}

func (ps *PipelineService) GetConfiguredPipelines() (models.PipelinesConf, error) {
	ps.l.Info("[GetConfiguredPipelines] Getting pipelines from pipeline service")
	pipelines, err := ps.fs.GetConfiguredPipelines()
	if err != nil {
		return nil, err
	}

	return pipelines, nil

}
