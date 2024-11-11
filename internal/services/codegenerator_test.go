package services

import (
    "context"
    "testing"
    "time"
    
    "github.com/google/uuid"
    "github.com/jackc/pgx/v5/pgtype"
    "github.com/stretchr/testify/assert"
    
    "lvm/database" // adjust import path
)

// MockQueries implements database.Queries for testing
type MockQueries struct {
    // Embed the Queries interface to implement all methods with defaults
    database.Queries

    // Override only the methods we need for testing
    mockGetSeed func(ctx context.Context, id pgtype.UUID) (database.Seed, error)
    mockGetExistingCodesForSeed func(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error)
}

// Override the specific methods we want to mock
func (m *MockQueries) GetSeed(ctx context.Context, id pgtype.UUID) (database.Seed, error) {
    return m.mockGetSeed(ctx, id)
}

func (m *MockQueries) GetExistingCodesForSeed(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error) {
    return m.mockGetExistingCodesForSeed(ctx, params)
}

func TestGenerateCode(t *testing.T) {
    tests := []struct {
        name          string
        seed          database.Seed
        existingCodes []pgtype.Text
        expectedCode  string
        getSeedErr    error
        getCodesErr   error
    }{
        {
            name: "successful generation - first code",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []pgtype.Text{},
            expectedCode: "R1",
        },
        {
            name: "successful generation - with gaps",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []pgtype.Text{
                {String: "R1", Valid: true},
                {String: "R3", Valid: true},
            },
            expectedCode: "R2",
        },
        {
            name: "handle invalid codes",
            seed: database.Seed{Name: pgtype.Text{String: "Radish", Valid: true}},
            existingCodes: []pgtype.Text{
                {String: "R1", Valid: true},
                {String: "", Valid: false},
                {String: "R3", Valid: true},
            },
            expectedCode: "R2",
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock queries
            mockQueries := &MockQueries{
                mockGetSeed: func(ctx context.Context, id pgtype.UUID) (database.Seed, error) {
                    return tt.seed, tt.getSeedErr
                },
                mockGetExistingCodesForSeed: func(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error) {
                    return tt.existingCodes, tt.getCodesErr
                },
            }

            // Create service with mock queries
            service := NewOrderService(mockQueries)

            // Test the service
            code, err := service.GenerateCode(
                uuid.New(),
                time.Now(),
                time.Now().Add(24 * time.Hour),
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

// If you still have the findNextAvailableNumber function
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