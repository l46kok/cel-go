package cel

import (
	"github.com/google/cel-go/common"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type ASTValidator struct {
	name     string
	issues   *Issues
	visitors *VisitorOverrides
}

type VisitorOverrides struct {
	visitConst VisitConster
	visitIdent VisitIdenter
	visitList  VisitLister
}

type VisitConster func(ast *Ast, constant *exprpb.Constant) *Issues
type VisitIdenter func(ast *Ast, ident *exprpb.Expr_Ident) *Issues
type VisitLister func(ast *Ast, list *exprpb.Expr_CreateList) *Issues

func NewASTValidator(name string, visitors *VisitorOverrides) *ASTValidator {
	if visitors.visitConst == nil {
		visitors.visitConst = func(ast *Ast, constant *exprpb.Constant) *Issues { return nil }
	}
	if visitors.visitIdent == nil {
		visitors.visitIdent = func(ast *Ast, ident *exprpb.Expr_Ident) *Issues { return nil }
	}
	if visitors.visitList == nil {
		visitors.visitList = func(ast *Ast, list *exprpb.Expr_CreateList) *Issues { return nil }
	}

	b := &ASTValidator{
		name:     name,
		issues:   &Issues{},
		visitors: visitors,
	}
	return b
}

func (b *ASTValidator) appendToIssues(iss *Issues) {
	b.issues = b.issues.Append(iss)
}

func (b *ASTValidator) Validate(ast *Ast) *Issues {
	b.issues = NewIssues(common.NewErrors(ast.Source()))
	b.visitExpr(ast, ast.Expr())
	return b.issues
}

func (b *ASTValidator) visitExpr(ast *Ast, e *exprpb.Expr) {
	switch e.GetExprKind().(type) {
	case *exprpb.Expr_ConstExpr:
		b.appendToIssues(b.visitors.visitConst(ast, e.GetConstExpr()))
	case *exprpb.Expr_IdentExpr:
		b.appendToIssues(b.visitors.visitIdent(ast, e.GetIdentExpr()))
	case *exprpb.Expr_ListExpr:
		b.visitList(ast, e.GetListExpr())
	}
}

func (b *ASTValidator) visitList(ast *Ast, list *exprpb.Expr_CreateList) {
	b.appendToIssues(b.visitors.visitList(ast, list))
	for _, e := range list.GetElements() {
		b.visitExpr(ast, e)
	}
}
