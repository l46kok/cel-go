package cel

import (
	"github.com/google/cel-go/common"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

func NewHomogeneousLiteralValidator() *ASTValidator {
	return NewASTValidator(
		"HomogeneousLiteralValidator",
		&VisitorOverrides{
			visitList: visitListOverride,
		})
}

func visitListOverride(ast *Ast, list *exprpb.Expr_CreateList) *Issues {
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

	return NewIssues(common.NewErrors(ast.Source()).Append(errors))
}
