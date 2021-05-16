package state

/*
获取 LuaStack 的栈顶值
*/
func (self *luaState) GetTop() int {
	return self.stack.top
}

/*
将相对索引转换成绝对索引
*/
func (self *luaState) AbsIndex(idx int) int {
	return self.stack.absIndex(idx)
}

/*
检查当前 LuaStack 是否可以容纳 n 个值，如果不能那么进行扩容
*/
func (self *luaState) CheckStack(n int) bool {
	self.stack.check(n)
	return true
}

/*
从 LuaStack 中弹出 n 个值
TODO：不需要接收弹出的值吗？
*/
func (self *luaState) Pop(n int) {
	for i := 0; i < n; i++ {
		self.stack.pop()
	}
}

/*
把值从 LuaStack 的 fromIdx 处复制一份到 toIdx 处
*/
func (self *luaState) Copy(fromIdx, toIdx int) {
	val := self.stack.get(fromIdx)
	self.stack.set(toIdx, val)
}

/*
把 LuaStack 中 idx 处的值推入栈顶
*/
func (self *luaState) PushValue(idx int) {
	val := self.stack.get(idx)
	self.stack.push(val)
}

/*
PushValue 的反操作，把栈顶的值弹出并写入 LuaStack 的 idx 处
*/
func (self *luaState) Replace(idx int) {
	val := self.stack.pop()
	self.stack.set(idx, val)
}

/*
将栈顶值弹出并插入指定位置，同时将其后的所有值上移
*/
func (self *luaState) Insert(idx int) {
	self.Rotate(idx, 1)
}

/*
移除指定位置上的值，并将其上的所有值下移
*/
func (self *luaState) Remove(idx int) {
	self.Rotate(idx, -1)
	self.stack.pop()
}

/*
将 [idx, top] 区间循环上移（如果 n 是正数）；
或者循环下移（如果 n 是负数）
*/
func (self *luaState) Rotate(idx, n int) {
	t := self.stack.top - 1
	p := self.stack.absIndex(idx) - 1

	var m int
	if n >= 0 {
		m = t - n
	} else {
		m = p - n - 1
	}

	self.stack.reverse(p, m)
	self.stack.reverse(m+1, t)
	self.stack.reverse(p, t)
}

/*
设置 LuaStack 的 top；
如果新 top 的值小于旧的，那么会弹出期间的差值；
否则会推入 nil
*/
func (self *luaState) SetTop(idx int) {
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("Stack Underflow")
	}

	n := self.stack.top - newTop
	if n > 0 {
		self.Pop(n)
	} else {
		for i := 0; i > n; i-- {
			self.stack.push(nil)
		}
	}
}
