package core

import (
	"bytes"
	"encoding/binary"
)

type ByteBuf struct {
	memory          *[]byte
	abs_read_index  int
	abs_write_index int
	capacity        int
	offset          int
}

func NewByteBuf(cap int) *ByteBuf {
	var data = make([]byte, cap)
	return &ByteBuf{memory: &data, capacity: cap}
}

func (buf *ByteBuf) HasReadableBytes() bool {
	return buf.abs_read_index < buf.abs_write_index
}

func (buf *ByteBuf) HasWritableBytes() bool {
	return buf.WritableBytes() > 0
}

func (buf *ByteBuf) WritableBytes() int {
	return buf.capacity - buf.WriteIndex()
}

func (buf *ByteBuf) ReadableBytes() int {
	return buf.abs_write_index - buf.abs_read_index
}

func (buf *ByteBuf) WriteIndex() int {
	return buf.abs_write_index - buf.offset
}

func (buf *ByteBuf) ReadIndex() int {
	return buf.abs_read_index - buf.offset
}

func (buf *ByteBuf) ix(index int) int {
	return buf.offset + index
}

func (buf *ByteBuf) GetMemory() *[]byte {
	return buf.memory
}

func (buf *ByteBuf) WriteBytes(data []byte) {
	var l = len(data)
	var fromIndex = buf.abs_write_index
	var toIndex = fromIndex + l
	copy((*buf.memory)[fromIndex:toIndex], data)
	buf.abs_write_index = toIndex
}

func (buf *ByteBuf) ReadShortLE() uint16 {
	var ret = binary.LittleEndian.Uint16((*buf.memory)[buf.abs_read_index : buf.abs_read_index+2])
	buf.abs_read_index += 2
	return ret
}

func (buf *ByteBuf) ReadBytes(dst []byte) {
	var l = len(dst)
	var fromIndex = buf.abs_read_index
	var toIndex = fromIndex + l
	copy(dst, (*buf.memory)[fromIndex:toIndex])
	buf.abs_read_index = toIndex
}

//func (buf *ByteBuf) SliceWrite() []byte {
//	return (*buf.memory)[buf.abs_write_index:buf.capacity]
//}

//func (buf *ByteBuf) SliceRead() []byte {
//	return (*buf.memory)[buf.abs_read_index:buf.abs_write_index]
//}

func (buf *ByteBuf) SkipWrite(n int) {
	buf.abs_write_index += n
}

func (buf *ByteBuf) ReadIntLE() int {
	var ret = binary.LittleEndian.Uint32((*buf.memory)[buf.abs_read_index : buf.abs_read_index+4])
	buf.abs_read_index += 4
	return int(ret)
}

func (buf *ByteBuf) SkipRead(i int) {
	buf.abs_read_index += i
}

func (buf *ByteBuf) Compact() {
	if buf.abs_read_index > 0 {
		var temp = (*buf.memory)[buf.abs_read_index:buf.abs_write_index]
		buf.Clear()
		buf.WriteBytes(temp)
	}
}

func (buf *ByteBuf) Clear() *ByteBuf {
	buf.abs_read_index = 0
	buf.abs_write_index = 0
	return buf
}

func (buf *ByteBuf) WriteIntLe(v uint32) {
	binary.LittleEndian.PutUint32((*buf.memory)[buf.abs_write_index:], v)
	buf.abs_write_index += 4
}

func (buf *ByteBuf) ReadLongLE() uint64 {
	var fromIndex = buf.abs_read_index
	var toIndex = fromIndex + 8
	buf.abs_read_index = toIndex
	return binary.LittleEndian.Uint64((*buf.memory)[fromIndex:toIndex])
}

func (buf *ByteBuf) WriteByte(b byte) {
	(*buf.memory)[buf.abs_write_index] = b
	buf.abs_write_index++
}

func (buf *ByteBuf) GetByteAbs(i int) byte {
	return (*buf.memory)[i]
}

func (buf *ByteBuf) LastIndexOf(b byte) int {
	return bytes.LastIndexByte((*buf.memory)[buf.abs_read_index:buf.abs_write_index], b)
}

func (buf *ByteBuf) IndexOf(b byte, from int, to int) int {
	return bytes.IndexByte((*buf.memory)[from:to], b)
}
