package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"test/internal/result"
)

func main() {
	mode := flag.String("mode", "all", "what to run: benchmark | result | all")
	flag.Parse()

	switch *mode {
	case "benchmark":
		runBenchmarks()
	case "result":
		err := runResult("result/result.txt", "./README.md")
		if err != nil {
			log.Fatal(err)
		}
	case "all":
		runBenchmarks()
		err := runResult("result/result.txt", "./README.md")
		if err != nil {
			log.Fatal(err)
		}
	default:
		log.Fatalf("unknown mode: %s (use: benchmark | generator | result | all)", *mode)
	}
}

func runBenchmarks() {
	fmt.Println("Running benchmarks...")

	if err := os.MkdirAll("result", 0755); err != nil {
		log.Fatalf("failed to create result dir: %v", err)
	}

	outFile, err := os.Create("result/result.txt")
	if err != nil {
		log.Fatalf("failed to create result file: %v", err)
	}
	defer outFile.Close()

	cmd := exec.Command("go", "test", "-bench=.", "-benchmem", "./")
	cmd.Stdout = outFile
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		log.Fatalf("failed to run benchmarks: %v", err)
	}
	fmt.Println("Benchmarks complete: result/result.txt")
}

func runResult(input, output string) error {
	err := result.GeneratePerformancePlots(input)
	if err != nil {
		return err
	}

	err = result.GenerateMarkdownSection(input, output, "benchmarks")
	if err != nil {
		return err
	}

	return nil
}
