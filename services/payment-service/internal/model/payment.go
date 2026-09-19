package model

import "time"

// Subscription represents a user's subscription plan.
// This model is owned exclusively by the Payment Service.
type Subscription struct {
	ID          string     `json:"id" bson:"_id,omitempty"`
	UserID      string     `json:"userId" bson:"user_id"`
	PlanID      string     `json:"planId" bson:"plan_id"`
	Status      string     `json:"status" bson:"status"`
	StartsAt    time.Time  `json:"startsAt" bson:"starts_at"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty" bson:"expires_at,omitempty"`
	CancelledAt *time.Time `json:"cancelledAt,omitempty" bson:"cancelled_at,omitempty"`
	CreatedAt   time.Time  `json:"createdAt" bson:"created_at"`
	UpdatedAt   time.Time  `json:"updatedAt" bson:"updated_at"`
}

// Subscription statuses
const (
	SubscriptionStatusActive    = "active"
	SubscriptionStatusCancelled = "cancelled"
	SubscriptionStatusExpired   = "expired"
	SubscriptionStatusPastDue   = "past_due"
)

// Plan represents a pricing plan.
type Plan struct {
	ID          string `json:"id" bson:"_id,omitempty"`
	Name        string `json:"name" bson:"name"`
	Description string `json:"description" bson:"description"`
	PriceMonthly int64 `json:"priceMonthly" bson:"price_monthly"` // In smallest currency unit (e.g., paise)
	PriceYearly  int64 `json:"priceYearly" bson:"price_yearly"`
	Features    []string `json:"features" bson:"features"`
	IsActive    bool   `json:"isActive" bson:"is_active"`
}
