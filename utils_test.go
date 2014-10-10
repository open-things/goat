package goat_test

import (
	check "gopkg.in/check.v1"

	"github.com/open-things/goat"
)

type UtilsReplaceAllSuite struct{}

var _ = check.Suite(&UtilsReplaceAllSuite{})

func (this *UtilsReplaceAllSuite) TestReplacesAllFoundPatterns(c *check.C) {
	expected := "abcdef"
	actual := goat.ReplaceAll(
		"{a}{b}{c}",
		[]string{"{a}", "{b}", "{c}", "{d}"},
		[]string{"ab", "cd", "ef", "gh"},
	)
	c.Assert(actual, check.Equals, expected)
}

func (this *UtilsReplaceAllSuite) TestMismatchingAmounts(c *check.C) {
	expected := "abcdcd"
	actual := goat.ReplaceAll(
		"{a}{b}{c}",
		[]string{"{a}", "{b}", "{c}", "{d}"},
		[]string{"ab", "cd"},
	)
	c.Assert(actual, check.Equals, expected)
}

type UtilsCopyMapSuite struct{}

var _ = check.Suite(&UtilsCopyMapSuite{})

func (this *UtilsCopyMapSuite) TestMakesACopy(c *check.C) {
	original := map[string]interface{}{
		"foo": "bar",
	}
	copy := goat.CopyMap(original)
	c.Assert(copy["foo"], check.Equals, "bar")
	copy["foo"] = "baz"
	c.Assert(original["foo"], check.Equals, "bar")
}

func (this *UtilsCopyMapSuite) TestCreatesEmptyMapEvenWhenSourceIsNil(c *check.C) {
	c.Assert(goat.CopyMap(nil), check.NotNil)
}
