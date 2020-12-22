package main

import "html/template"

func raw(str string) template.HTML {
	return template.HTML(str)
}

func concat(str ...string) string {
	res := ""
	for i := 0; i < len(str); i++ {
		res += str[i]
	}
	return res
}

func add(a int) int {
	return a + 1
}