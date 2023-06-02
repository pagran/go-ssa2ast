package ssa2ast

import (
	"go/ast"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_typeToExpr(t *testing.T) {
	f, _, info, _ := mustParseFile("testdata/types/types.go")
	name, structAst := findStruct(f, "exampleStruct")
	obj := info.Defs[name]
	fc := &typeConverter{resolver: defaultImportNameResolver}
	convAst, err := fc.Convert(obj.Type().Underlying())
	if err != nil {
		t.Fatal(err)
	}

	structConvAst := convAst.(*ast.StructType)
	if structDiff := cmp.Diff(structAst, structConvAst, astCmpOpt); structDiff != "" {
		t.Fatalf("struct not equals: %s", structDiff)
	}
}
