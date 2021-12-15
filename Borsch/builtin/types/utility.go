package types

import (
	"errors"
	"hash/fnv"
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
		return "абиякий"
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
	default:
		return "невідомий"
	}
}

func GetTypeHash(typeName string) uint64 {
	switch typeName {
	case "абиякий":
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
	default:
		return FunctionTypeHash + 1
	}
}

func IsBuiltinType(typeName string) bool {
	return GetTypeHash(typeName) <= FunctionTypeHash
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

func logicalAnd(l, r Type) (Type, error) {
	return BoolType{Value: l.AsBool() && r.AsBool()}, nil
}

func logicalOr(l, r Type) (Type, error) {
	return BoolType{Value: l.AsBool() || r.AsBool()}, nil
}

func newBinaryMethod(
	name string,
	parentTypeHash uint64,
	doc string,
	handler func(Type, Type) (Type, error),
) FunctionType {
	return NewFunctionType(
		name,
		[]FunctionArgument{
			{
				TypeHash:   parentTypeHash,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				TypeHash:   0,
				Name:       "інший",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []Type, _ map[string]Type) (Type, error) {
			return handler(args[0], args[1])
		},
		FunctionReturnType{
			TypeHash:   parentTypeHash,
			IsNullable: false,
		},
		nil,
		doc,
	)
}

func newUnaryMethod(
	name string,
	parentTypeHash uint64,
	doc string,
	handler func(Type) (Type, error),
) FunctionType {
	return NewFunctionType(
		name,
		[]FunctionArgument{
			{
				TypeHash:   parentTypeHash,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(args []Type, _ map[string]Type) (Type, error) {
			return handler(args[0])
		},
		FunctionReturnType{
			TypeHash:   parentTypeHash,
			IsNullable: false,
		},
		nil,
		doc,
	)
}
