package internal

import (
	"io"

	"github.com/alecthomas/participle"
	"github.com/alecthomas/participle/lexer"
	"github.com/alecthomas/participle/lexer/ebnf"
)

var (
	lexr = lexer.Must(ebnf.New(`
		Comment = "#" { "\u0000"…"\uffff"-"\n" } .
		Ident = (alpha | "_") { "_" | alpha | digit } .
		String = "\"" { "\u0000"…"\uffff"-"\""-"\\" | "\\" any } "\"" .
		Regexp = "/" { "\u0000"…"\uffff"-"/"-"\\" | "\\" any } "/" .
		Int = [ "-" | "+" ] digit { digit } .
		Whitespace = " " | "\t" | "\n" | "\r" .
		Punct = "!"…"/" | ":"…"@" | "["…` + "\"`\"" + ` | "{"…"~" .
		alpha = "a"…"z" | "A"…"Z" .
		digit = "0"…"9" .
		any = "\u0000"…"\uffff" .
	`))
	parser = participle.MustBuild(
		&Root{},
		participle.Lexer(lexr),
		participle.Unquote("String"),
		participle.Elide("Whitespace", "Comment"),
		removeRegexpSlashes("Regexp"),
	)
)

func Parse(r io.Reader) (Root, error) {
	var rules Root
	err := parser.Parse(r, &rules)
	return rules, err
}

func removeRegexpSlashes(types ...string) participle.Option {
	return participle.Map(func(t lexer.Token) (lexer.Token, error) {
		if len(t.Value) < 3 || t.Value[0] != '/' || t.Value[len(t.Value)-1] != '/' {
			return t, lexer.Errorf(t.Pos, "invalid regexp %s", t.Value)
		}
		t.Value = t.Value[1 : len(t.Value)-1]
		return t, nil
	}, types...)
}
