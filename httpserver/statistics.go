package httpserver

import (
	"encoding/json"
	"log"
	"time"
)

/*
Object that stores statisitcs for an http server
Total is the number of processed requests
Average is the average time in ms to process those requests
*/
type Stats struct {
	Total   float64 "json:\"total\""
	Average float64 "json:\"average\""
}

/*
New processes have processed 0 requests
Negative value of average implies no information yet
*/
func NewStats() *Stats {
	return &Stats{0, -1}
}

/*
Adjusts the average to include the newest process time
*/
func adjustAverage(total, average, new float64) float64 {
	newAvg := (total*average + new) / (total + 1)
	return newAvg
}

/*
Encodes the stats struct as a json string
*/
func (s Stats) Encode() string {
	j, err := json.Marshal(s)
	if err != nil {
		log.Println(err)
	}
	return string(j)
}

/*
Signals that a process request has completed in time pt

Increases the total processes by 1, and adjusts the average time
according to the process time
*/
func (s *Stats) UpdateStatistics(pt time.Duration) {
	ms := pt.Seconds() * 1e3

	if s.Average < 0 {
		//negative average implies that no requests have been processed yet
		s.Average = ms
	} else {
		s.Average = adjustAverage(s.Total, s.Average, ms)
	}

	s.Total++
}
