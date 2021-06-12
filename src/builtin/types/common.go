package types

import "errors"

const (
	NoneTypeHash = iota
	RealTypeHash
	IntegerTypeHash
	StringTypeHash
	BoolTypeHash
	ListTypeHash
	DictionaryTypeHash
)

type ValueType interface {
	String() string
	Representation() string
	TypeHash() int
	TypeName() string
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
	case NoneTypeHash:
		return "без типу"
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
	default:
		return "невідомий"
	}
}
