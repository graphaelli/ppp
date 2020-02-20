package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/pprof/profile"
)

func main() {
	p, err := convert(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}
	if err := p.Write(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

func convert(r io.Reader) (*profile.Profile, error) {
	p := &profile.Profile{
		SampleType: []*profile.ValueType{
			{
				Type: "samples",
				Unit: "count",
			},
		},
		Sample: make([]*profile.Sample, 0),
		Mapping: []*profile.Mapping{&profile.Mapping{
			ID: 1,
		}},
		Location: make([]*profile.Location, 0),
		Function: make([]*profile.Function, 0),
	}

	funLookup := make(map[string]*profile.Function)
	locLookup := make(map[string]*profile.Location)
	var id uint64 = 0

	b := bufio.NewReader(r)
	for {
		line, err := b.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				return p, nil
			}
			return nil, err
		}
		fields := strings.Fields(line)
		if len(fields) != 2 {
			return nil, fmt.Errorf("unexpected profile line: %s", line)
		}
		count, err := strconv.ParseInt(fields[1], 10, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to find count in %s: %w", fields[1], err)
		}

		sample := &profile.Sample{
			Location: make([]*profile.Location, 9),
			Value:    []int64{count},
		}

		for _, name := range strings.Split(fields[0], ";") {
			if _, ok := funLookup[name]; !ok {
				funLookup[name] = &profile.Function{
					ID:   id,
					Name: name,
				}
				locLookup[name] = &profile.Location{
					ID:      id,
					Mapping: p.Mapping[0],
					Line:    []profile.Line{{Function: funLookup[name]}},
				}
				id++

				p.Function = append(p.Function, funLookup[name])
				p.Location = append(p.Location, locLookup[name])
			}

			sample.Location = append(sample.Location, locLookup[name])
		}

		p.Sample = append(p.Sample, sample)
	}

}
