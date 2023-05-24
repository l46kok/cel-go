package validator

import (
	"fmt"

	"github.com/google/cel-go/cel"
)

type Validator struct {
	validators []*ASTVisitor
}

func NewValidator(opts ...ValidatorOption) (*Validator, error) {
	v := &Validator{
		validators: []*ASTVisitor{},
	}

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (v *Validator) Validate(ast *cel.Ast) *cel.Issues {
	var issues *cel.Issues

	for _, validator := range v.validators {
		fmt.Println("Running Validator: " + validator.name)
		issues = validator.Validate(ast)
		if len(issues.Errors()) > 0 {
			return issues
		}
	}
	return nil
}

func ASTVisitors(astVisitors ...*ASTVisitor) ValidatorOption {
	return func(v *Validator) error {
		for _, visitor := range astVisitors {
			v.validators = append(v.validators, visitor)
		}
		return nil
	}
}
