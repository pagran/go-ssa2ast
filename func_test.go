package ssa2ast

import (
	"go/ast"
	"go/importer"
	"go/printer"
	"go/types"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"golang.org/x/tools/go/ast/astutil"
	"golang.org/x/tools/go/ssa"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/ssa/ssautil"
)

func Test_convertSignature(t *testing.T) {
	conv := newFuncConverter(DefaultConfig())

	f, _, info, _ := mustParseFile("testdata/func/sig.go")
	for _, funcName := range []string{"plainStructFunc", "plainStructAnonFunc", "genericStructFunc", "plainFuncSignature", "genericFuncSignature"} {
		funcDecl := findFunc(f, funcName)
		funcDecl.Body = nil

		funcObj := info.Defs[funcDecl.Name].(*types.Func)
		funcDeclConverted, err := conv.convertSignatureToFuncDecl(funcObj.Name(), funcObj.Type().(*types.Signature))
		if err != nil {
			t.Fatal(err)
		}
		if structDiff := cmp.Diff(funcDecl, funcDeclConverted, astCmpOpt); structDiff != "" {
			t.Fatalf("method decl not equals: %s", structDiff)
		}
	}
}

func Test_Convert(t *testing.T) {
	const testFile = "testdata/func/main.go"
	runGoFile := func(f string) string {
		cmd := exec.Command("go", "run", f)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("compile failed: %v\n%s", err, string(out))
		}
		return string(out)
	}

	originalOut := runGoFile(testFile)

	file, fset, _, _ := mustParseFile(testFile)
	ssaPkg, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, types.NewPackage("test/main", ""), []*ast.File{file}, 0)
	if err != nil {
		panic(err)
	}

	for fIdx, decl := range file.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if !ok {
			continue
		}

		path, _ := astutil.PathEnclosingInterval(file, funcDecl.Pos(), funcDecl.Pos())
		ssaFunc := ssa.EnclosingFunction(ssaPkg, path)

		astFunc, err := Convert(ssaFunc, DefaultConfig())
		if err != nil {
			panic(err)
		}
		file.Decls[fIdx] = astFunc
	}

	convertedFile := filepath.Join(t.TempDir(), "main.go")
	f, err := os.Create(convertedFile)
	if err != nil {
		panic(err)
	}
	if err := printer.Fprint(f, fset, file); err != nil {
		panic(err)
	}
	_ = f.Close()

	convertedOut := runGoFile(convertedFile)

	if convertedOut != originalOut {
		t.Fatalf("Output not equals:\nOriginal: %s\n\nConverted: %s", originalOut, convertedOut)
	}
}
