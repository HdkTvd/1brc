package revisions

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/golobby/container/v3"
)

func init() {
	container.NamedSingleton("rev3", NewRev3)
}

type Rev3 struct{}

func NewRev3() Revision { return &Rev3{} }

func (rev Rev3) ProcessTemperatures(filepath string, output io.Writer) error {
	allStations := make(map[string]*City, 0)
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

		cityTemperatureInFloat, err := strconv.ParseFloat(temperature, 64)
		if err != nil {
			return err
		}

		city, ok := allStations[cityName]
		if !ok {
			allStations[cityName] = &City{
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
		}
	}

	sort.Strings(allCities)

	for _, cityName := range allCities {
		city := allStations[cityName]
		mean := city.sum / float64(city.count)
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", cityName, city.min, mean, city.max)
	}

	return nil
}
