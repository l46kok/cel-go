package validator

import (
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type HomogeneousLiteralValiadtor struct {
	ASTVisitor
}

func NewHomogeneousLiteralValidator() *HomogeneousLiteralValiadtor {
	h := &HomogeneousLiteralValiadtor{
		ASTVisitor: *NewASTVisitor(),
	}
	h.visitList = h.visitListOverride

	return h
}

func (h *HomogeneousLiteralValiadtor) Name() string {
	return "HomogeneousLiteralValidator"
}

func (h *HomogeneousLiteralValiadtor) visitListOverride(list *exprpb.Expr_CreateList) {
	var elementToCompare *expr.Type = nil
	for _, e := range list.GetElements() {
		h.visitExpr(e)
		currElementType := h.getType(e)
		if elementToCompare != nil {
			if elementToCompare != currElementType {
				h.reportError(e.Id, "Mismatch found")
			}
		}
		elementToCompare = currElementType
	}
}
