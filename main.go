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
	http.ListenAndServe(":8090", nil)
}

func thanks(respWriter http.ResponseWriter, req *http.Request) {
	req.ParseForm()

	for key, val := range req.Form {
		fontRatings := getValueOrBlank(key)
		vString := strings.Join(val, "")
		var points uint
		val, _ := strconv.Atoi(vString)
		switch val {
		case 1:
			fontRatings.FirstPlaceOccurances++
			points = 80
		case 2:
			fontRatings.SecondPlaceOccurances++
			points = 40
		case 3:
			fontRatings.ThirdPlaceOccurances++
			points = 20
		case 4:
			fontRatings.FourthPlaceOccurances++
			points = 10
		case 5:
			fontRatings.FifthPlaceOccurances++
			points = 5
		case 6:
			fontRatings.SixthPlaceOccurances++
			points = 0
		}
		fontRatings.TotalEntries++
		fontRatings.Points += points
		fontRatings.AveragePoints = fontRatings.Points / fontRatings.TotalEntries
		byteJSON, err := json.Marshal(fontRatings)
		check(err)
		db.Put([]byte(key), byteJSON, nil)
		fmt.Println(string(byteJSON))
	}
	respWriter.Header().Set("Content-Type", "application/json; charset=utf-8")
	fmt.Fprintf(respWriter, "Thanks!\n")
}

func returnJSON(respWriter http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(respWriter, getEntireDatabaseAsJSON())
}

func returnGraph(respWriter http.ResponseWriter, req *http.Request) {
	respWriter.Header().Set("Content-Type", "image/svg+xml; charset=utf-8")
	graph(respWriter)
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
	rand := rand.Perm(int(config.NumberOfFonts))
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

func check(e error) {
	if e != nil {
		log.Fatal(e)
		panic(e)
	}
}
