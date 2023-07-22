package bytebuf

import (
	"encoding/binary"
	"math"
)

type ByteBuffer struct {
	B []byte
}

func (b *ByteBuffer) Len() int {
	return len(b.B)
}

func (b *ByteBuffer) Bytes() []byte {
	return b.B
}

func (b *ByteBuffer) SetBytes(data []byte) {
	b.B = data
}

func (b *ByteBuffer) WriteInt64(num int64) {
	buf := make([]byte, binary.MaxVarintLen64)
	n := binary.PutVarint(buf, num)
	bytes := buf[:n]

	b.B = append(b.B, bytes...)
}

func (b *ByteBuffer) ReadInt64() int64 {
	num, len := binary.Varint(b.B)
	b.B = b.B[len:]
	return num
}

func (b *ByteBuffer) WriteFloat64(num float64) {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, math.Float64bits(num))
	b.B = append(b.B, buf...)
}

func (b *ByteBuffer) ReadFloat64() float64 {
	num := binary.BigEndian.Uint64(b.B)
	b.B = b.B[8:]
	return math.Float64frombits(num)
}

func (b *ByteBuffer) Write(p []byte) {
	b.WriteInt64(int64(len(p)))
	b.B = append(b.B, p...)
}

func (b *ByteBuffer) Read() []byte {
	size := b.ReadInt64()
	bytes := b.B[:size]
	b.B = b.B[size:]
	return bytes
}

func (b *ByteBuffer) WriteString(s string) {
	b.Write([]byte(s))
}

func (b *ByteBuffer) ReadString() string {
	return string(b.Read())
}

func (b *ByteBuffer) WriteBool(value bool) {
	if value {
		b.Write([]byte{1})
	} else {
		b.Write([]byte{0})
	}
}

func (b *ByteBuffer) ReadBool() bool {
	return b.Read()[0] == 1
}

func (b *ByteBuffer) Reset() {
	b.B = b.B[:0]
}
