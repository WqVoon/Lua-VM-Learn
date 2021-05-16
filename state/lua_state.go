package state

var STACK_SIZE = 20

/*
创建一个具有 STASK_SIZE 容量栈的 LuaState
*/
func New() *luaState {
	return &luaState{
		stack: newLuaStack(STACK_SIZE),
	}
}

/*
LuaState 结构体，用于描述 Lua 解释器的状态；
当前内部仅有一个 LuaStack
*/
type luaState struct {
	stack *luaStack
}
