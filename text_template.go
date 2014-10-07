package goat

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	tt "text/template"
)

const TextTemplateType string = "text"
const TextTemplateExtension string = "tt"

type TextTemplate struct {
	template *tt.Template
	group    string
	name     string
}

func NewTextTemplate(group, name string) *TextTemplate {
	return &TextTemplate{
		template: nil,
		group:    group,
		name:     name,
	}
}

func (this *TextTemplate) Template() *tt.Template {
	return this.template
}

func (this *TextTemplate) Group() string {
	return this.group
}

func (this *TextTemplate) Name() string {
	return this.name
}

func (this *TextTemplate) Type() string {
	return TextTemplateType
}

func (this *TextTemplate) Extension() string {
	return TextTemplateExtension
}

func (this *TextTemplate) LoadFiles(root string) error {
	var err error
	var patterns []string = GetPatterns(this)

	var template *tt.Template = tt.New(this.name)
	for _, pattern := range patterns {
		_, err = template.ParseGlob(filepath.Join(root, pattern))
		if err != nil && !strings.HasPrefix(err.Error(), "template: pattern matches no files:") {
			return err
		}
	}

	if len(template.Templates()) == 0 {
		return fmt.Errorf(
			"None of the patterns (%#q) matched anything in %#q",
			strings.Join(patterns, ", "),
			root,
		)
	}

	this.template = template
	return nil
}

func (this *TextTemplate) Funcs(funcMap map[string]interface{}) {
	this.template.Funcs(funcMap)
}

func (this *TextTemplate) Execute(writer io.Writer, data interface{}) error {
	return this.template.ExecuteTemplate(writer, this.name, data)
}
