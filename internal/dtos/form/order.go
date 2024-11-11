package fromDto

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

// Crop represents the data transfer object for a crop
type Order struct {
	CropID    uuid.UUID `json:"crop_id"`
	SeedID    uuid.UUID `json:"seed_id"`
	Date      time.Time `json:"date"`
	Time      time.Time `json:"time"`
	IsHarvest bool      `json:"is_harvest"`
	Yield     int       `json:"yield"` // Yield in grams
}

// NewCrop creates a new Crop instance
func NewOrder(seedID uuid.UUID, date, time time.Time, isHarvest bool, yield int) Order {
	return Order{
		CropID:    uuid.New(),
		SeedID:    seedID,
		Date:      date,
		Time:      time,
		IsHarvest: isHarvest,
		Yield:     yield,
	}
}

func (o Order) GetDateTime() time.Time {
	return time.Date(o.Date.Year(), o.Date.Month(), o.Date.Day(), o.Time.Hour(), o.Time.Minute(), o.Time.Second(), 0, time.Local)
}

// Validate checks if the Crop fields are valid
func (o *Order) Validate() error {
	if o.SeedID == uuid.Nil {
		return fmt.Errorf("seed ID cannot be nil")
	}
	if o.Date.IsZero() {
		return fmt.Errorf("date cannot be zero")
	}
	if o.Time.IsZero() {
		return fmt.Errorf("time cannot be zero")
	}
	if o.Yield <= 0 {
		return fmt.Errorf("yield must be greater than zero")
	}
	return nil
}
