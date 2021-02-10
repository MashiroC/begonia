package main

import (
	"html/template"
	"os"
	"strings"
)

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

func parseFullName(fullName string) (serviceName, dirPath string) {
	serviceName = fullName[strings.LastIndex(fullName, "_")+1:]
	dirPos := strings.LastIndex(fullName, string(os.PathSeparator))
	if dirPos == -1 {
		dirPath = "/"
	} else {
		dirPath = fullName[:dirPos]
	}
	return
}
