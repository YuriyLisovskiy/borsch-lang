package types

type NoneType struct {
}

func (t NoneType) String() string {
	return "NoneType{" + t.Representation() + "}"
}

func (t NoneType) Representation() string {
	return "Порожнеча"
}

func (t NoneType) TypeHash() int {
	return noneType
}

func (t NoneType) TypeName() string {
	return GetTypeName(t.TypeHash())
}
