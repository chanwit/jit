package jit

import "unsafe"

/*
extern int NativeMult(int a, int b);
*/
import "C"

//export NativeMult
func NativeMult(a, b C.int) C.int {
	return a * b
}

func NativeMultPtr() unsafe.Pointer { return C.NativeMult }
