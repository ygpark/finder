package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"regexp"
)

type extractor struct {
	name       string
	headerName string // CSV 헤더에 표시될 이름
	pattern    *regexp.Regexp
	group      int
}

func extractData(line string, extractors []extractor) map[string]string {
	extractedData := make(map[string]string)
	hasMatch := false

	for _, ex := range extractors {
		matches := ex.pattern.FindStringSubmatch(line)
		if len(matches) > ex.group {
			extractedData[ex.name] = matches[ex.group]
			hasMatch = true
		}
	}

	if hasMatch {
		return extractedData
	}
	return nil
}

func main() {
	inputFilePath := flag.String("i", "", "Input file path (required)")
	flag.Parse()

	if *inputFilePath == "" {
		fmt.Println("Error: Input file path is required.")
		flag.Usage()
		os.Exit(1)
	}

	file, err := os.Open(*inputFilePath)
	if err != nil {
		fmt.Println("Error opening input file:", err)
		os.Exit(1)
	}
	defer file.Close()

	writer := csv.NewWriter(os.Stdout)
	defer writer.Flush()

	extractors := []extractor{ // 슬라이스로 변경
		{
			name:       "date1",
			headerName: "Date (YYYY-MM-DD)", // 헤더 이름 명시적으로 지정
			pattern:    regexp.MustCompile(`\[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})\]`),
			group:      1,
		},
		{
			name:       "date2",
			headerName: "Date (Apache)", // 헤더 이름 명시적으로 지정
			pattern:    regexp.MustCompile(`(\d{4}-\d{2}-\d{2}( \d{2}:\d{2}:\d{2})?)`),
			group:      1,
		},
		{
			name:       "ip",
			headerName: "IP Address", // 헤더 이름 명시적으로 지정
			pattern:    regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`),
			group:      0,
		},
		{
			name:       "email",
			headerName: "Email Address", // 헤더 이름 명시적으로 지정
			pattern:    regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
			group:      0,
		},

		// {
		// 	name:       "userAgent",
		// 	headerName: "User Agent",
		// 	pattern:    regexp.MustCompile(`"([^"]*)"`), // 따옴표로 묶인 user agent 추출
		// 	group:      1,
		// },
	}

	var header []string
	for _, ex := range extractors { // 헤더 생성 로직 변경
		header = append(header, ex.headerName)
	}
	writer.Write(header)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		extractedData := extractData(line, extractors)

		if extractedData != nil {
			var row []string
			for _, ex := range extractors { // 데이터 추출 로직 변경
				row = append(row, extractedData[ex.name])
			}
			writer.Write(row)
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input file:", err)
		os.Exit(1)
	}
}
