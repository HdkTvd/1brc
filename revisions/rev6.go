package revisions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"

	"github.com/golobby/container/v3"
)

func init() {
	container.NamedSingleton("rev6", NewRev6)
}

type Rev6 struct{}

func NewRev6() Revision { return &Rev6{} }

func (rev Rev6) ProcessTemperatures(filepath string, output io.Writer) error {
	allStations := make(map[string]*CityStats, 0)
	allCities := make([]string, 0)

	file, err := os.Open(filepath)
	if err != nil {
		return err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		end := len(line)
		cityName := ""
		cityTemperatureInt32 := int32(line[end-1] - '0')    //last digit
		cityTemperatureInt32 += int32(line[end-3]-'0') * 10 //digit after decimal point
		// number, -ve sign, semicolon
		if line[end-4]-'0' >= '0' && line[end-4]-'0' <= '9' {
			cityTemperatureInt32 += int32(line[end-4]-'0') * 100 // digit at 10's place
			semicolon := end - 5
			if line[end-5] == '-' {
				cityTemperatureInt32 = -cityTemperatureInt32
				semicolon = end - 6
			}
			cityName = line[:semicolon]
		}

		// handling semicolon or -ve sign
		if line[end-4] == '-' {
			semicolon := end - 5
			cityTemperatureInt32 = -cityTemperatureInt32
			cityName = line[:semicolon]
		} else {
			semicolon := end - 4
			cityName = line[:semicolon]
		}

		city, ok := allStations[cityName]
		if !ok {
			allStations[cityName] = &CityStats{
				min:   cityTemperatureInt32,
				sum:   int64(cityTemperatureInt32),
				max:   cityTemperatureInt32,
				count: 1,
			}
			allCities = append(allCities, cityName)
		} else {
			city.min = min(city.min, cityTemperatureInt32)
			city.max = max(city.max, cityTemperatureInt32)
			city.sum += int64(cityTemperatureInt32)
			city.count++
		}
	}

	sort.Strings(allCities)

	for _, cityName := range allCities {
		city := allStations[cityName]
		mean := float64(city.sum) / float64(city.count) / 10
		fmt.Fprintf(output, "%s=%.1f/%.1f/%.1f, ", cityName, float64(city.min)/10, mean, float64(city.max)/10)
	}

	return nil
}
