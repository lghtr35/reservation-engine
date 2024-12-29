package models

import "time"

type ReadAllCustomers struct {
	Pagination Pagination `json:"pagination"`
	IDs        *[]string  `json:"ids"`
	Name       *string    `json:"name"`
}

type ReadAllSources struct {
	Pagination Pagination `json:"pagination"`
	IDs        *[]string  `json:"ids"`
	Name       *string    `json:"name"`
}

type ReadAllReservations struct {
	Pagination Pagination `json:"pagination"`
	IDs        *[]string  `json:"ids"`
	ReserverID *string    `json:"reserverId"`
	ReserveeID *string    `json:"reserveeId"`
	SourceID   *string    `json:"sourceId"`
}

type CreateCustomer struct {
	Name    string `json:"name"`
	Company string `json:"company"`
	Email   string `json:"email"`
}

type CreateSource struct {
	Name                string `json:"name"`
	MaxPossibleDuration string `json:"maxPossibleReservationDuration"`
	CustomerID          string `json:"customerId"`
}

type CreateReservation struct {
	From       time.Time `json:"from"`
	To         time.Time `json:"to"`
	ReserverID string    `json:"reserverId"`
	ReserveeID string    `json:"reserveeId"`
	SourceID   string    `json:"sourceId"`
}

type UpdateCustomer struct {
	ID             string  `json:"id" binding:"required"`
	Name           *string `json:"name"`
	Company        *string `json:"company"`
	Email          *string `json:"email"`
	MaxSourceLimit *int    `json:"maxSourceLimit"`
}

type UpdateSource struct {
	ID                  string  `json:"id" binding:"required"`
	Name                *string `json:"name"`
	MaxPossibleDuration *string `json:"maxPossibleReservationDuration"`
}

type UpdateReservation struct {
	ID   string     `json:"id" binding:"required"`
	From *time.Time `json:"from"`
	To   *time.Time `json:"to"`
}
