package asthelper

import (
	"go/ast"
	"go/token"
	"strconv"
)

func BlockStmt(stmts ...ast.Stmt) *ast.BlockStmt {
	return &ast.BlockStmt{List: stmts}
}

func CallExpr(fun ast.Expr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{Fun: fun, Args: args}
}

func CallExprByName(fun string, args ...ast.Expr) *ast.CallExpr {
	return CallExpr(ast.NewIdent(fun), args...)
}

func IndexExpr(xExpr, indexExpr ast.Expr) *ast.IndexExpr {
	return &ast.IndexExpr{X: xExpr, Index: indexExpr}
}

func AssignStmt(Lhs ast.Expr, Rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{Lhs},
		Tok: token.ASSIGN,
		Rhs: []ast.Expr{Rhs},
	}
}

func SelectExpr(x ast.Expr, sel *ast.Ident) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   x,
		Sel: sel,
	}
}

func AssignDefineStmt(Lhs ast.Expr, Rhs ast.Expr) *ast.AssignStmt {
	return &ast.AssignStmt{
		Lhs: []ast.Expr{Lhs},
		Tok: token.DEFINE,
		Rhs: []ast.Expr{Rhs},
	}
}

func IntLit(val int) *ast.BasicLit {
	return &ast.BasicLit{Kind: token.INT, Value: strconv.Itoa(val)}
}
