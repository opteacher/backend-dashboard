package utils

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
