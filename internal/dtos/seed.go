package dtos

import (
	"lvm/database"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/google/uuid"
)

type Seed struct {
	SeedID uuid.UUID
	Name   string
	Type   string
}

func (s Seed) ToDatabaseModel() database.Seed {
	seed := database.Seed{}
	seed.SeedID = pgtype.UUID{Bytes: s.SeedID, Valid: true}
	seed.Name = pgtype.Text{String: s.Name, Valid: true}
	seed.Type = pgtype.Text{String: s.Type, Valid: true}

	return seed
}

func (s Seed) FromDatabaseModel(seed database.Seed) Seed {
	s.SeedID = seed.SeedID.Bytes
	s.Name = seed.Name.String
	s.Type = seed.Type.String

	return s
}

func (s Seed) ToNewSeedParams() database.NewSeedParams {
	return database.NewSeedParams{
		SeedID: pgtype.UUID{Bytes: s.SeedID, Valid: true},
		Name: pgtype.Text{String: s.Name, Valid: true},
		Type: pgtype.Text{String: s.Type, Valid: true},
	}
}
