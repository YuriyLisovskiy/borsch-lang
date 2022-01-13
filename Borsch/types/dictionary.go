package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type DictionaryEntry struct {
	Key   Type
	Value Type
}

type DictionaryInstance struct {
	Object
	Map map[uint64]DictionaryEntry
}

func NewDictionaryInstance() DictionaryInstance {
	return DictionaryInstance{
		Map: map[uint64]DictionaryEntry{},
		Object: Object{
			typeName:    GetTypeName(DictionaryTypeHash),
			Attributes:  nil,
			callHandler: nil,
		},
	}
}

func (t DictionaryInstance) calcHash(obj interface{}) (uint64, error) {
	h := sha256.New()
	_, err := h.Write([]byte(fmt.Sprintf("%v", obj)))
	if err != nil {
		return 0, err
	}

	return binary.BigEndian.Uint64(h.Sum(nil)), nil
}

func (t DictionaryInstance) String() string {
	return t.Representation()
}

func (t DictionaryInstance) Representation() string {
	var strValues []string
	for _, value := range t.Map {
		strValues = append(
			strValues, fmt.Sprintf(
				"%s: %s", value.Key.Representation(), value.Value.Representation(),
			),
		)
	}

	return "{" + strings.Join(strValues, ", ") + "}"
}

func (t DictionaryInstance) GetTypeHash() uint64 {
	return t.GetClass().GetTypeHash()
}

func (t DictionaryInstance) AsBool() bool {
	return t.Length() != 0
}

func (t DictionaryInstance) SetAttribute(name string, _ Type) (Type, error) {
	if t.Object.HasAttribute(name) || t.GetClass().HasAttribute(name) {
		return nil, util.AttributeIsReadOnlyError(t.GetTypeName(), name)
	}

	return nil, util.AttributeNotFoundError(t.GetTypeName(), name)
}

func (t DictionaryInstance) GetAttribute(name string) (Type, error) {
	if attribute, err := t.Object.GetAttribute(name); err == nil {
		return attribute, nil
	}

	return t.GetClass().GetAttribute(name)
}

func (t DictionaryInstance) GetClass() *Class {
	return Dictionary
}

func (t DictionaryInstance) Length() int64 {
	return int64(len(t.Map))
}

func (t DictionaryInstance) GetElement(key Type) (Type, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	if value, ok := t.Map[keyHash]; ok {
		return value.Value, nil
	}

	return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String()))
}

func (t *DictionaryInstance) SetElement(key Type, value Type) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	t.Map[keyHash] = DictionaryEntry{Key: key, Value: value}
	return nil
}

func (t *DictionaryInstance) RemoveElement(key Type) (Type, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	value, ok := t.Map[keyHash]
	if !ok {
		return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String()))
	}

	delete(t.Map, keyHash)
	return value.Value, nil
}

func compareDictionaries(self Type, other Type) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case *DictionaryInstance, DictionaryInstance:
		return -2, util.RuntimeError(
			fmt.Sprintf(
				"непідтримувані типи операндів для оператора %s: '%s' і '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	default:
		return -2, errors.New(
			fmt.Sprintf(
				"неможливо застосувати оператор %s до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newDictionaryClass() *Class {
	attributes := mergeAttributes(
		map[string]Type{
			// TODO: add doc
			ops.ConstructorName: newBuiltinConstructor(DictionaryTypeHash, ToDictionary, ""),

			// TODO: add doc
			ops.LengthOperatorName: newLengthOperator(ListTypeHash, getLength, ""),
			"вилучити": NewFunctionInstance(
				"вилучити",
				[]FunctionArgument{
					{
						TypeHash:   DictionaryTypeHash,
						Name:       "я",
						IsVariadic: false,
						IsNullable: false,
					},
					{
						TypeHash:   AnyTypeHash,
						Name:       "ключ",
						IsVariadic: false,
						IsNullable: true,
					},
				},
				func(args *[]Type, _ *map[string]Type) (Type, error) {
					dict := (*args)[0].(DictionaryInstance)
					value, err := dict.RemoveElement((*args)[1])
					if err != nil {
						return nil, util.RuntimeError(err.Error())
					}

					return value, nil
				},
				[]FunctionReturnType{
					{
						TypeHash:   AnyTypeHash,
						IsNullable: false,
					},
				},
				true,
				nil,
				"", // TODO: add doc
			),
		},
		makeLogicalOperators(DictionaryTypeHash),
		makeComparisonOperators(DictionaryTypeHash, compareDictionaries),
	)
	return NewBuiltinClass(
		DictionaryTypeHash,
		BuiltinPackage,
		attributes,
		"", // TODO: add doc
		func() (Type, error) {
			return NewDictionaryInstance(), nil
		},
	)
}
