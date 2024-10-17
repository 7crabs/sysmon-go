package main

import (
	"fmt"
	"github.com/7crabs/sysmon-go/topfetcher"
	"os"
	"strconv"
)

func main() {
	if len(os.Args) != 4 {
		fmt.Println("Usage: smg <output-file> <interval> <count>")
		os.Exit(1)
	}

	outputFile := os.Args[1]
	interval, err := strconv.Atoi(os.Args[2])
	if err != nil || interval <= 0 {
		fmt.Println("Interval must be a positive integer.")
		os.Exit(1)
	}

	count, err := strconv.Atoi(os.Args[3])
	if err != nil || count <= 0 {
		fmt.Println("Count must be a positive integer.")
		os.Exit(1)
	}

	// topfetcherモジュールを使ってデータを取得
	data, err := topfetcher.FetchTopData(interval, count)
	if err != nil {
		fmt.Printf("Error fetching top data: %v\n", err)
		os.Exit(1)
	}

	// JSONに変換
	jsonData, err := topfetcher.ToJSON(data)
	if err != nil {
		fmt.Printf("Error converting data to JSON: %v\n", err)
		os.Exit(1)
	}

	// ファイルに保存
	err = os.WriteFile(outputFile, []byte(jsonData), 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Filtered top command output saved to %s\n", outputFile)
}

