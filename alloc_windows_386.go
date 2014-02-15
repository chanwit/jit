package jit

import "syscall"

const (
	MEM_COMMIT   = 0x1000
	MEM_RESERVE  = 0x2000
	MEM_DECOMMIT = 0x4000
	MEM_RELEASE  = 0x8000

	PAGE_EXECUTE_READWRITE = 0x40
)

var (
	kernel32     = syscall.MustLoadDLL("kernel32.dll")
	VirtualAlloc = kernel32.MustFindProc("VirtualAlloc")
	VirtualFree  = kernel32.MustFindProc("VirtualFree")
)

func MallocExec(size uint) uintptr {
	addr, _, _ := VirtualAlloc.Call(0, uintptr(size), MEM_RESERVE|MEM_COMMIT, PAGE_EXECUTE_READWRITE)
	return addr
}

func FreeExec(ptr uintptr, size uint) {
	if ptr == 0 {
		VirtualFree.Call(ptr, 0, MEM_RELEASE)
	}
}

func FlushExec(ptr uintptr, size uint) {
	// TODO must be writing in Assembly
}
