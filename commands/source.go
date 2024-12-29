/*
 * Everything involving a mutation belongs to the 'commands' package.
 */
package commands

import (
	"errors"
	"fmt"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateSourceCommand struct {
	db          *gorm.DB
	logger      *zerolog.Logger
	name        string
	maxDuration string
	customerId  string
}

func NewCreateSourceCommand(db *gorm.DB, logger *zerolog.Logger, name string, maxPossibleDuration string, customerId string) *CreateSourceCommand {
	return &CreateSourceCommand{db: db, logger: logger, name: name, maxDuration: maxPossibleDuration, customerId: customerId}
}

func (s *CreateSourceCommand) Execute() (string, error) {
	if s.name == "" {
		return "", errors.New("CreateSourceCommand: Tried creating with empty name")
	}
	s.logger.Debug().Msg("CreateSourceCommand: Started")

	var customer models.Customer
	res := s.db.Preload("Sources").First(&customer, s.customerId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("CreateSourceCommand: Could not find the customer with id: %s", s.customerId)
		}
		return "", res.Error
	}

	if customer.MaxSourceLimit == len(customer.Sources) {
		return "", fmt.Errorf("CreateSourceCommand: Customer with id %s, has already hit the limit for sources", s.customerId)
	}

	source := models.Source{
		Name:                s.name,
		MaxPossibleDuration: s.maxDuration,
		CustomerID:          s.customerId,
	}

	res = s.db.Create(&source)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("CreateSourceCommand: Finished with success")

	return source.ID, nil
}

type DeleteSourceCommand struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewDeleteSourceCommand(db *gorm.DB, logger *zerolog.Logger, id string) *DeleteSourceCommand {
	return &DeleteSourceCommand{db: db, logger: logger, id: id}
}

func (s *DeleteSourceCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("DeleteSourceCommand: Tried deleting with empty id")
	}
	s.logger.Debug().Msg("DeleteSourceCommand: Started")

	res := s.db.Delete(models.Source{}, s.id)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("DeleteSourceCommand: Finished with success")

	return s.id, nil
}

type UpdateSourceCommand struct {
	db          *gorm.DB
	logger      *zerolog.Logger
	id          string
	name        *string
	maxDuration *string
}

func NewUpdateSourceCommand(db *gorm.DB, logger *zerolog.Logger, id string, name, maxDuration *string) *UpdateSourceCommand {
	return &UpdateSourceCommand{db: db, logger: logger, id: id, name: name, maxDuration: maxDuration}
}

func (s *UpdateSourceCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("UpdateSourceCommand: Tried updating with empty id")
	}
	s.logger.Debug().Msg("UpdateSourceCommand: Started")

	var source models.Source
	res := s.db.First(&source, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("UpdateReservationCommand: Could not find the source with id: %s", s.id)
		}
		return "", res.Error
	}

	if s.name != nil && *s.name != "" {
		source.Name = *s.name
	}

	if s.maxDuration != nil && *s.maxDuration != "" {
		source.Name = *s.name
	}

	res = s.db.Save(&source)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("UpdateSourceCommand: Finished with success")

	return s.id, nil
}
