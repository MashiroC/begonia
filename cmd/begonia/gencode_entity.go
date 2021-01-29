package main

import (
	"html/template"
	"os"
	"strings"
)

func genEntity(objs map[string][]string) {
	for k, v := range objs {
		tmpl, err := template.New("entity").Parse(entityTmplStr)
		if err != nil {
			panic(err)
		}
		f, err := os.Create(strings.Join([]string{k, "call", "entity.begonia.go"}, string(os.PathSeparator)))
		if err != nil {
			panic(err)
		}
		err = tmpl.Execute(f, map[string]interface{}{
			"objs": v,
		})
		if err != nil {
			panic(err)
		}
		//fmt.Println(k)
		//fmt.Println(v)
	}
}
