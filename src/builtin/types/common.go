package types

import "errors"

const (
	noneType = iota
	realType
	integerType
	stringType
	boolType
	listType
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
}

func getIndex(index, length int64) (int64, error) {
	if index >= 0 && index < length {
		return index, nil
	} else if index < 0 && index >= -length {
		return length + index, nil
	}

	return 0, errors.New("індекс рядка за межами послідовності")
}

func GetTypeName(typeValue int) string {
	switch typeValue {
	case noneType:
		return "без типу"
	case realType:
		return "дійсний"
	case integerType:
		return "цілий"
	case stringType:
		return "рядок"
	case boolType:
		return "логічний"
	case listType:
		return "список"
	default:
		return "невідомий"
	}
}
