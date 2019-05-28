package utils

func Includes(a []string, t string) bool {
	for _, s := range a {
		if s == t {
			return true
		}
	}
	return false
}

func TakeProps(a []map[string]interface{}, pname string) (props []interface{}) {
	for _, item := range a {
		props = append(props, item[pname])
	}
	return
}
