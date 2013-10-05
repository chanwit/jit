package jit
/*
#cgo LDFLAGS: -lm -ljit
#include <stdio.h>
#include <jit/jit.h>
#include <jit/jit-dump.h>
*/
import "C"

import "unsafe"

const CDECL = C.jit_abi_cdecl

type Context struct {
    C C.jit_context_t
}

type Type struct {
    C C.jit_type_t
}

type Signature struct {
    C C.jit_type_t
}

type Function struct {
    C         C.jit_function_t
    retType   Type
    paramType []Type
}

type Value struct {
    C C.jit_value_t
}

var intType = Type{ C.jit_type_int }
func Int() Type {
    return intType
}

func NewContext() *Context {
    return &Context{ C.jit_context_create() }
}

func NewSignature(ret Type, params []Type) *Signature {
    signature := C.jit_type_create_signature(C.jit_abi_cdecl,
        C.jit_type_int,
        (*C.jit_type_t)(unsafe.Pointer(&params[0])),
        C.uint(len(params)), 1)
    return &Signature{signature}
}

func (c *Context) BuildStart() {
    C.jit_context_build_start(c.C)
}

func (c *Context) BuildEnd() {
    C.jit_context_build_end(c.C)
}

func (c *Context) Destroy() {
    C.jit_context_destroy(c.C)
}

func (c *Context) NewFunction(ret Type, params []Type) *Function {
    signature := NewSignature(ret, params)
    function := C.jit_function_create(c.C, signature.C)
    C.jit_type_free(signature.C)
    return &Function{C: function, retType: ret, paramType: params}
}

func (f *Function) Param(i int) *Value {
    return &Value{C.jit_value_get_param(f.C, C.uint(i))}
}

func (f *Function) Param2() (*Value, *Value) {
    return &Value{C.jit_value_get_param(f.C, C.uint(0))},
           &Value{C.jit_value_get_param(f.C, C.uint(1))}
}

func (f *Function) Param3() (*Value, *Value, *Value) {
    return &Value{C.jit_value_get_param(f.C, C.uint(0))},
           &Value{C.jit_value_get_param(f.C, C.uint(1))},
           &Value{C.jit_value_get_param(f.C, C.uint(2))}
}

func (f *Function) Mul(a, b *Value) *Value {
    return &Value{ C.jit_insn_mul(f.C, a.C, b.C) }
}

func (f *Function) Add(a, b *Value) *Value {
    return &Value{ C.jit_insn_add(f.C, a.C, b.C) }
}

func (f *Function) Return(ret *Value) {
    C.jit_insn_return(f.C, ret.C)
}

func (f *Function) Compile() {
    C.jit_function_compile(f.C)
}

func (f *Function) Run(values ...interface{}) interface{} {
    args := make([]unsafe.Pointer, len(values))
    for i := range values {
        switch f.paramType[i].C {
            case C.jit_type_int:
                x := C.jit_int(values[i].(int))
                args[i] = (unsafe.Pointer)(&x)
        }
    }

    switch f.retType.C {
        case C.jit_type_int:
            result := C.jit_int(0)
            C.jit_function_apply(f.C, (*unsafe.Pointer)(&args[0]), unsafe.Pointer(&result))
            return int(result)

    }

    return nil
}

func (f *Function) Dump() {
    C.jit_dump_function((*C.FILE)(C.stdout), f.C, C.CString("main"))
}
