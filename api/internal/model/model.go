package model

import "time"

type User struct {
	ID           string    `db:"id" json:"id"`
	Email        string    `db:"email" json:"email"`
	PasswordHash string    `db:"password_hash" json:"-"`
	Role         string    `db:"role" json:"role"`
	CreatedAt    time.Time `db:"created_at" json:"created_at"`
}

type Workspace struct {
	ID               string    `db:"id" json:"id"`
	Name             string    `db:"name" json:"name"`
	Plan             string    `db:"plan" json:"plan"`
	StripeCustomerID *string   `db:"stripe_customer_id" json:"stripe_customer_id,omitempty"`
	CreatedAt        time.Time `db:"created_at" json:"created_at"`
}