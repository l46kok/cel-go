package validator

import (
	"fmt"

	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	expr "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type ASTVisitor struct {
	errors     []common.Error
	visitConst VisitConster
	visitIdent VisitIdenter
	visitList  VisitLister
	ast        *cel.Ast
}

type VisitConster func(constant *exprpb.Constant)
type VisitIdenter func(constant *exprpb.Expr_Ident)
type VisitLister func(constant *exprpb.Expr_CreateList)

func NewASTVisitor() *ASTVisitor {
	b := &ASTVisitor{
		errors: []common.Error{},
	}
	b.visitConst = b.defaultVisitConstant
	b.visitIdent = b.defaultVisitIdent
	b.visitList = b.defaultVisitList
	return b
}

func (b *ASTVisitor) Validate(ast *cel.Ast) *cel.Issues {
	b.ast = ast
	b.visitExpr(ast.Expr())

	return cel.NewIssues(common.NewErrors(ast.Source()).Append(b.errors))
}

func (b *ASTVisitor) visitExpr(e *exprpb.Expr) {
	switch e.GetExprKind().(type) {
	case *exprpb.Expr_ConstExpr:
		b.visitConst(e.GetConstExpr())
	case *exprpb.Expr_IdentExpr:
		b.visitIdent(e.GetIdentExpr())
	case *exprpb.Expr_ListExpr:
		b.visitList(e.GetListExpr())
	}
}

func (b *ASTVisitor) reportError(exprId int64, msg string) {
	b.errors = append(b.errors, common.Error{Message: msg})
}

func (b *ASTVisitor) getType(e *exprpb.Expr) *expr.Type {
	return b.ast.GetType(e.Id)
}

func (b *ASTVisitor) defaultVisitConstant(constant *exprpb.Constant) {
	fmt.Println("Base Visit Constant")
}

func (b *ASTVisitor) defaultVisitIdent(ident *exprpb.Expr_Ident) {
	fmt.Println("Base Visit Ident")
}

func (b *ASTVisitor) defaultVisitList(list *exprpb.Expr_CreateList) {
	for _, e := range list.GetElements() {
		b.visitExpr(e)
	}
}
