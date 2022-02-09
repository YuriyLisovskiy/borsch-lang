package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (node *Package) String() string {
	panic("unreachable")
}

func (node *BlockStmts) String() string {
	return node.GetCurrentStmt().String()
}

func (node *Stmt) String() string {
	if node.Throw != nil {
		return node.Throw.String()
	} else if node.Unsafe != nil {
		return node.Unsafe.String()
	} else if node.IfStmt != nil {
		return node.IfStmt.String("")
	} else if node.LoopStmt != nil {
		return node.LoopStmt.String()
	} else if node.Block != nil {
		return node.Block.String()
	} else if node.FunctionDef != nil {
		return node.FunctionDef.String()
	} else if node.ClassDef != nil {
		return node.ClassDef.String()
	} else if node.ReturnStmt != nil {
		return node.ReturnStmt.String()
	} else if node.BreakStmt {
		return "перервати"
	} else if node.Assignment != nil {
		return node.Assignment.String() + ";"
	} else if node.Empty {
		return ";"
	}

	panic("unreachable")
}

func (node *Throw) String() string {
	return fmt.Sprintf("панікувати %s", node.Expression.String())
}

func (node *Unsafe) String() string {
	return "небезпечно"
}

func (node *Catch) String() string {
	return fmt.Sprintf("піймати (%s: %s)", node.ErrorVar, node.ErrorType.String())
}

func (node *Assignment) String() string {
	var lhs []string
	for _, expr := range node.Expressions {
		lhs = append(lhs, expr.String())
	}

	rhsStr := ""
	if node.Next != nil {
		var rhs []string
		for _, expr := range node.Next {
			rhs = append(rhs, expr.String())
		}

		rhsStr = " " + node.Op + " " + strings.Join(rhs, ", ")
	}

	return strings.Join(lhs, ", ") + rhsStr
}

func (node *Expression) String() string {
	return node.LogicalAnd.String()
}

func (node *LogicalAnd) String() string {
	return node.LogicalOr.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *LogicalOr) String() string {
	return node.LogicalNot.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *LogicalNot) String() string {
	return node.Comparison.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *Comparison) String() string {
	return node.BitwiseOr.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *BitwiseOr) String() string {
	return node.BitwiseXor.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *BitwiseXor) String() string {
	return node.BitwiseAnd.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *BitwiseAnd) String() string {
	return node.BitwiseShift.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *BitwiseShift) String() string {
	return node.Addition.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *Addition) String() string {
	return node.MultiplicationOrMod.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *MultiplicationOrMod) String() string {
	return node.Unary.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *Unary) String() string {
	if node.Exponent != nil {
		return node.Op + node.Exponent.String()
	}

	return node.Op + node.Next.String()
}

func (node *Exponent) String() string {
	return node.Primary.String() + nextOrEmpty(node.Op, node.Next)
}

func (node *Primary) String() string {
	switch {
	case node.Constant != nil:
		return node.Constant.String()
	case node.LambdaDef != nil:
		return node.LambdaDef.String()
	case node.AttributeAccess != nil:
		return node.AttributeAccess.String()
	case node.SubExpression != nil:
		return fmt.Sprintf("(%s)", node.SubExpression.String())
	default:
		panic("unreachable")
	}
}

func (node *Constant) String() string {
	switch {
	case node.Integer != nil:
		return fmt.Sprintf("%d", *node.Integer)
	case node.Real != nil:
		return strconv.FormatFloat(*node.Real, 'f', -1, 64)
	case node.Bool != nil:
		if *node.Bool {
			return "істина"
		}

		return "хиба"
	case node.StringValue != nil:
		return fmt.Sprintf("\"%s\"", *node.StringValue)
	case node.List != nil:
		var values []string
		for _, expr := range node.List {
			values = append(values, expr.String())
		}

		return "[" + strings.Join(values, ", ") + "]"
	case node.EmptyList == true:
		return "[]"
	case node.Dictionary != nil:
		var values []string
		for _, entry := range node.Dictionary {
			values = append(values, fmt.Sprintf("%s: %s", entry.Key.String(), entry.Value.String()))
		}

		return "{" + strings.Join(values, ", ") + "}"
	case node.EmptyDictionary == true:
		return "{}"
	default:
		panic("unreachable")
	}
}

func (node *DictionaryEntry) String() string {
	return fmt.Sprintf("%s: %s", node.Key.String(), node.Value.String())
}

func (node *LambdaDef) String() string {
	// TODO:
	return ""
}

func (node *AttributeAccess) String() string {
	str := node.SlicingOrSubscription.String()
	if node.AttributeAccess != nil {
		str += "." + node.AttributeAccess.String()
	}

	return str
}

func (node *IdentOrCall) String() string {
	str := ""
	if node.Call != nil {
		str = node.Call.String()
	} else {
		str = *node.Ident
	}

	if node.SlicingOrSubscription != nil {
		str += node.SlicingOrSubscription.String()
	}

	return str
}

func (node *SlicingOrSubscription) String() string {
	str := ""
	if len(node.Ranges) != 0 {
		for _, rng := range node.Ranges {
			str += rng.String()
		}
	}

	return str
}

func (node *Call) String() string {
	var args []string
	if len(node.Arguments) != 0 {
		for _, arg := range node.Arguments {
			args = append(args, arg.String())
		}
	}

	return node.Ident + "(" + strings.Join(args, ", ") + ")"
}

func (node *Range) String() string {
	rightBound := ""
	if node.IsSlicing {
		rightBound = ":"
		if node.RightBound != nil {
			rightBound += node.RightBound.String()
		}
	}

	return "[" + node.LeftBound.String() + rightBound + "]"
}

func (node *IfStmt) String(indent string) string {
	result := fmt.Sprintf("якщо (%s) {\n%s  ...\n%s}", node.Condition.String(), indent, indent)
	if len(node.ElseIfStmts) != 0 {
		for _, stmt := range node.ElseIfStmts {
			result += fmt.Sprintf("\n%s", stmt.String(indent))
		}
	}

	if node.Else != nil {
		result += fmt.Sprintf("\n%sінакше {\n%s  ...\n%s}", indent, indent, indent)
	}

	return result
}

func (node *ElseIfStmt) String(indent string) string {
	return fmt.Sprintf("%sінакше якщо (%s) {\n%s  ...\n%s}", indent, node.Condition.String(), indent, indent)
}

func (node *LoopStmt) String() string {
	result := "цикл "
	if node.RangeBasedLoop != nil {
		result += node.RangeBasedLoop.String()
	} else if node.ConditionalLoop != nil {
		result += node.ConditionalLoop.String()
	}

	return result
}

func (node *RangeBasedLoop) String() string {
	return fmt.Sprintf(
		"(%s : %s %s %s)",
		node.Variable,
		node.LeftBound.String(),
		node.Separator,
		node.RightBound.String(),
	)
}

func (node *ConditionalLoop) String() string {
	return fmt.Sprintf("(%s)", node.Condition.String())
}

func (node *FunctionDef) String() string {
	var returnTypes []string
	for _, returnType := range node.ReturnTypes {
		returnTypes = append(returnTypes, returnType.String())
	}

	returnTypesStr := ""
	if len(returnTypes) != 0 {
		returnTypesStr = fmt.Sprintf(": %s", strings.Join(returnTypes, ", "))
	}

	return fmt.Sprintf("функція %s(%s)%s", node.Name, node.ParametersSet.String(), returnTypesStr)
}

func (node *ParametersSet) String() string {
	var parameters []string
	for _, parameter := range node.Parameters {
		parameters = append(parameters, parameter.String())
	}

	return strings.Join(parameters, ", ")
}

func (node *Parameter) String() string {
	result := fmt.Sprintf("%s: %s", node.Name, node.Type)
	if node.IsNullable {
		result += "?"
	}

	return result
}

func (node *ReturnType) String() string {
	result := node.Name
	if node.IsNullable {
		result += "?"
	}

	return result
}

func (node *ReturnStmt) String() string {
	var expressions []string
	for _, expression := range node.Expressions {
		expressions = append(expressions, expression.String())
	}

	return strings.Join(expressions, ", ")
}

func (node *ClassDef) String() string {
	var bases []string
	for _, base := range node.Bases {
		bases = append(bases, base)
	}

	basesStr := ""
	if len(bases) != 0 {
		basesStr = fmt.Sprintf(": %s", strings.Join(bases, ", "))
	}

	final := ""
	if node.IsFinal {
		final = "заключний "
	}

	return fmt.Sprintf("клас %s %s%s", node.Name, basesStr, final)
}

func (node *ClassMember) String() string {
	if node.Variable != nil {
		return node.Variable.String()
	}

	if node.Method != nil {
		return node.Method.String()
	}

	if node.Class != nil {
		return node.Class.String()
	}

	panic("unreachable")
}

func nextOrEmpty(op string, next fmt.Stringer) string {
	if !reflect.ValueOf(next).IsNil() {
		return fmt.Sprintf(" %s %s", op, next.String())
	}

	return ""
}
