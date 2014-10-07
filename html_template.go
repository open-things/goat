package goat

import (
	"fmt"
	ht "html/template"
	"io"
	"path/filepath"
	"strings"
)

const HtmlTemplateType string = "html"
const HtmlTemplateExtension string = "ht"

type HtmlTemplate struct {
	template *ht.Template
	group    string
	name     string
}

func NewHtmlTemplate(group, name string) *HtmlTemplate {
	return &HtmlTemplate{
		template: nil,
		group:    group,
		name:     name,
	}
}

func (this *HtmlTemplate) Template() *ht.Template {
	return this.template
}

func (this *HtmlTemplate) Group() string {
	return this.group
}

func (this *HtmlTemplate) Name() string {
	return this.name
}

func (this *HtmlTemplate) Type() string {
	return HtmlTemplateType
}

func (this *HtmlTemplate) Extension() string {
	return HtmlTemplateExtension
}

func (this *HtmlTemplate) LoadFiles(root string) error {
	var err error
	var patterns []string = GetPatterns(this)

	var template *ht.Template = ht.New(this.name)
	for _, pattern := range patterns {
		_, err = template.ParseGlob(filepath.Join(root, pattern))
		if err != nil && !strings.HasPrefix(err.Error(), "html/template: pattern matches no files:") {
			return err
		}
	}

	if len(template.Templates()) == 1 && template.Templates()[0].Name() == template.Name() {
		return fmt.Errorf(
			"None of the patterns (%#q) matched anything in %#q",
			strings.Join(patterns, ", "),
			root,
		)
	}

	this.template = template
	return nil
}

func (this *HtmlTemplate) Funcs(funcMap map[string]interface{}) {
	this.template.Funcs(funcMap)
}

func (this *HtmlTemplate) Execute(writer io.Writer, data interface{}) error {
	return this.template.ExecuteTemplate(writer, this.name, data)
}
