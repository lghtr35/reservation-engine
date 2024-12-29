/*
 * Everything involving a mutation belongs to the 'commands' package.
 */
package commands

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
)

const CHECK_IF_INSERT_POSSIBLE_SQL string = `SELECT id FROM reservations r 
WHERE ((r.from  < @from  AND r.to > @from) OR (r.from  < @to  AND r.to > @to)) 
AND r.sourceId = @source
AND (r.reserveeId = @reservee OR r.reserverId = @reserver)`
const CHECK_IF_UPDATE_POSSIBLE_SQL string = `SELECT id FROM reservations r 
WHERE ((r.from  < @from  AND r.to > @from) OR (r.from  < @to  AND r.to > @to)) 
AND r.sourceId = @source
AND (r.reserveeId = @reservee OR r.reserverId = @reserver)
AND r.id != @id`

type CreateReservationCommand struct {
	db         *gorm.DB
	logger     *zerolog.Logger
	from       time.Time
	to         time.Time
	reserverId string
	reserveeId string
	sourceId   string
}

func NewCreateReservationCommand(db *gorm.DB, logger *zerolog.Logger, from time.Time, to time.Time, reserverId, reserveeId string) *CreateReservationCommand {
	return &CreateReservationCommand{db: db, logger: logger, from: from, to: to, reserverId: reserverId, reserveeId: reserveeId}
}

func (s *CreateReservationCommand) Execute() (string, error) {
	if s.reserveeId == "" || s.reserverId == "" {
		return "", errors.New("CreateReservationCommand: Tried creating with empty name")
	}
	s.logger.Debug().Msg("CreateReservationCommand: Started")

	var source models.Source
	res := s.db.First(&source, s.sourceId)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("CreateReservationCommand: Could not find the source with this id: %s", s.sourceId)
		}
		return "", res.Error
	}

	maxDurationForSource, err := time.ParseDuration(source.MaxPossibleDuration)
	if err != nil {
		return "", err
	}

	reservationDuration := s.to.Sub(s.from)

	if maxDurationForSource < reservationDuration {
		return "", errors.New("CreateReservationCommand: Tried creating a reservation longer than maximum for this source")
	}

	var countOfOverlaps int64
	res = s.db.Raw(CHECK_IF_INSERT_POSSIBLE_SQL,
		sql.Named("from", s.from),
		sql.Named("to", s.to),
		sql.Named("source", s.sourceId),
		sql.Named("reservee", s.reserveeId),
		sql.Named("reserver", s.reserverId),
	).Count(&countOfOverlaps)
	if res.Error != nil {
		return "", res.Error
	}

	if countOfOverlaps > 0 {
		return "", errors.New("CreateReservationCommand: Can not create reservation there are overlapping reservations")
	}

	reservation := models.Reservation{
		From:       s.from,
		To:         s.to,
		SourceID:   s.sourceId,
		ReserverID: s.reserverId,
		ReserveeID: s.reserveeId,
	}

	res = s.db.Create(&reservation)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("CreateReservationCommand: Finished with success")

	return reservation.ID, nil
}

type DeleteReservationCommand struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewDeleteReservationCommand(db *gorm.DB, logger *zerolog.Logger, id string) *DeleteReservationCommand {
	return &DeleteReservationCommand{db: db, logger: logger, id: id}
}

func (s *DeleteReservationCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("DeleteReservationCommand: Tried deleting with empty id")
	}
	s.logger.Debug().Msg("DeleteReservationCommand: Started")

	res := s.db.Delete(models.Reservation{}, s.id)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("DeleteReservationCommand: Finished with success")

	return s.id, nil
}

type UpdateReservationCommand struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
	from   *time.Time
	to     *time.Time
}

func NewUpdateReservationCommand(db *gorm.DB, logger *zerolog.Logger, id string, from, to *time.Time) *UpdateReservationCommand {
	return &UpdateReservationCommand{db: db, logger: logger, id: id, from: from, to: to}
}

func (s *UpdateReservationCommand) Execute() (string, error) {
	if s.id == "" {
		return "", errors.New("UpdateReservationCommand: Tried updating with empty id")
	}
	s.logger.Debug().Msg("UpdateReservationCommand: Started")

	var reservation models.Reservation
	res := s.db.First(&reservation, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("UpdateReservationCommand: Could not find the reservation with this id: %s", s.id)
		}
		return "", res.Error
	}

	var source models.Source
	res = s.db.First(&source, reservation.SourceID)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("UpdateReservationCommand: Could not find the source with this id: %s", reservation.SourceID)
		}
		return "", res.Error
	}

	if s.from != nil {
		reservation.From = *s.from
	}
	if s.to != nil {
		reservation.To = *s.to
	}

	maxDurationForSource, err := time.ParseDuration(source.MaxPossibleDuration)
	if err != nil {
		return "", err
	}

	reservationDuration := reservation.To.Sub(reservation.From)

	if maxDurationForSource < reservationDuration {
		return "", errors.New("CreateReservationCommand: Tried creating a reservation longer than maximum for this source")
	}

	var countOfOverlaps int64
	res = s.db.Raw(CHECK_IF_UPDATE_POSSIBLE_SQL,
		sql.Named("from", reservation.From),
		sql.Named("to", reservation.To),
		sql.Named("source", source.ID),
		sql.Named("reservee", reservation.ReserveeID),
		sql.Named("reserver", reservation.ReserverID),
		sql.Named("id", reservation.ID),
	).Count(&countOfOverlaps)
	if res.Error != nil {
		return "", res.Error
	}

	if countOfOverlaps > 0 {
		return "", errors.New("CreateReservationCommand: Can not create reservation there are overlapping reservations")
	}

	res = s.db.Save(&reservation)
	if res.Error != nil {
		return "", res.Error
	}

	s.logger.Debug().Msg("UpdateReservationCommand: Finished with success")

	return s.id, nil
}
