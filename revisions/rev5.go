package revisions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"

	"github.com/golobby/container/v3"
)

func init() {
	container.NamedSingleton("rev5", NewRev5)
}

type CityStats struct {
	min, max, count int32
	sum             int64
}

type Rev5 struct{}

func NewRev5() Revision { return &Rev5{} }

func (rev Rev5) ProcessTemperatures(filepath string, output io.Writer) error {
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
		cityName, temperature, hasSemi := strings.Cut(line, ";")
		if !hasSemi {
			continue
		}

		index := 0
		negative := false

		if temperature[0] == '-' {
			negative = true
			index++
		}

		cityTemperatureInt32 := int32(temperature[index] - '0')
		index++

		if temperature[index] != '.' {
			cityTemperatureInt32 = cityTemperatureInt32*10 + int32(temperature[index]-'0')
			index++
		}

		index++

		cityTemperatureInt32 = cityTemperatureInt32*10 + int32(temperature[index]-'0')
		if negative {
			cityTemperatureInt32 = -cityTemperatureInt32
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
