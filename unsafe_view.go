package echo

import "unsafe"

type stringView struct {
	p unsafe.Pointer
	n int
}

func str2raw(s string) (unsafe.Pointer, uint64) {
	v := (*stringView)(unsafe.Pointer(&s))
	return v.p, uint64(v.n)
}

func raw2str(p unsafe.Pointer, n uint64) string {
	v := stringView{p: p, n: int(n)}
	return *(*string)(unsafe.Pointer(&v))
}

type sliceView struct {
	p unsafe.Pointer
	n int
	m int
}

func bytes2raw(s []byte) (unsafe.Pointer, uint64) {
	v := (*sliceView)(unsafe.Pointer(&s))
	return v.p, uint64(v.n)
}

func raw2bytes(p unsafe.Pointer, n uint64) []byte {
	v := sliceView{p: p, n: int(n), m: int(n)}
	return *(*[]byte)(unsafe.Pointer(&v))
}

type anyView struct {
	ty unsafe.Pointer
	p  unsafe.Pointer
}

func any2raw(f any) (unsafe.Pointer, unsafe.Pointer) {
	v := *(*anyView)(unsafe.Pointer(&f))
	return v.ty, v.p
}

func raw2any(ty unsafe.Pointer, p unsafe.Pointer) any {
	v := anyView{ty: ty, p: p}
	return *(*any)(unsafe.Pointer(&v))
}

type ifaceView struct {
	itab unsafe.Pointer
	p    unsafe.Pointer
}

func iface2raw[T any](f T) (unsafe.Pointer, unsafe.Pointer) {
	v := *(*ifaceView)(unsafe.Pointer(&f))
	return v.itab, v.p
}

func raw2iface[T any](itab unsafe.Pointer, p unsafe.Pointer) T {
	v := ifaceView{itab: itab, p: p}
	return *(*T)(unsafe.Pointer(&v))
}
