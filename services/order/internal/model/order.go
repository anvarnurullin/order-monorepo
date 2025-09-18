package model

import "slices"

import "time"

type Order struct {
	ID        int64     `json:"id"`
	ProductID int64     `json:"product_id"`
	Quantity  int       `json:"quantity"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

var AllowedStatuses = []string{
	"pending",
	"paid",
	"shipped",
	"completed",
	"canceled",
}

var ValidTransitions = map[string][]string {
	"pending": {"paid", "cancelled"},
	"paid": {"shipped", "cancelled"},
	"shipped": {"completed"},
	"completed": {},
	"cancelled": {},
}

func IsValidStatusTransition(from string, to string) bool {
	allowed, ok := ValidTransitions[from]
	if !ok {
		return false
	}
	return slices.Contains(allowed, to)
}
