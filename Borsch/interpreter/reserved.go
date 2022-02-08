package interpreter

import (
	"sort"

	"github.com/YuriyLisovskiy/borsch-lang/Borsch/builtin"
)

var keywords = []string{
	"заключний",
	"клас",
	"нуль",
	"перервати",
	"повернути",
	"функція",
	"хиба",
	"цикл",
	"якщо",
	"інакше",
	"істина",
}

var builtinIds []string

func binSearchString(arr []string, start, end int, item string) int {
	if end < start {
		return -1
	}

	mid := start + (end-start)/2
	if arr[mid] == item {
		return mid
	}

	if item < arr[mid] {
		return binSearchString(arr, start, mid-1, item)
	}

	return binSearchString(arr, mid+1, end, item)
}

func isKeyword(word string) bool {
	return binSearchString(keywords, 0, len(keywords), word) != -1
}

func isBuiltin(name string) bool {
	return binSearchString(builtinIds, 0, len(builtinIds), name) != -1
}

func init() {
	if !sort.StringsAreSorted(keywords) {
		sort.Strings(keywords)
	}

	for key := range builtin.GlobalScope {
		builtinIds = append(builtinIds, key)
	}

	if !sort.StringsAreSorted(builtinIds) {
		sort.Strings(builtinIds)
	}
}
