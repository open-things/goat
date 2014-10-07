package goat

import (
	"path/filepath"
)

type Loader struct {
	templateRoot    string
	defaultKind     string
	patterns        []string
	execNamePattern string
	globalFuncs     map[string]interface{}
	globalValues    map[string]interface{}
}

func NewLoader(templateRoot string) *Loader {
	return &Loader{
		templateRoot: templateRoot,
		defaultKind:  "html",
		patterns: []string{
			"base.:extension:",
			":group:_base.:extension:",
			":group:_:name:.:extension:",
			":group:_:name:_*.:extension:",
		},
		execNamePattern: ":name:_main",
		globalFuncs:     make(map[string]interface{}),
		globalValues:    make(map[string]interface{}),
	}
}

func (this *Loader) TemplateRoot() string                      { return this.templateRoot }
func (this *Loader) SetTemplateRoot(root string)               { this.templateRoot = root }
func (this *Loader) DefaultKind() string                       { return this.defaultKind }
func (this *Loader) SetDefaultKind(kind string)                { this.defaultKind = kind }
func (this *Loader) Patterns() []string                        { return this.patterns }
func (this *Loader) SetPatterns(patterns []string)             { this.patterns = patterns }
func (this *Loader) ExecNamePattern() string                   { return this.execNamePattern }
func (this *Loader) SetExecNamePattern(execNamePattern string) { this.execNamePattern = execNamePattern }

func (this *Loader) AddPatterns(patterns ...string) {
	this.patterns = append(this.patterns, patterns...)
}

func (this *Loader) GetValue(name string) interface{} {
	value, ok := this.globalValues[name]
	if ok {
		return value
	}
	return nil
}

func (this *Loader) SetValue(name string, value interface{}) {
	this.globalValues[name] = value
}

func (this *Loader) Funcs(funcMap map[string]interface{}) {
	if funcMap == nil {
		return
	}
	for key, val := range funcMap {
		this.globalFuncs[key] = val
	}
}

func (this *Loader) Get(group, name string) (*Template, error) {
	return this.GetFull(this.defaultKind, group, name)
}

func (this *Loader) GetFull(kind, group, name string) (*Template, error) {
	extension, err := TemplateFactory.GetExtension(kind)
	if err != nil {
		return nil, err
	}
	patterns := this.getProcessedPatterns(kind, group, name, extension)
	execName := this.getProcessedExecName(kind, group, name, extension)
	template := NewTemplate(kind, group, name, patterns, execName, this.globalFuncs, this.globalValues)
	return template, nil
}

func processPattern(pattern, kind, group, name, extension string) string {
	match := []string{":kind:", ":group:", ":name:", ":extension:"}
	replace := []string{kind, group, name, extension}
	return ReplaceAll(pattern, match, replace)
}

func (this *Loader) getProcessedPatterns(kind, group, name, extension string) []string {
	result := []string{}
	for _, pattern := range this.patterns {
		processed := filepath.Join(this.templateRoot, processPattern(pattern, kind, group, name, extension))
		result = append(result, processed)
	}
	return result
}

func (this *Loader) getProcessedExecName(kind, group, name, extension string) string {
	return processPattern(this.execNamePattern, kind, group, name, extension)
}
