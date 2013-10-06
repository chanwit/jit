package main

import "fmt"
import . "github.com/chanwit/jit"

func compile_mul_add(function *Function) bool {
    fmt.Println(">> Recompile ...")
    x, y, z := function.Param3()
    t1 := function.Mul(x, y)
    t2 := function.Add(t1, z)
    function.Return(t2)

    // compile successfully
    return true
}

func main() {
    ctx := NewContext()
    defer ctx.Destroy()

    ctx.BuildStart()

    f := ctx.NewFunction(Int(), []Type{Int(), Int(), Int()})
    f.SetRecompilable()
    f.SetOnDemandCompiler(compile_mul_add)

    ctx.BuildEnd()

    // First-time on-demand compiler will be called
    fmt.Printf("mul_add(3, 5, 2) = %d\n", f.Run(3, 5, 2).(int))

    // Second-time, nothing compiled
    fmt.Printf("mul_add(13, 5, 7) = %d\n", f.Run(13, 5, 7).(int))

    ctx.BuildStart()
    f.GetOnDemandCompiler()(f) // return closure
    f.Compile()
    ctx.BuildEnd()

    fmt.Printf("mul_add(2, 18, -3) = %d\n", f.Run(2, 18, -3).(int))
}