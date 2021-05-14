package vm

/*
四种指令类型
*/
const (
	/*
		可以携带 3 个操作数，分别用 8、9、9 个比特表示，共计 39 条指令
	*/
	IABC = iota
	/*
		携带 A 和 Bx 两个操作数，分别占用 8 和 18 个比特，共计 3 条指令
	*/
	IABx
	/*
		携带 A 和 sBx 两个操作数，分别占用 8 和 18 个比特，共计 4 条指令
		**只有这种形式的指令中的 sBx 操作数会被解释成有符号整数**
	*/
	IAsBx
	/*
		携带一个操作数，占用 26 个比特，只有 1 条“指令”
		TODO：实际并不是真正的指令，只用来扩展其他指令操作数？
	*/
	IAx
)

/*
指令码
	所有指令均占用 32 个比特，其中低 6 比特用于操作码，高 26 比特用于操作数
	6 比特本身可表示 2^6=64 条指令，Lua 中共有 47 条指令，操作码从 0～46，详见下面的常量
*/
const (
	OP_MOVE = iota
	OP_LOADK
	OP_LOADKX
	OP_LOADBOOL
	OP_LOADNIL
	OP_GETUPVAL
	OP_GETTABUP
	OP_GETTABLE
	OP_SETTABUP
	OP_SETUPVAL
	OP_SETTABLE
	OP_NEWTABLE
	OP_SELF
	OP_ADD
	OP_SUB
	OP_MUL
	OP_MOD
	OP_POW
	OP_DIV
	OP_IDIV
	OP_BAND
	OP_BOR
	OP_BXOR
	OP_SHL
	OP_SHR
	OP_UNM
	OP_BNOT
	OP_NOT
	OP_LEN
	OP_CONCAT
	OP_JMP
	OP_EQ
	OP_LT
	OP_LE
	OP_TEST
	OP_TESTSET
	OP_CALL
	OP_TAILCALL
	OP_RETURN
	OP_FORLOOP
	OP_FORPREP
	OP_TFORCALL
	OP_TFORLOOP
	OP_SETLIST
	OP_CLOSURE
	OP_VARARG
	OP_EXTRAARG
)

/*
四种操作数类型
*/
const (
	/*
		操作数未被使用，如 MOVE 指令中的 C 操作数
	*/
	OpArgN = iota
	/*
		操作数已被使用，可能为布尔值、整数值、upvalue 索引、子函数索引等
	*/
	OpArgU
	/*
		操作数是寄存器（在 iABC 模式下）或一个跳转偏移（在 iAsBx 模式下）
	*/
	OpArgR
	/*
		操作数是常量或寄存器
			在 iABx 模式下表示常量表索引，比如 LOADK 指令
			在 iABC 模式下该操作数类型只能使用 9 个比特中的低 8 位，其中最高位为 1 表示常量表索引，否则为寄存器索引

		TODO：这样表示常量索引时应该只有 2^7=128 个值，但是实际测试时发现 Lua 对常量的数量似乎没有限制？
	*/
	OpArgK
)
