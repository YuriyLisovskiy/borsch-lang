package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type DictionaryEntry struct {
	Key   Type
	Value Type
}

// DictionaryType TODO: move methods impl to attributes
type DictionaryType struct {
	Object

	Map map[uint64]DictionaryEntry
	package_ *PackageType
}

func NewDictionaryType() *DictionaryType {
	return &DictionaryType{
		Map: map[uint64]DictionaryEntry{},
		Object: *newBuiltinObject(
			DictionaryTypeHash, map[string]Type{
				"__документ__": &NilType{}, // TODO: set doc
				"__пакет__":    BuiltinPackage,
			},
		),
		package_: BuiltinPackage,
	}
}

func (t DictionaryType) calcHash(obj interface{}) (uint64, error) {
	h := sha256.New()
	_, err := h.Write([]byte(fmt.Sprintf("%v", obj)))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(h.Sum(nil)), nil
}

func (t DictionaryType) String() string {
	return t.Representation()
}

func (t DictionaryType) Representation() string {
	var strValues []string
	for _, value := range t.Map {
		strValues = append(strValues, fmt.Sprintf(
			"%s: %s", value.Key.Representation(), value.Value.Representation(),
		))
	}

	return "{" + strings.Join(strValues, ", ") + "}"
}

func (t DictionaryType) AsBool() bool {
	return t.Length() != 0
}

func (t DictionaryType) Length() int64 {
	return int64(len(t.Map))
}

func (t DictionaryType) GetElement(key Type) (Type, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	if value, ok := t.Map[keyHash]; ok {
		return value.Value, nil
	}

	return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String()))
}

func (t *DictionaryType) SetElement(key Type, value Type) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	t.Map[keyHash] = DictionaryEntry{Key: key, Value: value}
	return nil
}

func (t *DictionaryType) RemoveElement(key Type) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	if _, ok := t.Map[keyHash]; !ok {
		return errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String()))
	}

	delete(t.Map, keyHash)
	return nil
}

func (t DictionaryType) SetAttribute(name string, _ Type) (Type, error) {
	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t DictionaryType) Pow(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) Plus() (Type, error) {
	return nil, nil
}

func (t DictionaryType) Minus() (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseNot() (Type, error) {
	return nil, nil
}

func (t DictionaryType) Mul(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) Div(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) Mod(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) Add(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) Sub(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseLeftShift(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseRightShift(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseAnd(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseXor(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseOr(Type) (Type, error) {
	return nil, nil
}

func (t DictionaryType) CompareTo(other Type) (int, error) {
	switch right := other.(type) {
	case NilType:
	case DictionaryType:
		return -2, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	default:
		return -2, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.GetTypeName(), right.GetTypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t DictionaryType) Not() (Type, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t DictionaryType) And(other Type) (Type, error) {
	return logicalAnd(t, other)
}

func (t DictionaryType) Or(other Type) (Type, error) {
	return logicalOr(t, other)
}
