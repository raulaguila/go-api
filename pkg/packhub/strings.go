package packhub

import "strings"

func Capitalize(s string) string {
	res := strings.TrimSpace(s)
	return strings.ToUpper(res[:1]) + strings.ToLower(res[1:])
}
