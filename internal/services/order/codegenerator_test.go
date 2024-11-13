package order

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/stretchr/testify/assert"

	"lvm/database" // adjust import path
	mocks "lvm/internal/mocks/repositories"
)

func TestGenerateCode(t *testing.T) {
    tests := []struct {
        name          string
        seed          database.Seed
        existingCodes []string
        expectedCode  string
        getSeedErr    error
        getCodesErr   error
    }{
        {
            name: "successful generation - first code",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []string{},
            expectedCode: "R1",
        },
        {
            name: "successful generation - with gaps",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []string{"R1", "R3"},
            expectedCode: "R2",
        },
        {
            name: "handle invalid codes",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []string{"R1", "", "R3"},
            expectedCode: "R2",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            seedId := uuid.New()
            startTime := time.Now()
            endTime := time.Now().Add(24 * time.Hour)
            // Create mock repos
            mockSeedRepo := &mocks.MockSeedRepository{}
            mockSeedRepo.On("GetSeed", context.Background(), seedId).Return(tt.seed, tt.getSeedErr)

            mockCropRepo := &mocks.MockCropRepository{}
            mockCropRepo.On("GetExistingCodesForSeed", context.Background(), tt.seed.Name.String, startTime, endTime, seedId).Return(tt.existingCodes, tt.getCodesErr)
            
            mockInstructionRepo := &mocks.MockSeedInstructionRepository{}
        
            // Create service with mock queries
            service := NewOrderService(mockSeedRepo, mockInstructionRepo, mockCropRepo)

            // Test the service
            code, err := service.GenerateCode(
                seedId,
                startTime,
                endTime,
            )

            if tt.getSeedErr != nil || tt.getCodesErr != nil {
                assert.Error(t, err)
                return
            }

            assert.NoError(t, err)
            assert.Equal(t, tt.expectedCode, code)
        })
    }
}

func TestFindNextAvailableNumber(t *testing.T) {
    tests := []struct {
        name     string
        codes    []string
        expected int
    }{
        {
            name:     "empty codes",
            codes:    []string{},
            expected: 1,
        },
        {
            name:     "sequential codes",
            codes:    []string{"R1", "R2", "R3"},
            expected: 4,
        },
        {
            name:     "gap in codes",
            codes:    []string{"R1", "R3", "R4"},
            expected: 2,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            result := findNextAvailableNumber(tt.codes)
            assert.Equal(t, tt.expected, result)
        })
    }
}