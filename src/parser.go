// 语法解析器
package src

import "github.com/pkg/errors"

// 解析变量名
func parseName(l *Lexer) (string, error) {
	_, name := l.NextTokenIs(TokenName)
	return name, nil
}

// 解析字符串字面量
// 空字符+空白字符 或 非空字符串 + 空白字符
func parseString(l *Lexer) (string, error) {
	switch l.LookAhead() {
	case TokenDoubleQuote: // 空字符串
		l.NextTokenIs(TokenDoubleQuote)
		l.skipIgnored()
		return "", nil
	case TokenQuote: // 非空字符串
		l.NextTokenIs(TokenQuote)
		s := l.scanString()
		l.NextTokenIs(TokenQuote)
		l.skipIgnored()
		return s, nil
	default:
		return "", errors.New("not a string")
	}
}

// 解析变量， 变量由 $+变量名组成
func parseVariable(l *Lexer) (*Variable, error) {
	lineNum := l.LineNum()
	l.NextTokenIs(TokenVarPrefix)
	name, err := parseName(l)
	if err != nil {
		return nil, err
	}
	// 跳过空白字符
	l.skipIgnored()
	return &Variable{LineNum: lineNum, Name: name}, nil
}

// 解析赋值语句
// 赋值语句由 $name 空格 = 空白 string 空白 构成
func parseAssignment(l *Lexer) (*Assignment, error) {
	lineNum := l.LineNum()
	v, err := parseVariable(l) // 解析变量
	if err != nil {
		return nil, err
	}
	l.skipIgnored()
	l.NextTokenIs(TokenEqual)
	l.skipIgnored()
	s, err := parseString(l)
	if err != nil {
		return nil, err
	}
	l.skipIgnored()
	return &Assignment{LineNum: lineNum, Variable: v, String: s}, nil
}

// 解析print语句
// print语句由 print( 空白 变量 空白 ）空白 构成
func parsePrint(l *Lexer) (*Print, error) {
	lineNum := l.LineNum()
	l.NextTokenIs(TokenPrint)
	l.NextTokenIs(TokenLeftParen)
	l.skipIgnored()
	v, err := parseVariable(l)
	if err != nil {
		return nil, err
	}
	l.skipIgnored()
	l.NextTokenIs(TokenRightParen)
	l.skipIgnored()
	return &Print{LineNum: lineNum, Variable: v}, nil
}



// 解析语句：包含赋值语句和print语句
func parseStatement(l *Lexer) (Statement, error) {
	l.skipIgnored()
	if t := l.LookAhead(); t == TokenPrint {
		return parsePrint(l)
	} else if t == TokenVarPrefix {
		return parseAssignment(l)
	}
	return nil, errors.New("not support statement")
}

func parseStatements(l *Lexer) ([]Statement, error) {
	res := make([]Statement, 0, 2)
	for l.LookAhead() != TokenEOF {
		s, err := parseStatement(l)
		if err != nil {
			return nil, err
		}
		res = append(res, s)
	}
	return res, nil
}

// 解析源代码, 源代码由一些语句构成
func parseSourceCode(l *Lexer) (*SourceCode, error) {
	lineNum := l.LineNum()
	ss, err := parseStatements(l)
	if err != nil {
		return nil, err
	}
	return &SourceCode{LineNum: lineNum, Statements: ss}, nil
}

// 总体解析器
func parse(code string) (*SourceCode, error) {
	l := NewLexer(code)
	s, err := parseSourceCode(l)
	if err != nil {
		return nil, err
	}
	l.NextTokenIs(TokenEOF)
	return s, nil
}


