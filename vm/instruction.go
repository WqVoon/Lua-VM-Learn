package vm

// Bx 操作数所能表示的最大数，其整体范围为 [0, 262143]
const MAXARG_Bx = 1<<18 - 1

// sBx 操作数所能表示的最大数，其整体范围为 [-131071, 131072]
const MAXARG_sBx = MAXARG_Bx >> 1

type Instruction uint32

/*
取指令的低 6 位，也就是操作码部分
*/
func (self Instruction) Opcode() int {
	return int(self & 0x3F)
}

/*
用于从 iABC 模式指令中提取出三个参数，分别占 8，9，9 个比特
*/
func (self Instruction) ABC() (a, b, c int) {
	a = int(self >> 6 & 0xFF)
	c = int(self >> 14 & 0x1FF)
	b = int(self >> 23 & 0x1FF)
	return
}

/*
用于从 ABx 模式指令中提取出两个参数，分别占 8, 18 个比特
*/
func (self Instruction) ABx() (a, bx int) {
	a = int(self >> 6 & 0xFF)
	bx = int(self >> 14)
	return
}

/*
用于从 iAsBx 模式指令中提取出两个参数，分别占 8，18 个比特
*/
func (self Instruction) AsBx() (a, sbx int) {
	a, bx := self.ABx()
	return a, bx - MAXARG_sBx
}

/*
用于从 iAx 模式指令中提取出参数，占据 26 个比特
*/
func (self Instruction) Ax() int {
	return int(self >> 6)
}

/*
用于获得当前指令的名字
*/
func (self Instruction) OpName() string {
	return opcodes[self.Opcode()].name
}

/*
用于获得当前指令的模式
*/
func (self Instruction) OpMode() byte {
	return opcodes[self.Opcode()].opMode
}

/*
用于获得 B 操作数的模式
*/
func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argBMode
}

/*
用于获得 C 操作数的模式
*/
func (self Instruction) CMode() byte {
	return opcodes[self.Opcode()].argCMode
}
