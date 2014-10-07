package goat

import (
	"fmt"
	"io"
)

type Template struct {
	kind         string
	group        string
	name         string
	patterns     []string
	execName     string
	globalFuncs  map[string]interface{}
	globalValues map[string]interface{}
	localFuncs   map[string]interface{}
	localValues  map[string]interface{}
	userFuncs    map[string]interface{}
	userValues   map[string]interface{}
}

func NewTemplate(kind, group, name string, patterns []string, execName string, globalFuncs map[string]interface{}, globalValues map[string]interface{}) *Template {
	template := &Template{
		kind:         kind,
		group:        group,
		name:         name,
		patterns:     patterns,
		execName:     execName,
		globalFuncs:  CopyMap(globalFuncs),
		globalValues: CopyMap(globalValues),
		localFuncs:   make(map[string]interface{}),
		localValues:  make(map[string]interface{}),
		userFuncs:    make(map[string]interface{}),
		userValues:   make(map[string]interface{}),
	}
	template.init()
	return template
}

func (this *Template) Kind() string       { return this.kind }
func (this *Template) Group() string      { return this.group }
func (this *Template) Name() string       { return this.name }
func (this *Template) Patterns() []string { return this.patterns }

func (this *Template) GetValue(name string) interface{} {
	value, ok := this.userValues[name]
	if ok {
		return value
	}
	value, ok = this.localValues[name]
	if ok {
		return value
	}
	value, ok = this.globalValues[name]
	if ok {
		return value
	}
	return nil
}

func (this *Template) SetValue(name string, value interface{}) {
	this.userValues[name] = value
}

func (this *Template) Funcs(funcMap map[string]interface{}) {
	if funcMap == nil {
		return
	}
	for key, val := range funcMap {
		this.userFuncs[key] = val
	}
}

func (this *Template) ParseGlob(pattern string) error {
	this.patterns = append(this.patterns, pattern)
	return nil
}

func (this *Template) ExecuteTemplate(writer io.Writer, name string, data interface{}) error {
	template, err := this.loadTemplate()
	if err != nil {
		return err
	}
	return template.ExecuteTemplate(writer, name, data)
}

func (this *Template) Execute(writer io.Writer, data interface{}) error {
	return this.ExecuteTemplate(writer, this.execName, data)
}

func (this *Template) init() {
	this.localFuncs["kind"] = this.Kind
	this.localFuncs["group"] = this.Group
	this.localFuncs["name"] = this.Name
	this.localFuncs["value"] = this.GetValue
}

func (this *Template) loadTemplate() (TemplateCore, error) {
	// construct TemplateCore object
	templateName := fmt.Sprintf(":%s:%s:%s:root", this.kind, this.group, this.name)
	constructor, err := TemplateFactory.GetConstructor(this.kind)
	if err != nil {
		return nil, err
	}
	template := constructor(templateName)

	// setup funcs
	if len(this.globalFuncs) > 0 {
		template.Funcs(this.globalFuncs)
	}
	if len(this.localFuncs) > 0 {
		template.Funcs(this.localFuncs)
	}
	if len(this.userFuncs) > 0 {
		template.Funcs(this.userFuncs)
	}

	// parse template files
	for _, pattern := range this.patterns {
		err := template.ParseGlob(pattern)
		if err != nil {
			return nil, err
		}
	}

	// all good!
	return template, nil
}
