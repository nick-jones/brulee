## Overview

Brulee is a basic string matching and scoring rules engine.

## Installation

```
go get github.com/nick-jones/brulee
```

## Example

The following is a minimal example where a variable is checked for equality against a string.

```go
package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nick-jones/brulee"
)

func main() {
	input := `
	  when
	    var(foo) == "bar"
	  then
	    score(foo) = 1
	  done
	`
	program := brulee.MustCompile(strings.NewReader(input))
	vars := map[string]string{"foo": "bar"}
	scores, err := program.Run(vars)
	if err != nil {
		panic(err)
	}
	for key, value := range scores {
		fmt.Printf("%s = %d\n", key, value)
	}
}
```

The above program results in the following output:

```
foo = 1
```

## Advanced Example

A more advanced example is contained with the [example directory](example).

## Internals

The input program is parsed and instructions are generated from it - these are then executed. The instructions can be
viewed by invoking the `Dump()` function against the compiled program. By way of example:


```go
package main

import (
	"os"
	"strings"

	"github.com/nick-jones/brulee"
)

func main() {
	input := `
	  when
	    (var(x) == "foo" or var(x) == "bar")
	    and var(y) matches /ba[rz]/
	  then
	    score(foo) = 1
	  done
	`
	program := brulee.MustCompile(strings.NewReader(input))
	program.Dump(os.Stdout)
}
```

The above outputs:

```
+-----+------------------+-----+------------+------------------+
| POS |        OP        | RET |  OPERAND1  |     OPERAND2     |
+-----+------------------+-----+------------+------------------+
|   0 | IS_EQUAL         | $3  | var(x)     | string("foo")    |
|   1 | JUMP_IF_NOT_ZERO | $2  | $3         | ->4              |
|   2 | IS_EQUAL         | $3  | var(x)     | string("bar")    |
|   3 | JUMP_IF_NOT_ZERO | $2  | $3         | ->4              |
|   4 | JUMP_IF_ZERO     | $1  | $2         | ->7              |
|   5 | MATCHES          | $2  | var(y)     | regexp(/ba[rz]/) |
|   6 | JUMP_IF_ZERO     | $1  | $2         | ->7              |
|   7 | NEGATE           | $1  | $1         |                  |
|   8 | JUMP_IF_ZERO     |     | $1         | ->10             |
|   9 | SET_SCORE        |     | score(foo) | int(1)           |
|  10 | NOOP             |     |            |                  |
+-----+------------------+-----+------------+------------------+
```
