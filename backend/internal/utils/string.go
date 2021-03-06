package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
)

func ToSingular(word string) string {
	l := len(word)
	switch {
	case l > 3 && word[l-3:] == "ies":
		return word[:l-3] + "y"
	case l > 3 && word[l-3:] == "ves":
		return word[:l-3] + "f"
	case l > 3 && word[l-3:] == "ses":
		return word[:l-2]
	case l > 1 && word[l-1:] == "s":
		return word[:l-1]
	default:
		return word
	}
}

func ToPlural(word string) string {
	l := len(word)
	if l == 0 {
		return ""
	}
	switch word[l-1] {
	case 'y':
		return word[:l-1] + "ies"
	case 'f':
		return word[:l-1] + "ves"
	case 's':
		return word + "es"
	default:
		return word + "s"
	}
}

func CamelToPascal(word string) string {
	str := ""
	for _, w := range word {
		if w >= 'A' && w <= 'Z' {
			str += fmt.Sprintf("_%s", string(w+32))
		} else {
			str += string(w)
		}
	}
	str = strings.TrimLeft(str, "_")
	return str
}

func PascalToCamel(word string) string {
	str := ""
	flag := false
	for _, w := range word {
		if w == '_' {
			flag = true
		} else if flag {
			flag = false
			if w >= 'a' && w <= 'z' {
				str += string(w-32)
			} else {
				str += string(w)
			}
		} else {
			str += string(w)
		}
	}
	return str
}

func Md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

func Capital(str string) string {
	if str[0] >= 'a' && str[0] <= 'z' {
		return string(str[0]-32) + str[1:]
	} else {
		return str
	}
}

func Uncapital(str string) string {
	if str[0] >= 'A' && str[0] <= 'Z' {
		return string(str[0]+32) + str[1:]
	} else {
		return str
	}
}

func AddSpacesBeforeRow(data string, n int) string {
	spaces := ""
	for i := 0; i < n; i++ {
		spaces += "\t"
	}
	ret := spaces + strings.Replace(data, "\n", "\n" + spaces, -1)
	return strings.TrimRight(ret, "\t")
}
