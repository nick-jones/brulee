package brulee

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

func ParseBytes(b []byte) (Root, error) {
	var rules Root
	err := parser.ParseBytes(b, &rules)
	return rules, err
}

func ParseString(s string) (Root, error) {
	var rules Root
	err := parser.ParseString(s, &rules)
	return rules, err
}
