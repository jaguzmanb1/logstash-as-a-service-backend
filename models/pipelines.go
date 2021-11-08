package models

import "github.com/breml/logstash-config/ast"

type PipelineConf struct {
	ID      string
	Path    string
	Plugins []ast.Plugin
}

type PipelinesConf []PipelineConf

type Dependences struct {
	Key   string
	Value string
}

type PluginSettings struct {
	Settings        string
	Type            string
	Required        bool
	HasValueList    bool
	RequiredDepends bool
	Dependence      []Dependences
	ValueList       []string
	SelectedValue   []bool
}

type Plugin struct {
	Type             string
	Name             string
	Description      string
	AvailableConfigs []PluginSettings
	Version          string
	ReleasedDate     string
	Configs          []PluginSettings
}
