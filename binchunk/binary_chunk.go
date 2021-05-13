package binchunk

/*
二进制 Chunk 结构体
*/
type binaryChunk struct {
	// 详见 header 结构体
	header
	// TODO：主函数的 Upvalue 数量？
	sizeUpvalues byte
	// 主函数的原型
	mainFunc *Prototype
}

/*
文件头中的信息，共占 17 + luaIntegerSize + luaNumberSize 个字节
在我的平台下共占 33 个字节，在 `hexdump -C` 的输出中一直到 0x21 处
*/
type header struct {
	// 签名（魔数），四个字节分别为 ESC, L, u, a，即 "\x1bLua"
	signature [4]byte
	// 版本号，Lua 版本 x.y.z 分别表示大版本，小版本和发布号，此处的值等于 大版本x16 + 小版本
	version byte
	// 格式号，官方实现的虚拟机其值为 0
	format byte
	// 被称为 LUAC_DATA，其中前两个字节为 0x1993，表示 lua1.0 发布的年份
	// 后四个字节依次为 0x0d, 0x0a, 0x1a, 0x0a，因此写成字面量为 "\x19\x93\r\n\x1a\n"
	luacData [6]byte
	// 后五个字节分别表示 cint, size_t, Lua 虚拟机指令、Lua 整数和 Lua 浮点数所占的字节数
	cintSize        byte
	sizeSize        byte
	instructionSize byte
	luaIntegerSize  byte
	luaNumberSize   byte
	// 用来存放 Lua 整数值的 0x5678，其所占大小实际等于上面的 luaIntegerSize
	// 主要用来检测大小端
	luacInt int64
	// 用来存放 Lua 浮点数 370.5，其所占大小实际等于上面的 luaNumberSize
	// 主要用来检测浮点数格式
	luacNum float64
}

/*
函数原型结构体
*/
type Prototype struct {
	// 函数来源，如果以 @ 开头表示从后面紧随的 lua 源文件编译而来
	// 如果以 = 开头则有特殊含义，比如通过 `luac -` 命令编译的内容此处为 "=stdin"
	// TODO：如果什么都没有，则说明该二进制 chunk 是从程序提供的字符串编译而来的？
	// 该信息不是运行必须的，因此如果使用 `luac -s` 会被剔除掉
	Source string
	// 起始行号，如果是 main 函数则为 0，和下面的结束行号均为 cint 型
	LineDefined uint32
	// 结束行号，如果是 main 函数则为 0
	LastLineDefined uint32
	// 固定参数的个数，如果函数是变长参数（如 main 函数），那么此值为 0
	NumParams byte
	// 是否为 Vararg 函数
	IsVararg byte
	// 寄存器数量，表示该函数所需的最大寄存器数量
	// TODO：由于实际采用的是个栈结构，因此叫做 MaxStackSize
	MaxStackSize byte
	// 指令表，Lua 虚拟机的每条指令占 4 个字节，表开头有一个 cint 表示表大小
	Code []uint32
	// 常量表，用于存放 Lua 代码中出现的常量
	// 其中每个常量都以一个 tag 开头，tag 的意义如下：
	//  0x00 nil 不存储
	//  0x01 boolean 存储 0 或 1
	//  0x03 number Lua 浮点数
	//  0x13 integer Lua 整数
	//  0x04 string 短字符串
	//  0x14 string 长字符串
	// 表开头有一个 cint 表示表大小
	Constants []interface{}
	// Upvalue 表，表开头有一个 cint 表示表大小
	Upvalues []Upvalue
	// 子函数原型表，表开头有一个 cint 表示表大小
	Protos []*Prototype
	// 行号表，与指令表中的内容一一对应，每一项用 cint 表示，表开头有一个 cint 表示表大小
	LineInfo []uint32
	// 局部变量表，用于记录局部变量名，表开头有一个 cint 表示表大小
	LocVars []LocVar
	// Upvalue 名表，与 Upvalue 表中的每一项对应，用于记录每一个 Upvalue 的名字
	// 表开头有一个 cint 表示表大小
	// TODO：有一个名为 `_ENV` 的 Upvalue 尚不知含义
	UpvalueNames []string
}

/*
Upvalue 结构体
TODO：当前仅知道该元素占用 2 个字节，具体内容待研究
*/
type Upvalue struct {
	Instack byte
	Idx     byte
}

/*
局部变量结构体
*/
type LocVar struct {
	// 变量名
	VarName string
	// 变量起始指令索引
	StartPc uint32
	// 变量结束指令索引
	EndPc uint32
}

const (
	// 此处的常量用于定义我的平台下 5.3.5 版本的 Lua 虚拟机头部信息
	LUA_SIGNATURE    = "\x1bLua"
	LUAC_VERSION     = 0x53
	LUAC_FORMAT      = 0
	LUAC_DATA        = "\x19\x93\r\n\x1a\n"
	CINT_SIZE        = 4
	CSIZET_SIZE      = 8
	INSTRUCTION_SIZE = 4
	LUA_INTEGER_SIZE = 8
	LUA_NUMBER_SIZE  = 8
	LUAC_INT         = 0x5678
	LUAC_NUM         = 370.5
)

const (
	// 此处的常量用于定义常量表中的 tag
	TAG_NIL       = 0x00
	TAG_BOOLEAN   = 0x01
	TAG_NUMBER    = 0x03
	TAG_INTEGER   = 0x13
	TAG_SHORT_STR = 0x04
	TAG_LONG_STR  = 0x14
)

/*
用来从 BinChunk 中读取数据并返回主函数的 Prototype
*/
func Undump(data []byte) *Prototype {
	r := reader{data}
	r.checkHeader()
	// 跳过主函数的 Upvalue 数量，因为这个值从 Prototype 中也可以拿到
	r.readByte()
	return r.readProto("")
}
