package ssa2ast

import (
	"go/ast"
	"go/importer"
	"go/parser"
	"go/token"
	"go/types"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/google/go-cmp/cmp/cmpopts"
)

var astCmpOpt = cmpopts.IgnoreTypes(token.NoPos, &ast.Object{})

func findStruct(file *ast.File, structName string) (name *ast.Ident, structType *ast.StructType) {
	ast.Inspect(file, func(node ast.Node) bool {
		if structType != nil {
			return false
		}

		typeSpec, ok := node.(*ast.TypeSpec)
		if !ok || typeSpec.Name == nil || typeSpec.Name.Name != structName {
			return true
		}
		typ, ok := typeSpec.Type.(*ast.StructType)
		if !ok {
			return true
		}
		structType = typ
		name = typeSpec.Name
		return true
	})

	if structType == nil {
		panic(structName + " not found")
	}
	return
}

func findFunc(file *ast.File, funcName string) (funcDecl *ast.FuncDecl) {
	ast.Inspect(file, func(node ast.Node) bool {
		if funcDecl != nil {
			return false
		}

		fDecl, ok := node.(*ast.FuncDecl)
		if !ok || fDecl.Name.Name != funcName {
			return true
		}
		funcDecl = fDecl
		return true
	})
	if funcDecl == nil {
		panic(funcName + " not found")
	}
	return
}

func mustParseFile(fileName string) (*ast.File, *token.FileSet, *types.Info, *types.Package) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, fileName, nil, 0)
	if err != nil {
		panic(err)
	}

	config := types.Config{Importer: importer.Default()}
	info := &types.Info{
		Types:      make(map[ast.Expr]types.TypeAndValue),
		Defs:       make(map[*ast.Ident]types.Object),
		Uses:       make(map[*ast.Ident]types.Object),
		Instances:  make(map[*ast.Ident]types.Instance),
		Implicits:  make(map[ast.Node]types.Object),
		Scopes:     make(map[ast.Node]*types.Scope),
		Selections: make(map[*ast.SelectorExpr]*types.Selection),
	}
	pkg, err := config.Check("test/main", fset, []*ast.File{f}, info)
	if err != nil {
		panic(err)
	}
	return f, fset, info, pkg
}

func mustParseDirectory(dir string) ([]*ast.File, []string, *token.FileSet) {
	var files []*ast.File
	var fileNames []string
	fset := token.NewFileSet()
	err := filepath.Walk(dir, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !strings.HasSuffix(path, ".go") {
			return nil
		}
		f, err := parser.ParseFile(fset, path, nil, 0)
		if err != nil {
			return err
		}
		files = append(files, f)
		fileNames = append(fileNames, filepath.Base(path))
		return nil
	})
	if err != nil {
		panic(err)
	}

	return files, fileNames, fset
}
