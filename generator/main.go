package main

import (
	"bufio"
	"fmt"
	"image/color"
	"log"
	"os"
	"regexp"
	"strconv"
	"strings"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
	"gonum.org/v1/plot/vg/draw"
)

func main() {
	filePath := "result/result.txt"

	f, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer f.Close()

	// Регулярные выражения для парсинга
	regexes := map[string]*regexp.Regexp{
		"StdLib":   regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/number_structure_(\d+)/encoding/json-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"EasyJSON": regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/number_structure_(\d+)/mailru/easyjson-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"Sonic":    regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/number_structure_(\d+)/bytedance/sonic-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"JsonIter": regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/number_structure_(\d+)/json-iterator/go-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
	}

	// Данные для графиков
	results := map[string]map[string][]float64{
		"Marshal":   {"StdLib": {}, "EasyJSON": {}, "Sonic": {}, "JsonIter": {}},
		"Unmarshal": {"StdLib": {}, "EasyJSON": {}, "Sonic": {}, "JsonIter": {}},
	}

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()

		for lib, regex := range regexes {
			matches := regex.FindStringSubmatch(line)
			if len(matches) == 4 {
				operation := matches[1]
				depthVal, timeVal, ok := parseMatches(matches[2], matches[3])
				if ok {
					results[operation][lib] = append(results[operation][lib], depthVal, timeVal)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("failed to read file: %v", err)
	}

	// Построение графиков
	createPlotWithColors(results["Marshal"], "Marshall Performance", "marshall_plot_colored.png")
	createPlotWithColors(results["Unmarshal"], "Unmarshall Performance", "unmarshall_plot_colored.png")
	fmt.Println("Графики сохранены в marshall_plot_colored.png и unmarshall_plot_colored.png")
}

func parseMatches(depthStr, timeStr string) (float64, float64, bool) {
	timeStr = strings.ReplaceAll(timeStr, ",", ".")
	depthVal, err := strconv.Atoi(depthStr)
	if err != nil {
		log.Printf("failed to parse depth: %v", err)
		return 0, 0, false
	}
	timeVal, err := strconv.ParseFloat(timeStr, 64)
	if err != nil {
		log.Printf("failed to parse ns/op: %v", err)
		return 0, 0, false
	}
	return float64(depthVal), timeVal, true
}

func createPlotWithColors(data map[string][]float64, title, outputPath string) {
	p := plot.New()

	p.Title.Text = title
	p.X.Label.Text = "Number of fields"
	p.Y.Label.Text = "ns/op"

	// Цвета для каждой библиотеки
	colors := map[string]color.RGBA{
		"StdLib":   {R: 255, G: 0, B: 0, A: 255},   // Красный
		"EasyJSON": {R: 0, G: 255, B: 0, A: 255},   // Зелёный
		"Sonic":    {R: 0, G: 0, B: 255, A: 255},   // Синий
		"JsonIter": {R: 255, G: 165, B: 0, A: 255}, // Оранжевый
	}

	// Добавляем линии
	for lib, values := range data {
		if len(values) == 0 {
			continue
		}
		depths, times := splitPoints(values)
		points := makePoints(depths, times)
		line, pointsPlot, err := plotter.NewLinePoints(points)
		if err != nil {
			log.Fatalf("failed to create line points: %v", err)
		}
		line.Color = colors[lib]
		pointsPlot.Shape = draw.CircleGlyph{}
		pointsPlot.Color = colors[lib]

		p.Add(line, pointsPlot)
		p.Legend.Add(lib, line)
	}

	// Сохраняем график в файл PNG
	if err := p.Save(9*vg.Inch, 6*vg.Inch, outputPath); err != nil {
		log.Fatalf("failed to save plot: %v", err)
	}
}

func splitPoints(data []float64) ([]float64, []float64) {
	var depths, times []float64
	for i := 0; i < len(data); i += 2 {
		depths = append(depths, data[i])
		times = append(times, data[i+1])
	}
	return depths, times
}

func makePoints(depths, times []float64) plotter.XYs {
	pts := make(plotter.XYs, len(depths))
	for i := range depths {
		pts[i].X = depths[i]
		pts[i].Y = times[i]
	}
	return pts
}
