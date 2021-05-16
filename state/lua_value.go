package state

import . "lua-vm/api"

/*
定义 Lua 中的数据类型，目前包括如下映射：
	Lua类型 Go类型
	nil     nil
	boolean bool
	integer int64
	float   float64
	string  string
*/
type luaValue interface{}

/*
返回值的类型对应的常量，常量定义在 api/consts.go 中
*/
func typeOf(val luaValue) LuaType {
	switch val.(type) {
	case nil:
		return LUA_TNIL
	case bool:
		return LUA_TBOOLEAN
	case int64:
		return LUA_TNUMBER
	case float64:
		return LUA_TNUMBER
	case string:
		return LUA_TSTRING
	default:
		panic("Todo")
	}
}
