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

	for _, ex := range extractors {
		matches := ex.pattern.FindStringSubmatch(line)
		if len(matches) > ex.group {
			extractedData[ex.name] = matches[ex.group]
		}
	}

	if len(extractedData) > 0 { // hasMatch 변수 제거, 맵 길이로 확인
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

	extractors := []extractor{
		{
			name:       "date1",
			headerName: "Date (Apache)", // 헤더 이름 수정
			pattern:    regexp.MustCompile(`\[(\d{2}/\w+/\d{4}:\d{2}:\d{2}:\d{2} \+\d{4})\]`),
			group:      1, //대괄호는 버리고 날짜만 추출하기 위해 그룹1번 선택
		},
		{
			name:       "date2",
			headerName: "Date and Time", // 헤더 이름 수정
			pattern:    regexp.MustCompile(`(\d{4}-\d{2}-\d{2}( \d{2}:\d{2}:\d{2})?)`),
			group:      1, //날짜와 시간정보가 함께 포함된 그룹1번 선택
		},
		{
			name:       "ip",
			headerName: "IP Address",
			pattern:    regexp.MustCompile(`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`),
			group:      0, // 매치된 전체 문자열을 나타내는 0번 그룹 선택
		},
		{
			name:       "email",
			headerName: "Email Address",
			pattern:    regexp.MustCompile(`[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}`),
			group:      0, // 매치된 전체 문자열을 나타내는 0번 그룹 선택
		},
	}

	var header []string
	for _, ex := range extractors {
		header = append(header, ex.headerName)
	}
	writer.Write(header)

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		extractedData := extractData(line, extractors)

		if extractedData != nil {
			var row []string
			for _, ex := range extractors {
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
