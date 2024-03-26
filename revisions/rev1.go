package revisions

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"sort"
	"strconv"

	"github.com/golobby/container/v3"
)

func init() {
	container.NamedSingleton("rev1", NewRev1)
}

type Rev1 struct{}

func NewRev1() Revision { return &Rev1{} }

func (rev Rev1) ProcessTemperatures(filepath string, output io.Writer) error {
	allStations := make(map[string]City, 0)
	allCities := make([]string, 0)

	file, err := os.Open(filepath)
	if err != nil {
		return err
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
				return err
			}
		}

		cityName := record[0]
		cityTemperature := record[1]
		cityTemperatureInFloat, err := strconv.ParseFloat(cityTemperature, 64)
		if err != nil {
			return err
		}

		city, ok := allStations[cityName]
		if !ok {
			allStations[cityName] = City{
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
		fmt.Printf("%s=%.1f/%.1f/%.1f, ", cityName, city.min, mean, city.max)
	}

	return nil
}
