package types

import (
	"errors"
	"fmt"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

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
	AsBool() bool
	GetAttr(string) (ValueType, error)
	SetAttr(string, ValueType) (ValueType, error)

	Pow(ValueType) (ValueType, error)

	Plus() (ValueType, error)
	Minus() (ValueType, error)
	BitwiseNot() (ValueType, error)

	Mul(ValueType) (ValueType, error)
	Div(ValueType) (ValueType, error)
	Mod(ValueType) (ValueType, error)

	Add(ValueType) (ValueType, error)
	Sub(ValueType) (ValueType, error)

	BitwiseLeftShift(ValueType) (ValueType, error)
	BitwiseRightShift(ValueType) (ValueType, error)

	BitwiseAnd(ValueType) (ValueType, error)

	BitwiseXor(ValueType) (ValueType, error)

	BitwiseOr(ValueType) (ValueType, error)

	CompareTo(ValueType) (int, error)

	Not() (ValueType, error)

	And(ValueType) (ValueType, error)

	Or(ValueType) (ValueType, error)
}

type SequentialType interface {
	Length() int64
	GetElement(int64) (ValueType, error)
	SetElement(int64, ValueType) (ValueType, error)
	Slice(int64, int64) (ValueType, error)
}

type ObjectType struct {
	typeHash int
	typeName string
	Attributes map[string]ValueType
}

func newObjectType(typeHash int, attributes map[string]ValueType) *ObjectType {
	return &ObjectType{
		typeHash:   typeHash,
		typeName:   GetTypeName(typeHash),
		Attributes: attributes,
	}
}

func (o ObjectType) makeAttributes() (*DictionaryType, error) {
	dict := NewDictionaryType()
	for key, val := range o.Attributes {
		err := dict.SetElement(NewStringType(key), val)
		if err != nil {
			return nil, err
		}
	}

	return dict, nil
}

func (o ObjectType) GetTypeHash() int {
	return o.typeHash
}

func (o ObjectType) GetTypeName() string {
	return o.typeName
}

func (o ObjectType) GetAttribute(name string) (ValueType, error) {
	if name == "__атрибути__" {
		return o.makeAttributes()
	}

	if val, ok := o.Attributes[name]; ok {
		return val, nil
	}

	return nil, util.AttributeError(o.GetTypeName(), name)
}

func (o ObjectType) SetAttribute(name string, value ValueType) error {
	if val, ok := o.Attributes[name]; ok {
		if val.TypeHash() == value.TypeHash() {
			o.Attributes[name] = value
			return nil
		}

		return util.RuntimeError(fmt.Sprintf(
			"неможливо записати значення типу '%s' у атрибут '%s' з типом '%s'",
			value.TypeName(), name, val.TypeName(),
		))
	}

	o.Attributes[name] = value
	return nil
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

func logicalAnd(l, r ValueType) (ValueType, error) {
	return BoolType{Value: l.AsBool() && r.AsBool()}, nil
}

func logicalOr(l, r ValueType) (ValueType, error) {
	return BoolType{Value: l.AsBool() || r.AsBool()}, nil
}

var BuiltinPackage = &PackageType{
	IsBuiltin:  true,
	Name:       "",
	Parent:     "",
	Object: newObjectType(PackageTypeHash, map[string]ValueType{}),
}
