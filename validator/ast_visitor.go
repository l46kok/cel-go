package validator

import (
	"github.com/google/cel-go/cel"
	"github.com/google/cel-go/common"
	exprpb "google.golang.org/genproto/googleapis/api/expr/v1alpha1"
)

type ASTVisitor struct {
	name    string
	issues  *cel.Issues
	visitor *CustomVisitor
}

type CustomVisitor struct {
	visitConst VisitConster
	visitIdent VisitIdenter
	visitList  VisitLister
}

type VisitConster func(ast *cel.Ast, constant *exprpb.Constant) *cel.Issues
type VisitIdenter func(ast *cel.Ast, ident *exprpb.Expr_Ident) *cel.Issues
type VisitLister func(ast *cel.Ast, list *exprpb.Expr_CreateList) *cel.Issues

func NewASTVisitor(name string, customVisitor *CustomVisitor) *ASTVisitor {
	if customVisitor.visitConst == nil {
		customVisitor.visitConst = func(ast *cel.Ast, constant *exprpb.Constant) *cel.Issues { return nil }
	}
	if customVisitor.visitIdent == nil {
		customVisitor.visitIdent = func(ast *cel.Ast, ident *exprpb.Expr_Ident) *cel.Issues { return nil }
	}
	if customVisitor.visitList == nil {
		customVisitor.visitList = func(ast *cel.Ast, list *exprpb.Expr_CreateList) *cel.Issues { return nil }
	}

	b := &ASTVisitor{
		name:    name,
		issues:  &cel.Issues{},
		visitor: customVisitor,
	}
	return b
}

func (b *ASTVisitor) appendToIssues(iss *cel.Issues) {
	b.issues = b.issues.Append(iss)
}

func (b *ASTVisitor) Validate(ast *cel.Ast) *cel.Issues {
	b.issues = cel.NewIssues(common.NewErrors(ast.Source()))
	b.visitExpr(ast, ast.Expr())
	return b.issues
}

func (b *ASTVisitor) visitExpr(ast *cel.Ast, e *exprpb.Expr) {
	switch e.GetExprKind().(type) {
	case *exprpb.Expr_ConstExpr:
		b.appendToIssues(b.visitor.visitConst(ast, e.GetConstExpr()))
	case *exprpb.Expr_IdentExpr:
		b.appendToIssues(b.visitor.visitIdent(ast, e.GetIdentExpr()))
	case *exprpb.Expr_ListExpr:
		b.defaultVisitList(ast, e.GetListExpr())
	}

}

func (b *ASTVisitor) defaultVisitList(ast *cel.Ast, list *exprpb.Expr_CreateList) {
	b.appendToIssues(b.visitor.visitList(ast, list))
	for _, e := range list.GetElements() {
		b.visitExpr(ast, e)
	}
}
