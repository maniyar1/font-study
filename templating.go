package main

import (
	"bytes"
	"io"
	"text/template"
)

func runHTMLTemplate(fileName string, data PageData, writer io.Writer) {
	tmpl := template.Must(template.ParseFiles(fileName))
	tmpl.Execute(writer, data)
}

func runCSSTemplate(fileName string, data *PageData) {
	tmpl := template.Must(template.ParseFiles(fileName))
	var buf bytes.Buffer
	check(tmpl.Execute(&buf, *data))
	data.CSS = buf.String()
}

func runResponseTemplate(fileName string, data Responses, writer io.Writer) {
	tmpl := template.Must(template.ParseFiles(fileName))
	check(tmpl.Execute(writer, data))
}
