// 定义了词法解析
package src

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// 语言的使用例子: 只包含赋值语句和打印语句
// $a = "hell0"
// print(a)

// 支持的token枚举

const (
	TokenEOF = iota  	// end of file
	TokenVarPrefix  	// $
	TokenLeftParen  	// (
	TokenRightParen 	// )
	TokenEqual     		// =
	TokenQuote	  		// " 单个双引号
	TokenDoubleQuote	// "" 两个双引号 空字符串
	TokenName			// 变量名称 ::=[_A-Za-z][_0-9A-Za-z]*
	TokenPrint			// print
	TokenIgnored		// 忽略的字符串（空白字符串）
)

var tokenNames = map[int]string {
	TokenEOF        : "EOF",
	TokenVarPrefix  : "$",
	TokenLeftParen  : "(",
	TokenRightParen : ")",
	TokenEqual      : "=",
	TokenQuote      : "\"",
	TokenDoubleQuote: "\"\"",
	TokenName       : "Name",
	TokenPrint      : "print",
	TokenIgnored    : "Ignored",
}

// 关键字
var keywords = map[string]int {
	"print": TokenPrint,
}
// 变量名的正则模式 匹配一个或多个字母，数字，下划线
var regexName = regexp.MustCompile(`^[_\d\w]+`)

// 词法解析器定义
type Lexer struct {
	sourceCode string  	// 源码
	lineNum int			// 行号
	nextToken string	// 下一个token
	nextTokenType int	// 下一个token类型
	nextTokenLineNum int // 下一个token所在行号, 该值如果非0，说明前看的时候已经加载过了，可以优化性能
}

func NewLexer(sourceCode string) *Lexer {
	return &Lexer{
		sourceCode:       sourceCode,
		lineNum:          1,
		nextToken:        "",
		nextTokenType:    0,
		nextTokenLineNum: 0,
	}
}

func (l *Lexer) LineNum() int {
	return l.lineNum
}

// 源代码向后跳n个字符
func (l *Lexer) skipSourceCode(n int){
	l.sourceCode = l.sourceCode[n:]
}

func (l *Lexer) nextSourceCodeIs(s string) bool{
	return strings.HasPrefix(l.sourceCode, s)
}

// 判断源码后续部分是否是忽略的空白字符，如果是就跳过
func (l *Lexer) isIgnored() bool{
	ignore := false
	isNewLine := func(c byte) bool {
		return c == '\r' || c == '\n'
	}
	isWhiteSpace := func(c byte) bool {
		return c == '\t' || c == '\n' || c == '\v' || c == '\f' || c == '\r' || c == ' '
	}

	for len(l.sourceCode) > 0 {
		if l.nextSourceCodeIs("\r\n") || l.nextSourceCodeIs("\n\r") {
			l.skipSourceCode(2)
			l.lineNum += 1
			ignore = true
		} else if isNewLine(l.sourceCode[0]) { // 先判断空行再判断空格符
			l.skipSourceCode(1)
			l.lineNum += 1
			ignore = true
		} else if isWhiteSpace(l.sourceCode[0]) {
			l.skipSourceCode(1)
			ignore = true
		} else {
			break
		}
	}
	return ignore
}

func isLetter(c byte) bool {
	return (c >='A' && c <= 'Z') || ( c >= 'a' && c <= 'z')
}

func (l *Lexer) scan(reg *regexp.Regexp) string {
	if token := reg.FindString(l.sourceCode); token != "" {
		l.skipSourceCode(len(token))
		return token
	}
	panic("unreachable")
	return ""
}


func (l *Lexer) scanName() string {
	return l.scan(regexName)
}

func (l *Lexer) MatchToken() (lineNum int, tokenType int, token string){
	if l.isIgnored() {
		return l.lineNum, TokenIgnored, tokenNames[TokenIgnored]
	}
	if len(l.sourceCode) == 0 {
		return l.lineNum, TokenEOF, tokenNames[TokenEOF]
	}
	// check token
	switch l.sourceCode[0] {
	case '$':
		l.skipSourceCode(1)
		return l.lineNum, TokenVarPrefix, "$"
	case '(':
		l.skipSourceCode(1)
		return l.lineNum, TokenLeftParen, "("
	case ')':
		l.skipSourceCode(1)
		return l.lineNum, TokenRightParen, ")"
	case '=':
		l.skipSourceCode(1)
		return l.lineNum, TokenEqual, "="
	case '"':
		// 判断是否是连续两个双引号
		if l.nextSourceCodeIs("\"\"") {
			l.skipSourceCode(2)
			return l.lineNum, TokenDoubleQuote, "\"\""
		}
		l.skipSourceCode(1)
		return l.lineNum, TokenQuote, "\""
	}
	// 判断是否以下划线或者字母开头
	if l.sourceCode[0] == '_' || isLetter(l.sourceCode[0]) {
		// 尝试解析变量名
		token := l.scanName()
		// print 关键字
		if tokenType, ok := keywords[token]; ok {
			return l.lineNum, tokenType, token
		}
		return l.lineNum, TokenName, token // 普通变量
	}
	err := errors.Errorf("unexpected symbol near '%q'.", l.sourceCode[0])
	panic(err)
	return
}

// 获取下一个token
func (l *Lexer) NextToken() (lineNum int, tokenType int, token string) {
	// 如果nextTokenNum 非0 说明前看时已经预加载了，只需要直接迭代之后返回就行
	if l.nextTokenLineNum > 0 {
		// 直接取值
		lineNum = l.nextTokenLineNum
		tokenType = l.nextTokenType
		token = l.nextToken

		// next属性reset
		l.lineNum = l.nextTokenLineNum
		l.nextTokenLineNum = 0
		return
	}
	return l.MatchToken()
}

// 验证下一个token是预期传入的类型，如果不是就panic，如果是，返回下一个token的行号和内容
func (l *Lexer) NextTokenIs(tokenType int) (lineNum int, token string) {
	nowLineNum, nowTokenType, nowToken := l.NextToken()
	if tokenType != nowTokenType {
		err := errors.Errorf("line:%v syntax error near '%s', expected token: {%s} but got {%s}",
			nowLineNum, nowToken, tokenNames[tokenType], tokenNames[nowTokenType])
		panic(err)
	}
	return nowLineNum, nowToken
}

func (l *Lexer) LookAheadAndSkip(expectedType int) {
	oldLine := l.lineNum
	newLine, tokenType, token := l.NextToken()
	// 回退
	if expectedType != tokenType {
		l.lineNum = oldLine
		l.nextTokenLineNum = newLine
		l.nextToken = token
		l.nextTokenType = tokenType
	}
}

func (l *Lexer) skipIgnored() {
	l.LookAheadAndSkip(TokenIgnored)
}

// 前看一个token，并且返回token类型
func (l *Lexer) LookAhead() int {
	// 已经前看过，直接返回
	if l.nextTokenLineNum > 0 {
		return l.nextTokenType
	}
	nowLine := l.lineNum
	nextLine, tokenType, token := l.NextToken()
	l.lineNum = nowLine
	l.nextTokenType = tokenType
	l.nextTokenLineNum = nextLine
	l.nextToken = token
	return tokenType
}

func (l *Lexer) scanString() string{
	hasQuote := false
	idx := 0
	for i , c := range l.sourceCode{
		if c == '"' {
			hasQuote = true
			idx = i
			break
		}
	}
	if !hasQuote {
		panic("string error")
	}
	s := l.sourceCode[:idx]
	l.skipSourceCode(idx)
	return s
}
