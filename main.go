package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"time"
)

type City struct {
	name  string
	min   float64
	sum   float64
	max   float64
	count int
}

func main() {
	fmt.Println("-- This is a 1BRC challenge --")
	start := time.Now()
	fmt.Printf("Time starts now. %v\n", start.String())

	allStations := make(map[string]City, 0)
	allCities := make([]string, 0)

	filepath := "./data/weather_stations.csv"
	file, err := os.Open(filepath)
	if err != nil {
		log.Fatalf("Error opening file [%v]", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ';'
	reader.FieldsPerRecord = 2

	// Hack: skipping first 2 rows
	_, _ = reader.Read()
	_, _ = reader.Read()

	var errRead error
	for errRead != io.EOF {
		record, err := reader.Read()
		if err != nil {
			if err == io.EOF {
				errRead = err
				continue
			} else {
				log.Fatalf("Error reading a line from the file [%v]", err)
			}
		}

		cityName := record[0]
		cityTemperature := record[1]
		cityTemperatureInFloat, err := strconv.ParseFloat(cityTemperature, 64)
		if err != nil {
			log.Fatalf("Error parsing string to float64 [%v]", err)
		}

		city, ok := allStations[cityName]
		if !ok {
			allStations[cityName] = City{
				name:  cityName,
				min:   cityTemperatureInFloat,
				sum:   cityTemperatureInFloat,
				max:   cityTemperatureInFloat,
				count: 1,
			}
			allCities = append(allCities, cityName)
		} else {
			city.min = min(city.min, cityTemperatureInFloat)
			city.max = max(city.max, cityTemperatureInFloat)
			city.sum += cityTemperatureInFloat
			city.count++

			allStations[cityName] = city
		}
	}

	sort.Strings(allCities)

	for _, cityName := range allCities {
		city := allStations[cityName]
		mean := city.sum / float64(city.count)
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", city.name, city.min, mean, city.max)
	}

	end := time.Until(start)
	fmt.Printf("\nChallenge completed. Time taken in seconds - [%v]\n", end.Seconds())
}

func min(x, y float64) float64 {
	if x > y {
		return y
	}
	return x
}

func max(x, y float64) float64 {
	if x < y {
		return y
	}
	return x
}
