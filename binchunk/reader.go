package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	// 内部记录将要被解析的字节流
	data []byte
}

/*
从当前数据中读取一个 byte 出来
*/
func (self *reader) readByte() byte {
	b := self.data[0]
	self.data = self.data[1:]
	return b
}

/*
从当前数据中读取 n 个 byte 出来
*/
func (self *reader) readBytes(n uint) []byte {
	bytes := self.data[:n]
	self.data = self.data[n:]
	return bytes
}

/*
从当前数据中读取一个 cint 出来
*/
func (self *reader) readUint32() uint32 {
	i := binary.LittleEndian.Uint32(self.data)
	self.data = self.data[4:]
	return i
}

/*
从当前数据中读取一个 size_t 出来
*/
func (self *reader) readUint64() uint64 {
	i := binary.LittleEndian.Uint64(self.data)
	self.data = self.data[8:]
	return i
}

/*
从当前数据中读取一个 Lua 整数出来
*/
func (self *reader) readLuaInteger() int64 {
	return int64(self.readUint64())
}

/*
从当前数据中读取一个 Lua 浮点数出来
*/
func (self *reader) readLuaNumber() float64 {
	return math.Float64frombits(self.readUint64())
}

/*
从当前数据中读取一个 string 出来
 BinChunk 中的字符串分为空字符串，长字符串和短字符串三种：
  对于 NULL 字符串，用 0x00 表示
  对于长度小于等于 253 的字符串，先用一个字节记录长度+1，后面跟着字节数组
  对于长度大于等于 254 的字符串，第一个字节是 0xFF，后面跟着 size_t 来记录长度+1，再跟着字节数组
*/
func (self *reader) readString() string {
	size := uint(self.readByte())
	if size == 0 {
		return ""
	}
	if size == 0xFF {
		size = uint(self.readUint64())
	}
	bytes := self.readBytes(size - 1)
	return string(bytes)
}
