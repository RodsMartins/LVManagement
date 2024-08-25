package db

import (
	"github.com/google/uuid"
	"database/sql"
	//"time"
)

type Seed struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type string    `json:"type"`
}

type SeedInstructions struct {
	ID               uuid.UUID `json:"id"`
	Seed             uuid.UUID `json:"seed"`
	SeedGrams        int       `json:"seed_grams"`
	SoakingHours     int       `json:"soaking_hours"`
	StackingHours    int       `json:"stacking_hours"`
	BlackoutHours    int       `json:"blackout_hours"`
	LightsHours      int       `json:"lights_hours"`
	YieldGrams       int       `json:"yield_grams"`
	SpecialTreatment string    `json:"special_treatment"`
}

type Crop struct {
	ID            uuid.UUID `json:"id"`
	Seed          uuid.UUID `json:"seed"`
	SoakingStart  sql.NullTime `json:"soaking_start"`
	SoakingEnd    sql.NullTime `json:"soaking_end"`
	StackingStart sql.NullTime `json:"stacking_start"`
	StackingEnd   sql.NullTime `json:"stacking_end"`
	LightsStart   sql.NullTime `json:"lights_start"`
	LightsEnd     sql.NullTime `json:"lights_end"`
}

// Watering represents a watering event.
type Watering struct {
	ID           uuid.UUID `json:"id"`
	Crop         uuid.UUID `json:"crop"`
	QuantityML   uuid.UUID `json:"quantity_ml"` // Assuming quantity_ml is a unique identifier
	FertilizerML int       `json:"fertilizer_ml"`
	Fertilizer   uuid.UUID `json:"fertilizer"`
}

// Fertilizer represents a fertilizer.
type Fertilizer struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// Supplier represents a supplier.
type Supplier struct {
	ID    uuid.UUID `json:"id"`
	Name  string    `json:"name"`
	URL   string    `json:"url"`
	Swiss bool      `json:"swiss"`
}

// Order represents an order.
type Order struct {
	ID       uuid.UUID `json:"id"`
	Total    float64   `json:"total"`
	Shipping float64   `json:"shipping"`
}

// OrderItem represents an item in an order.
type OrderItem struct {
	ID       uuid.UUID `json:"id"`
	Order    uuid.UUID `json:"order"`
	Type     uuid.UUID `json:"type"`
	Name     string    `json:"name"`
	Seed     uuid.UUID `json:"seed"`
	Material uuid.UUID `json:"material"`
	Supplier uuid.UUID `json:"supplier"`
	Price    uuid.UUID `json:"price"` // Assuming price is a unique identifier
}

// Material represents a material.
type Material struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
	Type uuid.UUID `json:"type"`
	Main bool      `json:"main"`
}

// MaterialType represents a material type.
type MaterialType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

// OrderType represents an order type.
type OrderType struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}
