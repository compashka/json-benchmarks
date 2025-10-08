package result

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type entry struct {
	cat  string
	impl string
	avg  float64
}

// GenerateMarkdownSection reads benchmark results from inputFile
// and overwrites the section in outFile defined by <!-- sectionName start --> and <!-- sectionName end -->.
// All old content inside the section is removed.
func GenerateMarkdownSection(inputFile, outFile, sectionName string) error {
	// Open benchmark file
	fl, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("failed to open input file: %v", err)
	}
	defer fl.Close()

	results := make(map[string][]float64)
	scanner := bufio.NewScanner(fl)
	for scanner.Scan() {
		text := scanner.Text()
		text, ok := strings.CutPrefix(text, "Benchmark")
		if !ok {
			continue
		}

		// Skip nested and number_structure benchmarks
		if strings.HasPrefix(text, "Marshal/nested_structure") || strings.HasPrefix(text, "Unmarshal/nested_structure") ||
			strings.HasPrefix(text, "Marshal/number_structure") || strings.HasPrefix(text, "Unmarshal/number_structure") {
			continue
		}

		spl := strings.Split(text, "\t")
		if len(spl) < 4 {
			continue
		}

		mbs, ok := strings.CutSuffix(spl[3], " MB/s")
		if !ok {
			continue
		}
		num, err := strconv.ParseFloat(strings.TrimSpace(mbs), 64)
		if err != nil {
			continue
		}
		results[spl[0]] = append(results[spl[0]], num)
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("failed to read benchmark file: %v", err)
	}

	// Prepare entries
	ents := make([]entry, 0, len(results))
	for name, result := range results {
		spl := strings.SplitN(name, "/", 3)
		if len(spl) != 3 {
			continue
		}
		var sum float64
		for _, f64 := range result {
			sum += f64
		}
		ents = append(ents, entry{
			cat:  fmt.Sprintf("%s - %s", spl[0], spl[1]),
			impl: spl[2],
			avg:  sum / float64(len(result)),
		})
	}

	// Sort entries
	slices.SortFunc(ents, func(a, b entry) int {
		if cmp := strings.Compare(a.cat, b.cat); cmp != 0 {
			return cmp
		}
		return strings.Compare(a.impl, b.impl)
	})

	// Generate Mermaid content
	var mermaid []string
	currentCat := ""
	for idx, ent := range ents {
		if idx == 0 || currentCat != ent.cat {
			currentCat = ent.cat
			mermaid = append(mermaid, "```mermaid")
			mermaid = append(mermaid, "gantt")
			mermaid = append(mermaid, fmt.Sprintf("title %s (MB/s - higher is better)", ent.cat))
			mermaid = append(mermaid, "dateFormat X")
			mermaid = append(mermaid, "axisFormat %s")
			mermaid = append(mermaid, "")
		}
		mermaid = append(mermaid, "section "+ent.impl)
		mermaid = append(mermaid, fmt.Sprintf("%d:0,%d", int(ent.avg), int(ent.avg)))
		if idx+1 >= len(ents) || ents[idx+1].cat != ent.cat {
			mermaid = append(mermaid, "```")
			mermaid = append(mermaid, "")
		}
	}

	newSection := strings.Join(mermaid, "\n")
	startMarker := fmt.Sprintf("<!-- %s start -->", sectionName)
	endMarker := fmt.Sprintf("<!-- %s end -->", sectionName)

	// Read existing README
	content, err := os.ReadFile(outFile)
	if err != nil {
		return fmt.Errorf("failed to read output file: %v", err)
	}
	lines := strings.Split(string(content), "\n")

	var out []string
	inSection := false
	for _, line := range lines {
		if strings.Contains(line, startMarker) {
			inSection = true
			out = append(out, line)       // keep start marker
			out = append(out, newSection) // insert new content
			continue
		}
		if strings.Contains(line, endMarker) {
			inSection = false
			out = append(out, line) // keep end marker
			continue
		}
		if !inSection {
			out = append(out, line)
		}
	}

	fmt.Println("Mermaid results in README generated")

	return os.WriteFile(outFile, []byte(strings.Join(out, "\n")), 0644)
}
