package validator

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func NewHomogeneousLiteralValidator() *ASTVisitor {
	h := NewASTVisitor(
		"HomogeneousLiteralValidator",
		&CustomVisitor{
			visitList: visitListOverride,
		})
	return h
}

func visitListOverride(ast *cel.Ast, list *exprpb.Expr_CreateList) *cel.Issues {
	var errors []common.Error = []common.Error{}
	var elementToCompare *expr.Type = nil
	for _, e := range list.GetElements() {
		currElementType := ast.GetType(e.Id)
		if elementToCompare != nil {
			if elementToCompare != currElementType {
				errors = append(errors, common.Error{
					Location: common.NoLocation,
					Message:  "Mismatch Found",
				})
			}
		}
		elementToCompare = currElementType
	}

	return cel.NewIssues(common.NewErrors(ast.Source()).Append(errors))
}
