package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/valyala/fasthttp"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	//	"strings"
	"time"
)

type Config struct {
	NumberOfFonts uint
}

var config Config = Config{NumberOfFonts: 20}
var (
	addr               = flag.String("addr", ":28892", "TCP address to listen to")
	compress           = flag.Bool("compress", false, "Whether to enable transparent response compression")
	dir                = flag.String("dir", "web/", "Directory to serve static files from")
	generateIndexPages = flag.Bool("generateIndexPages", false, "Whether to generate directory index pages")
	byteRange          = flag.Bool("byteRange", false, "Enables byte range requests if set to true")
)

func main() {
	openDB()
	defer closeDB()
	flag.Parse()
	fs := &fasthttp.FS{
		Root:               *dir,
		IndexNames:         []string{"index.html"},
		GenerateIndexPages: *generateIndexPages,
		Compress:           *compress,
		AcceptByteRange:    *byteRange,
	}
	fsHandler := fs.NewRequestHandler()

	requestHandler := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/":
			index(ctx)
		case "/font-study":
			index(ctx)
		case "/thanks":
			thanks(ctx)
		case "/data.json":
			returnJSON(ctx)
		case "/graph.svg":
			returnGraph(ctx)
		default:
			fsHandler(ctx)
		}
	}

	if err := fasthttp.ListenAndServe(*addr, requestHandler); err != nil {
		log.Fatalf("Error in ListenAndServe: %s", err)
	}
}

func thanks(ctx *fasthttp.RequestCtx) {
	responses := Responses{}

	ctx.SetContentType("text/html; charset=utf-8")
	form := ctx.PostArgs()
	form.VisitAll(func(key, val []byte) {
		fontRatings := getValueOrBlank(string(key))
		var points uint
		value, _ := strconv.Atoi(string(val))
		fmt.Println(value)

		switch value {
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
		default:
			log.Println("Non-valid number")
		}
		fontRatings.TotalEntries++
		fontRatings.Points += uint64(points)
		fontRatings.AveragePoints = float64(fontRatings.Points) / float64(fontRatings.TotalEntries)
		responses.Responses = append(responses.Responses, Response{Family: string(key), UserPoints: points, AveragePoints: fontRatings.AveragePoints})
		byteJSON, err := json.Marshal(fontRatings)
		check(err)
		db.Put(key, byteJSON, nil)
		fmt.Println(string(byteJSON))
	})
	runResponseTemplate("web/thanks.html", responses, ctx)
}

func returnJSON(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("application/json; charset=utf-8")
	fmt.Fprintf(ctx, getEntireDatabaseAsJSON())
}

func returnGraph(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("image/svg+xml; charset=utf-8")
	graph(ctx)
}

func index(ctx *fasthttp.RequestCtx) {
	ctx.SetContentType("text/html; charset=utf-8")
	paragraph := "If you know the enemy and know yourself, you need not fear the result of a hundred battles. If you know yourself but not the enemy, for every victory gained you will also suffer a defeat. If you know neither the enemy nor yourself, you will succumb in every battle. <br>"
	data := PageData{Options: createOptions(paragraph)}
	runCSSTemplate("web/template.css", &data)
	runHTMLTemplate("web/template.html", data, ctx)
}

func createOptions(pangram string) []Option {
	rand.Seed(time.Now().UnixNano())
	fonts := getJSON()
	rand := rand.Perm(int(config.NumberOfFonts))
	var sixOptions [6]Option
	for i, r := range rand[:6] {
		for len(fonts[r].Files.Regular[:]) == 0 || fonts[r].Files.Regular[len(fonts[r].Files.Regular)-3:] != "ttf" {
			r += 6
		}
		fonts[r].Files.Regular = fonts[r].Files.Regular[:4] + "s" + fonts[r].Files.Regular[4:]
		sixOptions[i] = Option{Number: i, Pangram: pangram, Font: fonts[r]}
	}
	return sixOptions[:]
}

func getJSON() []Font {
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
