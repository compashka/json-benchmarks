package result

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

// BenchData represents a single benchmark data point (X = depth or field count, Y = time in ns/op).
type BenchData struct {
	X float64
	Y float64
}

// GeneratePerformancePlots parses Go benchmark results from result.txt
// and generates performance plots for different libraries and structure types.
// It produces four PNG files:
// - marshal_number.png
// - unmarshal_number.png
// - marshal_nested.png
// - unmarshal_nested.png
func GeneratePerformancePlots(inputFile string) error {
	f, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open file: %v", err)
	}
	defer f.Close()

	// Regular expressions to extract benchmark results for each JSON library
	patterns := map[string]*regexp.Regexp{
		"StdLib":   regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/(number_structure|nested_structure)_(\d+)/encoding/json-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"EasyJSON": regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/(number_structure|nested_structure)_(\d+)/mailru/easyjson-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"Sonic":    regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/(number_structure|nested_structure)_(\d+)/bytedance/sonic-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
		"JsonIter": regexp.MustCompile(`Benchmark(Marshal|Unmarshal)/(number_structure|nested_structure)_(\d+)/json-iterator/go-\d+\s+\d+\s+([\d\.]+)\s+ns/op`),
	}

	// results[operation][structureType][library] = []BenchData
	results := map[string]map[string]map[string][]BenchData{
		"Marshal":   {"number_structure": {}, "nested_structure": {}},
		"Unmarshal": {"number_structure": {}, "nested_structure": {}},
	}

	for op := range results {
		for typ := range results[op] {
			for lib := range patterns {
				results[op][typ][lib] = []BenchData{}
			}
		}
	}

	// Parse benchmark file line by line
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		for lib, re := range patterns {
			m := re.FindStringSubmatch(line)
			if len(m) == 5 {
				op := m[1]
				structType := m[2]
				x := parseFloat(m[3])
				y := parseFloat(m[4])
				results[op][structType][lib] = append(results[op][structType][lib], BenchData{X: x, Y: y})
			}
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read benchmark file: %v", err)
	}

	// Generate plots for all benchmark types
	createPlot(results["Marshal"]["number_structure"],
		"Marshal — Number Structure", "Number of Fields", "Time per Operation (ns/op)", "result/marshal_number.png")
	createPlot(results["Unmarshal"]["number_structure"],
		"Unmarshal — Number Structure", "Number of Fields", "Time per Operation (ns/op)", "result/unmarshal_number.png")
	createPlot(results["Marshal"]["nested_structure"],
		"Marshal — Nested Structure", "Nesting Depth", "Time per Operation (ns/op)", "result/marshal_nested.png")
	createPlot(results["Unmarshal"]["nested_structure"],
		"Unmarshal — Nested Structure", "Nesting Depth", "Time per Operation (ns/op)", "result/unmarshal_nested.png")

	fmt.Println("Plots generated: marshal_number.png, unmarshal_number.png, marshal_nested.png, unmarshal_nested.png")
	return nil
}

// parseFloat safely converts a string into a float64.
func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", ".")
	v, _ := strconv.ParseFloat(s, 64)
	return v
}

// createPlot builds a single performance plot for all libraries.
func createPlot(data map[string][]BenchData, title, xlabel, ylabel, filename string) {
	p := plot.New()
	p.Title.Text = title
	p.X.Label.Text = xlabel
	p.Y.Label.Text = ylabel
	p.BackgroundColor = color.RGBA{245, 245, 245, 255} // light gray background
	p.Add(plotter.NewGrid())

	// Library colors
	colors := map[string]color.RGBA{
		"StdLib":   {R: 255, G: 99, B: 132, A: 255}, // red-pink
		"EasyJSON": {R: 75, G: 192, B: 192, A: 255}, // cyan
		"Sonic":    {R: 54, G: 162, B: 235, A: 255}, // blue
		"JsonIter": {R: 255, G: 206, B: 86, A: 255}, // yellow
	}

	for lib, vals := range data {
		if len(vals) == 0 {
			continue
		}
		points := make(plotter.XYs, len(vals))
		for i, v := range vals {
			points[i].X = v.X
			points[i].Y = v.Y
		}

		line, dots, err := plotter.NewLinePoints(points)
		if err != nil {
			log.Fatalf("failed to create line points: %v", err)
		}

		line.Color = colors[lib]
		dots.Color = colors[lib]
		dots.Shape = draw.CircleGlyph{}

		p.Add(line, dots)
		p.Legend.Add(lib, line)
	}

	p.Legend.Top = true
	p.Legend.XOffs = -5

	if err := p.Save(9*vg.Inch, 6*vg.Inch, filename); err != nil {
		log.Fatalf("failed to save plot: %v", err)
	}
}
