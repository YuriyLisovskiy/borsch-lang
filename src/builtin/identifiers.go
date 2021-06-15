package builtin

const (
	ConstantKeywordId = iota
	KeywordId
	FunctionId
)

var RegisteredIdentifiers = map[string]int{
	// Keyword names
	"пакет": ConstantKeywordId,
	"істина": ConstantKeywordId,
	"хиба": ConstantKeywordId,
	"якщо": KeywordId,
	"інакше": KeywordId,
	"для": KeywordId,

	// Types + cast functions
	"рядок": FunctionId,
	"цілий": FunctionId,
	"дійсний": FunctionId,
	"логічний": FunctionId,
	"список": FunctionId,
	"словник": FunctionId,

	// Functions
	"друк": FunctionId,
	"друкр": FunctionId,
	"ввід": FunctionId,
	"середовище": FunctionId,
	"паніка": FunctionId,
	"довжина": FunctionId,
	"вихід": FunctionId,
	"додати": FunctionId,
	"вилучити": FunctionId,
	"підтвердити": FunctionId,
}
