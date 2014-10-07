package goat

import (
	"io"
	"path/filepath"
	text_template "text/template"
)

type TextTemplate struct {
	template *text_template.Template
}

func NewTextTemplate(name string) TemplateCore {
	return &TextTemplate{
		template: text_template.New(name),
	}
}

func (this *TextTemplate) Funcs(funcMap map[string]interface{}) {
	this.template.Funcs(funcMap)
}

func (this *TextTemplate) ParseGlob(pattern string) error {
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

func (this *TextTemplate) ExecuteTemplate(writer io.Writer, name string, data interface{}) error {
	return this.template.ExecuteTemplate(writer, name, data)
}
