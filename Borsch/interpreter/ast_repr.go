package interpreter

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

func (s *Stmt) String() string {
	if s.IfStmt != nil {
		return "s.IfStmt."
	} else if s.LoopStmt != nil {
		return "s.LoopStmt."
	} else if s.Block != nil {
		return "s.Block."
	} else if s.FunctionDef != nil {
		return "s.FunctionDef."
	} else if s.ClassDef != nil {
		return "s.ClassDef."
	} else if s.ReturnStmt != nil {
		return "повернути ..."
	} else if s.BreakStmt {
		return "перервати"
	} else if s.Assignment != nil {
		return s.Assignment.String() + ";"
	} else if s.Empty {
		return ";"
	}

	panic("unreachable")
}

func (a *Assignment) String() string {
	var lhs []string
	for _, expr := range a.Expressions {
		lhs = append(lhs, expr.String())
	}

	rhsStr := ""
	if a.Next != nil {
		var rhs []string
		for _, expr := range a.Next {
			rhs = append(rhs, expr.String())
		}

		rhsStr = " " + a.Op + " " + strings.Join(rhs, ", ")
	}

	return strings.Join(lhs, ", ") + rhsStr
}

func (e *Expression) String() string {
	return e.LogicalAnd.String()
}

func (a *LogicalAnd) String() string {
	return a.LogicalOr.String() + nextOrEmpty(a.Op, a.Next)
}

func (o *LogicalOr) String() string {
	return o.LogicalNot.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *LogicalNot) String() string {
	return o.Comparison.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *Comparison) String() string {
	return o.BitwiseOr.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *BitwiseOr) String() string {
	return o.BitwiseXor.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *BitwiseXor) String() string {
	return o.BitwiseAnd.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *BitwiseAnd) String() string {
	return o.BitwiseShift.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *BitwiseShift) String() string {
	return o.Addition.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *Addition) String() string {
	return o.MultiplicationOrMod.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *MultiplicationOrMod) String() string {
	return o.Unary.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *Unary) String() string {
	if o.Exponent != nil {
		return o.Op + o.Exponent.String()
	}

	return o.Op + o.Next.String()
}

func (o *Exponent) String() string {
	return o.Primary.String() + nextOrEmpty(o.Op, o.Next)
}

func (o *Primary) String() string {
	switch {
	case o.Constant != nil:
		return o.Constant.String()
	case o.LambdaDef != nil:
		return o.LambdaDef.String()
	case o.AttributeAccess != nil:
		return o.AttributeAccess.String()
	case o.SubExpression != nil:
		return fmt.Sprintf("(%s)", o.SubExpression.String())
	default:
		panic("unreachable")
	}
}

func (o *Constant) String() string {
	switch {
	case o.Integer != nil:
		return fmt.Sprintf("%d", *o.Integer)
	case o.Real != nil:
		return strconv.FormatFloat(*o.Real, 'f', -1, 64)
	case o.Bool != nil:
		if *o.Bool {
			return "істина"
		}

		return "хиба"
	case o.StringValue != nil:
		return fmt.Sprintf("\"%s\"", *o.StringValue)
	case o.List != nil:
		var values []string
		for _, expr := range o.List {
			values = append(values, expr.String())
		}

		return "[" + strings.Join(values, ", ") + "]"
	case o.EmptyList == true:
		return "[]"
	case o.Dictionary != nil:
		var values []string
		for _, entry := range o.Dictionary {
			values = append(values, fmt.Sprintf("%s: %s", entry.Key.String(), entry.Value.String()))
		}

		return "{" + strings.Join(values, ", ") + "}"
	case o.EmptyDictionary == true:
		return "{}"
	default:
		panic("unreachable")
	}
}

func (o *LambdaDef) String() string {
	// TODO:
	return ""
}

func (o *AttributeAccess) String() string {
	str := o.SlicingOrSubscription.String()
	if o.AttributeAccess != nil {
		str += "." + o.AttributeAccess.String()
	}

	return str
}

func (o *SlicingOrSubscription) String() string {
	str := ""
	if o.Call != nil {
		str = o.Call.String()
	} else {
		str = *o.Ident
	}

	if len(o.Ranges) != 0 {
		for _, rng := range o.Ranges {
			str += rng.String()
		}
	}

	return str
}

func (o *Call) String() string {
	var args []string
	if len(o.Arguments) != 0 {
		for _, arg := range o.Arguments {
			args = append(args, arg.String())
		}
	}

	return o.Ident + "(" + strings.Join(args, ", ") + ")"
}

func (o *Range) String() string {
	rightBound := ""
	if o.IsSlicing {
		rightBound = ":"
		if o.RightBound != nil {
			rightBound += o.RightBound.String()
		}
	}

	return "[" + o.LeftBound.String() + rightBound + "]"
}

func nextOrEmpty(op string, next fmt.Stringer) string {
	if !reflect.ValueOf(next).IsNil() {
		return fmt.Sprintf(" %s %s", op, next.String())
	}

	return ""
}
