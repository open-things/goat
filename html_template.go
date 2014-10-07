package goat

import (
	html_template "html/template"
	"io"
	"path/filepath"
)

type HtmlTemplate struct {
	template *html_template.Template
}

func NewHtmlTemplate(name string) TemplateCore {
	return &HtmlTemplate{
		template: html_template.New(name),
	}
}

func (this *HtmlTemplate) Funcs(funcMap map[string]interface{}) {
	this.template.Funcs(funcMap)
}

func (this *HtmlTemplate) ParseGlob(pattern string) error {
	names, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	if len(names) > 0 {
		_, err := this.template.ParseGlob(pattern)
		if err != nil {
			return err
		}
	}
	return nil
}

func (this *HtmlTemplate) ExecuteTemplate(writer io.Writer, name string, data interface{}) error {
	return this.template.ExecuteTemplate(writer, name, data)
}
