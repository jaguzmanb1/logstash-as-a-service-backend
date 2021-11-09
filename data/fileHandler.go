package data

import (
	"bufio"
	"fmt"
	"log"
	"logstash-as-a-service-backend/models"
	"os"
	"regexp"
	"strings"
	"unicode/utf8"

	config "github.com/breml/logstash-config"
	"github.com/breml/logstash-config/ast"
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
				scanner.Scan()
				txt = scanner.Text()
				pl.Path = strings.Replace(strings.Replace(strings.Split(txt, ":")[1], " ", "", -1), "\"", "", -1)
				pipelines = append(pipelines, pl)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pipelines, nil
}

func (f *FileService) GetConfiguredPipelinesDetailed() (models.PipelinesConf, error) {
	f.l.Info("[GetConfiguredPipelinesDetailed] Start to configuration file ", "file path", f.ConfigPath)

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
				scanner.Scan()
				txt = scanner.Text()
				pl.Path = strings.Replace(strings.Replace(strings.Split(txt, ":")[1], " ", "", -1), "\"", "", -1)
				pl.Config = f.GetConfiguredPlugins(pl.Path)
				pipelines = append(pipelines, pl)
			}
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return pipelines, nil
}

func (f *FileService) GetPlugin(path string, ty string) ([]ast.PluginSection, error) {
	f.l.Info("[GetPlugin]")
	file, err := os.Open(path)
	if err != nil {
		return []ast.PluginSection{}, err
	}

	defer file.Close()
	//b, err := ioutil.ReadAll(file)

	c, err := config.ParseFile(path)
	if err != nil {
		return []ast.PluginSection{}, err
	}
	conf := c.(ast.Config)
	pls := []ast.PluginSection{}
	if len(conf.Input) > 0 {
		pls = conf.Input
	} else if len(conf.Filter) > 0 {
		pls = conf.Filter
	} else if len(conf.Output) > 0 {
		pls = conf.Output
	}

	return pls, nil
}

func (f *FileService) GetConfiguredPlugins(str string) (result ast.Config) {
	f.l.Info("[GetConfiguredPlugins]")
	replacer := strings.NewReplacer("{", "", "}", "")
	rex := regexp.MustCompile(`{.*}`)
	out := replacer.Replace(rex.FindAllStringSubmatch(str, -1)[0][0])

	plgsText := strings.Split(out, ",")
	conf := ast.Config{}
	for _, plg := range plgsText {
		ty := strings.Split(plg, "/")[0]
		sec := trimLastChar(ty)
		pl, err := f.GetPlugin(f.PipelinesPath+plg+".conf", sec)
		if sec == "input" {
			conf.Input = append(conf.Input, pl...)
		}
		if sec == "filter" {
			conf.Filter = append(conf.Filter, pl...)
		}
		if sec == "output" {
			conf.Output = append(conf.Filter, pl...)
		}

		if err != nil {
			f.l.Error("[GetConfiguredPlugins]", "error", err)
		}

	}

	return conf
}

func (f *FileService) CreatePipeline(config ast.Config, id string) error {
	f.l.Info("[CreatePipeline]")
	fmt.Printf(config.String())
	return nil
}

func trimLastChar(s string) string {
	r, size := utf8.DecodeLastRuneInString(s)
	if r == utf8.RuneError && (size == 0 || size == 1) {
		size = 0
	}
	return s[:len(s)-size]
}
