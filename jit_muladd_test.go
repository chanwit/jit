package jit

import (
    . "launchpad.net/gocheck"
    "testing"
)

// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type MySuite struct { }

var _ = Suite(&MySuite{})

func (s *MySuite) TestJitMulAdd(c *C) {
    // Create a context to hold the JIT's primary state
    // defer to Destroy the context at the end of func main()
    ctx := NewContext()
    defer ctx.Destroy()

    // Lock the context while we build and compile the function
    ctx.BuildStart()

    // Create the function object
    f := ctx.NewFunction(Int(), []Type{Int(), Int(), Int()})

    // Construct the function body
    x, y, z := f.Param3()
    t1 := f.Mul(x,  y)
    t2 := f.Add(t1, z)
    f.Return(t2)

    // Compile the function
    f.Compile()

    // Unlock the context
    ctx.BuildEnd()

    c.Check(f.Run(3, 5, 2), Equals, 17)
}
