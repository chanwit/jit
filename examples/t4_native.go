package main

/*
 * A libjit wrapper for Golang
 *
 * Copyright (c) 2013-2014 Chanwit Kaewkasi
 * Suranaree University of Technology
 *
 */
import "fmt"
import . "github.com/chanwit/jit"

/*
extern int NativeMult(int a,int b);
*/
import "C"

//export NativeMult
func NativeMult(a, b C.int) C.int {
	fmt.Println(">>>>Test")
	return a * b
}

func main() {

	// Create a context to hold the JIT's primary state
	// defer to Destroy the context at the end of func main()
	ctx := NewContext()
	defer ctx.Destroy()

	// Lock the context while we build and compile the function
	ctx.BuildStart()

	// Create the function object
	// void foo(int x, int y, int* result) {
	//   *result = NativeMult(x, y);
	// }
	f := ctx.NewFunction(Void(), []Type{Int(), Int(), VoidPtr()})

	// Construct the function body
	x, y, result := f.Param3()

	// This is native call
	sig := NewSignature(Int(), []Type{Int(), Int()})
	res := f.CallNative("NativeMult", C.NativeMult, sig, x, y)
	f.Store(x, res)
	f.StoreRelative(result, 0, x)

	f.Dump("foo [uncompiled]")

	// Compile the function
	f.Compile()

	// Dump the result to standard output
	f.Dump("foo [compiled]")

	// Unlock the context
	ctx.BuildEnd()

	// Execute the function
	var r int = 0
	f.Run(3, 5, &r)
	fmt.Println(r)
}
