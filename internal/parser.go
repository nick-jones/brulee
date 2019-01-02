package internal

import (
	"io"

	"github.com/alecthomas/participle"
)

var parser = participle.MustBuild(&Root{})

func Parse(r io.Reader) (Root, error) {
	var rules Root
	err := parser.Parse(r, &rules)
	return rules, err
}
