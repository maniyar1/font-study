package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type Config struct {
	NumberOfFonts uint
}

var config Config = Config{NumberOfFonts: 20}

func main() {
	openDB()
	defer closeDB()
	http.HandleFunc("/index", index)
	http.HandleFunc("/", index)
	http.HandleFunc("/thanks", thanks)
	http.HandleFunc("/data.json", returnJSON)
	http.HandleFunc("/graph.svg", returnGraph)
	http.HandleFunc("/graph.html", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/graph.html")
	})
	http.HandleFunc("/graph.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/graph.js")
	})
	http.HandleFunc("/dataValidator.js", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "web/dataValidator.js")
	})
	log.Println("Listening on :28892...")
	log.Println("Listening on :28892...")
	http.ListenAndServe(":28892", nil)
}

func thanks(respWriter http.ResponseWriter, req *http.Request) {
	req.ParseForm()
	responses := Responses{}

	for key, val := range req.Form {
		fontRatings := getValueOrBlank(key)
		vString := strings.Join(val, "")
		var points uint
		val, _ := strconv.Atoi(vString)
		switch val {
		case 1:
			fontRatings.FirstPlaceOccurances++
			points = 5
		case 2:
			fontRatings.SecondPlaceOccurances++
			points = 4
		case 3:
			fontRatings.ThirdPlaceOccurances++
			points = 3
		case 4:
			fontRatings.FourthPlaceOccurances++
			points = 2
		case 5:
			fontRatings.FifthPlaceOccurances++
			points = 1
		case 6:
			fontRatings.SixthPlaceOccurances++
			points = 0
		}
		fontRatings.TotalEntries++
		fontRatings.Points += points
		fontRatings.AveragePoints = fontRatings.Points / fontRatings.TotalEntries
		responses.Responses = append(responses.Responses, Response{Family: key, UserPoints: points, AveragePoints: fontRatings.AveragePoints})
		byteJSON, err := json.Marshal(fontRatings)
		check(err)
		db.Put([]byte(key), byteJSON, nil)
		fmt.Println(string(byteJSON))
	}
	runResponseTemplate("web/thanks.html", responses, respWriter)
}

func returnJSON(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(respWriter, getEntireDatabaseAsJSON())
}

func returnGraph(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	graph(respWriter)
}

func index(respWriter http.ResponseWriter, req *http.Request) {
	pangram := "Pack my box with five dozen liquor jugs."
	data := PageData{Options: createOptions(pangram)}
	runCSSTemplate("web/template.css", &data)
	runHTMLTemplate("web/template.html", data, respWriter)
}

func createOptions(pangram string) []Option {
	rand.Seed(time.Now().UnixNano())
	fonts := getJson()
	rand := rand.Perm(int(config.NumberOfFonts))
	var sixOptions [6]Option
	for i, r := range rand[:6] {
		if len(fonts[r].Files.Regular[:]) != 0 {
			fonts[r].Files.Regular = fonts[r].Files.Regular[:4] + "s" + fonts[r].Files.Regular[4:]
		}
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

func check(e error) {
	if e != nil {
		log.Println(e)
	}
}
