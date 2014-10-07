package goat_test

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	check "gopkg.in/check.v1"
)

func Test(t *testing.T) {
	check.TestingT(t)
}

func CreateStructure(c *check.C, structure map[string]map[string]string) (string, error) {
	var err error
	var rootDir string

	rootDir = c.MkDir()
	c.Logf("created %v", rootDir)

	for subDir, files := range structure {
		var dir string = filepath.Join(rootDir, subDir)
		if err = os.MkdirAll(dir, 0777); err != nil {
			return "", err
		}
		c.Logf("created %v", dir)

		for filename, content := range files {
			var full string = filepath.Join(dir, filename)
			var buf []byte = bytes.NewBufferString(content).Bytes()
			if err = ioutil.WriteFile(full, buf, 0777); err != nil {
				return "", err
			}
			c.Logf("created %v", full)
		}
	}

	return rootDir, nil
}
