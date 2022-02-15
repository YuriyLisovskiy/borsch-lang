package interpreter

import "github.com/alecthomas/participle/v2/lexer"

func (node *Throw) Position() lexer.Position {
	return node.Pos
}

func (node *Unsafe) Position() lexer.Position {
	return node.Pos
}

func (node *Catch) Position() lexer.Position {
	return node.Pos
}

func (node *Stmt) Position() lexer.Position {
	return node.Pos
}

func (node *BlockStmts) Position() lexer.Position {
	return node.Pos
}

func (node *FunctionBody) Position() lexer.Position {
	return node.Pos
}

func (node *ClassDef) Position() lexer.Position {
	return node.Pos
}

func (node *ClassMember) Position() lexer.Position {
	return node.Pos
}

func (node *Assignment) Position() lexer.Position {
	return node.Pos
}

func (node *IdentOrCall) Position() lexer.Position {
	return node.Pos
}
