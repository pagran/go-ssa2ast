package main

import "unsafe"

type genericStruct[T interface{}] struct{}
type plainStruct struct {
	Dummy struct{}
}

func (s *plainStruct) plainStructFunc() {

}

func (*plainStruct) plainStructAnonFunc() {

}

func (s *genericStruct[T]) genericStructFunc() {

}
func (s *genericStruct[T]) genericStructAnonFunc() (test T) {
	return
}

func plainFuncSignature(a int, b string, c struct{}, d struct{ string }, e interface{ Dummy() string }, pointer unsafe.Pointer) (i int, er error) {
	return
}

func genericFuncSignature[T interface{ interface{} | ~int64 | bool }, X interface{ comparable }](a T, b X, c genericStruct[struct{ a T }], d genericStruct[T]) (res T) {
	return
}
