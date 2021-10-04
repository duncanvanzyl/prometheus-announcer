package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"text/template"

	"github.com/kelseyhightower/envconfig"
)

// Flag usage messages
const (
	idUsage = `Unique ID for this announcer.
If not provided a random uuid will be used.`
	labelUsage = `Labels to tag metrics scraped from the supplied targets with.
Must be in the form: "name1:value1,name2:value2".
Label names may contain ASCII letters, numbers, as well as underscores and must 
match the regex "[a-zA-Z_][a-zA-Z0-9_]*".
Label values may contain any Unicode characters.`
	envFormat = `
The following environment variables can be used:
{{range .}}  {{.Key}}
  	{{usage_description .Tags}}
  	  [type]        {{usage_type .Field}}
  	  [default]     {{usage_default .Tags}}
{{end}}
`
)

func formatDescription(s string) string {
	indent := "\n    \t"
	if ss := strings.Split(s, "\n"); len(ss) == 1 {
		return s
	}

	return strings.ReplaceAll(s, "\n", indent)
}

func Usage(prefix string, s interface{}) func() {
	return func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [options] [target1 target2 ... targetN]:\n", os.Args[0])

		fmt.Fprintf(os.Stderr, "\nThe following flags can be used:\n")
		flag.PrintDefaults()

		tmpl, err := envTemplate(envFormat)
		if err != nil {
			log.Fatalf("could not create template: %v", err)
		}
		envconfig.Usaget(prefix, s, os.Stderr, tmpl)
	}
}

func envTemplate(format string) (*template.Template, error) {
	functions := template.FuncMap{
		"desc":              func(v interface{}) string { return fmt.Sprintf("%+v", v) },
		"usage_description": func(v reflect.StructTag) string { return formatDescription(v.Get("desc")) },
		"usage_type":        func(v reflect.Value) string { return strings.Title(v.Type().Name()) },
		"usage_default":     func(v reflect.StructTag) string { return v.Get("default") },
	}

	return template.New("envconfig").Funcs(functions).Parse(format)
}
