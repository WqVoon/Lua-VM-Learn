package binchunk

import (
	"encoding/binary"
	"math"
)

/*
用于读取并分析 BinChunk 中的字节流
*/
type reader struct {
	// 内部记录将要被解析的字节流
	data []byte
}

/*
检查头部是否和所需的数据一致，不一致则直接退出函数；
由于头部对于运行时没有什么作用，所以这里读取并判断后直接丢弃
*/
func (self *reader) checkHeader() {
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("Signature Error")
	} else if self.readByte() != LUAC_VERSION {
		panic("Luac Version Error")
	} else if self.readByte() != LUAC_FORMAT {
		panic("Luac Format Error")
	} else if string(self.readBytes(6)) != LUAC_DATA {
		panic("Luac Data Error")
	} else if self.readByte() != CINT_SIZE {
		panic("Cint Size Error")
	} else if self.readByte() != CSIZET_SIZE {
		panic("Size_t Size Error")
	} else if self.readByte() != INSTRUCTION_SIZE {
		panic("Instruction Size Error")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("Lua Integer Size Error")
	} else if self.readByte() != LUA_NUMBER_SIZE {
		panic("Lua Number Size Error")
	} else if self.readLuaInteger() != LUAC_INT {
		panic("Lua Integer Format Error")
	} else if self.readLuaNumber() != LUAC_NUM {
		panic("Lua Number Format Error")
	}
}

/*
递归读取函数 Prototype 并返回主函数
*/
func (self *reader) readProto(parentSource string) *Prototype {
	return &Prototype{}
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
