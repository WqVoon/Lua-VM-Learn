package main

import (
	"fmt"
	"io/ioutil"
	"lua-vm/binchunk"
	"os"
)

func main() {
	var filename string
	if len(os.Args) != 1 {
		filename = os.Args[1]
	} else {
		filename = "a.out"
	}

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	list(binchunk.Undump(data))
}

/*
递归输出 Prototype 结构体中的内容
*/
func list(f *binchunk.Prototype) {
	printHeader(f)
	printCode(f)
	printDetail(f)
	for _, p := range f.Protos {
		list(p)
	}
}

/*
打印 Prototype 的头部信息，格式基本等同于 `luac -l` 的反编译输出
*/
func printHeader(f *binchunk.Prototype) {
	funcType := "main"
	if f.LineDefined > 0 {
		funcType = "function"
	}

	varargFlag := ""
	if f.IsVararg == 1 {
		varargFlag = "+"
	}

	fmt.Printf("\n%s <%s: %d, %d> (%d instructions)\n",
		funcType, f.Source, f.LineDefined, f.LastLineDefined, len(f.Code))

	fmt.Printf("%d%s params, %d slots, %d upvalues, ",
		f.NumParams, varargFlag, f.MaxStackSize, len(f.Upvalues))

	fmt.Printf("%d locals, %d constants, %d functions\n",
		len(f.LocVars), len(f.Constants), len(f.Protos))
}

/*
打印代码段信息，当前仅输出其序号，行号和十六进制表示
*/
func printCode(f *binchunk.Prototype) {
	for pc, c := range f.Code {
		//TODO：尚不清楚这样写的原因
		line := "-"
		if len(f.LineInfo) > 0 {
			line = fmt.Sprintf("%d", f.LineInfo[pc])
		}
		fmt.Printf("\t%d\t[%s]\t0x%08X\n", pc+1, line, c)
	}
}

/*
打印常量表，局部变量表以及 Upvalue 表中的内容
*/
func printDetail(f *binchunk.Prototype) {
	fmt.Printf("constants (%d):\n", len(f.Constants))
	for i, k := range f.Constants {
		fmt.Printf("\t%d\t%s\n", i+1, constantToString(k))
	}

	fmt.Printf("locals (%d):\n", len(f.LocVars))
	for i, locVar := range f.LocVars {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, locVar.VarName, locVar.StartPc, locVar.EndPc)
	}

	fmt.Printf("upvalues (%d):\n", len(f.Upvalues))
	for i, upval := range f.Upvalues {
		fmt.Printf("\t%d\t%s\t%d\t%d\n", i, f.UpvalueNames[i], upval.Instack, upval.Idx)
	}
}

/*
类型断言 Constant 表中的内容，并返回对应的字符串形式
*/
func constantToString(k interface{}) string {
	switch k.(type) {
	case nil:
		return "nil"
	case bool:
		return fmt.Sprintf("%t", k)
	case float64:
		return fmt.Sprintf("%g", k)
	case int64:
		return fmt.Sprintf("%d", k)
	case string:
		return fmt.Sprintf("%q", k)
	default:
		return "?"
	}
}
