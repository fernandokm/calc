package lexer

import (
	"fmt"
	"strings"
	"unicode"
)

var (
	digit      = unicode.IsNumber
	underscore = "_"
	letter     = unicode.IsLetter
	signSymbol = "+-"
	boolean    = stringSet("true", "false")
	operator   = stringSet("+", "-", "*", "/", "%", "^", "!",
		"==", "!=", ">", ">=", "<", "<=", "&&", "||", "|",
		"=", "∑", "∏", "∪", "∩", "∁", "⊂", "⊄", "⊆", "⊈",
		"∈", "∉", "∀", "∃", "∄", "∃", "∃!", "∄!")
	whitespace     = " \t"
	newLine        = "\r\n"
	punctuation    = "(){},.[]"
	exponentSymbol = "eE"
	decimalPoint   = "."
)

var ruleMaxLengths = make(map[interface{}]int)

func stringSet(values ...string) map[string]struct{} {
	set := make(map[string]struct{})
	for _, val := range values {
		set[val] = struct{}{}
	}
	return set
}

func is(value interface{}, rules ...interface{}) bool {
	switch val := value.(type) {
	case rune:
		return runeIs(val, rules...)
	case uint8:
		return runeIs(rune(val), rules...)
	case string:
		return stringIs(val, rules...)
	default:
		panic("is expected a rune, uint8 or string")
	}
}

func runeIs(r rune, rules ...interface{}) bool {
	for _, ruleInterface := range rules {
		switch rule := ruleInterface.(type) {
		case string:
			if strings.IndexRune(rule, r) >= 0 {
				return true
			}
		case func(rune) bool:
			if rule(r) {
				return true
			}
		case map[string]struct{}:
			for k := range rule {
				if rune(k[0]) == r {
					return true
				}
			}
		default:
			panic(fmt.Sprintf("runeIs expected a string, func(rune) bool or map[string]struct{}, got %T", ruleInterface))
		}
	}
	return false
}

func stringIs(str string, rules ...interface{}) bool {
	for _, ruleInterface := range rules {
		switch rule := ruleInterface.(type) {
		case string:
			if strings.Index(rule, str) >= 0 {
				return true
			}
		case map[string]struct{}:
			if _, ok := rule[str]; ok {
				return true
			}
		default:
			panic(fmt.Sprintf("stringIs expected a string or map[string]struct{}, got %T", ruleInterface))
		}
	}
	return false
}

func stringRuleLength(str string, startIndex int, rules ...interface{}) int {
	maxLength := 0
	for _, rule := range rules {
		ruleMax := getRuleMaxLength(rule)
		if ruleMax > maxLength {
			maxLength = ruleMax
		}
	}
	if maxLength+startIndex > len(str) {
		maxLength = len(str) - startIndex
	}

	for length := maxLength; length > 0; length-- {
		if stringIs(str[startIndex:startIndex+length], rules...) {
			return length
		}
	}
	return 0
}

func getRuleMaxLength(rule interface{}) int {
	key := fmt.Sprintf("%+v", rule)
	if max, ok := ruleMaxLengths[key]; ok {
		return max
	}
	switch r := rule.(type) {
	case string:
		ruleMaxLengths[key] = len(r)
	case map[string]struct{}:
		max := 0
		for k := range r {
			if len(k) > max {
				max = len(k)
			}
		}
		ruleMaxLengths[key] = max
	default:
		panic(fmt.Sprintf("getRuleMaxLength expected string or map[string]struct{}, got %T", rule))
	}
	return ruleMaxLengths[key]
}
