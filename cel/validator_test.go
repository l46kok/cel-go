package cel

import (
	"fmt"
	"testing"
)

var testCases = []testInfo{
	// Const types
	{
		in: `A`,
	},
	{
		in: `[1,2,"a"]`,
	},
}

type testInfo struct {
	// in contains the expression to be parsed.
	in string
}

func Test(t *testing.T) {

	// Variables used within this expression environment.
	env, err := NewEnv(
		Variable("A", StringType),
		ASTValidators(NewHomogeneousLiteralValidator()),
	)
	if err != nil {
		t.Fatalf("environment creation error: %s\n", err)
	}

	validator, _ := env.NewValidator()
	for i, tst := range testCases {
		name := fmt.Sprintf("%d %s", i, tst.in)
		tc := tst
		t.Run(name, func(t *testing.T) {
			ast, iss := env.Compile(tc.in)
			if iss.Err() != nil {
				t.Fatal(iss.Err())
			}

			issues := validator.Validate(ast)
			fmt.Println(issues)
		})
	}
}
