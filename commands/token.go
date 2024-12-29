package commands

import (
	"errors"
	"fmt"
	"time"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/lghtr35/reservation-engine/util"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

type CreateApiTokenCommand struct {
	db         *gorm.DB
	logger     *zerolog.Logger
	hasher     *util.Hasher
	customerId string
	sourceId   string
}

func NewCreateApiTokenCommand(db *gorm.DB, logger *zerolog.Logger, hasher *util.Hasher, customerId, sourceId string) *CreateApiTokenCommand {
	return &CreateApiTokenCommand{db: db, logger: logger, hasher: hasher, customerId: customerId, sourceId: sourceId}
}

func (s *CreateApiTokenCommand) Execute() (string, error) {
	if s.customerId == "" || s.sourceId == "" {
		return "", errors.New("CreateApiTokenCommand: missing arguments")
	}
	s.logger.Debug().Msg("CreateApiTokenCommand: Started")

	var customer models.Customer
	res := s.db.Preload("Secret").First(&customer, s.customerId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("CreateApiTokenCommand: Could not find the customer with id: %s", s.customerId)
		}
		return "", res.Error
	}

	hashed, err := s.hasher.GetHash(fmt.Sprintf("%s:%s", customer.Secret.Value, util.GetRandString(8)))
	if err != nil {
		return "", err
	}

	oneYearLater := time.Now().AddDate(1, 0, 0)

	apiToken := models.ApiToken{
		SourceID:   s.sourceId,
		ValidUntil: oneYearLater,
		Token:      hashed,
	}
	res = s.db.Create(&apiToken)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("CreateApiTokenCommand: Finished with success")
	return apiToken.ID, nil
}
