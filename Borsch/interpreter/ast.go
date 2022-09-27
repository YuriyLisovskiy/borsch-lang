package interpreter

import (
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin/types"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
	"github.com/alecthomas/participle/v2/lexer"
)

type Package struct {
	Pos lexer.Position

	Stmts *BlockStmts `@@`
}

type Throw struct {
	Pos lexer.Position

	Expression *Expression `"панікувати" @@`
}

type Block struct {
	Pos lexer.Position

	Stmts       *BlockStmts `"блок" @@`
	CatchBlocks []*Catch    `[ @@ (@@)* ] "кінець"`
}

type Ident string

func (i *Ident) Capture(values []string) error {
	ident := values[0]
	if ident == "кінець" {
		// TODO: write ukr error!
		return utilities.SyntaxError("unexpected token 'кінець'")
	}

	*i = Ident(ident)
	return nil
}

func (i Ident) String() string {
	return string(i)
}

type Catch struct {
	Pos lexer.Position

	ErrorVar  Ident            `"піймати" "(" @Ident`
	ErrorType *AttributeAccess `":" @@ ")"`
	Stmts     *BlockStmts      `@@`
}

type ReturnStmt struct {
	Pos lexer.Position

	Expressions []*Expression `"повернути" (@@ ("," @@)*)?`
}

type LoopStmt struct {
	Pos lexer.Position

	Keyword         string           `"цикл"`
	RangeBasedLoop  *RangeBasedLoop  `["(" (@@ `
	ConditionalLoop *ConditionalLoop `|    @@) ")"]`
	Body            *BlockStmts      `@@ "кінець"`
}

// RangeBasedLoop is a loop with two bounds to
// iterate over.
//
// Example:
//   цикл (і : 1 .. 7)
//   {
//   }
type RangeBasedLoop struct {
	Pos lexer.Position

	Variable   Ident       `@Ident ":"`
	LeftBound  *Expression `@@`
	Separator  string      `@("."".")`
	RightBound *Expression `@@`
}

// ConditionalLoop
//
// Example:
//   цикл (умова_логічного_типу)
//   {
//   }
type ConditionalLoop struct {
	Pos lexer.Position

	Condition *Expression `@@`
}

type IfStmt struct {
	Pos lexer.Position

	Condition   *Expression   `"якщо" "(" @@ ")"`
	Body        *BlockStmts   `@@`
	ElseIfStmts []*ElseIfStmt `(@@ (@@)*)?`
	Else        *BlockStmts   `("інакше" @@)? "кінець"`
}

type ElseIfStmt struct {
	Condition *Expression `"інакше" "якщо" "(" @@ ")"`
	Body      *BlockStmts `@@`
}

type BlockStmts struct {
	Pos lexer.Position

	Stmts []*Stmt `@@*`

	stmtPos int
}

func (node *BlockStmts) GetCurrentStmt() *Stmt {
	return node.Stmts[node.stmtPos]
}

type Stmt struct {
	Pos lexer.Position

	Throw       *Throw       `(?!("піймати" | "інакше" | "кінець")) (@@`
	IfStmt      *IfStmt      `| @@ ";"`
	LoopStmt    *LoopStmt    `| @@ ";"`
	Block       *Block       `| @@ ";"`
	FunctionDef *FunctionDef `| @@ ";"`
	ClassDef    *ClassDef    `| @@ ";"`
	ReturnStmt  *ReturnStmt  `| @@ ";"`
	BreakStmt   bool         `| @"перервати" ";"`
	Assignment  *Assignment  `| (@@ ";")`
	Empty       bool         `| @";")`
}

type FunctionBody struct {
	Pos lexer.Position

	Stmts *BlockStmts `@@`
}

type FunctionDef struct {
	Pos lexer.Position

	Name          Ident          `"функція" @Ident`
	ParametersSet *ParametersSet `@@`
	ReturnTypes   []*ReturnType  `[":" (@@ | ("(" (@@ ("," @@)+ )? ")"))]`
	Body          *FunctionBody  `@@ "кінець"`
}

type ParametersSet struct {
	Pos lexer.Position

	Parameters []*Parameter `"(" (@@ ("," @@)* )? ")"`
}

type Parameter struct {
	Pos lexer.Position

	Name       Ident `@Ident ":"`
	TypeName   Ident `@Ident | @"?"`
	IsNullable bool  `@"?"?`
}

type ReturnType struct {
	Pos lexer.Position

	Name       Ident `@Ident`
	IsNullable bool  `@"?"?`
}

type ClassDef struct {
	Pos lexer.Position

	Name    Ident          `"клас" @Ident`
	IsFinal bool           `@"заключний"?`
	Bases   []Ident        `(":" @Ident ("," @Ident)*)?`
	Members []*ClassMember `(@@ ";")* "кінець"`

	operators []*types.Method
}

type ClassMember struct {
	Pos lexer.Position

	Method   *FunctionDef `  @@`
	Operator *OperatorDef `| @@`
	Class    *ClassDef    `| @@`
	Variable *Assignment  `| @@`
}

type OperatorDef struct {
	Pos lexer.Position

	Op            string         `"оператор" @("==" | "!=" | "<" | "<""=" | ">" | ">""=" | "+" | "-" | "/" | "*""*" | "*" | "%" | "<""<" | ">"">" | "|" | "^" | "&" | "~" | "&""&" | "|""|")`
	ParametersSet *ParametersSet `@@`
	ReturnTypes   []*ReturnType  `[":" (@@ | ("(" (@@ ("," @@)+ )? ")"))]`
	Body          *FunctionBody  `@@ "кінець"`
}

type Assignment struct {
	Pos lexer.Position

	Expressions []*Expression `@@ ("," @@)*`
	Op          string        `[     @"="`
	Next        []*Expression `@@ ("," @@)*]`
}

type Expression struct {
	Pos lexer.Position

	LogicalAnd *LogicalAnd `@@`
}

type LogicalAnd struct {
	Pos lexer.Position

	LogicalOr *LogicalOr  `@@`
	Op        string      `[ @("&""&")`
	Next      *LogicalAnd `  @@ ]`
}

type LogicalOr struct {
	Pos lexer.Position

	LogicalNot *LogicalNot `@@`
	Op         string      `[ @("|""|")`
	Next       *LogicalOr  `  @@ ]`
}

type LogicalNot struct {
	Pos lexer.Position

	Op         string      `  ( @"!"`
	Next       *LogicalNot `    @@ )`
	Comparison *Comparison `| @@`
}

type Comparison struct {
	Pos lexer.Position

	BitwiseOr *BitwiseOr  `@@`
	Op        string      `[ @(">""=" | ">" | "<""=" | "<" | "=""=" | "!""=")`
	Next      *Comparison `  @@ ]`
}

type BitwiseOr struct {
	Pos lexer.Position

	BitwiseXor *BitwiseXor `  @@`
	Op         string      `[ @("|")`
	Next       *BitwiseOr  `  @@ ]`
}

type BitwiseXor struct {
	Pos lexer.Position

	BitwiseAnd *BitwiseAnd `@@`
	Op         string      `[ @("^")`
	Next       *BitwiseXor `  @@ ]`
}

type BitwiseAnd struct {
	Pos lexer.Position

	BitwiseShift *BitwiseShift `@@`
	Op           string        `[ @("&")`
	Next         *BitwiseAnd   `  @@ ]`
}

type BitwiseShift struct {
	Pos lexer.Position

	Addition *Addition     `@@`
	Op       string        `[ @(">"">" | "<""<")`
	Next     *BitwiseShift `  @@ ]`
}

type Addition struct {
	Pos lexer.Position

	MultiplicationOrMod *MultiplicationOrMod `@@`
	Op                  string               `[ @("-" | "+")`
	Next                *Addition            `  @@ ]`
}

type MultiplicationOrMod struct {
	Pos lexer.Position

	Unary *Unary               `@@`
	Op    string               `[ @("/" | "*" | "%")`
	Next  *MultiplicationOrMod `  @@ ]`
}

type Unary struct {
	Pos lexer.Position

	Op       string    `  ( @("+" | "-" | "~")`
	Next     *Unary    `    @@ )`
	Exponent *Exponent `| @@`
}

type Exponent struct {
	Pos lexer.Position

	Primary *Primary  `@@`
	Op      string    `[ @("*""*")`
	Next    *Exponent `  @@ ]`
}

type Primary struct {
	Pos lexer.Position

	Literal         *Literal         `  @@`
	LambdaDef       *LambdaDef       `| @@`
	AttributeAccess *AttributeAccess `| @@`
	SubExpression   *Expression      `| "(" @@ ")"`
}

type Literal struct {
	Pos lexer.Position

	Nil             bool               `  @"нуль"`
	Integer         *string            `| @Int`
	Real            *string            `| @Float`
	Bool            *Boolean           `| @("істина" | "хиба")`
	StringValue     *string            `| @String`
	MultilineString *string            `| @RawString`
	List            []*Expression      `| "[" @@ ("," @@)* "]"`
	EmptyList       bool               `| @("[""]")`
	Dictionary      []*DictionaryEntry `| "{" @@ ("," @@)* "}"`
	EmptyDictionary bool               `| @("{""}")`
	// SubExpression   *Expression        `| "(" @@ ")"`
}

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "істина"
	return nil
}

type DictionaryEntry struct {
	Pos lexer.Position

	Key   *Expression `@@`
	Value *Expression `":" @@`
}

type LambdaDef struct {
	Pos lexer.Position

	ParametersSet        *ParametersSet `"лямбда" @@`
	ReturnTypes          []*ReturnType  `[":" (@@ | ("(" (@@ ("," @@)+ )? ")"))]`
	Body                 *FunctionBody  `@@ "кінець"`
	InstantCall          bool           `[ @"("`
	InstantCallArguments []*Expression  `[(@@ ("," @@)*)?] ")"]`
}

type AttributeAccess struct {
	Pos lexer.Position

	IdentOrCall     *IdentOrCall     `@@`
	AttributeAccess *AttributeAccess `("." @@)?`
}

type IdentOrCall struct {
	Pos lexer.Position

	Call                  *Call                  `( @@`
	Ident                 *Ident                 `| @Ident)`
	SlicingOrSubscription *SlicingOrSubscription `@@?`
}

type SlicingOrSubscription struct {
	Pos lexer.Position

	Ranges []*Range `@@+`
}

type Range struct {
	Pos lexer.Position

	LeftBound  *Expression `"[" @@`
	IsSlicing  bool        `[ @":"`
	RightBound *Expression ` @@?] "]"`
}

type Call struct {
	Pos lexer.Position

	Ident     Ident         `@Ident`
	Arguments []*Expression `"(" (@@ ("," @@)*)? ")"`
}
