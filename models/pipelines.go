package models

import (
	"encoding/json"

	"github.com/breml/logstash-config/ast"
)

type PipelineConf struct {
	ID     string
	Path   string
	Config ast.Config
}

type PipelinesConf []PipelineConf

type PluginType int

type Pos struct {
	Line   int
	Column int
	Offset int
}

type ElseBlock struct {
	Start         Pos
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

type BooleanOperator struct {
	Op    int
	Start Pos
}

type ElseIfBlock struct {
	Start         Pos
	Condition     Condition
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

type Condition struct {
	Expression []Expression
}

type Expression interface {
	Pos() Pos
	BoolOperator() BooleanOperator
	SetBoolOperator(BooleanOperator)
	expressionNode()
}

type CommentBlock []Comment

type Config struct {
	Input         []PluginSection
	Filter        []PluginSection
	Output        []PluginSection
	FooterComment CommentBlock
	Warnings      []string
}

type Comment struct {
	comment     string
	SpaceBefore bool
	SpaceAfter  bool
}

type PluginSection struct {
	Start           Pos
	PluginType      PluginType
	BranchOrPlugins []json.RawMessage
	CommentBlock    CommentBlock
	FooterComment   CommentBlock
}

type BranchOrPlugin interface {
	Pos() Pos
	branchOrPlugin()
}

type Plugin struct {
	Start         Pos
	name          string
	Attributes    []Attribute
	Comment       CommentBlock
	FooterComment CommentBlock
}

type Branch struct {
	IfBlock     IfBlock
	ElseIfBlock []ElseIfBlock
	ElseBlock   ElseBlock
}

func (p Plugin) Pos() Pos      { return p.Start }
func (Plugin) branchOrPlugin() {}

type IfBlock struct {
	Start         Pos
	Condition     Condition
	Block         []BranchOrPlugin
	Comment       CommentBlock
	FooterComment CommentBlock
}

func (Branch) branchOrPlugin() {}
func (b Branch) Pos() Pos      { return b.IfBlock.Start }

type Attribute interface {
	Name() string
	String() string
	ValueString() string
	CommentBlock() string
	Pos() Pos
	attributeNode()
}
