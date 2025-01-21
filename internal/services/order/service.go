package order

import (
	"context"
	"errors"
	"fmt"
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
	cropRepository            repositories.CropRepository
}

// NewOrderService creates a new OrderService
func NewOrderService(
	SeedRepository repositories.SeedRepository,
	SeedInstructionsRepository repositories.SeedInstructionRepository,
	CropRepository repositories.CropRepository,
) *OrderService {
	return &OrderService{
		codeGenerator:             newCodeGenerator(SeedRepository, CropRepository),
		seedRepository:            SeedRepository,
		seedInstructionRepository: SeedInstructionsRepository,
		cropRepository:            CropRepository,
	}
}

func (o OrderService) GenerateCropsFromOrder(order formDtos.Order) ([]dtos.Crop, error) {
	ctx := context.Background()

	crops := []dtos.Crop{}

	// Get seed instructions
	seedInstruction, err := o.seedInstructionRepository.GetSeedInstructionBySeedId(ctx, order.SeedID)
	if err != nil {
		return crops, errors.New("unable to get seed instructions: " + err.Error())
	}

	// Calculate number of trays needed
	yieldWithMargin := seedInstruction.YieldGrams - (seedInstruction.YieldGrams * YIELD_ERROR_MARGIN_PERCENTAGE / 100)
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
		lightsStart = harvestTime.Add(-time.Duration(seedInstruction.LightsHours) * time.Hour)

		// Create a separate variable for stacking calculation
		var stackingBase time.Time = lightsStart

		if seedInstruction.BlackoutHours > 0 {
			// Calculate blackout start time
			blackoutTime := lightsStart.Add(-time.Duration(seedInstruction.BlackoutHours) * time.Hour)
			blackoutStart = &blackoutTime
			// Use the blackout time for stacking calculations
			stackingBase = blackoutTime // Create a new value, not a pointer reference
		}

		// Calculate stacking start from stackingBase
		stackingStart = stackingBase.Add(-time.Duration(seedInstruction.StackingHours) * time.Hour)

		// If soaking is needed, calculate backwards from stacking
		if seedInstruction.SoakingHours > 0 {
			soakTime := stackingStart.Add(-time.Duration(seedInstruction.SoakingHours) * time.Hour)
			soakingStart = &soakTime
		}
	} else {
		// Working forwards from start date
		if seedInstruction.SoakingHours > 0 {
			soakTime := orderDateTime // Create new value
			soakingStart = &soakTime  // Store pointer to new value
			stackingStart = orderDateTime.Add(time.Duration(seedInstruction.SoakingHours) * time.Hour)
		} else {
			stackingStart = orderDateTime
		}

		if seedInstruction.BlackoutHours > 0 {
			blackoutTime := stackingStart.Add(time.Duration(seedInstruction.StackingHours) * time.Hour)
			blackoutStart = &blackoutTime // Store pointer to new value
			lightsStart = blackoutTime.Add(time.Duration(seedInstruction.BlackoutHours) * time.Hour)
		} else {
			lightsStart = stackingStart.Add(time.Duration(seedInstruction.StackingHours) * time.Hour)
		}

		// Calculate harvest time from end of lights period
		harvestTime = lightsStart.Add(time.Duration(seedInstruction.LightsHours) * time.Hour)
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
	return o.codeGenerator.GenerateCode(context.Background(), GenerateCodeInput{
		SeedID:    seedId,
		StartDate: start,
		EndDate:   end,
	})
}
