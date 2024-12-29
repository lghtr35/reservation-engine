/*
 * Everything involving a mutation belongs to the 'commands' package.
 */
package commands

import (
	"errors"
	"fmt"
	"regexp"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateCustomerCommand struct {
	db      *gorm.DB
	logger  *zerolog.Logger
	name    string
	company string
	email   string
}

func NewCreateCustomerCommand(db *gorm.DB, logger *zerolog.Logger, name, company, email string) *CreateCustomerCommand {
	return &CreateCustomerCommand{db: db, logger: logger, name: name, company: company, email: email}
}

func (s *CreateCustomerCommand) Execute() (string, error) {
	if s.name == "" {
		return "", errors.New("CreateCustomerCommand: Tried creating with empty name")
	}
	s.logger.Debug().Msg("CreateCustomerCommand: Started")

	if ok, err := regexp.Match("", []byte(s.email)); !ok || err != nil {
		if err != nil {
			return "", err
		}
		return "", errors.New("CreateCustomerCommand: Email is not in correct format")
	}

	customer := models.Customer{
		Name:           s.name,
		Company:        s.company,
		Email:          s.email,
		MaxSourceLimit: 1,
	}

	res := s.db.Create(&customer)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("CreateCustomerCommand: Finished with success")

	return customer.ID, nil
}

type DeleteCustomerCommand struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewDeleteCustomerCommand(db *gorm.DB, logger *zerolog.Logger, id string) *DeleteCustomerCommand {
	return &DeleteCustomerCommand{db: db, logger: logger, id: id}
}

func (s *DeleteCustomerCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("DeleteCustomerCommand: Tried deleting with empty id")
	}
	s.logger.Debug().Msg("DeleteCustomerCommand: Started")

	res := s.db.Delete(models.Customer{}, s.id)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("DeleteCustomerCommand: Finished with success")

	return s.id, nil
}

type UpdateCustomerCommand struct {
	db             *gorm.DB
	logger         *zerolog.Logger
	id             string
	name           *string
	email          *string
	company        *string
	maxSourceLimit *int
}

func NewUpdateCustomerCommand(db *gorm.DB, logger *zerolog.Logger, id string, name, email, company *string, maxSourceLimit *int) *UpdateCustomerCommand {
	return &UpdateCustomerCommand{db: db, logger: logger, id: id, name: name, email: email, company: company, maxSourceLimit: maxSourceLimit}
}

func (s *UpdateCustomerCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("UpdateCustomerCommand: Tried updating with empty id")
	}
	s.logger.Debug().Msg("UpdateCustomerCommand: Started")

	var customer models.Customer
	res := s.db.First(&customer, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("UpdateReservationCommand: Could not find the customer with id: %s", s.id)
		}
		return "", res.Error
	}

	if s.name != nil && *s.name != "" {
		customer.Name = *s.name
	}
	if s.email != nil && *s.email != "" {
		customer.Email = *s.email
	}
	if s.company != nil && *s.company != "" {
		customer.Company = *s.company
	}
	if s.maxSourceLimit != nil && *s.maxSourceLimit >= 0 {
		customer.MaxSourceLimit = *s.maxSourceLimit
	}

	res = s.db.Save(&customer)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("UpdateCustomerCommand: Finished with success")

	return s.id, nil
}
