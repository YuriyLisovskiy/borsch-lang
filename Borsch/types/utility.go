package types

import (
	"errors"
	"fmt"
	"hash/fnv"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

func hashObject(s string) uint64 {
	h := fnv.New64()
	_, _ = h.Write([]byte(s))
	return h.Sum64()
}

func getIndex(index, length int64) (int64, error) {
	if index >= 0 && index < length {
		return index, nil
	} else if index < 0 && index >= -length {
		return length + index, nil
	}

	return 0, errors.New("індекс за межами послідовності")
}

func GetTypeName(typeValue uint64) string {
	switch typeValue {
	case AnyTypeHash:
		return "довільний"
	case NilTypeHash:
		return "нульовий"
	case RealTypeHash:
		return "дійсний"
	case IntegerTypeHash:
		return "цілий"
	case StringTypeHash:
		return "рядок"
	case BoolTypeHash:
		return "логічний"
	case ListTypeHash:
		return "список"
	case DictionaryTypeHash:
		return "словник"
	case PackageTypeHash:
		return "пакет"
	case FunctionTypeHash:
		return "функція"
	case TypeClassTypeHash:
		return "тип"
	default:
		return "невідомий"
	}
}

func GetTypeHash(typeName string) uint64 {
	switch typeName {
	case "довільний":
		return AnyTypeHash
	case "нульовий":
		return NilTypeHash
	case "дійсний":
		return RealTypeHash
	case "цілий":
		return IntegerTypeHash
	case "рядок":
		return StringTypeHash
	case "логічний":
		return BoolTypeHash
	case "список":
		return ListTypeHash
	case "словник":
		return DictionaryTypeHash
	case "пакет":
		return PackageTypeHash
	case "функція":
		return FunctionTypeHash
	case "тип":
		return TypeClassTypeHash
	default:
		return hashObject(typeName)
	}
}

func IsBuiltinType(typeName string) bool {
	return GetTypeHash(typeName) <= FunctionTypeHash
}

func CheckResult(ctx common.Context, result common.Type, function *FunctionInstance) error {
	if len(function.ReturnTypes) == 1 {
		err := checkSingleResult(ctx, result, function.ReturnTypes[0], function.Name)
		if err != nil {
			return errors.New(fmt.Sprintf(err.Error(), ""))
		}

		return nil
	}

	switch value := result.(type) {
	case ListInstance:
		if int64(len(function.ReturnTypes)) != value.Length(ctx) {
			var expectedTypes []string
			for _, retType := range function.ReturnTypes {
				expectedTypes = append(expectedTypes, retType.String())
			}

			var typesGot []string
			for _, retType := range value.Values {
				typesGot = append(typesGot, retType.GetTypeName())
			}

			return util.RuntimeError(
				fmt.Sprintf(
					"'%s' повертає значення з типами (%s), отримано (%s)",
					function.Name,
					strings.Join(expectedTypes, ", "),
					strings.Join(typesGot, ", "),
				),
			)
		}

		// TODO: check values in list
		for i, returnType := range function.ReturnTypes {
			if err := checkSingleResult(ctx, value.Values[i], returnType, function.Name); err != nil {
				return errors.New(fmt.Sprintf(err.Error(), fmt.Sprintf(" на позиції %d", i+1)))
			}
		}
	default:
		var expectedTypes []string
		for _, retType := range function.ReturnTypes {
			expectedTypes = append(expectedTypes, retType.String())
		}

		return util.RuntimeError(
			fmt.Sprintf(
				"'%s()' повертає значення з типами '(%s)', отримано '%s'",
				function.Name,
				strings.Join(expectedTypes, ", "),
				value.GetTypeName(),
			),
		)
	}

	return nil
}

func makeFuncSignature(funcName string) string {
	if funcName == "" {
		return "лямбда-вираз"
	}

	return fmt.Sprintf("'%s()", funcName)
}

func checkSingleResult(
	ctx common.Context,
	result common.Type,
	returnType FunctionReturnType,
	funcName string,
) error {
	if result.(ObjectInstance).GetPrototype() == Nil {
		if returnType.Type != Nil && !returnType.IsNullable {
			return util.RuntimeError(
				fmt.Sprintf(
					"%s повертає ненульове значення%s, отримано '%s'",
					makeFuncSignature(funcName),
					"%s",
					result.String(ctx),
				),
			)
		}
	} else if returnType.Type != Any && result.(ObjectInstance).GetPrototype() != returnType.Type {
		return util.RuntimeError(
			fmt.Sprintf(
				"%s повертає значення типу '%s'%s, отримано значення з типом '%s'",
				makeFuncSignature(funcName), returnType.String(), "%s", result.GetTypeName(),
			),
		)
	}

	return nil
}

func CheckFunctionArguments(
	ctx common.Context,
	function *FunctionInstance,
	args *[]common.Type,
	_ *map[string]common.Type,
) error {
	parametersLen := len(*args)
	argsLen := len(function.Arguments)
	if argsLen > 0 && function.Arguments[argsLen-1].IsVariadic {
		argsLen--
		if parametersLen > argsLen {
			parametersLen = argsLen
		}
	}

	if parametersLen != argsLen {
		diffLen := argsLen - parametersLen
		if diffLen > 0 {
			end1 := "ій"
			end2 := "ий"
			end3 := ""
			if diffLen > 1 && diffLen < 5 {
				end1 = "і"
				end2 = "і"
				end3 = "и"
			} else if diffLen != 1 {
				end1 = "і"
				end2 = "их"
				end3 = "ів"
			}

			parametersStr := ""
			for c := parametersLen; c < argsLen; c++ {
				parametersStr += fmt.Sprintf("'%s'", function.Arguments[c].Name)
				if c < argsLen-2 {
					parametersStr += ", "
				} else if c < argsLen-1 {
					parametersStr += " та "
				}
			}

			return util.RuntimeError(
				fmt.Sprintf(
					"при виклику '%s()' відсутн%s %d необхідн%s параметр%s: %s",
					function.Name, end1, diffLen, end2, end3, parametersStr,
				),
			)
		} else {
			end1 := "ий"
			end2 := ""
			if argsLen > 1 && argsLen < 5 {
				end1 = "і"
				end2 = "и"
			} else if argsLen != 1 {
				end1 = "их"
				end2 = "ів"
			}

			return util.RuntimeError(
				fmt.Sprintf(
					"'%s()' приймає %d необхідн%s параметр%s, отримано %d",
					function.Name, argsLen, end1, end2, parametersLen,
				),
			)
		}
	}

	var c int
	for c = 0; c < argsLen; c++ {
		arg := (*args)[c]
		argPrototype := arg.(ObjectInstance).GetPrototype()
		if argPrototype == Nil {
			if function.Arguments[c].Type != Nil && !function.Arguments[c].IsNullable {
				return util.RuntimeError(
					fmt.Sprintf(
						"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
						function.Arguments[c].Name, arg.String(ctx),
					),
				)
			}
		} else if function.Arguments[c].Type != Any && argPrototype != function.Arguments[c].Type {
			return util.RuntimeError(
				fmt.Sprintf(
					"аргумент '%s' очікує параметр з типом '%s', отримано '%s'",
					function.Arguments[c].Name, function.Arguments[c].GetTypeName(), arg.GetTypeName(),
				),
			)
		}
	}

	if len(function.Arguments) > 0 {
		if lastArgument := function.Arguments[len(function.Arguments)-1]; lastArgument.IsVariadic {
			if len(*args)-parametersLen > 0 {
				parametersLen = len(*args)
				for k := c; k < parametersLen; k++ {
					arg := (*args)[k]
					argPrototype := arg.(ObjectInstance).GetPrototype()
					if argPrototype == Nil {
						if lastArgument.Type != Nil && !lastArgument.IsNullable {
							return util.RuntimeError(
								fmt.Sprintf(
									"аргумент '%s' очікує ненульовий параметр, отримано '%s'",
									lastArgument.Name, arg.String(ctx),
								),
							)
						}
					} else if lastArgument.Type != nil && argPrototype != lastArgument.Type {
						return util.RuntimeError(
							fmt.Sprintf(
								"аргумент '%s' очікує список параметрів з типом '%s', отримано '%s'",
								lastArgument.Name,
								lastArgument.GetTypeName(),
								argPrototype.GetTypeName(),
							),
						)
					}
				}
			}
		}
	}

	return nil
}

func boolToInt64(v bool) int64 {
	if v {
		return 1
	}

	return 0
}

func boolToFloat64(v bool) float64 {
	if v {
		return 1.0
	}

	return 0.0
}

func getAttributes(attributes map[string]common.Type) (DictionaryInstance, error) {
	dict := NewDictionaryInstance()
	for key, val := range attributes {
		err := dict.SetElement(NewStringInstance(key), val)
		if err != nil {
			return DictionaryInstance{}, err
		}
	}

	return dict, nil
}

func getLength(ctx common.Context, sequence common.Type) (int64, error) {
	switch self := sequence.(type) {
	case ListInstance:
		return self.Length(ctx), nil
	case DictionaryInstance:
		return self.Length(), nil
	}

	return 0, errors.New(fmt.Sprint("invalid type in length operator: ", sequence.GetTypeName()))
}

func mergeAttributes(a map[string]common.Type, b ...map[string]common.Type) map[string]common.Type {
	for _, m := range b {
		for key, val := range m {
			a[key] = val
		}
	}

	return a
}

func newBinaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.Context, common.Type, common.Type) (common.Type, error),
) *FunctionInstance {
	return NewFunctionInstance(
		name,
		[]FunctionArgument{
			{
				Type:       selfType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: true,
			},
		},
		func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			return handler(ctx, (*args)[0], (*args)[1])
		},
		[]FunctionReturnType{
			{
				Type:       returnType,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func newUnaryMethod(
	name string,
	selfType *Class,
	returnType *Class,
	doc string,
	handler func(common.Context, common.Type) (common.Type, error),
) *FunctionInstance {
	return NewFunctionInstance(
		name,
		[]FunctionArgument{
			{
				Type:       selfType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			return handler(ctx, (*args)[0])
		},
		[]FunctionReturnType{
			{
				Type:       returnType,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func makeComparisonOperator(
	operator ops.Operator,
	itemType *Class,
	doc string,
	comparator func(common.Context, common.Type, common.Type) (int, error),
	checker func(res int) bool,
) *FunctionInstance {
	return newBinaryMethod(
		operator.Caption(),
		itemType,
		Bool,
		doc,
		func(ctx common.Context, self common.Type, other common.Type) (common.Type, error) {
			res, err := comparator(ctx, self, other)
			if err != nil {
				return nil, err
			}

			return NewBoolInstance(checker(res)), nil
		},
	)
}

func makeComparisonOperators(
	itemType *Class,
	comparator func(common.Context, common.Type, common.Type) (int, error),
) map[string]common.Type {
	return map[string]common.Type{
		ops.EqualsOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.EqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0
			},
		),
		ops.NotEqualsOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.NotEqualsOp, itemType, "", comparator, func(res int) bool {
				return res != 0
			},
		),
		ops.GreaterOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.GreaterOp, itemType, "", comparator, func(res int) bool {
				return res == 1
			},
		),
		ops.GreaterOrEqualsOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.GreaterOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == 1
			},
		),
		ops.LessOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.LessOp, itemType, "", comparator, func(res int) bool {
				return res == -1
			},
		),
		ops.LessOrEqualsOp.Caption(): makeComparisonOperator(
			// TODO: add doc
			ops.LessOrEqualsOp, itemType, "", comparator, func(res int) bool {
				return res == 0 || res == -1
			},
		),
	}
}

func makeLogicalOperators(itemType *Class) map[string]common.Type {
	return map[string]common.Type{
		ops.NotOp.Caption(): newUnaryMethod(
			// TODO: add doc
			ops.NotOp.Caption(),
			itemType,
			Bool,
			"",
			func(ctx common.Context, self common.Type) (common.Type, error) {
				return NewBoolInstance(!self.AsBool(ctx)), nil
			},
		),
		ops.AndOp.Caption(): newBinaryMethod(
			// TODO: add doc
			ops.AndOp.Caption(),
			itemType,
			Bool,
			"",
			func(ctx common.Context, self common.Type, other common.Type) (common.Type, error) {
				return NewBoolInstance(self.AsBool(ctx) && other.AsBool(ctx)), nil
			},
		),
		ops.OrOp.Caption(): newBinaryMethod(
			// TODO: add doc
			ops.OrOp.Caption(),
			itemType,
			Bool,
			"",
			func(ctx common.Context, self common.Type, other common.Type) (common.Type, error) {
				return NewBoolInstance(self.AsBool(ctx) || other.AsBool(ctx)), nil
			},
		),
	}
}

func makeCommonOperators(itemType *Class) map[string]common.Type {
	return map[string]common.Type{
		// TODO: add doc
		ops.BoolOperatorName: newUnaryMethod(
			ops.BoolOperatorName, itemType, Bool, "",
			func(ctx common.Context, self common.Type) (common.Type, error) {
				return NewBoolInstance(self.AsBool(ctx)), nil
			},
		),
	}
}

func newBuiltinConstructor(
	itemType *Class,
	handler func(common.Context, ...common.Type) (common.Type, error),
	doc string,
) *FunctionInstance {
	return NewFunctionInstance(
		ops.ConstructorName,
		[]FunctionArgument{
			{
				Type:       itemType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "значення",
				IsVariadic: true,
				IsNullable: true,
			},
		},
		func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			self, err := handler(ctx, (*args)[1:]...)
			if err != nil {
				return nil, err
			}

			(*args)[0] = self
			return NewNilInstance(), nil
		},
		[]FunctionReturnType{
			{
				Type:       Nil,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}

func newLengthOperator(
	itemType *Class,
	handler func(common.Context, common.Type) (int64, error),
	doc string,
) *FunctionInstance {
	return NewFunctionInstance(
		ops.LengthOperatorName,
		[]FunctionArgument{
			{
				Type:       itemType,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(ctx common.Context, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
			length, err := handler(ctx, (*args)[0])
			if err != nil {
				return nil, err
			}

			return NewIntegerInstance(length), nil
		},
		[]FunctionReturnType{
			{
				Type:       Integer,
				IsNullable: false,
			},
		},
		true,
		nil,
		doc,
	)
}
