package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"math/rand"
	"net/http"
	"text/template"
	"time"
)

// JSON Structs
type JSONFormat struct {
	Items []Font
}

type Font struct {
	Kind         string
	Family       string
	Category     string
	Variants     []string
	Subsets      []string
	Version      string
	LastModified string
	Files        Files
}

type Files struct {
	Regular string
}

// Option is for HTML/CSS
type Option struct {
	Number  int
	Pangram string
	Font    Font
}

type PageData struct {
	CSS     string
	Options []Option
}

func main() {
	http.HandleFunc("/index", index)
	http.ListenAndServe(":8090", nil)
}

func index(respWriter http.ResponseWriter, req *http.Request) {
	pangram := "Pack my box with five dozen liquor jugs."
	data := PageData{Options: createOptions(pangram)}
	runCSSTemplate("template.css", &data)
	runHTMLTemplate("template.html", data, respWriter)
}

func createOptions(pangram string) []Option {
	rand.Seed(time.Now().UnixNano())
	fonts := getJson()
	rand := rand.Perm(50)
	var sixOptions [6]Option
	for i, r := range rand[:6] {
		sixOptions[i] = Option{Number: i, Pangram: pangram, Font: fonts[r]}
	}
	return sixOptions[:]
}

func getJson() []Font {
	client := http.Client{
		Timeout: time.Second * 2, // Maximum of 2 secs
	}
	url := "https://www.googleapis.com/webfonts/v1/webfonts?sort=popularity&key=AIzaSyDDxtMndiR4WRsqzWW3QNUZix2sbCVNOzI&"
	req, err := http.NewRequest("GET", url, nil)
	check(err)

	req.Header.Set("User-Agent", "font-stats")
	res, getErr := client.Do(req)
	check(getErr)

	body, readErr := ioutil.ReadAll(res.Body)
	check(readErr)

	var jsonResult JSONFormat
	json.Unmarshal(body, &jsonResult)
	return jsonResult.Items
}

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

func check(e error) {
	if e != nil {
		panic(e)
	}
}
