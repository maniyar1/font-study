package main

import (
	"io"

	"github.com/wcharczuk/go-chart"
)

func getChartValueArray() []chart.Value {
	var chartArray []chart.Value
	fonts := getEntireDatabaseAsSlice()
	for _, font := range fonts {
		chartArray = append(chartArray, chart.Value{Value: float64(font.AveragePoints), Label: font.Family})
	}
	return chartArray
}
func graph(writer io.Writer) {
	graph := chart.BarChart{
		Title: "Font Average Points",
		Background: chart.Style{
			Padding: chart.Box{
				Top:   40,
				Left:  20,
				Right: 20,
			},
		},
		Canvas: chart.Style{
			Padding: chart.Box{
				Left:  50,
				Right: 50,
			},
		},
		BarSpacing: 20,
		Height:     512,
		Width:      1080,
		BarWidth:   40,
		Bars:       getChartValueArray(),
	}

	graph.Render(chart.SVG, writer)
}
