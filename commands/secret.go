package commands

import (
	"errors"
	"fmt"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/lghtr35/reservation-engine/util"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateSecretCommand struct {
	db         *gorm.DB
	logger     *zerolog.Logger
	hasher     *util.Hasher
	customerId string
}

func NewCreateSecretCommand(db *gorm.DB, logger *zerolog.Logger, hasher *util.Hasher, customerId string) *CreateSecretCommand {
	return &CreateSecretCommand{db: db, logger: logger, hasher: hasher, customerId: customerId}
}

func (s *CreateSecretCommand) Execute() (string, error) {
	if s.customerId == "" {
		return "", errors.New("CreateSecretCommand: missing arguments")
	}
	s.logger.Debug().Msg("CreateSecretCommand: Started")

	var customer models.Customer
	res := s.db.First(&customer, s.customerId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("CreateSecretCommand: Could not find the customer with id: %s", s.customerId)
		}
		return "", res.Error
	}

	hashed, err := s.hasher.GetHash(fmt.Sprintf("%s:%s", customer.Company, util.GetRandString(16)))
	if err != nil {
		return "", err
	}

	apiSecret := models.Secret{
		CustomerID: s.customerId,
		Value:      hashed,
	}
	res = s.db.Create(&apiSecret)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("CreateSecretCommand: Finished with success")
	return apiSecret.ID, nil
}
