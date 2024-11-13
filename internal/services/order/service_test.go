package order

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"lvm/internal/db/repositories"
	"lvm/internal/dtos"
	formDtos "lvm/internal/dtos/form"
	mocks "lvm/internal/mocks/repositories"
)

func TestGenerateCropsFromOrder(t *testing.T) {
	seedId := uuid.New()
	tests := []struct {
		name              string
		order             formDtos.Order
		seedInstructions  dtos.SeedInstruction
		expectedCrops     int
		expectedCrop      dtos.Crop
		generateCodeError error
		expectError       bool
	}{
		{
			name: "harvest date with soaking and blackout",
			order: formDtos.Order{
				SeedID:    seedId,
				Yield:     1000,
				IsHarvest: true,
				Date:      time.Date(2024, 3, 20, 0, 0, 0, 0, time.Local),
				Time:      time.Date(0, 0, 0, 10, 0, 0, 0, time.Local),
			},
			seedInstructions: dtos.SeedInstruction{
				YieldGrams:    500,
				SoakingHours:  12,
				StackingHours: 48,
				BlackoutHours: 24,
				LightsHours:   120,
			},
			expectedCrops: 3,
			expectedCrop: dtos.Crop{
				SeedID:        seedId,
				Harvest:       time.Date(2024, 3, 20, 10, 0, 0, 0, time.Local),
				SoakingStart:  func() *time.Time { t := time.Date(2024, 3, 11, 22, 0, 0, 0, time.Local); return &t }(),
				StackingStart: time.Date(2024, 3, 12, 10, 0, 0, 0, time.Local),
				BlackoutStart: func() *time.Time { t := time.Date(2024, 3, 14, 10, 0, 0, 0, time.Local); return &t }(),
				LightsStart:   time.Date(2024, 3, 15, 10, 0, 0, 0, time.Local),
			},
		},
		   {
		       name: "start date without soaking or blackout",
		       order: formDtos.Order{
		           SeedID:    uuid.New(),
		           Yield:     500,
		           IsHarvest: false,
		           Date:      time.Date(2024, 3, 20, 0, 0, 0, 0, time.Local),
		           Time:      time.Date(0, 0, 0, 10, 0, 0, 0, time.Local),
		       },
		       seedInstructions: dtos.SeedInstruction{
		           YieldGrams:    400,
		           StackingHours: 24,
		           LightsHours:   120,
		       },
		       expectedCrops: 1,
			   expectedCrop: dtos.Crop{
				   SeedID:        seedId,
				   Harvest:       time.Date(2024, 3, 20, 10, 0, 0, 0, time.Local),
				   SoakingStart:  nil,
				   StackingStart: time.Date(2024, 3, 13, 10, 0, 0, 0, time.Local),
				   BlackoutStart: nil,
				   LightsStart:   time.Date(2024, 3, 15, 10, 0, 0, 0, time.Local),
			   },
		   },
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mocks
			mockSeedRepo := &mocks.MockSeedRepository{}
			mockSeedInstructionRepo := &mocks.MockSeedInstructionRepository{}
			mockCropRepo := &mocks.MockCropRepository{}

			// Configure mock behavior
			mockSeedInstructionRepo.On("GetSeedInstructionBySeedId", context.Background(), tt.order.SeedID).
				Return(tt.seedInstructions, nil)

			mockSeedRepo.On("GetSeed", context.Background(), tt.order.SeedID).
				Return(dtos.Seed{Name: "Radish"}, nil)

				mockCropRepo.On("GetExistingCodesForSeed", 
				context.Background(),
				mock.MatchedBy(func(input repositories.GetExistingCodesForSeedInput) bool {
					return input.FirstLetter == "R" && 
						   input.SeedID == tt.order.SeedID
				})).
				Return([]string{}, nil)

			service := NewOrderService(mockSeedRepo, mockSeedInstructionRepo, mockCropRepo)

			// Execute test
			crops, err := service.GenerateCropsFromOrder(tt.order)

			// Assert results
			if tt.expectError {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Len(t, crops, tt.expectedCrops)

			if len(crops) > 0 {
				firstCrop := crops[0]
				assert.Equal(t, tt.expectedCrop.SeedID, firstCrop.SeedID)
				assert.Equal(t, tt.expectedCrop.Harvest, firstCrop.Harvest)
				assert.Equal(t, tt.expectedCrop.SoakingStart, firstCrop.SoakingStart)
				assert.Equal(t, tt.expectedCrop.StackingStart, firstCrop.StackingStart)
				assert.Equal(t, tt.expectedCrop.BlackoutStart, firstCrop.BlackoutStart)
				assert.Equal(t, tt.expectedCrop.LightsStart, firstCrop.LightsStart)
			}

			mockSeedInstructionRepo.AssertExpectations(t)
			mockCropRepo.AssertExpectations(t)
			mockSeedRepo.AssertExpectations(t)
		})
	}
}
