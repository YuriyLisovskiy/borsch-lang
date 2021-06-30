package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/YuriyLisovskiy/borsch/Borsch/util"
	"strings"
)

type DictionaryEntry struct {
	Key   ValueType
	Value ValueType
}

type DictionaryType struct {
	Map map[uint64]DictionaryEntry
}

func NewDictionaryType() DictionaryType {
	return DictionaryType{
		Map: map[uint64]DictionaryEntry{},
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

func (t DictionaryType) TypeHash() int {
	return DictionaryTypeHash
}

func (t DictionaryType) TypeName() string {
	return GetTypeName(t.TypeHash())
}

func (t DictionaryType) AsBool() bool {
	return t.Length() != 0
}

func (t DictionaryType) Length() int64 {
	return int64(len(t.Map))
}

func (t DictionaryType) GetElement(key ValueType) (ValueType, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	if value, ok := t.Map[keyHash]; ok {
		return value.Value, nil
	}

	return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String()))
}

func (t *DictionaryType) SetElement(key ValueType, value ValueType) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	t.Map[keyHash] = DictionaryEntry{Key: key, Value: value}
	return nil
}

func (t *DictionaryType) RemoveElement(key ValueType) error {
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

func (t DictionaryType) GetAttr(name string) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t DictionaryType) SetAttr(name string, _ ValueType) (ValueType, error) {
	return nil, util.AttributeError(t.TypeName(), name)
}

func (t DictionaryType) Pow(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Plus() (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Minus() (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseNot() (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Mul(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Div(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Mod(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Add(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) Sub(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseLeftShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseRightShift(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseAnd(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseXor(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) BitwiseOr(ValueType) (ValueType, error) {
	return nil, nil
}

func (t DictionaryType) CompareTo(other ValueType) (int, error) {
	switch right := other.(type) {
	case NilType:
	case DictionaryType:
		return -2, util.RuntimeError(fmt.Sprintf(
			"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	default:
		return -2, errors.New(fmt.Sprintf(
			"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
			"%s", t.TypeName(), right.TypeName(),
		))
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func (t DictionaryType) Not() (ValueType, error) {
	return BoolType{Value: !t.AsBool()}, nil
}

func (t DictionaryType) And(other ValueType) (ValueType, error) {
	return logicalAnd(t, other)
}

func (t DictionaryType) Or(other ValueType) (ValueType, error) {
	return logicalOr(t, other)
}
