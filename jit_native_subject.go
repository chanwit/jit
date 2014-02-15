package jit

import "unsafe"

/*
extern int test_native_mult(int a, int b);
*/
import "C"

//export test_native_mult
func test_native_mult(a, b C.int) C.int {
	return a * b
}

func test_native_mult_ptr() unsafe.Pointer { return C.test_native_mult }
