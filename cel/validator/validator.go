package validator

import (
	"github.com/google/cel-go/cel"
)

type Validator struct {
	validators map[string]ASTValidator
}

func NewValidator(opts ...ValidatorOption) (*Validator, error) {
	v := &Validator{
		validators: map[string]ASTValidator{},
	}

	for _, opt := range opts {
		if err := opt(v); err != nil {
			return nil, err
		}
	}

	return v, nil
}

func (v *Validator) Validate(ast *cel.Ast) *cel.Issues {
	var issues *cel.Issues = &cel.Issues{}

	for _, validator := range v.validators {
		issues = validator.Validate(ast)
		if len(issues.Errors()) > 0 {
			return issues
		}
	}
	return nil
}

func ASTValidators(astValidators ...ASTValidator) ValidatorOption {
	return func(v *Validator) error {
		for _, validator := range astValidators {
			v.validators[validator.Name()] = validator
		}
		return nil
	}
}

type ASTValidator interface {
	Name() string
	Validate(ast *cel.Ast) *cel.Issues
}
