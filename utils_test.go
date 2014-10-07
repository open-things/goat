package goat_test

import (
	check "gopkg.in/check.v1"

	"github.com/Paulius-Maruska/goat"
)

type UtilsReplaceAllSuite struct{}

var _ = check.Suite(&UtilsReplaceAllSuite{})

func (this *UtilsReplaceAllSuite) TestReplacesAllFoundPatterns(c *check.C) {
	var expected string = "abcdef"
	var actual string = goat.ReplaceAll(
		"{a}{b}{c}",
		[]string{"{a}", "{b}", "{c}", "{d}"},
		[]string{"ab", "cd", "ef", "gh"},
	)
	c.Assert(actual, check.Equals, expected)
}

func (this *UtilsReplaceAllSuite) TestMismatchingAmounts(c *check.C) {
	var expected string = "abcdcd"
	var actual string = goat.ReplaceAll(
		"{a}{b}{c}",
		[]string{"{a}", "{b}", "{c}", "{d}"},
		[]string{"ab", "cd"},
	)
	c.Assert(actual, check.Equals, expected)
}

type UtilsGetPatternsSuite struct {
	originalTemplatePatterns []string
}

var _ = check.Suite(&UtilsGetPatternsSuite{})

func (this *UtilsGetPatternsSuite) SetUpTest(c *check.C) {
	this.originalTemplatePatterns = goat.TemplatePatterns
}

func (this *UtilsGetPatternsSuite) TearDownTest(c *check.C) {
	goat.TemplatePatterns = this.originalTemplatePatterns
}

type DummyTemplate struct {
	g, n, t, e string
}

func NewDummyTemplate(g, n, t, e string) *DummyTemplate {
	return &DummyTemplate{g: g, n: n, t: t, e: e}
}
func (this *DummyTemplate) Group() string     { return this.g }
func (this *DummyTemplate) Name() string      { return this.n }
func (this *DummyTemplate) Type() string      { return this.t }
func (this *DummyTemplate) Extension() string { return this.e }

func (this *UtilsGetPatternsSuite) TestReplacesFieldsWithValues(c *check.C) {
	goat.TemplatePatterns = []string{
		"{type}-{group}/{name}.{extension}",
		"{type}-{group}/{group}-{name}_*.{extension}",
	}
	var expected []string = []string{
		"t-g/n.e",
		"t-g/g-n_*.e",
	}
	var actual []string = goat.GetPatterns(NewDummyTemplate("g", "n", "t", "e"))
	c.Assert(len(expected), check.Equals, len(actual))
	c.Assert(expected, check.DeepEquals, actual)
}
