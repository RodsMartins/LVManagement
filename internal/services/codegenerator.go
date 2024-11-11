package services

import (
	"context"
	"fmt"
	"lvm/database"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// GenerateCodeInput represents the input parameters for code generation
type GenerateCodeInput struct {
    SeedID    uuid.UUID
    StartDate time.Time
    EndDate   time.Time
}

// CodeGenerator handles the business logic for generating codes
type CodeGenerator struct {
    queries interface {
        GetSeed(ctx context.Context, id pgtype.UUID) (database.Seed, error)
        GetExistingCodesForSeed(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error)
    }
}

// NewCodeGenerator creates a new CodeGenerator
func NewCodeGenerator(queries interface {
    GetSeed(ctx context.Context, id pgtype.UUID) (database.Seed, error)
    GetExistingCodesForSeed(ctx context.Context, params database.GetExistingCodesForSeedParams) ([]pgtype.Text, error)
}) *CodeGenerator {
    return &CodeGenerator{queries: queries}
}

// GenerateCode generates a unique code for a crop based on the seed type and date range
func (g *CodeGenerator) GenerateCode(ctx context.Context, input GenerateCodeInput) (string, error) {
    seed, err := g.queries.GetSeed(ctx, pgtype.UUID{Bytes: input.SeedID, Valid: true})
    if err != nil {
        return "", fmt.Errorf("get seed: %w", err)
    }

    if seed.Name.String == "" {
        return "", fmt.Errorf("invalid seed name")
    }

    params := database.GetExistingCodesForSeedParams{
        FirstLetter: pgtype.Text{String: string(seed.Name.String[0]), Valid: true},
        StartDate:  pgtype.Timestamp{Time: input.StartDate, Valid: true},
        EndDate:    pgtype.Timestamp{Time: input.EndDate, Valid: true},
        SeedID:    pgtype.UUID{Bytes: input.SeedID, Valid: true},
    }

    codes, err := g.queries.GetExistingCodesForSeed(ctx, params)
    if err != nil {
        return "", fmt.Errorf("get existing codes: %w", err)
    }

    stringCodes := make([]string, 0, len(codes))
    for _, code := range codes {
        if code.Valid {
            stringCodes = append(stringCodes, code.String)
        }
    }

    nextNum := findNextAvailableNumber(stringCodes)
    return fmt.Sprintf("%c%d", seed.Name.String[0], nextNum), nil
}

// findNextAvailableNumber finds the next available number for a code
func findNextAvailableNumber(codes []string) int {
    if len(codes) == 0 {
        return 1
    }

    numbers := make(map[int]bool)
    pattern := regexp.MustCompile(`^[A-Z](\d+)$`)

    for _, code := range codes {
        matches := pattern.FindStringSubmatch(code)
        if len(matches) == 2 {
            if num, err := strconv.Atoi(matches[1]); err == nil {
                numbers[num] = true
            }
        }
    }

    for i := 1; ; i++ {
        if !numbers[i] {
            return i
        }
    }
}