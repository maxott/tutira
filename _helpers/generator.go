package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"io/ioutil"
	"os"
	"text/template"
)

func main() {
	var tmplFile string
	var outFile string
	var skipFormatting bool

	flag.StringVar(&tmplFile, "template", "", "Path to template file")
	flag.StringVar(&outFile, "out", "-", "Path to where to write generated code ('-' = stdout)")
	flag.BoolVar(&skipFormatting, "skip-formatting", false, "Do not format or syntax check generated code")
	flag.Parse()

	content, err := ioutil.ReadFile(tmplFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal: opening file '%s' - %v\n", tmplFile, err)
		os.Exit(-1)
	}

	funcMap := map[string]interface{}{"list": mkList}
	t := template.Must(template.New("generator").Funcs(funcMap).Parse(string(content)))
	var buf bytes.Buffer
	t.Execute(&buf, nil)

	if skipFormatting {
		writeRes(buf.Bytes(), outFile) // .Print(buf.String())
	} else {
		if src, err := format.Source(buf.Bytes()); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: formatting error in '%s' - %v\n", tmplFile, err)
			os.Exit(-1)
		} else {
			writeRes(src, outFile)
			//fmt.Print(string(src))
		}
	}
}

func writeRes(content []byte, fileName string) {
	if fileName == "-" {
		fmt.Fprint(os.Stdout, string(content))
	} else {
		if err := ioutil.WriteFile(fileName, content, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal: writing generated code to '%s' - %v\n", fileName, err)
			os.Exit(-1)
		}
	}
}

func mkList(args ...interface{}) []interface{} {
	return args
}
