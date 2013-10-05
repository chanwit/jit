package main
/*
 * A libjit wrapper for Golang
 *
 * Copyright (c) 2013 Chanwit Kaewkasi
 * Suranaree University of Technology
 *
 */
import "fmt"
import . "github.com/chanwit/jit"

func main() {
    ctx := NewContext(); defer ctx.Destroy()

    ctx.BuildStart()

    f := ctx.NewFunction(Int(), []Type{Int(), Int(), Int()})

    x, y, z := f.Param3()
    t1 := f.Mul(x,  y)
    t2 := f.Add(t1, z)
    f.Return(t2)

    f.Compile()
    f.Dump()

    ctx.BuildEnd()

    fmt.Printf("mul_add(3, 5, 2): result = %d\n", f.Run(3, 5, 2).(int))
    fmt.Printf("mul_add(5, 5, 5): result = %d\n", f.Run(5, 5, 5).(int))
}

