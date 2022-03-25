// 编译器后端，直接解析抽象语法树，并执行
package src

import (
	"fmt"
	"github.com/pkg/errors"
)

// 变量列表列表, key变量名为string，val字符串字面量
type Vars map[string]string

func NewVars() Vars{
	return make(Vars)
}

func resolveAssignment(vars Vars, v *Assignment) error {
	vName := v.Variable.Name
	if vName == "" {
		return errors.New("variable name is empty")
	}
	vars[vName] = v.String
	return nil
}

func resolvePrint(vars Vars, p *Print) error {
	vName := p.Variable.Name
	if vName == "" {
		return errors.New("param of print is empty")
	}
	// 查找变量是否存在，如果已经定义则答应，否则提示未定义
	v, ok := vars[vName]
	if !ok {
		return errors.Errorf("variable: %s is not assignment", vName)
	}
	fmt.Println(v)
	return nil
}

func resolveStatement(vars Vars, s Statement) error {
	if assignment, ok := s.(*Assignment); ok {
		return resolveAssignment(vars, assignment)
	} else if p, ok := s.(*Print); ok {
		return resolvePrint(vars, p)
	}
	return errors.New("wrong statement")
}

func resolveAST(vars Vars, s *SourceCode) error {
	if len(s.Statements) == 0 {
		return errors.New("no statement to execute")
	}
	for _, v := range s.Statements {
		if err := resolveStatement(vars, v); err != nil {
			return err
		}
	}
	return nil
}

func Execute(code string) {
	vars := NewVars()
	// 语法解析
	ast, err := parse(code)
	if err != nil {
		panic(err)
	}
	// 解析抽象语法树(执行脚本)
	err = resolveAST(vars, ast)
	if err != nil {
		panic(err)
	}
}
