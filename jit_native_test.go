package jit

import . "launchpad.net/gocheck"

func (s *MySuite) TestJitCallingNative(c *C) {
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
	res := f.CallNative("NativeMult", test_native_mult_ptr(), sig, x, y)
	// f.Store(x, res)
	f.StoreRelative(result, 0, res)

	// Compile the function
	f.Compile()

	// Dump the result to standard output
	// f.Dump("foo")

	// Unlock the context
	ctx.BuildEnd()

	// Execute the function
	var r int = 0
	f.Run(3, 5, &r)

	c.Check(r, Equals, 15)
}
