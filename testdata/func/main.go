package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"sort"
	"strconv"
	"sync"
	"time"
	"unsafe"
)

func main() {
	methodOps()
	slicesOps()
	iterAndMapsOps()
	chanOps()
	flowOps()
	typeOps()
}

func makeSprintf(tag string) func(vals ...interface{}) {
	i := 0
	return func(vals ...interface{}) {
		fmt.Printf("%s(%d): %v\n", tag, i, vals)
		i++
	}
}

func return42() int {
	return 42
}

type arrayOfInts []int

type structOfArraysOfInts struct {
	a arrayOfInts
	b arrayOfInts
}

func slicesOps() {
	sprintf := makeSprintf("slicesOps")

	slice := [...]int{1, 2}
	sprintf(slice[0:1:2])
	// *ssa.IndexAddr
	sprintf(slice)
	slice[0] += 1
	sprintf(slice)

	sprintf(slice[:1])
	sprintf(slice[slice[0]:])
	sprintf(slice[0:2])

	sprintf((*[2]int)(slice[:])[return42()%2]) // *ssa.SliceToArrayPointer

	sprintf("test"[return42()%3]) // *ssa.Index

	structOfArrays := structOfArraysOfInts{a: slice[1:], b: slice[:1]}
	sprintf(structOfArrays.a[:1])
	sprintf(structOfArrays.b[:1])

	slice2 := make([]string, return42(), return42()*2)
	slice2[return42()-1] = "test"
	sprintf(slice2)

	return
}

func iterAndMapsOps() {
	sprintf := makeSprintf("iterAndMapsOps")

	// *ssa.MakeMap + *ssa.MapUpdate
	mmap := map[string]time.Month{
		"April":    time.April,
		"December": time.December,
		"January":  time.January,
	}

	var vals []string
	for k := range mmap {
		vals = append(vals, k)
	}
	for _, v := range mmap {
		vals = append(vals, v.String())
	}
	sort.Strings(vals) // Required. Order of map iteration not guaranteed
	sprintf(vals)

	if v, ok := mmap["?"]; ok {
		panic("unreachable: " + v.String())
	}
	for idx, s := range "hello world" {
		sprintf(idx, s)
	}

	sprintf(mmap["April"].String())
	return
}

type interfaceCalls interface {
	Return1() string
}

type structCalls struct {
}

func (r structCalls) Return1() string {
	return "Return1"
}

func (r *structCalls) Return2() string {
	return "Return2"
}

func multiOutputRes() (int, string) {
	return 42, "24"
}

func returnInterfaceCalls() interfaceCalls {
	return structCalls{}
}

func methodOps() {
	sprintf := makeSprintf("methodOps")

	defer func() {
		sprintf("from defer")
	}()
	defer sprintf("from defer 2")

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		sprintf("from go")
		wg.Done()
	}()
	wg.Wait()

	i, s := multiOutputRes()
	sprintf(strconv.Itoa(i))

	var strct structCalls

	strct.Return1()
	strct.Return2()

	intrfs := returnInterfaceCalls()
	intrfs.Return1()

	sprintf(strconv.Itoa(len(s)))

	strconv.Itoa(binary.Size(4))
	sprintf(binary.LittleEndian.AppendUint32(nil, 42))

	if len(s) == 0 {
		panic("unreachable")
	}

	sprintf(*unsafe.StringData(s))

	thunkMethod1 := structCalls.Return1
	sprintf(thunkMethod1(strct))

	thunkMethod2 := (*structCalls).Return2
	sprintf(thunkMethod2(&strct))

	closureVar := "c " + s
	anonFnc := func(n func(structCalls) string) string {
		return n(structCalls{}) + "anon" + closureVar
	}

	sprintf(anonFnc(structCalls.Return1))
}

func chanOps() {
	sprintf := makeSprintf("chanOps")

	a := make(chan string)
	b := make(chan string)
	c := make(chan string)
	d := make(chan string)

	select {
	case r1, ok := <-a:
		sprintf(r1, ok)
	case r2 := <-b:
		sprintf(r2)
	case <-c:
		sprintf("r3")
	case d <- "test":
		sprintf("d triggered")
	default:
		sprintf("default")
	}

	e := make(chan string, 1)
	e <- "hi"

	sprintf(<-e)

	close(a)
	val, ok := <-a

	sprintf(val, ok)
	return
}

func flowOps() {
	sprintf := makeSprintf("flowOps")
	i := 1
	if return42()%2 == 0 {
		sprintf("a")
		i++
	} else {
		sprintf("b")
	}
	sprintf(i)

	switch return42() {
	case 1:
		sprintf("1")
	case 2:
		sprintf("2")
	case 3:
		sprintf("3")
	case 42:
		sprintf("42")
	}
}

type interfaceB interface {
}

type testStruct struct {
	A, B int
}

func typeOps() {
	sprintf := makeSprintf("typeOps")

	// *ssa.ChangeType
	var interA interfaceCalls
	sprintf(interA)

	// *ssa.ChangeInterface
	var interB interfaceB = struct{}{}
	var inter0 interface{} = interB
	sprintf(inter0)

	// *ssa.Convert
	var f float64 = 1.0
	sprintf(int(f))

	casted, ok := inter0.(interfaceB)
	sprintf(casted, ok)

	casted2 := inter0.(interfaceB)
	sprintf(casted2)

	strc := testStruct{return42(), return42() + 2}
	strc.B += strc.A
	sprintf(strc)

	// Access to unexported structure
	discard := io.Discard
	if return42() == 0 {
		sprintf(discard) // Trigger phi block
	}
	_, _ = discard.Write([]byte("test"))
}
