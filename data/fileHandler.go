package data

import (
	"bufio"
	"fmt"
	"log"
	"logstash-as-a-service-backend/models"
	"os"
	"strings"

	"github.com/hashicorp/go-hclog"
)

type FileService struct {
	l             hclog.Logger
	PipelinesPath string
	ConfigPath    string
}

func NewFileService(l hclog.Logger, pipelinesPath, configPath string) *FileService {
	return &FileService{
		l:             l,
		PipelinesPath: pipelinesPath,
		ConfigPath:    configPath,
	}
}

func (f *FileService) GetConfiguredPipelines() (models.PipelinesConf, error) {
	f.l.Info("[GetConfiguredPipelines] Start to configuration file ", "file path", f.ConfigPath)

	pipelines := models.PipelinesConf{}
	file, err := os.Open(f.ConfigPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		txt := scanner.Text()
		if !strings.HasPrefix(txt, "#") {
			if strings.HasPrefix(txt, "-") {
				pl := models.PipelineConf{}
				pl.ID = strings.Replace(strings.Split(txt, ":")[1], " ", "", -1)
				pipelines = append(pipelines, &pl)
			}
			fmt.Println(scanner.Text())
		}

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pipelines, nil
}
