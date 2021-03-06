package core

import (
	"logstash-as-a-service-backend/data"
	"logstash-as-a-service-backend/models"

	"github.com/breml/logstash-config/ast"
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

func (ps *PipelineService) GetConfiguredPipelinesDetailed() (models.PipelinesConf, error) {
	ps.l.Info("[GetConfiguredPipelines] Getting pipelines from pipeline service")
	pipelines, err := ps.fs.GetConfiguredPipelinesDetailed()

	if err != nil {
		return nil, err
	}

	return pipelines, nil

}

func (ps *PipelineService) CreatePipeline(conf ast.Config, user string) {
	ps.l.Info("[CreatePipeline]")
	err := ps.fs.CreatePipeline(conf, user)
	if err != nil {
		ps.l.Error("[CreatePipeline] Error creating pipeline", "error", err)
	}
}
