package state

/*
顾名思义
*/
func (self *luaState) PushNil() {
	self.stack.push(nil)
}

/*
顾名思义
*/
func (self *luaState) PushBoolean(b bool) {
	self.stack.push(b)
}

/*
顾名思义
*/
func (self *luaState) PushInteger(n int64) {
	self.stack.push(n)
}

/*
顾名思义
*/
func (self *luaState) PushNumber(n float64) {
	self.stack.push(n)
}

/*
顾名思义
*/
func (self *luaState) PushString(s string) {
	self.stack.push(s)
}
