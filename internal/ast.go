// participle tags upset govet
// nolint:govet
package internal

type Root struct {
	Statements []Statement `@@*`
}

type Statement struct {
	Rule        *Rule        `@@`
	ScoreChange *ScoreChange `| @@`
	Exit        bool         `| @( "exit" )`
}

type Rule struct {
	Expression   Expression   `"when" @@`
	Consequences Consequences `"then" @@ "done"`
}

type Expression struct {
	Or []OrExpression `@@ { "or" @@ }`
}

type OrExpression struct {
	And []ConditionOrExpression `@@ { "and" @@ }`
}

type ConditionOrExpression struct {
	Condition  *Condition  `@@ `
	Expression *Expression `| "(" @@ ")"`
}

type Condition struct {
	ScalarCondition *ScalarCondition `@@`
	ListCondition   *ListCondition   `| @@`
}

type ScalarCondition struct {
	LeftValue  MixedValue `@@`
	Op         string     `@( "<" { "=" } | ">" { "=" } | "=" "=" | "!" "=" | "contains" | "matches" | "does" "not" ( "match" | "contain" ) )`
	RightValue MixedValue `@@`
}

type ListCondition struct {
	LeftValue   MixedValue   `@@`
	Op          string       `@( { "not" } "in" )`
	RightValues []MixedValue `"[" @@ { "," @@ } "]"`
}

type MixedValue struct {
	Var    *string `"var" "(" @Ident ")"`
	String *string `| @String`
	Int    *int    `| @Int`
	Score  *Score  `| @@`
	Regexp *string `| @Regexp`
}

type Consequences struct {
	Consequences []Statement `@@*`
}

type ScoreChange struct {
	Score    Score    `@@`
	Operator string   `@( "=" | "+" "=" | "-" "=" )`
	Value    IntValue `@@`
}

type Score struct {
	Name string `"score" "(" @Ident ")"`
}

type IntValue struct {
	Int   *int   `@Int`
	Score *Score `| @@`
}
