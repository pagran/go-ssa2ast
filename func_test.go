package ssa2ast

import (
	"go/ast"
	"go/importer"
	"go/printer"
	"go/types"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
	"golang.org/x/tools/go/ssa"
	"golang.org/x/tools/go/ssa/ssautil"
)

const funcSignatureSrc = `
package main

type genericStruct[T interface{}] struct {}
type plainStruct struct {
	Dummy struct{}
}

func (s *plainStruct) plainStructFunc()  {   
 
}

func (*plainStruct) plainStructAnonFunc()  {   
 
}

func (s *genericStruct[T]) genericStructFunc()  {   
 
}
func (s *genericStruct[T]) genericStructAnonFunc() (test T) {   
 return
}

func plainFuncSignature(a int, b string, c struct{}, d struct{string}, e interface { Dummy() string }) (i int, er error) {
	return
}

func genericFuncSignature[T interface { interface{} | ~int64 | bool }, X interface{ comparable }](a T, b X, c genericStruct[struct{ a T }], d genericStruct[T]) (res T) {
	return
}
`

func Test_convertSignature(t *testing.T) {
	conv := NewFuncConverter(DefaultConfig())

	f, info, _, _ := mustParseFile(funcSignatureSrc)
	for _, funcName := range []string{"plainStructFunc", "plainStructAnonFunc", "genericStructFunc", "plainFuncSignature", "genericFuncSignature"} {
		funcDecl := findFunc(f, funcName)
		funcDecl.Body = nil

		funcObj := info.Defs[funcDecl.Name].(*types.Func)
		funcDeclConverted, err := conv.convertSignature(funcObj.Name(), funcObj.Type().(*types.Signature))
		if err != nil {
			t.Fatal(err)
		}
		if structDiff := cmp.Diff(funcDecl, funcDeclConverted, astCmpOpt); structDiff != "" {
			t.Fatalf("method decl not equals: %s", structDiff)
		}
	}
}

const funcSrc = `
package main

import (
"encoding/binary"
"time"


)

func dummy(int, int) int  {
 return 0
}

func dummy1() int  {
 return 0
}

type interfaceA interface {
	Test() bool
}
type interfaceB interface {
}

func testInterface(interfaceB)  {
 
}

type testStruct struct {
	A, B string
}

func (testStruct) Test(arg string) string  {
	return "hello" + arg
}

var VariableStr string = "test"

func main() {
	binary.Size(0)

	var intrf interfaceA
	println(intrf.Test())
	
	tst := testStruct{}
	println(dummy1())
	println(tst.Test("hello"))
	go tst.Test("go")
	defer tst.Test("defer")
	println(binary.LittleEndian)
	x := &VariableStr
	println(*x)
	
	// *ssa.Alloc
	alloc := 3
	println(alloc)
	
	var localVar int
	chan1 := make(chan string, 3) 
	_, ok := <-chan1
	chan1 <- "hi"
	println(ok)
	
	// *ssa.If + *ssa.Jump
	if alloc%2 == 0 {
		// *ssa.BinOp 
		localVar = 1
	}
	localVar2 := alloc * localVar + 1
	
	// *ssa.Builtin
	println(localVar2)

	// *ssa.Call
	_ = dummy(localVar, localVar2)
	
	for i, s := range "test" {
	 print(i, ":", s, " ")
	}

	// *ssa.IndexAddr
	slice := [...]int{1,2}
	slice[0] += 1
	println(slice)
	// *ssa.SliceToArrayPointer
	println((*[2]int)(slice[:]))

	// *ssa.MakeMap + *ssa.MapUpdate
	mmap := map[string]time.Month{
		"April":    time.April,
		"December": time.December,
		"January":  time.January,
	}
	for k, v := range mmap {
		println(k, v)
		slice[0]++
	}
	if v, ok := mmap["?"]; ok {
		println(v)
	}
	println(mmap["April"])
	
	// *ssa.ChangeType
	var interA interfaceA
	testInterface(interA)
	
	// *ssa.ChangeInterface
	var interB interfaceB
	var inter0 interface{} = interB
	print(inter0)
	
	// *ssa.Convert
	var f float64 = 1.0
	println(int(f))	

	if f > 2.2 {
		//  *ssa.Panic
		panic(interA)
	}

	// *ssa.Index
	print("test"[slice[0]])

	// *ssa.TypeAssert
	{
		var typ interface{} = "str"
		switch x := typ.(type) {
			case string:
				println(x, "string")
			case int:
				println(x, "int")
			case float32:
				println(x, "float32")
		}
		
		
		println(typ.(string))
	}

	// *ssa.Defer
	defer println("defer 1", f)
	// *ssa.Go
	go println("go 1", f)

	
	strc := testStruct{"A", "B"}
	print(strc.A + strc.B)
	

	a := make(chan string)
	b := make(chan string)
	c := make(chan string)
	d := make(chan string)
	
	select{
		case r1 := <-a:
			println("a triggered", r1)
		case r2 := <-b:
			println("b triggered", r2)
		case r3 := <-c:
			println("c triggered", r3)
		case d<-"test":
			println("d triggered")
		default:
			println("default")
	}
}
`

func Test_Convert(t *testing.T) {
	f, _, fset, _ := mustParseFile(funcSrc)

	ssaPkg, _, err := ssautil.BuildPackage(&types.Config{Importer: importer.Default()}, fset, types.NewPackage("test/main", ""), []*ast.File{f}, ssa.NaiveForm)
	if err != nil {
		panic(err)
	}

	ssaFunc := ssaPkg.Func("main")
	ssaFunc.WriteTo(os.Stdout)

	conv := NewFuncConverter(DefaultConfig())
	astFunc, err := conv.Convert(ssaFunc)
	if err != nil {
		panic(err)
	}

	println("\n\n")

	err = printer.Fprint(os.Stdout, fset, astFunc)
	if err != nil {
		panic(err)
	}
}
