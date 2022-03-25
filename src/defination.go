// 一些类型定义
package src

// 变量由行号和名称组成
type Variable struct {
	LineNum int
	Name string
}

// 赋值语句，由行号，变量和字符串字面量组成
// eg  &a = "hello"
type Assignment struct {
	LineNum int
	Variable *Variable
	String string
}

// print语句由行号和变量组成
// eg print(a)
type Print struct {
	LineNum int
	Variable *Variable
}

// 声明语句定义为空接口，可以是赋值语句或者print语句
type Statement interface {}

var _ Statement = (*Assignment)(nil)
var _ Statement = (*Print)(nil)

// 源码由行号和一些语句构成
type SourceCode struct {
	LineNum int
	Statements []Statement
}

