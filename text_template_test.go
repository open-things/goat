package goat_test

import (
	check "gopkg.in/check.v1"

	"github.com/Paulius-Maruska/goat"
)

type TextTemplateSuite struct{}

var _ = check.Suite(&TextTemplateSuite{})

func (this *TextTemplateSuite) TestNewAndGetters(c *check.C) {
	var textTemplate *goat.TextTemplate = goat.NewTextTemplate("g", "n")
	c.Assert(textTemplate.Template(), check.IsNil)
	c.Assert(textTemplate.Group(), check.Equals, "g")
	c.Assert(textTemplate.Name(), check.Equals, "n")
	c.Assert(textTemplate.Type(), check.Equals, goat.TextTemplateType)
	c.Assert(textTemplate.Extension(), check.Equals, goat.TextTemplateExtension)
}

func (this *TextTemplateSuite) TestLoadFilesParsesEverything(c *check.C) {
	var structure map[string]map[string]string = make(map[string]map[string]string)
	structure["text/gr"] = make(map[string]string)
	structure["text/gr"]["templ.tt"] = "{define \"main\"}main{template \"sub1\"}{template \"sub2\"}{end}"
	structure["text/gr"]["templ_sub1.tt"] = "{define \"sub1\"}-sub1{end}"
	structure["text/gr"]["templ_sub2.tt"] = "{define \"sub2\"}-sub2{end}"

	var err error
	var root string
	root, err = CreateStructure(c, structure)
	c.Assert(err, check.IsNil)

	var textTemplate *goat.TextTemplate = goat.NewTextTemplate("gr", "templ")
	c.Assert(textTemplate.LoadFiles(root), check.IsNil)
	c.Assert(len(textTemplate.Template().Templates()), check.Equals, 3)
}

func (this *TextTemplateSuite) TestLoadFilesWorksWhenPatternMatchesNothing(c *check.C) {
	var structure map[string]map[string]string = make(map[string]map[string]string)
	structure["text/gr"] = make(map[string]string)
	structure["text/gr"]["templ_main.tt"] = "{define \"main\"}main{template \"sub1\"}{template \"sub2\"}{end}"
	structure["text/gr"]["templ_sub1.tt"] = "{define \"sub1\"}-sub1{end}"
	structure["text/gr"]["templ_sub2.tt"] = "{define \"sub2\"}-sub2{end}"

	var err error
	var root string
	root, err = CreateStructure(c, structure)
	c.Assert(err, check.IsNil)

	var textTemplate *goat.TextTemplate = goat.NewTextTemplate("gr", "templ")
	c.Assert(textTemplate.LoadFiles(root), check.IsNil)
	c.Assert(len(textTemplate.Template().Templates()), check.Equals, 3)
}

func (this *TextTemplateSuite) TestLoadFilesErrorsWhenNoFilesMatched(c *check.C) {
	var root string = "."

	var textTemplate *goat.TextTemplate = goat.NewTextTemplate("goat", "template")
	c.Assert(textTemplate.LoadFiles(root), check.NotNil)
	c.Assert(textTemplate.Template(), check.IsNil)
}
