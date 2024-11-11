package dtos

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"lvm/database"
	"time"
)

// Crop represents the data transfer object for a crop
type Crop struct {
	CropID        uuid.UUID  `json:"crop_id"`
	SeedID        uuid.UUID  `json:"seed_id"`
	SoakingStart  *time.Time `json:"soaking_start,omitempty"`
	StackingStart time.Time  `json:"stacking_start"`
	BlackoutStart *time.Time `json:"blackout_start,omitempty"`
	LightsStart   time.Time  `json:"lights_start"`
	Harvest       time.Time  `json:"harvest"`
	Code          string     `json:"code"`
	YieldGrams    *int       `json:"YieldGrams"` // YieldGrams in grams
}

// NewCrop creates a new Crop instance
func NewCrop(
	seedID uuid.UUID,
	soakingStart *time.Time,
	stackingStart time.Time,
	blackoutStart *time.Time,
	lightsStart time.Time,
	harvest time.Time,
	code string,
	yieldGrams *int,
) Crop {
	return Crop{
		CropID:        uuid.New(),
		SeedID:        seedID,
		SoakingStart:  soakingStart,
		StackingStart: stackingStart,
		BlackoutStart: blackoutStart,
		LightsStart:   lightsStart,
		Code:          code,
		Harvest:       harvest,
		YieldGrams:    yieldGrams,
	}
}

func (c *Crop) ToDatabaseModel() database.Crop {
	// Handle nullable SoakingStart
	soakingStart := pgtype.Timestamp{
		Valid: false,
	}
	if c.SoakingStart != nil {
		soakingStart.Time = *c.SoakingStart
		soakingStart.Valid = true
	}
	// Handle nullable blackoutStart
	blackoutStart := pgtype.Timestamp{
		Valid: false,
	}
	if c.BlackoutStart != nil {
		blackoutStart.Time = *c.BlackoutStart
		blackoutStart.Valid = true
	}

	// Handle nullable YieldGrams
	yieldGrams := pgtype.Int4{
		Valid: false,
	}
	if c.YieldGrams != nil {
		yieldGrams.Int32 = int32(*c.YieldGrams)
		yieldGrams.Valid = true
	}

	return database.Crop{
		CropID:        pgtype.UUID{Bytes: c.CropID, Valid: true},
		SeedID:        pgtype.UUID{Bytes: c.SeedID, Valid: true},
		SoakingStart:  soakingStart,
		StackingStart: pgtype.Timestamp{Time: c.StackingStart, Valid: true},
		BlackoutStart: blackoutStart,
		LightsStart:   pgtype.Timestamp{Time: c.LightsStart, Valid: true},
		Harvest:       pgtype.Timestamp{Time: c.Harvest, Valid: true},
		Code:          pgtype.Text{String: c.Code, Valid: true},
		YieldGrams:    yieldGrams,
	}
}

func (crop *Crop) GetNewCropParams() database.NewCropParams {
	var soakingStart pgtype.Timestamp
	if crop.SoakingStart != nil {
		soakingStart = pgtype.Timestamp{
			Time:  *crop.SoakingStart,
			Valid: true,
		}
	} else {
		soakingStart = pgtype.Timestamp{Valid: false}
	}

	var blackoutStart pgtype.Timestamp
	if crop.BlackoutStart != nil {
		blackoutStart = pgtype.Timestamp{
			Time:  *crop.BlackoutStart,
			Valid: true,
		}
	} else {
		blackoutStart = pgtype.Timestamp{Valid: false}
	}

	return database.NewCropParams{
		CropID: pgtype.UUID{
			Bytes: crop.CropID,
			Valid: true,
		},
		SeedID: pgtype.UUID{
			Bytes: crop.SeedID,
			Valid: true,
		},
		SoakingStart:  soakingStart,
		BlackoutStart: blackoutStart,
		StackingStart: pgtype.Timestamp{
			Time:  crop.StackingStart,
			Valid: !crop.StackingStart.IsZero(),
		},
		LightsStart: pgtype.Timestamp{
			Time:  crop.LightsStart,
			Valid: !crop.LightsStart.IsZero(),
		},
		Harvest: pgtype.Timestamp{
			Time:  crop.Harvest,
			Valid: !crop.Harvest.IsZero(),
		},
		Code: pgtype.Text{
			String: crop.Code,
			Valid:  true,
		},
	}
}

func CropFromDatabaseModel(crop database.Crop) Crop {
	// Handle nullable YieldGrams
	var yieldGrams *int
	if crop.YieldGrams.Valid {
		value := int(crop.YieldGrams.Int32)
		yieldGrams = &value
	}

	// Handle nullable SoakingStart
	var soakingStart *time.Time
	if crop.SoakingStart.Valid {
		value := crop.SoakingStart.Time
		soakingStart = &value
	}

	return Crop{
		CropID:        crop.CropID.Bytes,
		SeedID:        crop.SeedID.Bytes,
		SoakingStart:  soakingStart,
		StackingStart: crop.StackingStart.Time,
		LightsStart:   crop.LightsStart.Time,
		Harvest:       crop.Harvest.Time,
		YieldGrams:    yieldGrams,
		Code:          crop.Code.String,
	}
}
