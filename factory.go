package goat

import (
	"fmt"
	"io"
)

type TemplateCore interface {
	Funcs(funcMap map[string]interface{})
	ParseGlob(pattern string) error
	ExecuteTemplate(writer io.Writer, name string, data interface{}) error
}

type TemplateCoreConstructor func(name string) TemplateCore

type templateInfo struct {
	Constructor TemplateCoreConstructor
	Extension   string
}

type templateFactory struct {
	mapping map[string]*templateInfo
}

var TemplateFactory *templateFactory = newTemplateFactory()

func (this *templateFactory) Register(kind string, constructor TemplateCoreConstructor, extension string) error {
	ti := templateInfo{
		Constructor: constructor,
		Extension:   extension,
	}
	this.mapping[kind] = &ti
	return nil
}

func (this *templateFactory) Unregister(kind string) error {
	_, ok := this.mapping[kind]
	if ok {
		delete(this.mapping, kind)
	}
	return nil
}

func (this *templateFactory) GetConstructor(kind string) (TemplateCoreConstructor, error) {
	ti, err := this.getInfo(kind)
	if err != nil {
		return nil, err
	}
	return ti.Constructor, nil
}

func (this *templateFactory) GetExtension(kind string) (string, error) {
	ti, err := this.getInfo(kind)
	if err != nil {
		return "", err
	}
	return ti.Extension, nil
}

func (this *templateFactory) getInfo(kind string) (*templateInfo, error) {
	ti, ok := this.mapping[kind]
	if ok {
		return ti, nil
	}
	return nil, fmt.Errorf("goat: no match found")
}

func newTemplateFactory() *templateFactory {
	factory := &templateFactory{
		mapping: make(map[string]*templateInfo),
	}
	// register defaults
	factory.Register("text", NewTextTemplate, "txt")
	factory.Register("html", NewHtmlTemplate, "html")
	return factory
}
