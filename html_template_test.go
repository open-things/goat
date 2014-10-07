package goat_test

import (
	check "gopkg.in/check.v1"

	"github.com/Paulius-Maruska/goat"
)

type HtmlTemplateSuite struct{}

var _ = check.Suite(&HtmlTemplateSuite{})

func (this *HtmlTemplateSuite) TestNewAndGetters(c *check.C) {
	var htmlTemplate *goat.HtmlTemplate = goat.NewHtmlTemplate("g", "n")
	c.Assert(htmlTemplate.Template(), check.IsNil)
	c.Assert(htmlTemplate.Group(), check.Equals, "g")
	c.Assert(htmlTemplate.Name(), check.Equals, "n")
	c.Assert(htmlTemplate.Type(), check.Equals, goat.HtmlTemplateType)
	c.Assert(htmlTemplate.Extension(), check.Equals, goat.HtmlTemplateExtension)
}

func (this *HtmlTemplateSuite) TestLoadFilesParsesEverything(c *check.C) {
	var structure map[string]map[string]string = make(map[string]map[string]string)
	structure["html/gr"] = make(map[string]string)
	structure["html/gr"]["templ.ht"] = "{define \"main\"}main{template \"sub1\"}{template \"sub2\"}{end}"
	structure["html/gr"]["templ_sub1.ht"] = "{define \"sub1\"}-sub1{end}"
	structure["html/gr"]["templ_sub2.ht"] = "{define \"sub2\"}-sub2{end}"

	var err error
	var root string
	root, err = CreateStructure(c, structure)
	c.Assert(err, check.IsNil)

	var htmlTemplate *goat.HtmlTemplate = goat.NewHtmlTemplate("gr", "templ")
	c.Assert(htmlTemplate.LoadFiles(root), check.IsNil)
	c.Assert(len(htmlTemplate.Template().Templates()), check.Equals, 4)
}

func (this *HtmlTemplateSuite) TestLoadFilesWorksWhenPatternMatchesNothing(c *check.C) {
	var structure map[string]map[string]string = make(map[string]map[string]string)
	structure["html/gr"] = make(map[string]string)
	structure["html/gr"]["templ_main.ht"] = "{define \"main\"}main{template \"sub1\"}{template \"sub2\"}{end}"
	structure["html/gr"]["templ_sub1.ht"] = "{define \"sub1\"}-sub1{end}"
	structure["html/gr"]["templ_sub2.ht"] = "{define \"sub2\"}-sub2{end}"

	var err error
	var root string
	root, err = CreateStructure(c, structure)
	c.Assert(err, check.IsNil)

	var htmlTemplate *goat.HtmlTemplate = goat.NewHtmlTemplate("gr", "templ")
	c.Assert(htmlTemplate.LoadFiles(root), check.IsNil)
	c.Assert(len(htmlTemplate.Template().Templates()), check.Equals, 4)
}

func (this *HtmlTemplateSuite) TestLoadFilesErrorsWhenNoFilesMatched(c *check.C) {
	var root string = "."

	var htmlTemplate *goat.HtmlTemplate = goat.NewHtmlTemplate("goat", "template")
	c.Assert(htmlTemplate.LoadFiles(root), check.NotNil)
	c.Assert(htmlTemplate.Template(), check.IsNil)
}
