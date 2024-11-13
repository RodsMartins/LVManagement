package dtos

import (
	"lvm/database"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

type SeedInstruction struct {
	SeedInstructionID uuid.UUID
	SeedID            uuid.UUID
	SeedGrams         int
	SoakingHours      int
	StackingHours     int
	BlackoutHours     int
	LightsHours       int
	YieldGrams        int
	SpecialTreatment  string
}

func (s SeedInstruction) ToDatabaseModel() database.SeedInstruction {
	seedInstruction := database.SeedInstruction{}
	seedInstruction.SeedInstructionID = pgtype.UUID{Bytes: s.SeedInstructionID, Valid: true}
	seedInstruction.SeedID = pgtype.UUID{Bytes: s.SeedID, Valid: true}
	seedInstruction.SeedGrams = pgtype.Int4{Int32: int32(s.SeedGrams), Valid: true}
	seedInstruction.SoakingHours = pgtype.Int4{Int32: int32(s.SoakingHours), Valid: true}
	seedInstruction.StackingHours = pgtype.Int4{Int32: int32(s.StackingHours), Valid: true}
	seedInstruction.BlackoutHours = pgtype.Int4{Int32: int32(s.BlackoutHours), Valid: true}
	seedInstruction.LightsHours = pgtype.Int4{Int32: int32(s.LightsHours), Valid: true}
	seedInstruction.YieldGrams = pgtype.Int4{Int32: int32(s.YieldGrams), Valid: true}
	seedInstruction.SpecialTreatment = pgtype.Text{String: s.SpecialTreatment, Valid: true}

	return seedInstruction
}

func (s SeedInstruction) FromDatabaseModel(seedInstruction database.SeedInstruction) SeedInstruction {
	s.SeedInstructionID = seedInstruction.SeedInstructionID.Bytes
	s.SeedID = seedInstruction.SeedID.Bytes
	s.SeedGrams = int(seedInstruction.SeedGrams.Int32)
	s.SoakingHours = int(seedInstruction.SoakingHours.Int32)
	s.StackingHours = int(seedInstruction.StackingHours.Int32)
	s.BlackoutHours = int(seedInstruction.BlackoutHours.Int32)
	s.LightsHours = int(seedInstruction.LightsHours.Int32)
	s.YieldGrams = int(seedInstruction.YieldGrams.Int32)
	s.SpecialTreatment = seedInstruction.SpecialTreatment.String

	return s
}

func (s SeedInstruction) ToNewSeedInstructionParams() database.NewSeedInstructionParams {
	return database.NewSeedInstructionParams{
		SeedInstructionID: pgtype.UUID{Bytes: s.SeedInstructionID, Valid: true},
		SeedID:            pgtype.UUID{Bytes: s.SeedID, Valid: true},
		SeedGrams:         pgtype.Int4{Int32: int32(s.SeedGrams), Valid: true},
		SoakingHours:      pgtype.Int4{Int32: int32(s.SoakingHours), Valid: true},
		StackingHours:     pgtype.Int4{Int32: int32(s.StackingHours), Valid: true},
		BlackoutHours:     pgtype.Int4{Int32: int32(s.BlackoutHours), Valid: true},
		LightsHours:       pgtype.Int4{Int32: int32(s.LightsHours), Valid: true},
		YieldGrams:        pgtype.Int4{Int32: int32(s.YieldGrams), Valid: true},
		SpecialTreatment:  pgtype.Text{String: s.SpecialTreatment, Valid: true},
	}
}

func (s SeedInstruction) ToUpdateSeedInstructionParams() database.UpdateSeedInstructionParams {
	return database.UpdateSeedInstructionParams{
		SeedInstructionID: pgtype.UUID{Bytes: s.SeedInstructionID, Valid: true},
		SeedID:            pgtype.UUID{Bytes: s.SeedID, Valid: true},
		SeedGrams:         pgtype.Int4{Int32: int32(s.SeedGrams), Valid: true},
		SoakingHours:      pgtype.Int4{Int32: int32(s.SoakingHours), Valid: true},
		StackingHours:     pgtype.Int4{Int32: int32(s.StackingHours), Valid: true},
		BlackoutHours:     pgtype.Int4{Int32: int32(s.BlackoutHours), Valid: true},
		LightsHours:       pgtype.Int4{Int32: int32(s.LightsHours), Valid: true},
		YieldGrams:        pgtype.Int4{Int32: int32(s.YieldGrams), Valid: true},
		SpecialTreatment:  pgtype.Text{String: s.SpecialTreatment, Valid: true},
	}
}
