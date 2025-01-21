package dtos

import (
	"fmt"
	"lvm/database"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type CropStage int

const (
    SoakingStage CropStage = iota
    StackingStage
    BlackoutStage
    LightsStage
    HarvestStage
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

func (c *Crop) GetStageStartDate(stage CropStage) (time.Time, error) {
    switch stage {
    case SoakingStage:
        if c.SoakingStart == nil {
            return time.Time{}, fmt.Errorf("soaking start is not set")
        }
        return *c.SoakingStart, nil
    case StackingStage:
        return c.StackingStart, nil
    case BlackoutStage:
        if c.BlackoutStart == nil {
            return time.Time{}, fmt.Errorf("blackout start is not set")
        }
        return *c.BlackoutStart, nil
    case LightsStage:
        return c.LightsStart, nil
    case HarvestStage:
        return c.Harvest, nil
    default:
        return time.Time{}, fmt.Errorf("invalid stage")
    }
}

func (c *Crop) GetStageEndDate(stage CropStage) (time.Time, error) {
    switch stage {
    case SoakingStage:
        if c.SoakingStart == nil {
            return time.Time{}, fmt.Errorf("soaking start is not set")
        }
        return c.StackingStart, nil
    case StackingStage:
        if c.BlackoutStart == nil {
            // If no blackout stage, use lights start as end
            return c.LightsStart, nil
        }
        return *c.BlackoutStart, nil
    case BlackoutStage:
        if c.BlackoutStart == nil {
            return time.Time{}, fmt.Errorf("blackout start is not set")
        }
        return c.LightsStart, nil
    case LightsStage:
        return c.Harvest, nil
    case HarvestStage:
        return c.Harvest, nil
    default:
        return time.Time{}, fmt.Errorf("invalid stage")
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

func (c *Crop) GetActiveStages(date time.Time) []CropStage {
    var stages []CropStage
    
    // Convert input date to UTC
    checkDateUTC := date.UTC()
    
    // Helper function to check if date is between stage transitions
    isActiveDuring := func(stageStart, nextStageStart *time.Time) bool {
        if stageStart == nil {
            return false
        }
        
        // Convert stage date to UTC
        stageStartUTC := stageStart.UTC()
        
        // Check if checkDate is after or equal to stage start
        if !checkDateUTC.Before(stageStartUTC) {
            // If there's no next stage, this stage is still active
            if nextStageStart == nil {
                return true
            }
            
            // Convert next stage date to UTC
            nextStageStartUTC := nextStageStart.UTC()
            
            // Check if checkDate is before next stage
            return checkDateUTC.Before(nextStageStartUTC)
        }
        
        return false
    }
    
    // Check Soaking Stage
    if isActiveDuring(c.SoakingStart, &c.StackingStart) {
        stages = append(stages, SoakingStage)
    }
    
    // Check Stacking Stage
    var nextStage *time.Time
    if c.BlackoutStart != nil {
        nextStage = c.BlackoutStart
    } else {
        nextStage = &c.LightsStart
    }
    if isActiveDuring(&c.StackingStart, nextStage) {
        stages = append(stages, StackingStage)
    }
    
    // Check Blackout Stage
    if c.BlackoutStart != nil && isActiveDuring(c.BlackoutStart, &c.LightsStart) {
        stages = append(stages, BlackoutStage)
    }
    
    // Check Lights Stage
    if isActiveDuring(&c.LightsStart, &c.Harvest) {
        stages = append(stages, LightsStage)
    }
    
    // Check Harvest Stage (assuming it's the final stage with no end date)
    if isActiveDuring(&c.Harvest, nil) {
        stages = append(stages, HarvestStage)
    }
    
    return stages
}

func (c *Crop) GetNewStages(date time.Time) []CropStage {
    var stages []CropStage
    
    // Convert input date to UTC, keeping the same instant in time
    checkDateUTC := date.UTC()
    
    // Helper function to normalize stage dates to UTC and check if they occur on the same day
    isSameDay := func(stageDate *time.Time) bool {
        if stageDate == nil {
            return false
        }
        
        // Convert stage date to UTC
        stageDateUTC := stageDate.UTC()
        
        // Check if dates are on the same day in UTC
        return checkDateUTC.Year() == stageDateUTC.Year() &&
            checkDateUTC.Month() == stageDateUTC.Month() &&
            checkDateUTC.Day() == stageDateUTC.Day()
    }
    
    // Check each stage using the helper function
    if c.SoakingStart != nil && isSameDay(c.SoakingStart) {
        stages = append(stages, SoakingStage)
    }
    
    if isSameDay(&c.StackingStart) {
        stages = append(stages, StackingStage)
    }
    
    if c.BlackoutStart != nil && isSameDay(c.BlackoutStart) {
        stages = append(stages, BlackoutStage)
    }
    
    if isSameDay(&c.LightsStart) {
        stages = append(stages, LightsStage)
    }
    
    if isSameDay(&c.Harvest) {
        stages = append(stages, HarvestStage)
    }
    
    return stages
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
