package state

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

func typeOf(val luaValue) {}
