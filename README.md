# my_pineapple
实现了一个pineapple语言编译器
参照了https://github.com/karminski/pineapple 的代码

代码架构:
lexer: 负责解析一个Token
parser: 迭代解析token，以解析出一个个表达式，所有表达式构成程序，最终形成AST
backend: 对AST进行后续遍历并利用栈式计算机即可生成汇编代码
pineapple语言只支持赋值和打印，因此没有使用树结构存储表达式，而是使用了slice简单表示，因此重复声明(将字符串绑定到字面量)不会报错

该语言只支持字符串字面量的赋值和打印
例子:
$a="Hello"
print($a)
