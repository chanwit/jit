package jit
/*
#cgo LDFLAGS: -lm -ljit
#include <stdio.h>
#include <jit/jit.h>
#include <jit/jit-dump.h>

extern int on_demand_compile(jit_function_t);

static void SetOnDemandCompileFunction(jit_function_t f) {
    jit_function_set_on_demand_compiler(f, on_demand_compile);
}

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

type Label struct {
    C C.jit_label_t
}

var intType = Type{ C.jit_type_int }
func Int() Type {
    return intType
}

var uintType = Type{ C.jit_type_uint }
func UInt() Type {
    return uintType
}

func NewContext() *Context {
    return &Context{ C.jit_context_create() }
}

func NewSignature(ret Type, params []Type) *Signature {
    signature := C.jit_type_create_signature(C.jit_abi_cdecl,
        ret.C,
        (*C.jit_type_t)(unsafe.Pointer(&params[0])),
        C.uint(len(params)), 1)
    return &Signature{signature}
}

func NewLabel() *Label {
    return &Label{ C.jit_label_undefined }
}

// ========== Context =============

func (c *Context) NewFunction(ret Type, params []Type) *Function {
    signature := NewSignature(ret, params)
    function := C.jit_function_create(c.C, signature.C)
    C.jit_type_free(signature.C)
    return &Function{C: function, retType: ret, paramType: params}
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

// ========== Context =============


// ========== Function =============

func (f *Function) Param(i int) *Value {
    return &Value{ C.jit_value_get_param(f.C, C.uint(i)) }
}

func (f *Function) Param2() (*Value, *Value) {
    return &Value{ C.jit_value_get_param(f.C, C.uint(0)) },
           &Value{ C.jit_value_get_param(f.C, C.uint(1)) }
}

func (f *Function) Param3() (*Value, *Value, *Value) {
    return &Value{ C.jit_value_get_param(f.C, C.uint(0)) },
           &Value{ C.jit_value_get_param(f.C, C.uint(1)) },
           &Value{ C.jit_value_get_param(f.C, C.uint(2)) }
}

func (f *Function) Mul(a, b *Value) *Value {
    return &Value{ C.jit_insn_mul(f.C, a.C, b.C) }
}

func (f *Function) Add(a, b *Value) *Value {
    return &Value{ C.jit_insn_add(f.C, a.C, b.C) }
}

func (f *Function) Sub(a, b *Value) *Value {
    return &Value{ C.jit_insn_sub(f.C, a.C, b.C) }
}

func (f *Function) Return(ret *Value) {
    C.jit_insn_return(f.C, ret.C)
}

func (f *Function) Eq(a, b *Value) *Value {
    return &Value{ C.jit_insn_eq(f.C, a.C, b.C) }
}

func (f *Function) BranchIfNot(v *Value, label *Label) {
    C.jit_insn_branch_if_not(f.C, v.C, (*C.jit_label_t)(unsafe.Pointer(label)))
}

func (f *Function) Label(label *Label) {
    C.jit_insn_label(f.C, (*C.jit_label_t)(unsafe.Pointer(label)))
}

func (f *Function) LessThan(a, b *Value) *Value {
    return &Value{ C.jit_insn_lt(f.C, a.C, b.C) }
}

func (f *Function) Call(target *Function, values []*Value) *Value {
    args := make([]C.jit_value_t, len(values))
    for i := range values {
        args[i] = values[i].C
    }
    return &Value{ C.jit_insn_call(f.C,
                    C.CString("noname"),
                    target.C, nil, (*C.jit_value_t)(&args[0]), C.uint(len(args)), C.int(0)) }
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

            case C.jit_type_uint:
                x := C.jit_uint(values[i].(uint))
                args[i] = (unsafe.Pointer)(&x)
        }
    }

    switch f.retType.C {
        case C.jit_type_int:
            result := C.jit_int(0)
            C.jit_function_apply(f.C, (*unsafe.Pointer)(&args[0]), unsafe.Pointer(&result))
            return int(result)

        case C.jit_type_uint:
            result := C.jit_uint(0)
            C.jit_function_apply(f.C, (*unsafe.Pointer)(&args[0]), unsafe.Pointer(&result))
            return uint(result)
    }

    return nil
}

func (f *Function) Dump(name string) {
    C.jit_dump_function((*C.FILE)(C.stdout), f.C, C.CString(name))
}

func (f *Function) SetRecompilable() {
    C.jit_function_set_recompilable(f.C)
}

type compileFunction struct {
    F           *Function
    compileFunc func(*Function)bool
}

var registry = make(map[C.jit_function_t]*compileFunction)

func (f *Function) SetOnDemandCompiler(function func(f *Function)bool) {
    registry[f.C] = &compileFunction{f, function}
    C.SetOnDemandCompileFunction(f.C)
}

func (f *Function) GetOnDemandCompiler() func(*Function)bool {
    return registry[f.C].compileFunc
}

// ========== Function =============

//export on_demand_compile
func on_demand_compile(f C.jit_function_t) C.int {
    cf := registry[f]
    result := cf.compileFunc(cf.F)
    if(result) {
        return C.int(1)
    }
    return C.int(0)
}
