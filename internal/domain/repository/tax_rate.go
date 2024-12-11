package repository

import "time"

type TaxRate interface {
	GetRateByDate(date time.Time) (float64, error)
}
