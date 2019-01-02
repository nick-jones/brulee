package internal

import (
	"github.com/alecthomas/participle"
	"io"
)

var parser = participle.MustBuild(&Root{})

func Parse(r io.Reader) (Root, error) {
	var rules Root
	err := parser.Parse(r, &rules)
	return rules, err
}