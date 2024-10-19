package topfetcher

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

type ProcessData struct {
	Values map[string]string `json:"values"`
}

type TimeSeriesData struct {
	Time string        `json:"time"`
	Data []ProcessData `json:"data"`
}

// FetchTopData topコマンドのデータを収集し、JSONで返す
func FetchTopData(interval, count int) ([]TimeSeriesData, error) {
	cmd := exec.Command("top", "-b", "-d", strconv.Itoa(interval), "-n", strconv.Itoa(count))
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("error executing top command: %w", err)
	}

	var timeSeriesData []TimeSeriesData
	var headers []string
	headerParsed := false
	var currentTime string
	var currentProcesses []ProcessData

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()

		if strings.HasPrefix(line, "top -") {
			if headerParsed && len(currentProcesses) > 0 {
				timeSeriesData = append(timeSeriesData, TimeSeriesData{
					Time: currentTime,
					Data: currentProcesses,
				})
				currentProcesses = []ProcessData{}
			}
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				currentTime = parts[2]
			}
			continue
		}

		if strings.HasPrefix(line, "  PID") {
			headerParsed = true
			headers = strings.Fields(line)
			continue
		}

		if strings.HasPrefix(line, "Tasks:") || strings.HasPrefix(line, "%Cpu(s):") || strings.HasPrefix(line, "MiB Mem :") || strings.HasPrefix(line, "MiB Swap:") {
			continue
		}

		if headerParsed {
			if len(line) == 0 {
				continue
			}

			fields := strings.Fields(line)
			if len(fields) >= len(headers) {
				processData := ProcessData{
					Values: make(map[string]string),
				}
				for i, header := range headers {
					processData.Values[header] = fields[i]
				}
				currentProcesses = append(currentProcesses, processData)
			}
		}
	}

	if len(currentProcesses) > 0 {
		timeSeriesData = append(timeSeriesData, TimeSeriesData{
			Time: currentTime,
			Data: currentProcesses,
		})
	}

	return timeSeriesData, nil
}

// ToJSON JSONに変換して返す
func ToJSON(data []TimeSeriesData) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling data to JSON: %w", err)
	}
	return string(jsonData), nil
}
