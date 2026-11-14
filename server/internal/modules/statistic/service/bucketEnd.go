package service

import "time"

// bucketEnd — последний момент бакета: перед началом следующего, но не позже конца периода
func (c seriesCalculator) bucketEnd(index int) time.Time {
	end := c.endOfPeriod
	if index+1 < len(c.buckets) && c.buckets[index+1].Before(end) {
		end = c.buckets[index+1]
	}
	return end.Add(-time.Nanosecond)
}
