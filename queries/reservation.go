/*
 * Any operation that does not mutate the database belongs to 'queries'.
 */
package queries

import (
	"errors"
	"fmt"

	"github.com/lghtr35/reservation-engine/models"
	"github.com/rs/zerolog"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type FilterReservationsQuery struct {
	db         *gorm.DB
	logger     *zerolog.Logger
	ids        *[]string
	reserverID *string
	reserveeID *string
	sourceID   *string
	models.Pagination
}

func NewFilterReservationsQuery(db *gorm.DB, logger *zerolog.Logger, ids *[]string, reserveeID, reserverID, sourceID *string, pagination models.Pagination) *FilterReservationsQuery {
	return &FilterReservationsQuery{db: db, logger: logger, ids: ids, reserverID: reserverID, reserveeID: reserveeID, sourceID: sourceID, Pagination: pagination}
}

func (s *FilterReservationsQuery) Execute() (any, error) {
	s.logger.Debug().Msg("FilterReservationsQuery: Started")
	q := s.db.Model(models.Reservation{})
	if s.ids != nil && len(*s.ids) > 0 {
		q = q.Where("id IN ?", *s.ids)
	}
	if s.reserverID != nil && *s.reserverID != "" {
		q = q.Where("reserverId = ?", *s.reserverID)
	}
	if s.reserveeID != nil && *s.reserveeID != "" {
		q = q.Where("reserveeId = ?", *s.reserveeID)
	}
	if s.sourceID != nil && *s.sourceID != "" {
		q = q.Where("sourceId = ?", *s.sourceID)
	}

	offset := s.Pagination.Offset()

	var reservations []models.Reservation
	res := q.Offset(offset).Limit(int(s.Size)).Find(&reservations)
	if res.Error != nil {
		return models.NewPaginationResponse(reservations, 0, 0), res.Error
	}

	var totalCount int64
	res = q.Count(&totalCount)
	if res.Error != nil {
		return models.NewPaginationResponse(reservations, 0, 0), res.Error
	}

	s.logger.Debug().Msg("FilterReservationsQuery: Finished with success")
	return models.NewPaginationResponse(reservations, totalCount, s.Page), nil

}

type ReadReservationQuery struct {
	db     *gorm.DB
	logger *zerolog.Logger
	id     string
}

func NewReadReservationQuery(db *gorm.DB, logger *zerolog.Logger, id string) *ReadReservationQuery {
	return &ReadReservationQuery{db: db, logger: logger, id: id}
}

func (s *ReadReservationQuery) Execute() (any, error) {
	if s.id == "" {
		return models.Reservation{}, errors.New("ReadReservationQuery: Tried to read one with empty id")
	}
	s.logger.Debug().Msg("ReadReservationQuery: ReadOne started")

	var reservation models.Reservation
	res := s.db.Model(models.Reservation{}).Preload(clause.Associations).First(&reservation, s.id)
	if res.Error != nil {
		if res.Error == gorm.ErrRecordNotFound {
			return "", fmt.Errorf("ReadReservationQuery: Could not find the reservation with this id: %s", s.id)
		}
		return "", res.Error
	}

	s.logger.Debug().Msg("ReadReservationQuery: ReadOne finished with success")
	return reservation, nil
}
