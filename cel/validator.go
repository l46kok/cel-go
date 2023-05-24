package cel

import "fmt"

type Validator interface {
	Validate(ast *Ast) *Issues
}

type validator struct {
	env *Env
}

func ASTValidators(astValidators ...*ASTValidator) EnvOption {
	return func(e *Env) (*Env, error) {
		for _, visitor := range astValidators {
			e.validators = append(e.validators, visitor)
		}
		return e, nil
	}
}

func (v *validator) Validate(ast *Ast) *Issues {
	var issues *Issues

	for _, validator := range v.env.validators {
		fmt.Println("Running Validator: " + validator.name)
		issues = validator.Validate(ast)
		if len(issues.Errors()) > 0 {
			return issues
		}
	}
	return nil
}

func (e *Env) NewValidator(opts ...ValidatorOption) (Validator, error) {
	v := &validator{
		env: e,
	}

	var err error
	for _, opt := range opts {
		v, err = opt(v)
		if err != nil {
			return nil, err
		}
	}

	return v, nil
}
