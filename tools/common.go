package tools

import (
	"gii/glog"
	"io"
	"unicode"
)

func Close(c io.Closer) {
	if err := c.Close(); err != nil {
		glog.Error(err)
	}
}

// 驼峰单词转下划线单词

func CamelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}

		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}
