package order

import (
	"context"
	"fmt"
	"lvm/internal/db/repositories"
	"regexp"
	"strconv"
	"time"

	"github.com/google/uuid"
)

// GenerateCodeInput represents the input parameters for code generation
type GenerateCodeInput struct {
	SeedID    uuid.UUID
	StartDate time.Time
	EndDate   time.Time
}

// CodeGenerator handles the business logic for generating codes
type CodeGenerator struct {
	seedRepository repositories.SeedRepository
	cropRepository repositories.CropRepository
}

// newCodeGenerator creates a new CodeGenerator
func newCodeGenerator(
	seedRepo repositories.SeedRepository,
	cropRepo repositories.CropRepository,
) *CodeGenerator {
	return &CodeGenerator{
		seedRepository: seedRepo,
		cropRepository: cropRepo,
	}
}

// GenerateCode generates a unique code for a crop based on the seed type and date range
func (g *CodeGenerator) GenerateCode(ctx context.Context, input GenerateCodeInput) (string, error) {
	seed, err := g.seedRepository.GetSeed(ctx, input.SeedID)
	if err != nil {
		return "", fmt.Errorf("get seed: %w", err)
	}

	if seed.Name == "" {
		return "", fmt.Errorf("invalid seed name")
	}

	codes, err := g.cropRepository.GetExistingCodesForSeed(ctx, repositories.GetExistingCodesForSeedInput{
		FirstLetter: string(seed.Name[0]),
		SeedID:      input.SeedID,
		StartDate:   input.StartDate,
		EndDate:     input.EndDate,
	})

	if err != nil {
		return "", fmt.Errorf("get existing codes: %w", err)
	}

	stringCodes := make([]string, 0, len(codes))
	stringCodes = append(stringCodes, codes...)

	nextNum := findNextAvailableNumber(stringCodes)
	return fmt.Sprintf("%c%d", seed.Name[0], nextNum), nil
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
