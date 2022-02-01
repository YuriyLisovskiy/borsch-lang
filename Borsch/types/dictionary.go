package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/ops"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/util"
)

type DictionaryEntry struct {
	Key   common.Type
	Value common.Type
}

type DictionaryInstance struct {
	BuiltinObject
	Map map[uint64]DictionaryEntry
}

func NewDictionaryInstance() DictionaryInstance {
	return DictionaryInstance{
		Map: map[uint64]DictionaryEntry{},
		BuiltinObject: BuiltinObject{
			CommonObject{
				Object: Object{
					typeName:    common.DictionaryTypeName,
					Attributes:  nil,
					callHandler: nil,
				},
				prototype: Dictionary,
			},
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

func (t DictionaryInstance) String(state common.State) string {
	return t.Representation(state)
}

func (t DictionaryInstance) Representation(state common.State) string {
	var strValues []string
	for _, value := range t.Map {
		strValues = append(
			strValues, fmt.Sprintf(
				"%s: %s", value.Key.Representation(state), value.Value.Representation(state),
			),
		)
	}

	return "{" + strings.Join(strValues, ", ") + "}"
}

func (t DictionaryInstance) AsBool(state common.State) bool {
	return t.Length(state) != 0
}

func (t DictionaryInstance) Length(common.State) int64 {
	return int64(len(t.Map))
}

func (t DictionaryInstance) GetElement(state common.State, key common.Type) (common.Type, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	if value, ok := t.Map[keyHash]; ok {
		return value.Value, nil
	}

	return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String(state)))
}

func (t *DictionaryInstance) SetElement(key common.Type, value common.Type) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	t.Map[keyHash] = DictionaryEntry{Key: key, Value: value}
	return nil
}

func (t *DictionaryInstance) RemoveElement(state common.State, key common.Type) (common.Type, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	value, ok := t.Map[keyHash]
	if !ok {
		return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", key.String(state)))
	}

	delete(t.Map, keyHash)
	return value.Value, nil
}

func compareDictionaries(_ common.State, self common.Type, other common.Type) (int, error) {
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
				"неможливо застосувати оператор '%s' до значень типів '%s' та '%s'",
				"%s", self.GetTypeName(), right.GetTypeName(),
			),
		)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func newDictionaryClass() *Class {
	initAttributes := func() map[string]common.Type {
		return mergeAttributes(
			map[string]common.Type{
				// TODO: add doc
				ops.ConstructorName: newBuiltinConstructor(Dictionary, ToDictionary, ""),

				// TODO: add doc
				ops.LengthOperatorName: newLengthOperator(List, getLength, ""),
				"вилучити": NewFunctionInstance(
					"вилучити",
					[]FunctionArgument{
						{
							Type:       Dictionary,
							Name:       "я",
							IsVariadic: false,
							IsNullable: false,
						},
						{
							Type:       nil,
							Name:       "ключ",
							IsVariadic: false,
							IsNullable: true,
						},
					},
					func(state common.State, args *[]common.Type, _ *map[string]common.Type) (common.Type, error) {
						dict := (*args)[0].(DictionaryInstance)
						value, err := dict.RemoveElement(state, (*args)[1])
						if err != nil {
							return nil, util.RuntimeError(err.Error())
						}

						return value, nil
					},
					[]FunctionReturnType{
						{
							Type:       Any,
							IsNullable: false,
						},
					},
					true,
					nil,
					"", // TODO: add doc
				),
			},
			makeLogicalOperators(Dictionary),
			makeComparisonOperators(Dictionary, compareDictionaries),
			makeCommonOperators(Dictionary),
		)
	}

	return NewBuiltinClass(
		common.DictionaryTypeName,
		BuiltinPackage,
		initAttributes,
		"", // TODO: add doc
		func() (common.Type, error) {
			return NewDictionaryInstance(), nil
		},
	)
}
