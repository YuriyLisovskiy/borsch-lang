package types

import "errors"

const (
	AnyTypeHash = iota
	NilTypeHash
	RealTypeHash
	IntegerTypeHash
	StringTypeHash
	BoolTypeHash
	ListTypeHash
	DictionaryTypeHash
	PackageTypeHash
	FunctionTypeHash
)

type ValueType interface {
	String() string
	Representation() string
	TypeHash() int
	TypeName() string
	GetAttr(string) (ValueType, error)
	SetAttr(string, ValueType) (ValueType, error)
	CompareTo(ValueType) (int, error)
}

type SequentialType interface {
	Length() int64
	GetElement(int64) (ValueType, error)
	SetElement(int64, ValueType) (ValueType, error)
	Slice(int64, int64) (ValueType, error)
}

func getIndex(index, length int64) (int64, error) {
	if index >= 0 && index < length {
		return index, nil
	} else if index < 0 && index >= -length {
		return length + index, nil
	}

	return 0, errors.New("індекс за межами послідовності")
}

func GetTypeName(typeValue int) string {
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

func GetTypeHash(typeName string) int {
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
		return -1
	}
}

func IsBuiltinType(typeName string) bool {
	return GetTypeHash(typeName) != -1
}
