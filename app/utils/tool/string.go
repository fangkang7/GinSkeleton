package tool

import (
	"bytes"
	"strconv"
	"strings"
	"unicode/utf8"
)

// 字符串拼接
func StringBuild(temstr ...string) string {
	var strBuilder strings.Builder
	for _, str := range temstr {
		strBuilder.WriteString(str)
	}
	return strBuilder.String()
}

// Implode implode()
func Implode(glue string, pieces []string) string {
	var buf bytes.Buffer
	l := len(pieces)
	for _, str := range pieces {
		buf.WriteString(str)
		if l--; l > 0 {
			buf.WriteString(glue)
		}
	}
	return buf.String()
}

// Chr chr()
func Chr(ascii int) string {
	return string(rune(ascii))
}

// Ord ord()
func Ord(char string) int {
	r, _ := utf8.DecodeRune([]byte(char))
	return int(r)
}

// stripslashes 类似于 PHP 的 stripslashes 函数，用于去除字符串中的反斜杠转义字符
func Stripslashes(s string) string {
	res, err := strconv.Unquote(`"` + s + `"`)
	if err != nil {
		return s
	}
	// res = strings.ReplaceAll(res, `\n`, "\n")
	// res = strings.ReplaceAll(res, `\r`, "\r")
	// res = strings.ReplaceAll(res, `\t`, "\t")
	return res
}
