package interpreter

import "github.com/alecthomas/participle/v2/lexer"

type Package struct {
	Pos lexer.Position

	Stmts []*Stmt `@@*`
}

type ReturnStmt struct {
	Pos lexer.Position

	Expressions []*Expression `"повернути" (@@ ("," @@)*)? ";"`
}

type LoopStmt struct {
	Pos lexer.Position

	Keyword         string           `"цикл"`
	RangeBasedLoop  *RangeBasedLoop  `["(" (@@ `
	ConditionalLoop *ConditionalLoop `|    @@) ")"]`
	Body            *BlockStmts      `"{" @@ "}"`
}

// RangeBasedLoop is a loop with two bounds to
// iterate over.
//
//   цикл (і : 1 .. 7)
//   {
//   }
type RangeBasedLoop struct {
	Pos lexer.Position

	Variable   string      `@Ident ":"`
	LeftBound  *Expression `@@`
	Separator  string      `@("."".")`
	RightBound *Expression `@@`
}

type ConditionalLoop struct {
	Pos lexer.Position

	Condition *Expression `@@`
}

type IfStmt struct {
	Pos lexer.Position

	Condition   *Expression   `"якщо" "(" @@ ")"`
	Body        *BlockStmts   `"{" @@ "}"`
	ElseIfStmts []*ElseIfStmt `(@@ (@@)* )?`
	Else        *BlockStmts   `("інакше" "{" @@ "}")?`
}

type ElseIfStmt struct {
	Condition *Expression `"інакше" "якщо" "(" @@ ")"`
	Body      *BlockStmts `"{" @@ "}"`
}

type BlockStmts struct {
	Pos lexer.Position

	Stmts []*Stmt `@@*`
}

type Stmt struct {
	Pos lexer.Position

	IfStmt      *IfStmt      `  @@`
	LoopStmt    *LoopStmt    `| @@`
	Block       *BlockStmts  `| "{" @@ "}"`
	FunctionDef *FunctionDef `| @@`
	ClassDef    *ClassDef    `| @@`
	ReturnStmt  *ReturnStmt  `| @@`
	BreakStmt   bool         `| @"перервати"`
	Assignment  *Assignment  `| (@@ ";")`
	Empty       bool         `| @";"`
}

type FunctionBody struct {
	Pos lexer.Position

	Stmts *BlockStmts `@@`
}

type FunctionDef struct {
	Pos lexer.Position

	Name          string         `"функція" @Ident`
	ParametersSet *ParametersSet `@@`
	ReturnTypes   []*ReturnType  `[":" (@@ | ("(" (@@ ("," @@)+ )? ")"))]`
	Body          *FunctionBody  `"{" @@ "}"`
}

type ParametersSet struct {
	Pos lexer.Position

	Parameters []*Parameter `"(" (@@ ("," @@)* )? ")"`
}

type Parameter struct {
	Pos lexer.Position

	Name       string `@Ident ":"`
	Type       string `@Ident`
	IsNullable bool   `@"?"?`
}

type ReturnType struct {
	Pos lexer.Position

	Name       string `@Ident`
	IsNullable bool   `@"?"?`
}

type ClassDef struct {
	Pos lexer.Position

	Name    string         `"клас" @Ident`
	IsFinal bool           `@"заключний"?`
	Bases   []string       `[":" (@Ident)+]`
	Members []*ClassMember `"{" @@* "}"`
}

type ClassMember struct {
	Pos lexer.Position

	Variable *Assignment  ` (@@ ";")`
	Method   *FunctionDef `| @@`
	Class    *ClassDef    `| @@`
}

type Assignment struct {
	Pos lexer.Position

	Expressions []*Expression ` @@ ("," @@)*`
	Op          string        `[@"="`
	Next        []*Expression ` @@ ("," @@)*]`
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

	BitwiseXor *BitwiseXor `@@`
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

	Constant        *Constant        `  @@`
	LambdaDef       *LambdaDef       `| @@`
	AttributeAccess *AttributeAccess `| @@`
	SubExpression   *Expression      `| "(" @@ ")"`
}

type Constant struct {
	Pos lexer.Position

	Integer         *int64             `  @Int`
	Real            *float64           `| @Float`
	Bool            *Boolean           `| @("істина" | "хиба")`
	StringValue     *string            `| @String`
	List            []*Expression      `| "[" @@ ("," @@)* "]"`
	EmptyList       bool               `| @("[""]")`
	Dictionary      []*DictionaryEntry `| "{" @@ ("," @@)* "}"`
	EmptyDictionary bool               `| @("{""}")`
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

	ParametersSet        *ParametersSet `@@`
	ReturnTypes          []*ReturnType  `[":" (@@ | ("(" (@@ ("," @@)+ )? ")"))]`
	Body                 *FunctionBody  `"="">" "{" @@ "}"`
	InstantCall          bool           `[ @"("`
	InstantCallArguments []*Expression  `[(@@ ("," @@)*)?] ")"]`
}

type AttributeAccess struct {
	Pos lexer.Position

	SlicingOrSubscription *SlicingOrSubscription `@@`
	AttributeAccess       *AttributeAccess       `("." @@)?`
}

type SlicingOrSubscription struct {
	Pos lexer.Position

	Call   *Call    `( @@`
	Ident  *string  `| @Ident)`
	Ranges []*Range `@@*`
}

type Range struct {
	Pos lexer.Position

	LeftBound  *Expression `"[" @@`
	IsSlicing  bool        `[ @":"`
	RightBound *Expression ` @@?] "]"`
}

type Call struct {
	Pos lexer.Position

	Ident     string        `@Ident`
	Arguments []*Expression `"(" (@@ ("," @@)*)? ")"`
}
