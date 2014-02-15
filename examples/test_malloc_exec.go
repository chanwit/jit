package main

import (
	"../../jit"
	"fmt"
	"log"
	"syscall"
	"unsafe"
)

const SIZE = 64 * 1024

type FunctionBuilder struct {
	size uint32
	pc   uintptr
	buf  uintptr
}

type X86Reg byte

const (
	X86_EAX = X86Reg(iota)
	X86_ECX
	X86_EDX
	X86_EBX
	X86_ESP
	X86_EBP
	X86_ESI
	X86_EDI
	X86_NREG
)

func (fb *FunctionBuilder) Nop() {
	b := (*[SIZE]byte)(unsafe.Pointer(fb.buf))
	b[fb.pc] = 0x90
	fb.pc++
}

func (fb *FunctionBuilder) PushReg(reg X86Reg) {
	b := (*[1]byte)(unsafe.Pointer(fb.buf + fb.pc))
	b[0] = byte(0x50) + byte(reg)
	fb.pc++
}

func (fb *FunctionBuilder) PopReg(reg X86Reg) {
	b := (*[1]byte)(unsafe.Pointer(fb.buf + fb.pc))
	b[0] = byte(0x58) + byte(reg)
	fb.pc++
}

func (fb *FunctionBuilder) MovRegReg(dreg X86Reg, reg X86Reg, size byte) {
	b := (*[1]byte)(unsafe.Pointer(fb.buf + fb.pc))
	switch size {
	case 1:
		b[0] = byte(0x8a)
	case 2:
		b[0] = byte(0x66)
		fallthrough
	case 4:
		b[0] = byte(0x8b)
	default:
		panic("JIT Assert")
	}
	fb.pc++
	fb.RegEmit(dreg, reg)
}

func (fb *FunctionBuilder) RegEmit(dreg X86Reg, reg X86Reg) {
	fb.AddressByte(3, dreg, reg)
}

func (fb *FunctionBuilder) AddressByte(m byte, o X86Reg, r X86Reg) {
	b := (*[1]byte)(unsafe.Pointer(fb.buf + fb.pc))
	b[0] = ((((m) & 0x03) << 6) | ((byte(o) & 0x07) << 3) | (byte(r) & 0x07))
	fb.pc++
}

func (fb *FunctionBuilder) Mov_EBP_Param_EAX(i int) {
	b := (*[SIZE]byte)(unsafe.Pointer(fb.buf))
	b[fb.pc] = 0x8b
	fb.pc++
	b[fb.pc] = 0x45
	fb.pc++
	b[fb.pc] = byte((i + 2) * 4)
	fb.pc++
}

func (fb *FunctionBuilder) IMul_EBP_Param_EAX(i int) {
	b := (*[SIZE]byte)(unsafe.Pointer(fb.buf))
	b[fb.pc] = 0x0f
	fb.pc++
	b[fb.pc] = 0xaf
	fb.pc++
	b[fb.pc] = 0x45
	fb.pc++
	b[fb.pc] = byte((i + 2) * 4)
	fb.pc++
}

func (fb *FunctionBuilder) Add_EBP_Param_EAX(i int) {
	b := (*[SIZE]byte)(unsafe.Pointer(fb.buf))
	b[fb.pc] = 0x03
	fb.pc++
	b[fb.pc] = 0x45
	fb.pc++
	b[fb.pc] = byte((i + 2) * 4)
	fb.pc++
}

func (fb *FunctionBuilder) Ret() {
	b := (*[SIZE]byte)(unsafe.Pointer(fb.buf))
	b[fb.pc] = 0xc3
	fb.pc++
}

func mkprog() error {

	addr := jit.MallocExec(SIZE)
	if addr == 0 {
		return nil
	}

	fb := &FunctionBuilder{SIZE, 0, addr}

	fb.PushReg(X86_EBP)
	fb.MovRegReg(X86_EBP, X86_ESP, 4)
	fb.Mov_EBP_Param_EAX(0)
	fb.IMul_EBP_Param_EAX(1)
	fb.Add_EBP_Param_EAX(2)
	fb.PopReg(X86_EBP)
	fb.Ret()
	fb.Nop()

	r1, r2, e := syscall.Syscall(addr, 3, 3, 5, 2)

	fmt.Printf("%d\n", r1)
	fmt.Printf("%d\n", r2)
	fmt.Printf("%s\n", e)
	return nil
}

func main() {
	err := mkprog()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("HELLO\n")
}
