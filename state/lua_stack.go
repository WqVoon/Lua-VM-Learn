package state

/*
用于创建指定容量的栈
*/
func newLuaStack(size int) *luaStack {
	return &luaStack{
		slots: make([]luaValue, size),
		top:   0,
	}
}

/*
LuaStack 结构体，用于存放值，下标正向从 1 开始，反向从 -1 开始；
对于 Lua 来说，top 在栈空时无意义，栈内有元素后指向当前的栈顶元素处；
对于 Golang 来说，top 始终指向栈顶元素的下一个值
*/
type luaStack struct {
	slots []luaValue
	top   int
}

/*
检查当前的 LuaStack 是否还可以容纳 n 个值；
如果不能，那么为其扩容至可以为止
*/
func (self *luaStack) check(n int) {
	free := len(self.slots) - self.top
	for i := free; i < n; i++ {
		self.slots = append(self.slots, nil)
	}
}

/*
向 LuaStack 中压入一个值
*/
func (self *luaStack) push(val luaValue) {
	if self.top == len(self.slots) {
		panic("Stack Overflow")
	}
	self.slots[self.top] = val
	self.top++
}

/*
将 LuaStack 中的值弹出
*/
func (self *luaStack) pop() luaValue {
	if self.top < 1 {
		panic("Stack Underflow")
	}
	self.top--
	val := self.slots[self.top]
	self.slots[self.top] = nil
	return val
}

/*
把索引转换成绝对索引（在 Lua 视角下）
*/
func (self *luaStack) absIndex(idx int) int {
	if idx >= 0 {
		return idx
	}
	return idx + self.top + 1
}

/*
检查绝对索引是否有效（在 Lua 视角下）
*/
func (self *luaStack) isAbsValid(absIdx int) bool {
	return absIdx > 0 && absIdx <= self.top
}

/*
根据索引从栈中取值，如果索引无效那么返回 nil
*/
func (self *luaStack) get(idx int) luaValue {
	absIdx := self.absIndex(idx)
	if self.isAbsValid(absIdx) {
		return self.slots[absIdx-1]
	}
	return nil
}

/*
根据索引向栈中压值，如果索引无效则直接 panic
*/
func (self *luaStack) set(idx int, val luaValue) {
	absIndex := self.absIndex(idx)
	if self.isAbsValid(absIndex) {
		self.slots[absIndex-1] = val
		return
	}
	panic("Invalid Index")
}
