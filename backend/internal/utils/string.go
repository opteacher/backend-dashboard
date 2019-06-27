package utils

import (
	"fmt"
	"strings"
	"crypto/md5"
    "encoding/hex"
)

func ToSingular(word string) string {
	l := len(word)
	switch {
	case l > 3 && word[l-3:] == "ies":
		return word[:l-3] + "y"
	case l > 3 && word[l-3:] == "ves":
		return word[:l-3] + "f"
	case l > 2 && word[l-2:] == "es":
		return word[:l-2]
	case l > 1 && word[l-1:] == "s":
		return word[:l-1]
	default:
		return word
	}
}

func CamelToPascal(word string) string {
	str := ""
	for _, w := range word {
		if w >= 65 && w <= 90 {
			str += fmt.Sprintf("_%s", string(w+32))
		} else {
			str += string(w)
		}
	}
	str = strings.TrimLeft(str, "_")
	return str
}

func Md5Hex(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}