package types

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"strings"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/common"
	"github.com/YuriyLisovskiy/borsch-lang/Borsch/utilities"
)

type DictionaryEntry struct {
	Key   common.Value
	Value common.Value
}

type DictionaryInstance struct {
	BuiltinInstance
	Map map[uint64]DictionaryEntry
}

func NewDictionaryInstance() DictionaryInstance {
	return DictionaryInstance{
		BuiltinInstance: BuiltinInstance{
			ClassInstance{
				class:      Dictionary,
				attributes: map[string]common.Value{},
				address:    "",
			},
		},
		Map: map[uint64]DictionaryEntry{},
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

func (t DictionaryInstance) String(state common.State) (string, error) {
	return t.Representation(state)
}

func (t DictionaryInstance) Representation(state common.State) (string, error) {
	var strValues []string
	for _, value := range t.Map {
		keyRepresentation, err := value.Key.Representation(state)
		if err != nil {
			return "", err
		}

		valueRepresentation, err := value.Value.Representation(state)
		if err != nil {
			return "", err
		}

		strValues = append(strValues, fmt.Sprintf("%s: %s", keyRepresentation, valueRepresentation))
	}

	return "{" + strings.Join(strValues, ", ") + "}", nil
}

func (t DictionaryInstance) AsBool(state common.State) (bool, error) {
	return t.Length(state) != 0, nil
}

func (t DictionaryInstance) Length(common.State) int64 {
	return int64(len(t.Map))
}

func (t DictionaryInstance) GetElement(state common.State, key common.Value) (common.Value, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	if value, ok := t.Map[keyHash]; ok {
		return value.Value, nil
	}

	keyStr, err := key.String(state)
	if err != nil {
		return nil, err
	}

	return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", keyStr))
}

func (t *DictionaryInstance) SetElement(key common.Value, value common.Value) error {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return err
	}

	t.Map[keyHash] = DictionaryEntry{Key: key, Value: value}
	return nil
}

func (t *DictionaryInstance) RemoveElement(state common.State, key common.Value) (common.Value, error) {
	keyHash, err := t.calcHash(key)
	if err != nil {
		return nil, err
	}

	value, ok := t.Map[keyHash]
	if !ok {
		keyStr, err := key.String(state)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(fmt.Sprintf("значення за ключем '%s' не існує", keyStr))
	}

	delete(t.Map, keyHash)
	return value.Value, nil
}

func toDictionary(state common.State, args ...common.Value) (common.Value, error) {
	dict := NewDictionaryInstance()
	if len(args) == 0 {
		return dict, nil
	}

	if len(args) != 2 {
		return nil, errors.New(
			fmt.Sprintf(
				"функція 'словник()' приймає два аргументи, або жодного (отримано %d)", len(args),
			),
		)
	}

	switch keys := args[0].(type) {
	case ListInstance:
		switch values := args[1].(type) {
		case ListInstance:
			if keys.Length(state) != values.Length(state) {
				return nil, errors.New(
					fmt.Sprintf(
						"довжина списку ключів має співпадати з довжиною списку значень",
					),
				)
			}

			length := keys.Length(state)
			for i := int64(0); i < length; i++ {
				err := dict.SetElement(keys.Values[i], values.Values[i])
				if err != nil {
					return nil, err
				}
			}

			return dict, nil
		default:
			return nil, errors.New(
				fmt.Sprintf(
					"функція 'словник()' другим аргументом приймає список значень",
				),
			)
		}
	default:
		return nil, errors.New(
			fmt.Sprintf(
				"функція 'словник()' першим аргументом приймає список ключів",
			),
		)
	}
}

func compareDictionaries(_ common.State, op common.Operator, self common.Value, other common.Value) (int, error) {
	switch right := other.(type) {
	case NilInstance:
	case *DictionaryInstance, DictionaryInstance:
		return -2, utilities.OperandsNotSupportedError(op, self.GetTypeName(), right.GetTypeName())
	default:
		return -2, utilities.OperatorNotSupportedError(op, self, right)
	}

	// -2 is something other than -1, 0 or 1 and means 'not equals'
	return -2, nil
}

func dictionaryMethod_Remove(name string) common.Value {
	return NewFunctionInstance(
		name,
		[]FunctionParameter{
			{
				Type:       Dictionary,
				Name:       "я",
				IsVariadic: false,
				IsNullable: false,
			},
			{
				Type:       Any,
				Name:       "ключ",
				IsVariadic: false,
				IsNullable: false,
			},
		},
		func(state common.State, args *[]common.Value, _ *map[string]common.Value) (common.Value, error) {
			dict := (*args)[0].(DictionaryInstance)
			_, err := dict.RemoveElement(state, (*args)[1])
			return nil, err
		},
		[]FunctionReturnType{
			{
				Type:       Nil,
				IsNullable: false,
			},
		},
		true,
		nil,
		"", // TODO: add doc
	)
}

func newDictionaryClass() *Class {
	return &Class{
		Name:    common.DictionaryTypeName,
		IsFinal: true,
		Bases:   []*Class{},
		Parent:  BuiltinPackage,
		AttrInitializer: func(attrs *map[string]common.Value) {
			*attrs = MergeAttributes(
				map[string]common.Value{
					// TODO: add doc
					common.ConstructorName: makeVariadicConstructor(Dictionary, toDictionary, ""),

					// TODO: add doc
					common.LengthOperatorName: makeLengthOperator(List, ""),

					"вилучити": dictionaryMethod_Remove("вилучити"),
				},
				MakeLogicalOperators(Dictionary),
				MakeComparisonOperators(Dictionary, compareDictionaries),
				MakeCommonOperators(Dictionary),
			)
		},
		GetEmptyInstance: func() (common.Value, error) {
			return NewDictionaryInstance(), nil
		},
	}
}
