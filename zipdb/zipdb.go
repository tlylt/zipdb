package zipdb

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/tlylt/zipdb/domain"
)

func readLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		text := scanner.Text()
		if strings.Trim(text, " ") == "" {
			continue
		}
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func ParseLocations(path string) []domain.Location {
	content, err := readLines(path)
	if err != nil {
		panic(err)
	}
	locations := make([]domain.Location, len(content))

	for idx, line := range content {
		rawLocation := strings.Split(line, "\t")
		lat, err := strconv.ParseFloat(rawLocation[9], 64)
		if err != nil {
			fmt.Println(err)
		}
		long, err := strconv.ParseFloat(rawLocation[10], 64)
		if err != nil {
			fmt.Println(err)
		}
		locations[idx] = domain.Location{
			Country:   rawLocation[0],
			Zip:       rawLocation[1],
			City:      rawLocation[2],
			StateLong: rawLocation[3],
			State:     rawLocation[4],
			County:    rawLocation[5],
			Lat:       lat,
			Long:      long,
		}
		// fmt.Printf("%d,%v", idx, location)
	}
	return locations
}

func LoadZipDB() map[string]domain.Location {
	content, err := readLines("US.txt")
	if err != nil {
		panic(err)
	}
	locationMap := make(map[string]domain.Location)

	for _, line := range content {
		rawLocation := strings.Split(line, "\t")
		lat, err := strconv.ParseFloat(rawLocation[9], 64)
		if err != nil {
			fmt.Println(err)
		}
		long, err := strconv.ParseFloat(rawLocation[10], 64)
		if err != nil {
			fmt.Println(err)
		}
		location := domain.Location{
			Country:   rawLocation[0],
			Zip:       rawLocation[1],
			City:      rawLocation[2],
			StateLong: rawLocation[3],
			State:     rawLocation[4],
			County:    rawLocation[5],
			Lat:       lat,
			Long:      long,
		}
		locationMap[location.Zip] = location
	}
	return locationMap
}
