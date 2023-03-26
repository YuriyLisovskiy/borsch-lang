package interpreter

import "sort"

var keywords = []string{
	"блок",
	"заключний",
	"клас",
	"кінець",
	"лямбда",
	"небезпечно",
	"нуль",
	"панікувати",
	"перервати",
	"повернути",
	"піймати",
	"функція",
	"хиба",
	"цикл",
	"якщо",
	"інакше",
	"істина",
}

var builtinIds []string

func binSearchString(arr []string, low, high int, item string) int {
	if high < low {
		return -1
	}

	mid := low + (high-low)/2
	if arr[mid] == item {
		return mid
	}

	if item < arr[mid] {
		return binSearchString(arr, low, mid-1, item)
	}

	return binSearchString(arr, mid+1, high, item)
}

func isKeyword(word string) bool {
	return binSearchString(keywords, 0, len(keywords)-1, word) != -1
}

func init() {
	if !sort.StringsAreSorted(keywords) {
		sort.Strings(keywords)
	}

	for key := range GlobalScope {
		builtinIds = append(builtinIds, key)
	}

	if !sort.StringsAreSorted(builtinIds) {
		sort.Strings(builtinIds)
	}
}
