package main

import "fmt"
import . "../../jit"

func main() {

	ctx := NewContext()
	defer ctx.Destroy()

	ctx.BuildStart()

	f := ctx.NewFunction(UInt(), []Type{UInt(), UInt()})

	x, y := f.Param2()

	label1 := NewLabel()
	label2 := NewLabel()

	// Check the condition "if(x == y)"
	t1 := f.Eq(x, y)
	f.BranchIfNot(t1, label1)

	// Implement "return x"
	f.Return(x)

	f.Label(label1)
	// if (x < y)
	t2 := f.LessThan(x, y)
	f.BranchIfNot(t2, label2)

	// return gcd(x, y-x)
	t3 := f.TailCall(f, x, f.Sub(y, x))
	f.Return(t3)

	f.Label(label2)
	// return gcd(x-y, y)
	t4 := f.TailCall(f, f.Sub(x, y), y)
	f.Return(t4)

	f.Compile()

	ctx.BuildEnd()

	f.Dump("gcd")

	// Execute the function and print the result
	result := f.Run(uint(27), uint(14))
	fmt.Printf("gcd(27, 14) = %d\n", result.(uint))

}
