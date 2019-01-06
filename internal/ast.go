package internal

type Root struct {
	Statements []Statement `@@*`
}

type Statement struct {
	Rule        *Rule        `@@`
	ScoreChange *ScoreChange `| @@`
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
	LeftValue  MixedValue `@@`
	Op         string     `@( "<" { "=" } | ">" { "=" } | "=" "=" | "!" "=" | "contains" | "matches" | "does" "not" ( "match" | "contain" ) )`
	RightValue MixedValue `@@`
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
