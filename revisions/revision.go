package revisions

import "io"

type City struct {
	min   float64
	sum   float64
	max   float64
	count int
}

type Revision interface {
	ProcessTemperatures(filepath string, output io.Writer) error
}
