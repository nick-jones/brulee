package brulee

type Root struct {
	Statements []RootStatement `@@*`
}

type RootStatement struct {
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
	Op         string     `@( "<" { "=" } | ">" { "=" } | "=" "=" | "!" "=" | "does" "not" "contain" | "contains" | "does" "not" "match" | "matches" )`
	RightValue MixedValue `@@`
}

type MixedValue struct {
	Var    *string `"var" "(" @Ident ")"`
	String *string `| @String`
	Int    *int    `| @Int`
	Score  *Score  `| @@`
}

type Consequences struct {
	Consequences []Consequent `@@*`
}

type Consequent struct {
	ScoreChange *ScoreChange `@@`
	SubRule     *Rule        `| @@`
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
