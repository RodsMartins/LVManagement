package services

import (
	"context"
	"errors"
	"fmt"
	"log"
	"lvm/internal/db/repositories"
	"lvm/internal/dtos"
	formDtos "lvm/internal/dtos/form"
	"math"
	"time"

	"github.com/google/uuid"
)

const YIELD_ERROR_MARGIN_PERCENTAGE = 10

type OrderService struct {
	codeGenerator             *CodeGenerator
	seedRepository            repositories.SeedRepository
	seedInstructionRepository repositories.SeedInstructionRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(
	SeedRepository repositories.SeedRepository,
	SeedInstructionsRepository repositories.SeedInstructionRepository,
) *OrderService {
	return &OrderService{
		codeGenerator:  NewCodeGenerator(SeedRepository, SeedInstructionsRepository),
		seedRepository: SeedRepository,
	}
}

func (o OrderService) GenerateCropsFromOrder(order formDtos.Order) ([]dtos.Crop, error) {
	ctx := context.Background()

	crops := []dtos.Crop{}

	// Get seed instructions
	seedInstructions, err := o.seedInstructionRepository.GetSeedInstructionsBySeedId(ctx, order.SeedID)
	if err != nil {
		return crops, errors.New("unable to get seed instructions")
	}
	seedInstructionsModel := dtos.SeedInstruction{}.FromDatabaseModel(seedInstructions)

	// Calculate number of trays needed
	yieldWithMargin := seedInstructionsModel.YieldGrams - (seedInstructionsModel.YieldGrams * YIELD_ERROR_MARGIN_PERCENTAGE / 100)
	traysToPlant := int(math.Ceil(float64(order.Yield) / float64(yieldWithMargin)))

	// Calculate dates based on whether it's a harvest date or start date
	orderDateTime := order.GetDateTime()
	var (
		harvestTime   time.Time
		lightsStart   time.Time
		stackingStart time.Time
		soakingStart  *time.Time
		blackoutStart *time.Time
	)

	if order.IsHarvest {
		// Working backwards from harvest date
		harvestTime = orderDateTime

		// Start from harvest and work backwards through lights period
		lightsStart = harvestTime.Add(-time.Duration(seedInstructionsModel.LightsHours) * time.Hour)

		// Create a separate variable for stacking calculation
		var stackingBase time.Time = lightsStart

		if seedInstructionsModel.BlackoutHours > 0 {
			// Calculate blackout start time
			blackoutTime := lightsStart.Add(-time.Duration(seedInstructionsModel.BlackoutHours) * time.Hour)
			blackoutStart = &blackoutTime
			// Use the blackout time for stacking calculations
			stackingBase = blackoutTime // Create a new value, not a pointer reference
		}

		// Calculate stacking start from stackingBase
		stackingStart = stackingBase.Add(-time.Duration(seedInstructionsModel.StackingHours) * time.Hour)

		// If soaking is needed, calculate backwards from stacking
		if seedInstructionsModel.SoakingHours > 0 {
			soakTime := stackingStart.Add(-time.Duration(seedInstructionsModel.SoakingHours) * time.Hour)
			soakingStart = &soakTime
		}
		// Print all variables using log
		log.Printf("harvestTime: %v\n", harvestTime)
		log.Printf("lightsStart: %v\n", lightsStart)
		log.Printf("stackingStart: %v\n", stackingStart)
		if soakingStart != nil {
			log.Printf("soakingStart: %v\n", *soakingStart)
		} else {
			log.Println("soakingStart: nil")
		}
		if blackoutStart != nil {
			log.Printf("blackoutStart: %v\n", *blackoutStart)
		} else {
			log.Println("blackoutStart: nil")
		}
	} else {
		// Working forwards from start date
		// Handle soaking phase if needed
		if seedInstructionsModel.SoakingHours > 0 {
			// Create a new time value for soaking
			soakTime := orderDateTime // Create new value
			soakingStart = &soakTime  // Store pointer to new value
			stackingStart = orderDateTime.Add(time.Duration(seedInstructionsModel.SoakingHours) * time.Hour)
		} else {
			stackingStart = orderDateTime
		}

		if seedInstructionsModel.BlackoutHours > 0 {
			blackoutTime := stackingStart.Add(time.Duration(seedInstructionsModel.StackingHours) * time.Hour)
			blackoutStart = &blackoutTime // Store pointer to new value
			lightsStart = blackoutTime.Add(time.Duration(seedInstructionsModel.BlackoutHours) * time.Hour)
		} else {
			lightsStart = stackingStart.Add(time.Duration(seedInstructionsModel.StackingHours) * time.Hour)
		}

		// Calculate harvest time from end of lights period
		harvestTime = lightsStart.Add(time.Duration(seedInstructionsModel.LightsHours) * time.Hour)
	}

	codeStart := stackingStart
	if soakingStart != nil {
		codeStart = *soakingStart
	}

	code, err := o.GenerateCode(order.SeedID, codeStart, harvestTime)
	if err != nil {
		return crops, fmt.Errorf("error generating code: %w", err)
	}

	// Generate crops
	for i := 0; i < traysToPlant; i++ {
		crop := dtos.NewCrop(
			uuid.UUID(order.SeedID),
			soakingStart,
			stackingStart,
			blackoutStart,
			lightsStart,
			harvestTime,
			code,
			nil,
		)

		crops = append(crops, crop)
	}

	return crops, nil
}

// GenerateCode is now a thin wrapper around the CodeGenerator
func (o *OrderService) GenerateCode(seedId uuid.UUID, start time.Time, end time.Time) (string, error) {
	return o.codeGenerator.GenerateCode(context.Background(), struct {
		SeedID    uuid.UUID
		StartDate time.Time
		EndDate   time.Time
	}{
		SeedID:    seedId,
		StartDate: start,
		EndDate:   end,
	})
}
