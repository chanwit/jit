package jit

import . "launchpad.net/gocheck"

func (s *MySuite) TestJitCallingJit(c *C) {

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
    t3 := f.Call(f, x, f.Sub(y, x))
    f.Return(t3)

    f.Label(label2)
    // return gcd(x-y, y)
    t4 := f.Call(f, f.Sub(x, y), y)
    f.Return(t4)

    f.Compile()

    ctx.BuildEnd()

    c.Check(f.Run(uint(27), uint(14)), Equals, uint(1))
    c.Check(f.Run(uint(15), uint(30)), Equals, uint(15))
}
