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
	source := self.readString()
	// 只有最顶层的 Prototype 才会获得 Source
	// 子 Prototype 可以继承父 Prototype 的值
	if source == "" {
		source = parentSource
	}
	return &Prototype{
		Source:          source,
		LineDefined:     self.readUint32(),
		LastLineDefined: self.readUint32(),
		NumParams:       self.readByte(),
		IsVararg:        self.readByte(),
		MaxStackSize:    self.readByte(),
		Code:            self.readCode(),
		Constants:       self.readConstants(),
		Upvalues:        self.readUpvalues(),
		Protos:          self.readProtos(source),
		LineInfo:        self.readLineInfo(),
		LocVars:         self.readLocVars(),
		UpvalueNames:    self.readUpvalueNames(),
	}
}

/*
读取所有的指令
*/
func (self *reader) readCode() []uint32 {
	code := make([]uint32, self.readUint32())
	for i := range code {
		code[i] = self.readUint32()
	}
	return code
}

/*
读取一个常量
*/
func (self *reader) readConstant() interface{} {
	tag := self.readByte()
	switch tag {
	case TAG_NIL:
		return nil
	case TAG_BOOLEAN:
		return self.readByte()
	case TAG_NUMBER:
		return self.readLuaNumber()
	case TAG_INTEGER:
		return self.readLuaInteger()
	case TAG_SHORT_STR:
		return self.readString()
	case TAG_LONG_STR:
		return self.readString()
	default:
		panic("Tag Error")
	}
}

/*
读取所有的常量
*/
func (self *reader) readConstants() []interface{} {
	constants := make([]interface{}, self.readUint32())
	for i := range constants {
		constants[i] = self.readConstant()
	}
	return constants
}

/*
读取所有的 Upvalue
*/
func (self *reader) readUpvalues() []Upvalue {
	upvalues := make([]Upvalue, self.readUint32())
	for i := range upvalues {
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

/*
读取所有的子 Prototype
*/
func (self *reader) readProtos(source string) []*Prototype {
	protos := make([]*Prototype, self.readUint32())
	for i := range protos {
		protos[i] = self.readProto(source)
	}
	return protos
}

/*
读取行号表
*/
func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32, self.readUint32())
	for i := range lineInfo {
		lineInfo[i] = self.readUint32()
	}
	return lineInfo
}

/*
读取局部变量表
*/
func (self *reader) readLocVars() []LocVar {
	locVars := make([]LocVar, self.readUint32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPc: self.readUint32(),
			EndPc:   self.readUint32(),
		}
	}
	return locVars
}

/*
读取 Upvalue 名表
*/
func (self *reader) readUpvalueNames() []string {
	upValueNames := make([]string, self.readUint32())
	for i := range upValueNames {
		upValueNames[i] = self.readString()
	}
	return upValueNames
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
