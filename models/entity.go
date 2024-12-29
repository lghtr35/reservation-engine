package models

import "time"

type Base struct {
	ID        string    `gorm:"primarykey;type:uuid;default:gen_random_uuid()" json:"id"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

type Source struct {
	Base
	Name                string        `gorm:"type:nvarchar(256)" json:"name"`
	Tokens              []ApiToken    `json:"tokens"`
	Reservations        []Reservation `json:"reservations"`
	MaxPossibleDuration string        `json:"maxPossibleReservationDuration"`
	CustomerID          string        `json:"customerId"`
}

type ApiToken struct {
	Base
	CustomerID string    `gorm:"type:uuid" json:"customerId"`
	SourceID   string    `gorm:"type:uuid" json:"sourceId"`
	Token      string    `gorm:"type:nvarchar(64)" json:"token"`
	ValidUntil time.Time `json:"validUntil"`
}

type Reservation struct {
	Base
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	ReserverID string    `json:"reserverId"`
	ReserveeID string    `json:"reserveeId"`
	SourceID   string    `json:"sourceId"`
}

type Customer struct {
	Base
	Name           string     `gorm:"type:nvarchar(128)" json:"name"`
	Company        string     `gorm:"type:nvarchar(64)" json:"company"`
	Email          string     `gorm:"type:nvarchar(128)" json:"email"`
	Sources        []string   `json:"sources"`
	ApiTokens      []ApiToken `json:"apiTokens"`
	Secret         Secret     `json:"secret"`
	MaxSourceLimit int        `json:"maxSourceLimit"`
}

type Secret struct {
	Base
	CustomerID string `gorm:"type:uuid" json:"customerId"`
	Value      string `gorm:"type:nvarchar(64)" json:"secret"`
}
