package main

import (
	"1brc/revisions"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"github.com/golobby/container/v3"
)

func main() {
	filepath := flag.String("filepath", "./data/weather_stations.csv", "Enter a filepath of the csv data file.")
	cpuprofile := flag.String("cpuprofile", "./result/cpuprofile.prof", "Write cpu profile to file.")
	prevResults := flag.String("results", "./result/prevResults.csv", "Filepath of csv file storing results of previous runs.")
	rev := flag.String("revision", "rev1", "Select a revision to apply.")
	flag.Parse()

	var revision revisions.Revision
	if err := container.NamedResolve(&revision, *rev); err != nil {
		log.Fatalf("Error resolving revision implementation [%v]", err)
	}

	f, _ := os.Create(*cpuprofile)

	start := time.Now()

	pprof.StartCPUProfile(f)

	if err := revision.ProcessTemperatures(*filepath, os.Stdout); err != nil {
		log.Fatalf("Error processing temperatures - [%v]", err)
	}

	end := time.Until(start)

	pprof.StopCPUProfile()

	fmt.Printf("\nTime taken in seconds - [%v]\n", end.Seconds())

	resultsFile, err := os.OpenFile(*prevResults, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("Error opening/creating prev results file - [%v]", err)
	}
	defer resultsFile.Close()

	csvwriter := csv.NewWriter(resultsFile)
	csvwriter.Write([]string{*rev, end.String()})
	csvwriter.Flush()
}
