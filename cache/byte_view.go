package cache

// ByteView 只读数据
type ByteView struct {
	b []byte
}

func (b ByteView) Len() uint {
	return uint(len(b.b))
}

func (b ByteView) String() string {
	return string(b.b)
}

func cloneBytes(val []byte) (b []byte) {
	b = make([]byte, len(val))
	copy(b, val)
	return
}

// ByteSlice 副本，提供外界访问
func (b ByteView) ByteSlice() []byte {
	return cloneBytes(b.b)
}
