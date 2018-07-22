package store

import "time"

type Price struct {
	Item      string
	Currency  string
	Cost      float64
	ActualFor time.Time
}
