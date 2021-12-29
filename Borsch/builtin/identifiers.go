package builtin

const (
	ConstantKeywordId = iota
	KeywordId
	FunctionId
	TypeId
)

var RegisteredIdentifiers = map[string]int{
	// Keyword names
	"пакет":  ConstantKeywordId,
	"істина": ConstantKeywordId,
	"хиба":   ConstantKeywordId,
	"якщо":   KeywordId,
	"інакше": KeywordId,
	"для":    KeywordId,
	"клас":   KeywordId,

	// Types + cast functions
	"абиякий":  TypeId,
	"рядок":    TypeId,
	"цілий":    TypeId,
	"дійсний":  TypeId,
	"логічний": TypeId,
	"список":   TypeId,
	"словник":  TypeId,

	// Functions
	"друк":            FunctionId,
	"друкр":           FunctionId,
	"ввід":            FunctionId,
	"середовище":      FunctionId,
	"паніка":          FunctionId,
	"довжина":         FunctionId,
	"вихід":           FunctionId,
	"додати":          FunctionId,
	"підтвердити":     FunctionId,
	"авторське_право": FunctionId,
	"ліцензія":        FunctionId,
	"допомога":        FunctionId,
	"тип":             FunctionId,
}
